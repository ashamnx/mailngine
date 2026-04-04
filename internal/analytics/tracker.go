package analytics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/hellomail/hellomail/internal/db/sqlcdb"
)

// 1x1 transparent GIF
var transparentGIF = []byte{
	0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00, 0x01, 0x00,
	0x80, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x21,
	0xf9, 0x04, 0x01, 0x00, 0x00, 0x00, 0x00, 0x2c, 0x00, 0x00,
	0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x02, 0x02, 0x44,
	0x01, 0x00, 0x3b,
}

type Tracker struct {
	queries *sqlcdb.Queries
	cache   *redis.Client
	logger  zerolog.Logger
}

func NewTracker(db *pgxpool.Pool, cache *redis.Client, logger zerolog.Logger) *Tracker {
	return &Tracker{
		queries: sqlcdb.New(db),
		cache:   cache,
		logger:  logger,
	}
}

func (t *Tracker) TrackOpen(w http.ResponseWriter, r *http.Request) {
	trackingID := chi.URLParam(r, "id")
	emailID, err := uuid.Parse(trackingID)
	if err != nil {
		w.Header().Set("Content-Type", "image/gif")
		w.Write(transparentGIF)
		return
	}

	// Deduplicate opens
	cacheKey := fmt.Sprintf("open:%s", emailID.String())
	if t.cache.SetNX(r.Context(), cacheKey, "1", 24*time.Hour).Val() {
		// First open — record event
		email, err := t.queries.GetEmail(r.Context(), sqlcdb.GetEmailParams{ID: emailID, OrgID: uuid.Nil})
		if err == nil {
			t.queries.CreateEmailEvent(r.Context(), sqlcdb.CreateEmailEventParams{
				EmailID:   emailID,
				OrgID:     email.OrgID,
				EventType: "opened",
				Recipient: "",
				UserAgent: pgtype.Text{String: r.UserAgent(), Valid: true},
				IpAddress: r.RemoteAddr,
			})
		}
	}

	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	w.Write(transparentGIF)
}

func (t *Tracker) TrackClick(w http.ResponseWriter, r *http.Request) {
	trackingID := chi.URLParam(r, "id")
	targetURL := r.URL.Query().Get("url")

	if targetURL == "" {
		http.Error(w, "missing url", http.StatusBadRequest)
		return
	}

	emailID, err := uuid.Parse(trackingID)
	if err != nil {
		http.Redirect(w, r, targetURL, http.StatusFound)
		return
	}

	// Record click
	cacheKey := fmt.Sprintf("click:%s:%s", emailID.String(), targetURL)
	if t.cache.SetNX(r.Context(), cacheKey, "1", 24*time.Hour).Val() {
		metadata, _ := json.Marshal(map[string]string{"url": targetURL})
		email, err := t.queries.GetEmail(r.Context(), sqlcdb.GetEmailParams{ID: emailID, OrgID: uuid.Nil})
		if err == nil {
			t.queries.CreateEmailEvent(r.Context(), sqlcdb.CreateEmailEventParams{
				EmailID:   emailID,
				OrgID:     email.OrgID,
				EventType: "clicked",
				Recipient: "",
				Metadata:  metadata,
				UserAgent: pgtype.Text{String: r.UserAgent(), Valid: true},
				IpAddress: r.RemoteAddr,
			})
		}
	}

	http.Redirect(w, r, targetURL, http.StatusFound)
}
