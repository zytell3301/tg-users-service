USE tg;

CREATE TABLE IF NOT EXISTS users_pk_username
(
    username VARCHAR,
    id       UUID,
    PRIMARY KEY (username)
) WITH ID = '3c3a13ad-9172-46f9-bf6d-a5652cb3609c';