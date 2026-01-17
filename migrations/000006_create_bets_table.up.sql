CREATE TABLE IF NOT EXISTS bets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    market_id UUID NOT NULL REFERENCES markets(id),
    outcome_id UUID NOT NULL REFERENCES outcomes(id),
    amount INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(market_id, user_id)
);

CREATE INDEX bets_user_id_idx on bets(user_id);
CREATE INDEX bets_market_id_idx on bets(market_id);