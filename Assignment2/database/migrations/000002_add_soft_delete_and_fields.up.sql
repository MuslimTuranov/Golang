alter table users
  add column if not exists email varchar(255) unique,
  add column if not exists age int,
  add column if not exists created_at timestamptz not null default now(),
  add column if not exists updated_at timestamptz not null default now(),
  add column if not exists deleted_at timestamptz null;