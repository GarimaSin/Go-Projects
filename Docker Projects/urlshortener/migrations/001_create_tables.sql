CREATE SEQUENCE IF NOT EXISTS short_id_seq START 1;

CREATE TABLE IF NOT EXISTS urls (
    short TEXT PRIMARY KEY,
    destination TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    expires_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE TABLE IF NOT EXISTS url_clicks_hourly (
    short TEXT NOT NULL,
    hour TIMESTAMP WITH TIME ZONE NOT NULL,
    clicks BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (short, hour)
);
CREATE INDEX IF NOT EXISTS idx_clicks_hourly_short ON url_clicks_hourly(short);
