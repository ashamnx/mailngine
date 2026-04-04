package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hibiken/asynq"
	"github.com/hellomail/hellomail/internal/api/handler"
	"github.com/hellomail/hellomail/internal/api/middleware"
	"github.com/hellomail/hellomail/internal/auth"
	"github.com/hellomail/hellomail/internal/config"
	"github.com/hellomail/hellomail/internal/domain"
	"github.com/hellomail/hellomail/internal/email"
	"github.com/hellomail/hellomail/internal/analytics"
	"github.com/hellomail/hellomail/internal/audit"
	"github.com/hellomail/hellomail/internal/billing"
	"github.com/hellomail/hellomail/internal/inbox"
	"github.com/hellomail/hellomail/internal/suppression"
	"github.com/hellomail/hellomail/internal/team"
	hmtemplate "github.com/hellomail/hellomail/internal/template"
	"github.com/hellomail/hellomail/internal/webhook"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func NewRouter(cfg *config.Config, db *pgxpool.Pool, cache *redis.Client, asynqClient *asynq.Client, logger zerolog.Logger) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.RealIP)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.MaxBodySize(10 << 20))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging(logger))
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID", "X-Idempotency-Key"},
		ExposedHeaders:   []string{"X-Request-ID", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Auth dependencies
	jwtMgr := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiry)
	googleOAuth := auth.NewGoogleOAuth(cfg.Google.ClientID, cfg.Google.ClientSecret, cfg.Google.RedirectURL, cache)

	// Services
	domainSvc := domain.NewService(db, logger)
	emailSvc := email.NewService(db, asynqClient, logger)
	inboxSvc := inbox.NewService(db, logger)
	suppressionSvc := suppression.NewService(db, cache, logger)
	webhookSvc := webhook.NewService(db, asynqClient, logger)
	auditSvc := audit.NewService(db, logger)
	analyticsSvc := analytics.NewService(db, cache, logger)
	tracker := analytics.NewTracker(db, cache, logger)
	templateSvc := hmtemplate.NewService(db, logger)
	billingSvc := billing.NewService(db, logger)
	teamSvc := team.NewService(db, logger)

	// Handlers
	healthHandler := handler.NewHealthHandler(db, cache)
	authHandler := handler.NewAuthHandler(cfg, db, cache, googleOAuth, jwtMgr)
	apiKeyHandler := handler.NewAPIKeyHandler(db)
	domainHandler := handler.NewDomainHandler(domainSvc)
	emailHandler := handler.NewEmailHandler(emailSvc)
	inboxHandler := handler.NewInboxHandler(inboxSvc)
	suppressionHandler := handler.NewSuppressionHandler(suppressionSvc)
	webhookHandler := handler.NewWebhookHandler(webhookSvc)
	auditHandler := handler.NewAuditHandler(auditSvc)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsSvc)
	templateHandler := handler.NewTemplateHandler(templateSvc)
	billingHandler := handler.NewBillingHandler(billingSvc)
	teamHandler := handler.NewTeamHandler(teamSvc)

	// Middleware
	authenticate := middleware.Authenticate(jwtMgr, db, cache)
	rateLimit := middleware.RateLimit(cache, middleware.RateLimitConfig{
		APIKeyRate:  100,
		SessionRate: 60,
		Window:      1 * time.Second,
	})

	// Health & metrics (no auth)
	r.Get("/health", healthHandler.Check)
	r.Handle("/metrics", promhttp.Handler())

	// Tracking endpoints (public, no auth)
	r.Get("/t/o/{id}", tracker.TrackOpen)
	r.Get("/t/c/{id}", tracker.TrackClick)

	// API v1
	r.Route("/v1", func(r chi.Router) {
		// Auth routes (no auth required)
		r.Route("/auth", func(r chi.Router) {
			r.Get("/google", authHandler.GoogleRedirect)
			r.Get("/google/callback", authHandler.GoogleCallback)

			// Authenticated auth routes
			r.Group(func(r chi.Router) {
				r.Use(authenticate)
				r.Post("/logout", authHandler.Logout)
				r.Get("/me", authHandler.Me)
			})
		})

		// Protected routes (auth required)
		r.Group(func(r chi.Router) {
			r.Use(authenticate)
			r.Use(rateLimit)

			// API Keys
			r.Route("/api-keys", func(r chi.Router) {
				r.Use(middleware.RequireRole("member"))
				r.Post("/", apiKeyHandler.Create)
				r.Get("/", apiKeyHandler.List)
				r.Delete("/{id}", apiKeyHandler.Revoke)
			})

			// Domains
			r.Route("/domains", func(r chi.Router) {
				r.Use(middleware.RequireRole("member"))
				r.Post("/analyze", domainHandler.Analyze)
				r.Post("/", domainHandler.Create)
				r.Get("/", domainHandler.List)
				r.Get("/{id}", domainHandler.Get)
				r.Patch("/{id}", domainHandler.Update)
				r.Delete("/{id}", domainHandler.Delete)
				r.Post("/{id}/verify", domainHandler.Verify)
				r.Get("/{id}/connect-url", domainHandler.GetConnectURL)
			})

			// Emails
			r.Route("/emails", func(r chi.Router) {
				r.Post("/", emailHandler.Send)
				r.Get("/", emailHandler.List)
				r.Get("/{id}", emailHandler.Get)
			})

			// Inbox
			r.Route("/inbox", func(r chi.Router) {
				r.Get("/threads", inboxHandler.ListThreads)
				r.Get("/threads/{id}", inboxHandler.GetThread)
				r.Delete("/threads/{id}", inboxHandler.DeleteThread)

				r.Get("/messages/{id}", inboxHandler.GetMessage)
				r.Patch("/messages/{id}", inboxHandler.UpdateMessage)
				r.Delete("/messages/{id}", inboxHandler.DeleteMessage)
				r.Post("/messages/{id}/labels", inboxHandler.AddLabel)
				r.Delete("/messages/{id}/labels/{labelId}", inboxHandler.RemoveLabel)

				r.Get("/labels", inboxHandler.ListLabels)
				r.Post("/labels", inboxHandler.CreateLabel)
				r.Delete("/labels/{id}", inboxHandler.DeleteLabel)

				r.Get("/search", inboxHandler.SearchMessages)
			})

			// Suppressions
			r.Route("/suppressions", func(r chi.Router) {
				r.Get("/", suppressionHandler.List)
				r.Post("/", suppressionHandler.Create)
				r.Delete("/{id}", suppressionHandler.Delete)
			})

			// Webhooks
			r.Route("/webhooks", func(r chi.Router) {
				r.Use(middleware.RequireRole("admin"))
				r.Post("/", webhookHandler.Create)
				r.Get("/", webhookHandler.List)
				r.Get("/{id}", webhookHandler.Get)
				r.Patch("/{id}", webhookHandler.Update)
				r.Delete("/{id}", webhookHandler.Delete)
				r.Get("/{id}/deliveries", webhookHandler.ListDeliveries)
			})

			// Analytics
			r.Route("/analytics", func(r chi.Router) {
				r.Get("/overview", analyticsHandler.Overview)
				r.Get("/timeseries", analyticsHandler.Timeseries)
				r.Get("/events", analyticsHandler.Events)
			})

			// Audit Logs
			r.Route("/audit-logs", func(r chi.Router) {
				r.Use(middleware.RequireRole("admin"))
				r.Get("/", auditHandler.List)
				r.Get("/{id}", auditHandler.Get)
			})

			// Templates
			r.Route("/templates", func(r chi.Router) {
				r.Use(middleware.RequireRole("member"))
				r.Post("/", templateHandler.Create)
				r.Get("/", templateHandler.List)
				r.Get("/{id}", templateHandler.Get)
				r.Patch("/{id}", templateHandler.Update)
				r.Delete("/{id}", templateHandler.Delete)
				r.Post("/{id}/preview", templateHandler.Preview)
			})

			// Billing
			r.Route("/billing", func(r chi.Router) {
				r.Get("/usage", billingHandler.Usage)
				r.Get("/usage/history", billingHandler.History)
				r.Get("/plan", billingHandler.Plan)
			})

			// Organization / Team
			r.Get("/org", teamHandler.GetOrg)
			r.With(middleware.RequireRole("admin")).Patch("/org", teamHandler.UpdateOrg)
			r.Route("/org/members", func(r chi.Router) {
				r.Get("/", teamHandler.ListMembers)
				r.With(middleware.RequireRole("admin")).Post("/invite", teamHandler.InviteMember)
				r.With(middleware.RequireRole("admin")).Patch("/{id}", teamHandler.UpdateRole)
				r.With(middleware.RequireRole("admin")).Delete("/{id}", teamHandler.RemoveMember)
			})
		})
	})

	return r
}
