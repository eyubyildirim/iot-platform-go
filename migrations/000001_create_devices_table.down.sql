-- This reverses the first migration by dropping the trigger, the table, and the function in the correct order.
DROP TRIGGER IF EXISTS update_devices_updated_at ON devices;
DROP TABLE IF EXISTS devices;
DROP FUNCTION IF EXISTS update_updated_at_column();
