CREATE TABLE hosts (
  id SERIAL PRIMARY KEY,
  hostname VARCHAR(255) UNIQUE NOT NULL,
  last_seen TIMESTAMP NOT NULL
);

CREATE TABLE metrics (
  id SERIAL PRIMARY KEY,
  host_id INTEGER REFERENCES hosts(id) ON DELETE CASCADE,
  name VARCHAR(50) NOT NULL,
  value DOUBLE PRECISION NOT NULL,
  recorded_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_metrics_host_recorded 
  ON metrics(host_id, recorded_at DESC);

CREATE INDEX idx_metrics_name 
  ON metrics(name);
