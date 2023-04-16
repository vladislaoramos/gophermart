-- +goose Up
-- +goose StatementBegin
select 'up SQL query';
create table if not exists public.balance (
    id serial primary key,
    balance float,
    withdrawal float,
    user_id int unique,
    constraint FK_balance_user foreign key (user_id) references public.users (id)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
drop table public.balance;
-- +goose StatementEnd
