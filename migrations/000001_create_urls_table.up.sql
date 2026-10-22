create table if not exists urls (
    id bigserial primary key,
    short_id text not null unique,
    original_url text not null,
    created_at timestamptz not null default now()
);