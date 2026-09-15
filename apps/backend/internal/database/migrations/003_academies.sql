CREATE TYPE academy_member_role AS ENUM ('owner', 'instructor', 'front_desk', 'student');
CREATE TYPE academy_member_status AS ENUM ('active', 'suspended', 'removed');

CREATE TABLE academies (
    id         TEXT PRIMARY KEY,           -- pending: match users.id type
    name       TEXT NOT NULL,
    slug       TEXT NOT NULL UNIQUE,        -- URL-safe identifier, e.g. academy subdomain or /academies/:slug
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE academy_members (
    id         TEXT PRIMARY KEY,
    academy_id TEXT NOT NULL REFERENCES academies(id) ON DELETE CASCADE,
    user_id    TEXT NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
    role       academy_member_role   NOT NULL DEFAULT 'student',
    status     academy_member_status NOT NULL DEFAULT 'active',
    joined_via TEXT NOT NULL DEFAULT 'manual',  -- 'manual' | 'invite' | 'request' — plain text for now, see note below
    joined_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_academy_members_academy_user ON academy_members (academy_id, user_id);
CREATE INDEX idx_academy_members_user ON academy_members (user_id);

---- create above / drop below ----

DROP TABLE IF EXISTS academy_members;
DROP TABLE IF EXISTS academies;
DROP TYPE IF EXISTS academy_member_status;
DROP TYPE IF EXISTS academy_member_role;
