CREATE TABLE IF NOT EXISTS events(
    event_id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    event_date DATE NOT NULL,
    events TEXT NOT NULL
);