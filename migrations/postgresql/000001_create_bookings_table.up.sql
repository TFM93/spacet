CREATE OR REPLACE FUNCTION updated_at_refresh()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TYPE gender AS ENUM ('UNSPECIFIED', 'MALE', 'FEMALE', 'NON_BINARY', 'OTHER');

CREATE TABLE IF NOT EXISTS bookings(
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   first_name VARCHAR(25) NOT NULL,
   last_name VARCHAR(25)  NOT NULL,
   gender gender DEFAULT 'UNSPECIFIED',
   birthday DATE,
   launch_date DATE NOT NULL,
   launchpad_id VARCHAR(25),

   created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
   updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE INDEX idx_bookings_updated_at_id ON bookings (updated_at DESC, id DESC);
CREATE INDEX idx_bookings_launch_date_id ON bookings (launch_date DESC, id DESC);

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON bookings
FOR EACH ROW
EXECUTE FUNCTION updated_at_refresh();