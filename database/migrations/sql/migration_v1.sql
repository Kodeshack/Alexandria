CREATE TABLE schema_version (
    version INT
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    admin INT,
    email TEXT UNIQUE,
    display_name TEXT,
    creation_date INT,
    password TEXT,
    salt TEXT,
    argon2_key_len INT,
    argon2_memory INT,
    argon2_threads INT,
    argon2_time INT,
    argon2_version INT
)
