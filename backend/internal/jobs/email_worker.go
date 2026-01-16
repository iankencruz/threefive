// backend/internal/jobs/email_worker.go
package jobs

import (
	"context"
	"log"
	"time"

	"github.com/iankencruz/threefive/internal/contacts"
)

type EmailWorker struct {
	contactService *contacts.Service
	ticker         *time.Ticker
	done           chan bool
}

func NewEmailWorker(contactService *contacts.Service) *EmailWorker {
	return &EmailWorker{
		contactService: contactService,
		done:           make(chan bool),
	}
}

// Start begins the email retry worker
func (w *EmailWorker) Start(ctx context.Context) {
	log.Println("[Email Worker] Starting email retry worker (runs every 1 hour)")

	// Run immediately on startup
	w.runRetry(ctx)

	// Setup ticker for hourly runs
	w.ticker = time.NewTicker(1 * time.Hour)

	go func() {
		for {
			select {
			case <-w.ticker.C:
				w.runRetry(ctx)
			case <-w.done:
				log.Println("[Email Worker] Worker stopped")
				return
			case <-ctx.Done():
				log.Println("[Email Worker] Context cancelled, stopping worker")
				return
			}
		}
	}()
}

// Stop gracefully stops the worker
func (w *EmailWorker) Stop() {
	if w.ticker != nil {
		w.ticker.Stop()
	}
	close(w.done)
}

func (w *EmailWorker) runRetry(ctx context.Context) {
	log.Println("[Email Worker] Starting email retry run")
	startTime := time.Now()

	err := w.contactService.RetryFailedEmails(ctx)
	if err != nil {
		log.Printf("[Email Worker] Error during retry: %v", err)
		return
	}

	duration := time.Since(startTime)
	log.Printf("[Email Worker] Email retry completed in %v", duration)
}
