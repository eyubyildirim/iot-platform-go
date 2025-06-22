-- This table stores the rules defined by the user.
CREATE TABLE alert_rules (
    id UUID PRIMARY KEY,
    
    -- A user-friendly name for the rule, e.g., "Living Room High Temp Alert".
    name VARCHAR(255) NOT NULL,

    -- Optionally link a rule to a specific device. If NULL, the rule applies to ALL devices.
    device_id UUID REFERENCES devices(id) ON DELETE CASCADE,
    
    -- We use JSONB to store the flexible rule definition. This is extremely powerful.
    -- Example: {"metric": "temperature", "condition": "gt", "value": 30, "duration_minutes": 5}
    rule_definition JSONB NOT NULL,
    
    -- A simple flag to enable or disable a rule without deleting it.
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- This table will log every single time a rule's conditions are met.
CREATE TABLE alerts_triggered (
    id UUID PRIMARY KEY,
    
    -- The rule that was triggered.
    rule_id UUID NOT NULL REFERENCES alert_rules(id) ON DELETE CASCADE,
    
    -- The specific device that triggered the alert.
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    
    -- The actual value from the sensor that caused the trigger.
    triggered_value NUMERIC(10, 2) NOT NULL,
    
    -- A snapshot of the rule definition and other context at the time of the alert.
    details JSONB,
    
    -- The timestamp when the alert was created.
    triggered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add an index to quickly find all active rules.
CREATE INDEX idx_alert_rules_is_active ON alert_rules (is_active);
