-- +goose Up
-- +goose StatementBegin
select 'up SQL query';
create table if not exists public.users (
    id serial primary key,
    login varchar(64) unique,
    password_hash varchar(64)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
drop table public.users;
-- +goose StatementEnd
