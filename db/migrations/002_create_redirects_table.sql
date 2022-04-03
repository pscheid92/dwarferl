-- Write your migrate up statements here
create table "redirects" (
    short text primary key,
    url text not null,
    user_id text not null references "users" (id) on delete cascade,
    created_at timestamptz not null
);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
drop table if exists "redirects";
