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

func (r *Repository) SaveMetrics(ctx context.Context, metrics model.ServerMetrics) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var hostID int

	reqID := `INSERT INTO hosts (hostname, last_seen)
		VALUES ($1, $2)
		ON CONFLICT (hostname)
		DO UPDATE SET last_seen = EXCLUDED.last_seen
		RETURNING id;`

	reqInsert := `INSERT INTO metrics (host_id, name, value, recorded_at)
		VALUES ($1, $2, $3, $4)`

	row := tx.QueryRowContext(ctx, reqID, metrics.Hostname, time.Now())
	err = row.Scan(&hostID)
	if err != nil {
			return err
	}
	for i := 0; i < len(metrics.Metrics); i++ {
		_, err = tx.ExecContext(
			ctx, reqInsert, hostID,
			metrics.Metrics[i].Name, metrics.Metrics[i].Value, time.Now(),
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
