USE tg;

CREATE TABLE IF NOT EXISTS users
(
    id            UUID PRIMARY KEY,
    name          VARCHAR,
    lastname      VARCHAR,
    bio           VARCHAR,
    username      VARCHAR,
    phone         VARCHAR,
    online_status BOOLEAN,
    created_at    TIMESTAMP
);