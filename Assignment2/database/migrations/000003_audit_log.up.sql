create table if not exists audit_logs (
  id serial primary key,
  action varchar(64) not null,
  user_id int,
  payload jsonb,
  created_at timestamptz not null default now()
);