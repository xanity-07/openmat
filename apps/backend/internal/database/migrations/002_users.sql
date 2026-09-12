-- Write your migrate up statements here
CREATE  TABLE users (
    id VARCHAR(11) PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role user_roles DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

---- create above / drop below ----
DROP TABLE IF EXISTS users;
