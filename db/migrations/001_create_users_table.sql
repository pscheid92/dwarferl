-- Write your migrate up statements here
create table "users" (
    id text primary key,
    email text not null unique
);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
drop table if exists "users";
