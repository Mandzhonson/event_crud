CREATE TABLE IF NOT EXISTS events(
    user_id INTEGER,
    event_date DATE,
    events TEXT NOT NULL,
    PRIMARY KEY(user_id,event_date)
);