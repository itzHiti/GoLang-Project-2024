CREATE TABLE IF NOT EXISTS users
(
    id
    bigserial
    PRIMARY
    KEY,
    username
    citext
    NOT
    NULL
    UNIQUE,
    email
    citext
    NOT
    NULL
    UNIQUE,
    password
    bytea
    NOT
    NULL,
    token_hash
    text
    NOT
    NULL,
    activated
    bool
    NOT
    NULL
    DEFAULT
    false
);
CREATE TABLE IF NOT EXISTS admins
(
    id
    bigserial
    PRIMARY
    KEY,
    user_id
    bigint,
    FOREIGN
    KEY
(
    user_id
)
    REFERENCES users
(
    id
)
    );
CREATE TABLE IF NOT EXISTS bans
(
    id
    bigserial
    PRIMARY
    KEY,
    user_id
    bigint
    NOT
    NULL
    UNIQUE,
    expiry
    timestamp
(
    0
) with time zone NOT NULL,
      FOREIGN KEY (user_id)
    REFERENCES users
(
    id
)
  ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS verifications
(
    code
    bytea
    PRIMARY
    KEY,
    user_id
    bigint,
    expiry
    timestamp
(
    0
) with time zone NOT NULL,
      FOREIGN KEY (user_id) REFERENCES users
(
    id
)
  on delete CASCADE
    );
