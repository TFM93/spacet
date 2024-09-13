CREATE TYPE launchDomain AS ENUM ('SPACEX', 'SPACET');


CREATE TABLE IF NOT EXISTS launches(
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   external_id VARCHAR(25),
   domain launchDomain,
   launch_name VARCHAR(50) NOT NULL,
   date_utc TIMESTAMPTZ NOT NULL,
   launchpad_id VARCHAR(25) REFERENCES launchpads (id),
   destination VARCHAR(25),
   created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
   updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

ALTER TABLE bookings ADD COLUMN launch_id UUID REFERENCES launches (id);
ALTER TABLE bookings DROP COLUMN launch_date;

CREATE OR REPLACE FUNCTION trunc_date(timestamptz) RETURNS date AS $$
    SELECT date_trunc('day', $1)::date;
$$ LANGUAGE SQL IMMUTABLE;

CREATE UNIQUE INDEX idx_unique_launch_per_day_per_launchpad
ON launches (trunc_date(date_utc), launchpad_id);



CREATE TRIGGER set_updated_at
BEFORE UPDATE ON launches
FOR EACH ROW
EXECUTE FUNCTION updated_at_refresh();