-- This function will be triggered on every update to a row.
-- It automatically sets the 'updated_at' column to the current time.
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW(); 
   RETURN NEW;
END;
$$ language 'plpgsql';

-- Create the main 'devices' table
CREATE TABLE devices (
    -- 'UUID' is the best data type for universally unique identifiers.
    id UUID PRIMARY KEY,
    
    -- 'VARCHAR(255)' is a standard choice for names. 'NOT NULL' ensures it's always present.
    name VARCHAR(255) NOT NULL,
    
    -- 'VARCHAR(100)' for the device type, maps to your 'Kind' field.
    kind VARCHAR(100) NOT NULL,
    
    -- 'api_key' must be unique to ensure we can identify each device correctly.
    api_key VARCHAR(128) UNIQUE NOT NULL,
    
    -- 'TIMESTAMPTZ' stores the timestamp with timezone information, which is best practice.
    -- 'DEFAULT NOW()' automatically sets the creation time on insert.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Apply the trigger to the 'devices' table.
-- This will call the function before any row is updated.
CREATE TRIGGER update_devices_updated_at
BEFORE UPDATE ON devices
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Optional: Add an index on the api_key for very fast lookups during authentication.
CREATE INDEX idx_devices_api_key ON devices (api_key);
