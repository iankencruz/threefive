// backend/internal/jobs/page_cleanup.go
package jobs

import (
	"context"
	"log"
	"time"

	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// PageCleanupWorker handles automatic purging of old deleted pages
type PageCleanupWorker struct {
	queries       *sqlc.Queries
	retentionDays int
	ticker        *time.Ticker
	done          chan bool
}

// NewPageCleanupWorker creates a new page cleanup worker
func NewPageCleanupWorker(queries *sqlc.Queries, retentionDays int) *PageCleanupWorker {
	return &PageCleanupWorker{
		queries:       queries,
		retentionDays: retentionDays,
		done:          make(chan bool),
	}
}

// Start begins the cleanup worker
// It runs immediately once, then repeats every 24 hours
func (w *PageCleanupWorker) Start(ctx context.Context) {
	log.Printf("[Page Cleanup] Starting auto-purge worker (retention: %d days)", w.retentionDays)

	// Run immediately on startup
	w.runCleanup(ctx)

	// Setup ticker for daily runs at the same time
	w.ticker = time.NewTicker(24 * time.Hour)

	go func() {
		for {
			select {
			case <-w.ticker.C:
				w.runCleanup(ctx)
			case <-w.done:
				log.Println("[Page Cleanup] Worker stopped")
				return
			case <-ctx.Done():
				log.Println("[Page Cleanup] Context cancelled, stopping worker")
				return
			}
		}
	}()
}

// Stop gracefully stops the cleanup worker
func (w *PageCleanupWorker) Stop() {
	if w.ticker != nil {
		w.ticker.Stop()
	}
	close(w.done)
}

// runCleanup executes the cleanup logic
func (w *PageCleanupWorker) runCleanup(ctx context.Context) {
	startTime := time.Now()
	cutoffDate := time.Now().AddDate(0, 0, -w.retentionDays)

	log.Printf("[Page Cleanup] Starting purge of pages deleted before %s", cutoffDate.Format("2006-01-02 15:04:05"))

	// Convert time.Time to pgtype.Timestamptz
	cutoffTimestamp := pgtype.Timestamptz{
		Time:  cutoffDate,
		Valid: true,
	}

	rowsDeleted, err := w.queries.PurgeOldDeletedPages(ctx, cutoffTimestamp)
	if err != nil {
		log.Printf("[Page Cleanup] ERROR: Failed to purge pages: %v", err)
		return
	}

	duration := time.Since(startTime)
	log.Printf("[Page Cleanup] Successfully purged %d page(s) in %v", rowsDeleted, duration)

	// Log if no pages were purged
	if rowsDeleted == 0 {
		log.Printf("[Page Cleanup] No pages to purge (none older than %d days)", w.retentionDays)
	}
}

// RunManualCleanup allows manual triggering of cleanup (useful for testing or admin endpoints)
func (w *PageCleanupWorker) RunManualCleanup(ctx context.Context) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -w.retentionDays)

	cutoffTimestamp := pgtype.Timestamptz{
		Time:  cutoffDate,
		Valid: true,
	}

	return w.queries.PurgeOldDeletedPages(ctx, cutoffTimestamp)
}
