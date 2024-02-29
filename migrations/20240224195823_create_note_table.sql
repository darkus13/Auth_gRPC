-- +goose Up
CREATE TABLE users (
    id serial primary key,
    name text not null,
    email text,
    password text not null,
    password_confirm text not null,
    role text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

-- +goose Down
DROP TABLE users;