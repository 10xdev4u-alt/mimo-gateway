package handlers

import (
	"fmt"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mimo-gateway/apps/api/internal/files"
	"mimo-gateway/apps/api/internal/jobs"
	"mimo-gateway/apps/api/internal/models"
	"mimo-gateway/apps/api/internal/storage"
)

// AllowedMimeTypes defines which file types can be uploaded.
var AllowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	"video/mp4":       true,
	"video/webm":      true,
	"video/quicktime": true,
	"application/pdf": true,
	"text/plain":      true,
	"text/csv":        true,
	"application/json": true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}

// MaxUploadSize is the maximum file size (50 MB).
const MaxUploadSize = 50 << 20

// UploadHandler handles file upload endpoints.
type UploadHandler struct {
	DB      *gorm.DB
	Storage *storage.Storage
	Jobs    *jobs.Client
}

// Create handles file upload via multipart form.
//
// Query params (v3.31.30):
//   accepts   — comma-separated list of CLI accept aliases
//               (image, video, pdf, doc, excel, csv, zip, archive, all).
//               When present, validates the upload's MIME against the
//               alias set. Absent = fall back to the global allowlist.
//   max_size  — per-field byte cap. Overrides MaxUploadSize when set
//               (e.g. video fields raise it to 300MB).
//
// Response: a files.FileRef directly under data so the frontend can
// store it verbatim in form state, no shape massaging needed.
func (h *UploadHandler) Create(c *gin.Context) {
	if h.Storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{
				"code":    "STORAGE_UNAVAILABLE",
				"message": "File storage is not configured",
			},
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_FILE",
				"message": "No file provided",
			},
		})
		return
	}
	defer file.Close()

	// Per-field accept list. Comma-separated aliases.
	var acceptsList []string
	if a := c.Query("accepts"); a != "" {
		for _, s := range strings.Split(a, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				acceptsList = append(acceptsList, s)
			}
		}
	}

	// Per-field max size override. Bytes.
	maxSize := int64(MaxUploadSize)
	if m := c.Query("max_size"); m != "" {
		if parsed, perr := strconv.ParseInt(m, 10, 64); perr == nil && parsed > 0 {
			maxSize = parsed
		}
	} else if len(acceptsList) > 0 {
		// No explicit max_size, but field type is known — use the
		// default-for-accepts (5MB for most, 300MB for video).
		maxSize = files.DefaultMaxSizeBytes(acceptsList)
	}

	if header.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "FILE_TOO_LARGE",
				"message": fmt.Sprintf("File size exceeds maximum of %d MB", maxSize/(1<<20)),
			},
		})
		return
	}

	mimeType := header.Header.Get("Content-Type")

	// If accepts was provided, validate against the per-field allow set.
	// Otherwise fall back to the global allowlist (backwards-compat).
	allowed := false
	if len(acceptsList) > 0 {
		allowed = files.AllowsMIME(acceptsList, mimeType)
	} else {
		allowed = AllowedMimeTypes[mimeType]
	}
	if !allowed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_FILE_TYPE",
				"message": "File type not allowed",
			},
		})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), strings.TrimSuffix(filepath.Base(header.Filename), ext), ext)
	key := fmt.Sprintf("uploads/%s/%s", time.Now().Format("2006/01"), filename)

	// Upload to storage
	if err := h.Storage.Upload(c.Request.Context(), key, file, mimeType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "UPLOAD_FAILED",
				"message": "Failed to upload file",
			},
		})
		return
	}

	userID, _ := c.Get("user_id")

	upload := models.Upload{
		Filename:     filename,
		OriginalName: header.Filename,
		MimeType:     mimeType,
		Size:         header.Size,
		Path:         key,
		URL:          h.Storage.GetURL(key),
		UserID:       userID.(string),
	}

	if err := h.DB.Create(&upload).Error; err != nil {
		_ = h.Storage.Delete(c.Request.Context(), key)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to save upload record",
			},
		})
		return
	}

	// Enqueue image processing job. Width / height are written back to
	// the upload row by the worker; for now we return what we have and
	// the frontend can refetch the FileRef later if it needs dimensions.
	if h.Jobs != nil && storage.IsImageMimeType(mimeType) {
		_ = h.Jobs.EnqueueProcessImage(c.Request.Context(), upload.ID, key, mimeType, jobs.EnqueueOption{
			IdempotencyKey: "image:process:" + upload.ID,
		})
	}

	// Dimensions / duration aren't extracted synchronously -- the
	// image-processing worker populates ThumbnailURL asynchronously
	// and the frontend can re-fetch the record if it needs them later.
	ref := files.FileRef{
		URL:          upload.URL,
		Key:          upload.Path,
		Name:         upload.OriginalName,
		MIME:         upload.MimeType,
		Size:         upload.Size,
		ThumbnailURL: upload.ThumbnailURL,
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    ref,
		"message": "File uploaded successfully",
	})
}

// Stats returns aggregate storage usage across the uploads table.
// Surfaces total count, total bytes, and a per-kind breakdown
// (image / video / audio / document / other) so the storage admin
// page can show usage at a glance. v3.31.32.
func (h *UploadHandler) Stats(c *gin.Context) {
	type kindRow struct {
		Kind  string `gorm:"column:kind" json:"kind"`
		Count int64  `gorm:"column:count" json:"count"`
		Size  int64  `gorm:"column:size" json:"size"`
	}

	var total int64
	if err := h.DB.Model(&models.Upload{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to compute stats"},
		})
		return
	}

	var totalSize int64
	h.DB.Model(&models.Upload{}).Select("COALESCE(SUM(size), 0)").Scan(&totalSize)

	// Bucket by MIME kind. SUBSTR + CASE in raw SQL keeps this a single
	// scan regardless of DB engine (works on Postgres + SQLite).
	rows := []kindRow{}
	bucketExpr := `CASE
		WHEN mime_type LIKE 'image/%' THEN 'image'
		WHEN mime_type LIKE 'video/%' THEN 'video'
		WHEN mime_type LIKE 'audio/%' THEN 'audio'
		WHEN mime_type = 'application/pdf' THEN 'pdf'
		WHEN mime_type LIKE '%spreadsheet%' OR mime_type LIKE '%excel%' OR mime_type = 'text/csv' THEN 'spreadsheet'
		WHEN mime_type LIKE '%wordprocessing%' OR mime_type = 'application/msword' THEN 'document'
		ELSE 'other'
	END`
	h.DB.Model(&models.Upload{}).
		Select(bucketExpr+" AS kind, COUNT(*) AS count, COALESCE(SUM(size), 0) AS size").
		Group("kind").
		Scan(&rows)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total_count": total,
			"total_size":  totalSize,
			"by_kind":     rows,
		},
	})
}

// List returns a paginated list of uploads.
func (h *UploadHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := h.DB.Model(&models.Upload{})

	// Filter by MIME type
	if mimeType := c.Query("mime_type"); mimeType != "" {
		query = query.Where("mime_type LIKE ?", mimeType+"%")
	}

	var total int64
	query.Count(&total)

	var uploads []models.Upload
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&uploads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch uploads",
			},
		})
		return
	}

	pages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{
		"data": uploads,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     pages,
		},
	})
}

// GetByID returns a single upload by ID.
func (h *UploadHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	var upload models.Upload
	if err := h.DB.First(&upload, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Upload not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": upload,
	})
}

// Delete removes an upload and its stored file.
func (h *UploadHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	var upload models.Upload
	if err := h.DB.First(&upload, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Upload not found",
			},
		})
		return
	}

	// Delete from storage
	if h.Storage != nil {
		_ = h.Storage.Delete(c.Request.Context(), upload.Path)
		// Also delete thumbnail if it exists
		if upload.ThumbnailURL != "" {
			thumbKey := strings.Replace(upload.Path, "uploads/", "thumbnails/", 1)
			_ = h.Storage.Delete(c.Request.Context(), thumbKey)
		}
	}

	if err := h.DB.Delete(&upload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to delete upload",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Upload deleted successfully",
	})
}

// Presign generates a presigned PUT URL for direct browser-to-storage upload.
func (h *UploadHandler) Presign(c *gin.Context) {
	if h.Storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{"code": "STORAGE_UNAVAILABLE", "message": "File storage is not configured"},
		})
		return
	}

	var req struct {
		Filename    string `json:"filename" binding:"required"`
		ContentType string `json:"content_type" binding:"required"`
		FileSize    int64  `json:"file_size" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	if !AllowedMimeTypes[req.ContentType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_FILE_TYPE", "message": "File type not allowed"},
		})
		return
	}

	if req.FileSize > MaxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "FILE_TOO_LARGE", "message": fmt.Sprintf("File size exceeds maximum of %d MB", MaxUploadSize/(1<<20))},
		})
		return
	}

	ext := filepath.Ext(req.Filename)
	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), strings.TrimSuffix(filepath.Base(req.Filename), ext), ext)
	key := fmt.Sprintf("uploads/%s/%s", time.Now().Format("2006/01"), filename)

	presignedURL, err := h.Storage.PresignPutURL(c.Request.Context(), key, req.ContentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "PRESIGN_FAILED", "message": "Failed to generate upload URL"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"presigned_url": presignedURL,
			"key":           key,
			"public_url":    h.Storage.GetURL(key),
		},
	})
}

// CompleteUpload records a file that was uploaded directly to storage via presigned URL.
func (h *UploadHandler) CompleteUpload(c *gin.Context) {
	var req struct {
		Key         string `json:"key" binding:"required"`
		Filename    string `json:"filename" binding:"required"`
		ContentType string `json:"content_type" binding:"required"`
		Size        int64  `json:"size" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	userID, _ := c.Get("user_id")

	upload := models.Upload{
		Filename:     filepath.Base(req.Key),
		OriginalName: req.Filename,
		MimeType:     req.ContentType,
		Size:         req.Size,
		Path:         req.Key,
		URL:          h.Storage.GetURL(req.Key),
		UserID:       userID.(string),
	}

	if err := h.DB.Create(&upload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to save upload record"},
		})
		return
	}

	// Enqueue image processing job if it's an image.
	// IdempotencyKey = upload.ID so a client retry of the same upload
	// (rare but possible after a network drop) doesn't re-process.
	if h.Jobs != nil && storage.IsImageMimeType(req.ContentType) {
		_ = h.Jobs.EnqueueProcessImage(c.Request.Context(), upload.ID, req.Key, req.ContentType, jobs.EnqueueOption{
			IdempotencyKey: "image:process:" + upload.ID,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    upload,
		"message": "Upload recorded successfully",
	})
}
