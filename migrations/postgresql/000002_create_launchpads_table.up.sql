CREATE TABLE IF NOT EXISTS launchpads(
   id VARCHAR(25) PRIMARY KEY,
   pad_name VARCHAR(50) NOT NULL,
   locality VARCHAR(50)  NOT NULL,
   region VARCHAR(50)  NOT NULL,
   timezone VARCHAR(50)  NOT NULL,
   pad_status VARCHAR(25),

   created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
   updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);


CREATE TRIGGER set_updated_at
BEFORE UPDATE ON launchpads
FOR EACH ROW
EXECUTE FUNCTION updated_at_refresh();