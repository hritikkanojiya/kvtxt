package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/hritikkanojiya/kvtxt/internal/storage"
)

func StartCleanupWorker(
	ctx context.Context,
	store *storage.Storage,
	interval time.Duration,
) {

	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				slog.Info("cleanup worker stopped")
				return

			case <-ticker.C:
				now := time.Now().Unix()

				deleted, err := store.DeleteExpired(now)
				if err != nil {
					slog.Error("cleanup failed", "error", err)
					continue
				}

				if deleted > 0 {
					slog.Info("expired entries cleaned",
						"count", deleted,
					)
				}
			}
		}
	}()
}
