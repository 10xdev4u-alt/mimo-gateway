package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"mimo-gateway/apps/api/internal/models"
	"mimo-gateway/apps/api/internal/services"
)

type FormShareHandler struct {
	DB *gorm.DB
}

// ── Admin endpoints ────────────────────────────────────────────────────

// List paginates form shares for the admin dashboard.
func (h *FormShareHandler) List(c *gin.Context) {
	var shares []models.FormShare
	q := h.DB.Order("created_at DESC")
	if rn := c.Query("resource_name"); rn != "" {
		q = q.Where("resource_name = ?", rn)
	}
	if err := q.Find(&shares).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": shares,
		"meta": gin.H{"total": len(shares), "page": 1, "page_size": len(shares), "pages": 1},
	})
}

// Create generates a new share for a resource. Optional password is
// bcrypt-hashed before storage. A 32-char URL-safe token is generated
// automatically.
func (h *FormShareHandler) Create(c *gin.Context) {
	var req struct {
		ResourceName string `json:"resource_name" binding:"required"`
		Label        string `json:"label"`
		Password     string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	token, err := randomToken(24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "token generation failed"},
		})
		return
	}

	share := models.FormShare{
		ResourceName: req.ResourceName,
		Token:        token,
		Label:        req.Label,
		Enabled:      true,
	}
	if userID, ok := c.Get("user_id"); ok {
		if s, ok := userID.(string); ok {
			share.CreatedByUserID = s
		}
	}
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "INTERNAL_ERROR", "message": "password hash failed"},
			})
			return
		}
		share.PasswordHash = string(hash)
	}

	if err := h.DB.Create(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	share.HasPassword = share.PasswordHash != ""
	c.JSON(http.StatusCreated, gin.H{"data": share, "message": "Share created"})
}

// Update toggles enabled/label/password. Pass password="" to leave
// unchanged; pass password="-" to remove an existing password.
func (h *FormShareHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var share models.FormShare
	if err := h.DB.First(&share, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "Share not found"},
		})
		return
	}

	var req struct {
		Label    *string `json:"label"`
		Enabled  *bool   `json:"enabled"`
		Password *string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}
	if req.Label != nil {
		share.Label = *req.Label
	}
	if req.Enabled != nil {
		share.Enabled = *req.Enabled
	}
	if req.Password != nil {
		if *req.Password == "-" {
			share.PasswordHash = ""
		} else if *req.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{"code": "INTERNAL_ERROR", "message": "password hash failed"},
				})
				return
			}
			share.PasswordHash = string(hash)
		}
	}
	if err := h.DB.Save(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	share.HasPassword = share.PasswordHash != ""
	c.JSON(http.StatusOK, gin.H{"data": share, "message": "Share updated"})
}

// Delete soft-deletes a share (token stops working).
func (h *FormShareHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.FormShare{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Share deleted"})
}

// ── Public endpoints ───────────────────────────────────────────────────

// PublicGet returns the resource name + has_password so the public
// web page can decide whether to show a password gate. Does NOT expose
// the hash, the operator label, or submission stats.
func (h *FormShareHandler) PublicGet(c *gin.Context) {
	token := c.Param("token")
	var share models.FormShare
	if err := h.DB.First(&share, "token = ? AND enabled = ?", token, true).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "Link not found or disabled"},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"resource_name": share.ResourceName,
			"has_password":  share.PasswordHash != "",
			"label":         share.Label,
		},
	})
}

// PublicSubmit accepts the form payload, verifies the password (when
// required), and dispatches to the resource's create service. Returns
// the new record's ID + an opaque label.
func (h *FormShareHandler) PublicSubmit(c *gin.Context) {
	token := c.Param("token")
	var share models.FormShare
	if err := h.DB.First(&share, "token = ? AND enabled = ?", token, true).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "Link not found or disabled"},
		})
		return
	}

	var body struct {
		Password string                 `json:"_password"`
		Fields   map[string]interface{} `json:"fields" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	if share.PasswordHash != "" {
		if body.Password == "" || bcrypt.CompareHashAndPassword([]byte(share.PasswordHash), []byte(body.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{"code": "PASSWORD_REQUIRED", "message": "Password incorrect or missing"},
			})
			return
		}
	}

	out, err := services.SubmitSharedForm(h.DB, share.ResourceName, body.Fields)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "SUBMISSION_FAILED", "message": err.Error()},
		})
		return
	}

	// Bump submission count (best-effort — failure here doesn't
	// retroactively invalidate the user's submission).
	h.DB.Model(&share).UpdateColumns(map[string]interface{}{
		"submission_count": share.SubmissionCount + 1,
		"updated_at":       time.Now(),
	})

	// v3.31.25 — write the audit row. Best-effort; failure here means
	// the visitor still gets their record, the admin just misses one
	// line in the trail. We truncate UA at 500 chars (column width).
	ua := c.GetHeader("User-Agent")
	if len(ua) > 500 {
		ua = ua[:500]
	}
	_ = h.DB.Create(&models.FormSubmission{
		ShareID:      share.ID,
		ResourceName: share.ResourceName,
		RecordID:     out.ID,
		IP:           c.ClientIP(),
		UserAgent:    ua,
	}).Error

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id":    out.ID,
			"label": out.Label,
		},
		"message": "Submitted",
	})
}

// ListSubmissions returns audit-log rows for one or more shares.
// Filterable by share_id and resource_name; defaults to the 100 most
// recent rows. v3.31.25.
func (h *FormShareHandler) ListSubmissions(c *gin.Context) {
	var rows []models.FormSubmission
	q := h.DB.Order("created_at DESC").Limit(100)

	if shareID := c.Query("share_id"); shareID != "" {
		q = q.Where("share_id = ?", shareID)
	}
	if resourceName := c.Query("resource_name"); resourceName != "" {
		q = q.Where("resource_name = ?", resourceName)
	}

	if err := q.Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": rows,
		"meta": gin.H{"total": len(rows), "page": 1, "page_size": len(rows), "pages": 1},
	})
}

// randomToken returns a URL-safe base64 string of about `byteLen*4/3`
// characters. 24 bytes → ~32 chars, plenty of entropy against brute
// force at the public endpoint (paired with Sentinel rate limits).
func randomToken(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "="), nil
}
