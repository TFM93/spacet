CREATE TYPE launchDomain AS ENUM ('SPACEX', 'SPACET');


CREATE TABLE IF NOT EXISTS launches(
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   external_id VARCHAR(25) UNIQUE,
   domain launchDomain,
   launch_name VARCHAR(50) NOT NULL,
   date_utc TIMESTAMPTZ NOT NULL,
   launchpad_id VARCHAR(25) REFERENCES launchpads (id),
   booking_id UUID UNIQUE REFERENCES bookings(id),
   destination VARCHAR(25),
   lstatus VARCHAR(20) NOT NULL DEFAULT 'scheduled',
   created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
   updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE INDEX idx_launches_with_bookings_by_date ON launches (booking_id, date_utc)
WHERE booking_id IS NOT NULL;

CREATE OR REPLACE FUNCTION trunc_date(timestamptz) RETURNS date AS $$
    SELECT date_trunc('day', $1)::date;
$$ LANGUAGE SQL IMMUTABLE;

CREATE UNIQUE INDEX idx_unique_launch_per_day_per_launchpad
ON launches (trunc_date(date_utc), launchpad_id) WHERE lstatus != 'cancelled' AND domain = 'SPACET';


CREATE TRIGGER set_updated_at
BEFORE UPDATE ON launches
FOR EACH ROW
EXECUTE FUNCTION updated_at_refresh();