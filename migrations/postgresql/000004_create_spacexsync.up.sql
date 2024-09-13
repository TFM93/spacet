CREATE TABLE IF NOT EXISTS sync_info(
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   resource_name TEXT UNIQUE,
   last_sync TIMESTAMP NOT NULL
);
