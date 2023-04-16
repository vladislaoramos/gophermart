-- +goose Up
-- +goose StatementBegin
select 'up SQL query';
-- +goose StatementEnd
create table if not exists public.withdrawal (
    id serial primary key,
    order_number text,
    sum_number decimal,
    updated_at timestamp default now(),
    user_id int,
    constraint FK_withdrawal_user foreign key (user_id) references public.users (id)
);
-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
drop table public.withdrawal;
-- +goose StatementEnd
