USE tg;

CREATE TABLE IF NOT EXISTS users_pk_phone
(
    phone    VARCHAR,
    id       UUID,
    lastname VARCHAR,
    bio      VARCHAR,
    username VARCHAR,
    PRIMARY KEY (phone)
) WITH ID = 'a6199f05-ef65-4780-b934-bfd4c8ccd52c';