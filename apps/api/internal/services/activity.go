package services

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mimo-gateway/apps/api/internal/models"
)

// ActivityService centralises the "log a semantic event" call so
// handlers don't repeat IP-pull / user-agent / actor lookup boilerplate.
//
// Usage from a handler:
//
//	services.LogActivity(h.DB, c, services.ActivityArgs{
//	    Action:       "ticket.create",
//	    Severity:     "info",
//	    Summary:      fmt.Sprintf("Opened ticket %q", ticket.Subject),
//	    ResourceType: "ticket",
//	    ResourceID:   ticket.ID,
//	})
//
// Errors are logged, not returned — losing an audit row should never
// fail a real request. If you care about guaranteed delivery, queue
// these via asynq instead of writing inline.
type ActivityArgs struct {
	// UserID overrides the actor when set. Auth handlers (login, register)
	// pass it explicitly because the auth middleware hasn't yet populated
	// the gin context with "user_id" — without this override every auth
	// event would be attributed to the empty system actor.
	UserID       string
	Action       string
	Severity     string                 // info | warn | critical
	Summary      string
	ResourceType string
	ResourceID   string
	Metadata     map[string]interface{} // optional JSON-encodable extras
}

// LogActivity writes a UserActivity row. Picks actor + IP + user-agent
// from the request context automatically, falling back to args.UserID
// when the caller is in an unauthenticated handler (auth flows).
func LogActivity(db *gorm.DB, c *gin.Context, args ActivityArgs) {
	userID := args.UserID
	if userID == "" {
		if v, ok := c.Get("user_id"); ok {
			if s, ok := v.(string); ok {
				userID = s
			}
		}
	}

	var metaJSON string
	if args.Metadata != nil {
		if b, err := json.Marshal(args.Metadata); err == nil {
			metaJSON = string(b)
		}
	}

	row := models.UserActivity{
		UserID:       userID,
		Action:       args.Action,
		Severity:     args.Severity,
		Summary:      args.Summary,
		ResourceType: args.ResourceType,
		ResourceID:   args.ResourceID,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.GetHeader("User-Agent"),
		Metadata:     metaJSON,
	}

	if err := db.Create(&row).Error; err != nil {
		// Audit failures are non-fatal but worth knowing about — log and
		// keep moving. In production, wire a metric to alert on a sudden
		// surge in these (suggests DB write pressure).
		log.Printf("activity: failed to write %s: %v", args.Action, err)
	}
}

// Convenience helpers for the most common events. Use these in auth
// handlers + middleware so the dotted action names stay consistent.

func LogLogin(db *gorm.DB, c *gin.Context, userID, email string) {
	LogActivity(db, c, ActivityArgs{
		UserID:       userID,
		Action:       "auth.login",
		Severity:     "info",
		Summary:      email + " signed in",
		ResourceType: "user",
		ResourceID:   userID,
	})
}

// LogLoginFailed intentionally leaves UserID empty — the failed attempt
// might be an unknown email or a real account being brute-forced; either
// way, attributing it to a specific actor is misleading. The summary
// captures the attempted email for audit.
func LogLoginFailed(db *gorm.DB, c *gin.Context, email string) {
	LogActivity(db, c, ActivityArgs{
		Action:   "auth.login_failed",
		Severity: "warn",
		Summary:  "Failed sign-in attempt for " + email,
	})
}

func LogLogout(db *gorm.DB, c *gin.Context, userID, email string) {
	LogActivity(db, c, ActivityArgs{
		UserID:   userID,
		Action:   "auth.logout",
		Severity: "info",
		Summary:  email + " signed out",
		ResourceType: "user", ResourceID: userID,
	})
}

func LogRegister(db *gorm.DB, c *gin.Context, userID, email string) {
	LogActivity(db, c, ActivityArgs{
		UserID:       userID,
		Action:       "auth.register",
		Severity:     "info",
		Summary:      email + " created an account",
		ResourceType: "user",
		ResourceID:   userID,
	})
}
