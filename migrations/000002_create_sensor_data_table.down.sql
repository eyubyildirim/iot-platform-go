-- This reverses the second migration by dropping the sensor_data table.
-- Indexes associated with the table are dropped automatically.
DROP TABLE IF EXISTS sensor_data;
