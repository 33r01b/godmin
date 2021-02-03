CREATE TABLE users
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    encrypted_password TEXT NOT NULL
);
