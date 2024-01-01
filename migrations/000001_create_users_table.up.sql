create table users
(
    id            bigserial primary key,
    created_at    timestamp with time zone,
    updated_at    timestamp with time zone,
    deleted_at    timestamp with time zone,
    username      text unique,
    email         text unique,
    password      text,
    last_login_at timestamp with time zone
);

alter table users
    owner to myuser;

create index idx_users_deleted_at
    on users (deleted_at);