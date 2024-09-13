ALTER TABLE bookings DROP COLUMN launch_id;
ALTER TABLE bookings ADD COLUMN launch_date DATE NOT NULL;
DROP TABLE IF EXISTS launches;
DROP TYPE IF EXISTS launchDomain;

