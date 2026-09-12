-- Write your migrate up statements here
CREATE TYPE user_roles AS ENUM (
    'user',
    'admin',
    'instructor',
    'student'
);
---- create above / drop below ----
DROP TYPE IF EXISTS user_roles;