CREATE TABLE IF NOT EXISTS markets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    category TEXT NOT NULL,
    resolver_type TEXT NOT NULL,
    resolver_ref UUID REFERENCES users(id),
    expires_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_market_status ON markets(status);
CREATE INDEX idx_markets_expires_at ON markets(expires_at);