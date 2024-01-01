create table todos
(
    id          uuid not null primary key,
    title       text,
    description text,
    completed   boolean,
    user_id     bigint
        constraint fk_todos_user references users,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone
);

alter table todos
    owner to myuser;

create index idx_todos_deleted_at
    on todos (deleted_at);