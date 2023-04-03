-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE public.order (
    id serial primary key,
    order_number text,
    status varchar(32) default 'NEW',
    accrual decimal,
    uploaded_at timestamp default now(),
    user_id int,
    constraint FK_order_user foreign key (user_id) references public.user (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE public.order;
-- +goose StatementEnd
