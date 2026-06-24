package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/MUKE-coder/gin-docs/gindocs"
	"github.com/MUKE-coder/gorm-studio/studio"
	"github.com/MUKE-coder/pulse/pulse"
	"github.com/MUKE-coder/sentinel"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mimo-gateway/apps/api/internal/ai"
	"mimo-gateway/apps/api/internal/cache"
	"mimo-gateway/apps/api/internal/config"
	"mimo-gateway/apps/api/internal/handlers"
	"mimo-gateway/apps/api/internal/mail"
	"mimo-gateway/apps/api/internal/middleware"
	"mimo-gateway/apps/api/internal/models"
	"mimo-gateway/apps/api/internal/jobs"
	"mimo-gateway/apps/api/internal/realtime"
	"mimo-gateway/apps/api/internal/services"
	"mimo-gateway/apps/api/internal/storage"
	"mimo-gateway/apps/api/internal/flags"
	"mimo-gateway/apps/api/internal/sync"
	"mimo-gateway/apps/api/internal/webhooks"
)

// Services holds all Phase 4 services for dependency injection.
type Services struct {
	Cache   *cache.Cache
	Storage *storage.Storage
	Mailer  *mail.Mailer
	AI      *ai.AI
	Jobs    *jobs.Client
	// SecObsBridge talks to Sentinel + Pulse over loopback so the
	// in-app Security/Observability dashboards can show summary cards
	// without iframing. Nil when Sentinel/Pulse are both disabled.
	SecObs  *services.SecObsBridge
}

// Setup configures all routes and returns the Gin engine.
func Setup(db *gorm.DB, cfg *config.Config, svc *Services) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global middleware
	r.Use(middleware.Maintenance())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.MaxBodySize(10 << 20)) // 10MB max request body
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(cfg.CORSOrigins))
	r.Use(middleware.Gzip())

	// CSRF defence — only enforces on cookie-authenticated mutations.
	// Bearer (mobile/desktop) flows pass through with no header required.
	// Pairs with services.AuthService.SetAuthCookies (the HttpOnly cookie
	// path documented in /docs/backend/authentication).
	r.Use(middleware.AutoCSRF())

	// Idempotent retries for unsafe methods. Activates only when the client
	// sends an Idempotency-Key header; cached for 24h on 2xx responses.
	r.Use(middleware.Idempotency(svc.Cache))

	// Mount Sentinel security suite (WAF, rate limiting, auth shield, anomaly detection)
	if cfg.SentinelEnabled {
		// In development, use relaxed rate limits so devs don't get blocked while testing
		isDev := cfg.AppEnv == "development"
		ipLimit := &sentinel.Limit{Requests: 100, Window: 1 * time.Minute}
		routeLimits := map[string]sentinel.Limit{
			"/api/auth/login":    {Requests: 5, Window: 15 * time.Minute},
			"/api/auth/register": {Requests: 3, Window: 15 * time.Minute},
		}
		if isDev {
			ipLimit = &sentinel.Limit{Requests: 1000, Window: 1 * time.Minute}
			routeLimits = map[string]sentinel.Limit{
				"/api/auth/login":    {Requests: 100, Window: 1 * time.Minute},
				"/api/auth/register": {Requests: 100, Window: 1 * time.Minute},
			}
		}

		// Sentinel v2.0.1 — use MountE so we can recover gracefully on
		// misconfiguration in dev instead of log.Fatalf-ing the host.
		if err := sentinel.MountE(r, db, sentinel.Config{
			Dashboard: sentinel.DashboardConfig{
				Username:               cfg.SentinelUsername,
				Password:               cfg.SentinelPassword,
				SecretKey:              cfg.SentinelSecretKey,
				// v2.0 refuses default credentials in gin.ReleaseMode;
				// opt-in only for dev so prod can't ship forgeable JWTs.
				AllowInsecureDefaults:  isDev,
			},
			WAF: sentinel.WAFConfig{
				Enabled: true,
				Mode: func() sentinel.WAFMode {
					if isDev { return sentinel.ModeLog }
					return sentinel.ModeBlock
				}(),
				// v2.0 X-Forwarded-For trust closed. Empty list = ignore
				// XFF entirely (the safe default). Operators behind a known
				// reverse proxy should populate via SENTINEL_TRUSTED_PROXIES.
				TrustedProxies:        cfg.SentinelTrustedProxies,
				// 1 MB cap covers richtext admin payloads — Tiptap blog
				// bodies with embedded inline images comfortably exceed
				// the prior 64 KB ceiling. Bump higher if your content
				// embeds large base64 images.
				MaxBodyBytes:          1 * 1024 * 1024,
				RejectOversizedBody:   true,
				// Authenticated admin write endpoints handle their own
				// HTML/richtext payloads via Tiptap. The WAF's XSS detection
				// otherwise flags every <p>/<strong>/<img> tag in a blog
				// body as a payload. These routes still pass through auth
				// + RBAC + binding validation; WAF is just stepped aside
				// for their bodies.
				ExcludeRoutes: []string{
					"/api/blogs",
					"/api/blogs/:id",
					"/api/posts",
					"/api/posts/:id",
					"/api/articles",
					"/api/articles/:id",
					"/api/uploads",
					// v3.31.20 — public form-share submissions. Auth is
					// the share's bcrypt password (optional) and the
					// token itself; Sentinel rate-limits the path.
					"/api/public/forms/:token",
					"/api/public/forms/:token/submit",
				},
			},
			RateLimit: sentinel.RateLimitConfig{
				Enabled: !isDev,
				ByIP:    ipLimit,
				ByRoute: routeLimits,
			},
			AuthShield: sentinel.AuthShieldConfig{
				Enabled:    !isDev,
				LoginRoute: "/api/auth/login",
				// v2.0 CAPTCHA tier sits between soft and hard thresholds.
				// Wire a provider by setting CaptchaProvider in your app code.
			},
			Anomaly: sentinel.AnomalyConfig{Enabled: !isDev},
			Geo:     sentinel.GeoConfig{Enabled: !isDev},
		}); err != nil {
			log.Printf("Warning: Sentinel mount failed: %v", err)
		} else {
			log.Println("Sentinel v2.0 mounted at /sentinel")
		}
	}

	// Mount GORM Studio
	if cfg.GORMStudioEnabled {
		studioCfg := studio.Config{
			Prefix: "/studio",
		}
		if cfg.GORMStudioUsername != "" && cfg.GORMStudioPassword != "" {
			studioCfg.AuthMiddleware = gin.BasicAuth(gin.Accounts{
				cfg.GORMStudioUsername: cfg.GORMStudioPassword,
			})
		}
		studio.Mount(r, db, []interface{}{&models.User{}, &models.Upload{}, &models.Blog{}, /* grit:studio */}, studioCfg)
		log.Println("GORM Studio mounted at /studio")
	}

	// API Documentation (gin-docs — auto-generated from routes + models)
	gindocs.Mount(r, db, gindocs.Config{
		Title:       cfg.AppName + " API",
		Description: "REST API built with [Grit](https://gritframework.dev) — Go + React meta-framework.",
		Version:     "1.0.0",
		UI:          gindocs.UIScalar,
		ScalarTheme: "kepler",
		Models:      []interface{}{&models.User{}, &models.Upload{}, &models.Blog{}},
		Auth: gindocs.AuthConfig{
			Type:         gindocs.AuthBearer,
			BearerFormat: "JWT",
		},
	})
	log.Println("API docs available at /docs")

	// Mount Pulse observability (request tracing, DB monitoring, runtime metrics, error tracking)
	if cfg.PulseEnabled {
		// Pulse v1.0 uses functional options + a context. The context
		// drives clean shutdown of the dashboard's WebSocket + background
		// samplers; we hand it the request context so a server shutdown
		// also unwinds Pulse.
		pulseOpts := []pulse.Option{
			pulse.WithAppName(cfg.AppName),
			pulse.WithCredentials(cfg.PulseUsername, cfg.PulsePassword),
			pulse.WithExcludePaths("/studio/*", "/sentinel/*", "/docs/*", "/pulse/*"),
			pulse.WithPrometheus(),
		}
		if cfg.IsDevelopment() {
			pulseOpts = append(pulseOpts, pulse.WithDevMode())
		}
		// Pulse v1.0 SQLite-backed storage — request/query/error data
		// survives a restart. Stay on the in-memory ring buffer for peak
		// write throughput.
		if cfg.PulseStorage == "sqlite" && cfg.PulseStorageDSN != "" {
			pulseOpts = append(pulseOpts, pulse.WithSQLite(cfg.PulseStorageDSN))
		}
		p := pulse.Mount(context.Background(), r, db, pulseOpts...)

		// Register health checks for connected services
		if svc.Cache != nil {
			p.AddHealthCheck(pulse.HealthCheck{
				Name:     "redis",
				Type:     "redis",
				Critical: false,
				CheckFunc: func(ctx context.Context) error {
					return svc.Cache.Client().Ping(ctx).Err()
				},
			})
		}

		log.Println("Pulse observability mounted at /pulse")
	}

	// Auth service
	authService := &services.AuthService{
		Secret:        cfg.JWTSecret,
		AccessExpiry:  cfg.JWTAccessExpiry,
		RefreshExpiry: cfg.JWTRefreshExpiry,
	}

	// Handlers
	authHandler := &handlers.AuthHandler{
		DB:          db,
		AuthService: authService,
		Config:      cfg,
	}
	userHandler := &handlers.UserHandler{
		DB:          db,
		AuthService: authService,
	}
	uploadHandler := &handlers.UploadHandler{
		DB:      db,
		Storage: svc.Storage,
		Jobs:    svc.Jobs,
	}
	aiHandler := &handlers.AIHandler{
		AI: svc.AI,
	}
	jobsHandler := &handlers.JobsHandler{
		RedisURL: cfg.RedisURL,
	}
	cronHandler := &handlers.CronHandler{}
	blogHandler := handlers.NewBlogHandler(db)
	uiRegistryHandler := handlers.NewUIRegistryHandler(db, cfg.AppURL)
	totpHandler := &handlers.TOTPHandler{
		DB:          db,
		AuthService: authService,
		Issuer:      cfg.TOTPIssuer,
	}
	activityHandler := handlers.NewActivityHandler(db)
	webhookHandler := handlers.NewWebhookHandler(db)
	webhooks.Setup(db)
	realtimeHub := realtime.NewHub()
	flagsEngine := flags.New(db, realtimeHub)
	featureFlagHandler := handlers.NewFeatureFlagHandler(db, flagsEngine)
	realtimeHandler := handlers.NewRealtimeHandler(realtimeHub, authService)
	_ = realtimeHub // available to handlers/services that want to push events

	// In-app Security + Observability dashboards — read from Sentinel/Pulse APIs
	// over loopback. notificationHandler powers the admin bell.
	notificationHandler := &handlers.NotificationHandler{DB: db}
	securityHandler := &handlers.SecurityHandler{Bridge: svc.SecObs}
	observabilityHandler := &handlers.ObservabilityHandler{Bridge: svc.SecObs}

	// v3.30 — semantic activity log + ticket system. Mailer is optional;
	// when nil the ticket handler skips email-out and only writes the row
	// + admin notifications.
	userActivityHandler := &handlers.UserActivityHandler{DB: db}
	ticketHandler := &handlers.TicketHandler{DB: db, Mail: svc.Mailer}
	// v3.31.20 — public form sharing (Phase 2)
	formShareHandler := &handlers.FormShareHandler{DB: db}

	// Sync registry — list every model that should be syncable from
	// offline-first desktop clients. The resource generator injects
	// new resources at the marker below.
	syncRegistry := sync.NewRegistry()
	syncRegistry.Register("users", &models.User{})
	syncRegistry.Register("uploads", &models.Upload{})
	syncRegistry.Register("blogs", &models.Blog{})
	// grit:sync
	syncHandler := handlers.NewSyncHandler(db, syncRegistry)
	// grit:handlers

	// Health check
	// /api/health probes every infrastructure dependency the dashboard's
	// System Health page wants to render. Each probe is bounded by a 500ms
	// timeout so a hung dependency doesn't pile up health requests; failing
	// probes mark themselves down and the overall status downgrades to
	// "degraded" rather than failing the endpoint.
	r.GET("/api/health", func(c *gin.Context) {
		type compStatus struct {
			OK         bool   `json:"ok"`
			LatencyMS  int64  `json:"latency_ms,omitempty"`
			Tables     int    `json:"tables,omitempty"`
			QueueKeys  int    `json:"queue_keys,omitempty"`
			Configured bool   `json:"configured,omitempty"`
			Error      string `json:"error,omitempty"`
		}

		// Database ping + table count. We probe with a 500ms deadline so a
		// blocked write loop can't hang the health check.
		dbStatus := compStatus{OK: true}
		dbStart := time.Now()
		if sqlDB, err := db.DB(); err == nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
			defer cancel()
			if err := sqlDB.PingContext(ctx); err != nil {
				dbStatus.OK = false
				dbStatus.Error = err.Error()
			}
		}
		dbStatus.LatencyMS = time.Since(dbStart).Milliseconds()
		if dbStatus.OK {
			// Best-effort table count — failure is non-fatal and just drops
			// the "tables: N" tooltip on the health card.
			var count int
			db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = current_schema()").Scan(&count)
			dbStatus.Tables = count
		}

		// Redis ping. Reuse the same cache client the rest of the app uses
		// rather than opening a new connection — that way "Redis healthy"
		// on the dashboard means the same Redis the cache + jobs use.
		redisStatus := compStatus{}
		if svc.Cache != nil {
			redisStart := time.Now()
			ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
			defer cancel()
			if err := svc.Cache.Client().Ping(ctx).Err(); err != nil {
				redisStatus.OK = false
				redisStatus.Error = err.Error()
			} else {
				redisStatus.OK = true
			}
			redisStatus.LatencyMS = time.Since(redisStart).Milliseconds()
		}

		// Background-jobs queue — count active asynq keys as a liveness
		// signal. If asynq isn't wired (Jobs == nil), report unconfigured
		// rather than "down" so the dashboard distinguishes the cases.
		jobsStatus := compStatus{}
		if svc.Jobs != nil && svc.Cache != nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
			defer cancel()
			n, err := svc.Cache.Client().Eval(ctx,
				"local total = 0\nfor _, k in ipairs(redis.call('keys', 'asynq:*')) do total = total + 1 end\nreturn total",
				[]string{}).Int()
			if err == nil {
				jobsStatus.OK = true
				jobsStatus.QueueKeys = n
			} else {
				// Fall back to a simple ping so a "no keys yet" install still
				// reports OK rather than down.
				if perr := svc.Cache.Client().Ping(ctx).Err(); perr == nil {
					jobsStatus.OK = true
				}
			}
		}

		// Email is "configured" when Resend key is set + non-default. The
		// dashboard treats unconfigured as "—" not "down".
		mailStatus := compStatus{
			Configured: cfg.ResendAPIKey != "" && cfg.ResendAPIKey != "re_your_api_key",
			OK:         cfg.ResendAPIKey != "" && cfg.ResendAPIKey != "re_your_api_key",
		}

		// Overall status — ok if every wired-up component is up. Components
		// that aren't configured (e.g. Redis off in a single-binary dev
		// run) don't drag the overall status down.
		overall := "ok"
		if !dbStatus.OK || (svc.Cache != nil && !redisStatus.OK) {
			overall = "degraded"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   overall,
			"version":  "0.1.0",
			"database": dbStatus,
			"redis":    redisStatus,
			"api":      compStatus{OK: true},
			"jobs":     jobsStatus,
			"email":    mailStatus,
		})
	})

	// WebSocket: realtime hub. Auth via ?token=<jwt> on the handshake
	// because browsers can't set custom headers on WS upgrade.
	r.GET("/api/ws", realtimeHandler.Connect)

	// Public webhook receiver — no auth on the route itself; each
	// provider's signature verification is the real auth boundary.
	// POST /webhooks/:provider routes to whatever was registered via
	// webhooks.Register(...) at app boot.
	r.POST("/webhooks/:provider", webhookHandler.Receive)

	// Public Grit UI component registry (shadcn-compatible)
	r.GET("/r.json", uiRegistryHandler.GetRegistry)
	r.GET("/r/:name", uiRegistryHandler.GetComponent)

	// Public blog routes (no auth required)
	blogs := r.Group("/api/blogs")
	{
		blogs.GET("", blogHandler.ListPublished)
		blogs.GET("/:slug", blogHandler.GetBySlug)
	}

	// Public auth routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	// OAuth2 social login
	oauth := auth.Group("/oauth")
	{
		oauth.GET("/:provider", authHandler.OAuthBegin)
		oauth.GET("/:provider/callback", authHandler.OAuthCallback)
	}

	// TOTP verification (public — uses pending tokens, not JWT)
	auth.POST("/totp/verify", totpHandler.Verify)
	auth.POST("/totp/backup-codes/verify", totpHandler.VerifyBackupCode)

	// MiMo Proxy routes (OpenAI-compatible)
	mimoProxy := handlers.NewMiMoProxyHandler()
	v1 := r.Group("/v1")
	{
		v1.POST("/chat/completions", mimoProxy.HandleChat)
		v1.POST("/chat/completions/stream", mimoProxy.HandleStreamChat)
		v1.GET("/models", mimoProxy.HandleModels)
	}
	r.GET("/health", mimoProxy.HandleHealth)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.Auth(db, authService))
	// Activity logger writes one row per successful authenticated mutation.
	// Records who/what/when/where for audit. Read-only — see admin/activity.
	protected.Use(middleware.ActivityLogger(db))
	{
		protected.GET("/auth/me", authHandler.Me)
		protected.POST("/auth/logout", authHandler.Logout)

		// Two-Factor Authentication (TOTP)
		protected.POST("/auth/totp/setup", totpHandler.Setup)
		protected.POST("/auth/totp/enable", totpHandler.Enable)
		protected.POST("/auth/totp/disable", totpHandler.Disable)
		protected.GET("/auth/totp/status", totpHandler.Status)
		protected.POST("/auth/totp/backup-codes", totpHandler.RegenerateBackupCodes)
		protected.DELETE("/auth/totp/trusted-devices", totpHandler.RevokeTrustedDevices)

		// User routes (authenticated)
		protected.GET("/users/:id", userHandler.GetByID)

		// File uploads
		protected.POST("/uploads", uploadHandler.Create)
		protected.POST("/uploads/presign", uploadHandler.Presign)
		protected.POST("/uploads/complete", uploadHandler.CompleteUpload)
		protected.GET("/uploads", uploadHandler.List)
		protected.GET("/uploads/stats", uploadHandler.Stats)
		protected.GET("/uploads/:id", uploadHandler.GetByID)
		protected.DELETE("/uploads/:id", uploadHandler.Delete)

		// Offline-first sync — desktop clients call these to flush their
		// local outbox and pull server-side updates.
		protected.POST("/sync/push", syncHandler.Push)
		protected.GET("/sync/pull", syncHandler.Pull)

		// AI
		protected.POST("/ai/complete", aiHandler.Complete)
		protected.POST("/ai/chat", aiHandler.Chat)
		protected.POST("/ai/stream", aiHandler.Stream)

		// Grit UI component registry (authenticated browse)
		protected.GET("/ui-components", uiRegistryHandler.ListComponents)
		protected.GET("/ui-components/:name", uiRegistryHandler.GetComponentDetail)

		// In-app notification bell — every authenticated user. Pulls
		// from a single Notification table that the SecObs poller
		// writes into when Sentinel/Pulse fires a high-severity event.
		protected.GET("/notifications", notificationHandler.List)
		protected.POST("/notifications/:id/read", notificationHandler.MarkRead)
		protected.POST("/notifications/read-all", notificationHandler.MarkAllRead)

		// v3.30 — tickets. Any authenticated user can open + reply; the
		// handler scopes List/Get visibility to the caller unless they're
		// ADMIN/EDITOR (then they see the full queue).
		protected.POST("/tickets", ticketHandler.Create)
		protected.GET("/tickets", ticketHandler.List)
		protected.GET("/tickets/:id", ticketHandler.Get)
		protected.POST("/tickets/:id/reply", ticketHandler.Reply)
		protected.PATCH("/tickets/:id/close", ticketHandler.Close)
		protected.PATCH("/tickets/:id/reopen", ticketHandler.Reopen)
		protected.PATCH("/tickets/:id/assign", ticketHandler.Assign) // admin-gated inside the handler

		// grit:routes:protected
	}

	// Profile routes (any authenticated user)
	profile := protected.Group("/profile")
	{
		profile.GET("", userHandler.GetProfile)
		profile.PUT("", userHandler.UpdateProfile)
		profile.DELETE("", userHandler.DeleteProfile)
	}

	// Admin routes
	admin := r.Group("/api")
	admin.Use(middleware.Auth(db, authService))
	admin.Use(middleware.RequireRole("ADMIN"))
	{
		admin.GET("/users", userHandler.List)
		admin.POST("/users", userHandler.Create)
		admin.PUT("/users/:id", userHandler.Update)
		admin.DELETE("/users/:id", userHandler.Delete)

		// Activity audit log + tamper-evident chain verification
		admin.GET("/admin/activity", activityHandler.List)
		admin.GET("/admin/activity/integrity", activityHandler.VerifyIntegrity)

		// v3.30 — semantic user activity dashboard (action + IP + severity).
		// Separate from /admin/activity above which is the HTTP audit log.
		admin.GET("/user-activity", userActivityHandler.List)
		admin.GET("/user-activity/stats", userActivityHandler.Stats)

		// Webhook receiver admin (review + replay failed events)
		admin.GET("/admin/webhooks", webhookHandler.List)
		admin.POST("/admin/webhooks/:id/replay", webhookHandler.Replay)

		// Feature flags + A/B testing
		admin.GET("/admin/flags", featureFlagHandler.List)
		admin.POST("/admin/flags", featureFlagHandler.Create)
		admin.PUT("/admin/flags/:id", featureFlagHandler.Update)
		admin.DELETE("/admin/flags/:id", featureFlagHandler.Delete)
		admin.GET("/admin/flags/:id/exposures", featureFlagHandler.Exposures)

		// Admin system routes
		admin.GET("/admin/jobs/stats", jobsHandler.Stats)
		admin.GET("/admin/jobs/:status", jobsHandler.ListByStatus)
		admin.POST("/admin/jobs/:id/retry", jobsHandler.Retry)
		admin.DELETE("/admin/jobs/queue/:queue", jobsHandler.ClearQueue)
		admin.GET("/admin/cron/tasks", cronHandler.ListTasks)

		// Blog management (admin)
		admin.GET("/admin/blogs", blogHandler.List)
		admin.GET("/admin/blogs/:id", blogHandler.GetByID)
		admin.POST("/admin/blogs", blogHandler.Create)
		admin.PUT("/admin/blogs/:id", blogHandler.Update)
		admin.DELETE("/admin/blogs/:id", blogHandler.Delete)

		// Grit UI component registry (admin management)
		admin.POST("/admin/ui-components", uiRegistryHandler.CreateComponent)
		admin.PUT("/admin/ui-components/:name", uiRegistryHandler.UpdateComponent)
		admin.DELETE("/admin/ui-components/:name", uiRegistryHandler.DeleteComponent)

		// In-app Security dashboard — aggregates Sentinel APIs into one
		// envelope so the React page does a single round-trip. Operators
		// who want to dig deeper open /sentinel/ui directly.
		admin.GET("/admin/security/summary", securityHandler.Summary)
		// In-app Observability dashboard — same pattern against Pulse.
		// Operators who want a flame graph or the full SLO timeline open
		// /pulse/ui directly.
		admin.GET("/admin/observability/summary", observabilityHandler.Summary)

		// v3.31.20 — public form sharing admin
		admin.GET("/admin/form-shares", formShareHandler.List)
		admin.POST("/admin/form-shares", formShareHandler.Create)
		admin.PATCH("/admin/form-shares/:id", formShareHandler.Update)
		admin.DELETE("/admin/form-shares/:id", formShareHandler.Delete)
		// v3.31.25 — audit log of public submissions
		admin.GET("/admin/form-submissions", formShareHandler.ListSubmissions)

		// grit:routes:admin
	}

	// Public form-sharing endpoints. NO auth, NO CSRF — Sentinel rate
	// limits each token aggressively. The dispatch service is the
	// security boundary (whitelists which resources are reachable).
	publicForms := r.Group("/api/public/forms")
	{
		publicForms.GET("/:token", formShareHandler.PublicGet)
		publicForms.POST("/:token/submit", formShareHandler.PublicSubmit)
	}

	// Custom role-restricted routes
	// grit:routes:custom

	return r
}
