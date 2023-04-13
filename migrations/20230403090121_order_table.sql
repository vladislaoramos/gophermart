-- +goose Up
-- +goose StatementBegin
select 'up SQL query';
create table if not exists public.orders (
    id serial primary key,
    order_number text,
    status varchar(32) default 'NEW',
    accrual decimal,
    uploaded_at timestamp default now(),
    user_id int,
    constraint FK_order_user foreign key (user_id) references public.users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
select 'down SQL query';
drop table public.orders;
-- +goose StatementEnd
