package repository

import (
	"context"
	"database/sql"
	"time"

	"monitoring-system/internal/model"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveMetrics(ctx context.Context, hostname string, metrics model.ServerMetrics) error {

	tx, err := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	reqID := `INSERT INTO hosts (hostname, last_seen)
		VALUES ($1, $2)
		ON CONFLICT (hostname)
		DO UPDATE SET last_seen = EXCLUDED.last_seen
		RETURNING id;`

	reqInsert := `INSERT INTO metrics (host_id, name, value, recorded_at)
		VALUES ($1, $2, $3, $4)`

	db.BeginTx(
		hostID, err := db.ExecContext(ctx, reqID, metrics.Hostname, time.Now())
		if err != nil {
			return err
		}
		for i := 0; i < len(metrics.Metrics); i++ {
			_, err = db.ExecContext(
				ctx, reqInsert, hostID,
				metrics.Metrics[i].Name, metrics.Metrics[i].Value, time.Now(),
			)
			if err != nil {
				return err
			}
		}
	)

	return nil
}
