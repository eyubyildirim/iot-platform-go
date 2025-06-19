-- Create the 'sensor_data' table for storing time-series measurements.
CREATE TABLE sensor_data (
    -- 'BIGSERIAL' is an auto-incrementing 8-byte integer, which maps directly to Go's 'int64'.
    id BIGSERIAL PRIMARY KEY,
    
    -- This is the foreign key linking to the 'devices' table. It must be 'UUID' to match 'devices.id'.
    device_id UUID NOT NULL,
    
    -- Name of the metric being recorded, e.g., 'temperature', 'humidity'.
    metric_name VARCHAR(100) NOT NULL,
    
    -- 'NUMERIC(10, 2)' is a good choice for sensor values, allowing for numbers up to 99,999,999.99.
    -- It is an exact decimal type, avoiding floating-point inaccuracies. It maps to 'float64' in Go.
    metric_value NUMERIC(10, 2) NOT NULL,
    
    -- The timestamp when the data was recorded.
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- This constraint ensures that every 'device_id' in this table corresponds to a real device.
    -- 'ON DELETE CASCADE' means that if a device is deleted, all of its sensor data is automatically deleted too.
    -- This prevents orphaned data and maintains database integrity.
    CONSTRAINT fk_device
        FOREIGN KEY(device_id) 
        REFERENCES devices(id)
        ON DELETE CASCADE
);

-- === Indexes for Performance ===
-- Indexes are CRITICAL for time-series data. Without them, queries will become very slow as the table grows.

-- This is the most important index. It allows for very fast lookups of recent data for a specific device.
-- e.g., "Get the last 24 hours of temperature data for device X".
CREATE INDEX idx_sensor_data_device_id_timestamp ON sensor_data (device_id, timestamp DESC);

-- This index can speed up queries that look at all data within a time range, regardless of the device.
-- e.g., "How many data points did we receive in the last hour?"
CREATE INDEX idx_sensor_data_timestamp ON sensor_data (timestamp DESC);
