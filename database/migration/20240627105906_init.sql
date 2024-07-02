-- +goose Up
CREATE TABLE habit (
	id              INTEGER PRIMARY KEY AUTOINCREMENT,
	title           TEXT NOT NULL,
	day             TEXT NOT NULL,
	is_completed    INTEGER NOT NULL DEFAULT 0,
	updated_at      TEXT
);

CREATE INDEX idx_day ON habit (day);

-- +goose Down
DROP INDEX IF EXIST idx_day;
DROP TABLE IF EXIST habit;
