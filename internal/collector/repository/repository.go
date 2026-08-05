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

func (r *Repository) CleanOldMetrics(ctx context.Context, threshold time.Time) error {
	reqClean := `DELETE FROM metrics WHERE recorded_at < $1`

	_, err := r.db.ExecContext(ctx, reqClean, threshold)
	return err
}

func (r *Repository) SelectLastMetrics(ctx context.Context, hostID int, count int) ([]model.Metric, error) {
	reqSelect := `SELECT name, value FROM metrics WHERE host_id = $1 ORDER BY recorded_at DESC LIMIT $2`

	rows, err := r.db.QueryContext(ctx, reqSelect, hostID, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []model.Metric
	
	for rows.Next() {
		var metric model.Metric
		err = rows.Scan(&metric.Name, &metric.Value)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (r *Repository) GetAllHosts(ctx context.Context) ([]model.Host, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, hostname FROM hosts ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []model.Host
	for rows.Next() {
		var h model.Host
		if err := rows.Scan(&h.ID, &h.Hostname); err != nil {
			return nil, err
		}
		hosts = append(hosts, h)
	}
	return hosts, nil
}
