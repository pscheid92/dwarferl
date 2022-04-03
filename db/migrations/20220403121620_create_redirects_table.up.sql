create table "redirects" (
    short text primary key,
    url text not null,
    user_id text not null references "users" (id) on delete cascade,
    created_at timestamptz not null
);
