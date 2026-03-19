CREATE TABLE IF NOT EXISTS customers(
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    height_cm numeric,
    height_in numeric GENERATED ALWAYS AS (height_cm / 2.54),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
