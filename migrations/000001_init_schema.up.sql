CREATE TABLE IF NOT EXISTS events(
    event_id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    event_date DATE NOT NULL,
    event TEXT NOT NULL
);