CREATE TABLE IF NOT EXISTS outcomes (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       market_id UUID NOT NULL REFERENCES markets(id),
       label TEXT NOT NULL,
       created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
       UNIQUE (market_id, label)
);

CREATE INDEX outcomes_market_id_idx on outcomes(market_id);