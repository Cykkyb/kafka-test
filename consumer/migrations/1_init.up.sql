CREATE TABLE IF NOT EXISTS messages
(
    id           SERIAL PRIMARY KEY,
    creation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    message_data      TEXT NOT NULL
);