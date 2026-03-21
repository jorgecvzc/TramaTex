-- 026: Replace JSONB mes_work_refs with relational join tables
-- Rationale: name and observations (description) live exclusively in MES work_setups.
-- Sales only stores FKs + sequence.

BEGIN;

-- 1. Create quote_work_setups
CREATE TABLE IF NOT EXISTS quote_work_setups (
    id            UUID PRIMARY KEY,
    quote_id      UUID NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    work_setup_id UUID NOT NULL REFERENCES work_setups(id),
    sequence      INT  NOT NULL DEFAULT 1,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX idx_quote_work_setups_quote ON quote_work_setups(quote_id);

-- 2. Create order_work_setups
CREATE TABLE IF NOT EXISTS order_work_setups (
    id            UUID PRIMARY KEY,
    order_id      UUID NOT NULL REFERENCES sales_orders(id) ON DELETE CASCADE,
    work_setup_id UUID NOT NULL REFERENCES work_setups(id),
    work_order_id UUID REFERENCES mes_works(id),
    sequence      INT  NOT NULL DEFAULT 1,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX idx_order_work_setups_order ON order_work_setups(order_id);

-- 3. Migrate observations from JSONB into MES work_setups.description
--    (only where the work_setup currently has an empty description)
UPDATE work_setups ws
SET description = sub.observations
FROM (
    SELECT DISTINCT ON ((ref->>'work_setup_id')::UUID)
           (ref->>'work_setup_id')::UUID AS ws_id,
           ref->>'observations'          AS observations
    FROM (
        SELECT jsonb_array_elements(mes_work_refs) AS ref
        FROM quotes
        WHERE mes_work_refs IS NOT NULL
        UNION ALL
        SELECT jsonb_array_elements(mes_work_refs) AS ref
        FROM sales_orders
        WHERE mes_work_refs IS NOT NULL
    ) combined
    WHERE ref->>'work_setup_id' IS NOT NULL
      AND ref->>'work_setup_id' != ''
      AND COALESCE(ref->>'observations', '') != ''
    ORDER BY (ref->>'work_setup_id')::UUID
) sub
WHERE ws.id = sub.ws_id
  AND COALESCE(ws.description, '') = '';

-- 4. Migrate quote JSONB data → quote_work_setups (only rows with valid work_setup_id)
INSERT INTO quote_work_setups (id, quote_id, work_setup_id, sequence)
SELECT
    (ref->>'id')::UUID,
    q.id,
    (ref->>'work_setup_id')::UUID,
    COALESCE((ref->>'sequence')::INT, rn.seq)
FROM quotes q,
     LATERAL jsonb_array_elements(q.mes_work_refs) WITH ORDINALITY AS t(ref, seq),
     LATERAL (SELECT t.seq AS seq) rn
WHERE q.mes_work_refs IS NOT NULL
  AND ref->>'work_setup_id' IS NOT NULL
  AND ref->>'work_setup_id' != ''
ON CONFLICT (id) DO NOTHING;

-- 5. Migrate order JSONB data → order_work_setups
INSERT INTO order_work_setups (id, order_id, work_setup_id, work_order_id, sequence)
SELECT
    (ref->>'id')::UUID,
    o.id,
    (ref->>'work_setup_id')::UUID,
    CASE
        WHEN ref->>'work_order_id' IS NOT NULL AND ref->>'work_order_id' != ''
        THEN (ref->>'work_order_id')::UUID
        ELSE NULL
    END,
    COALESCE((ref->>'sequence')::INT, rn.seq)
FROM sales_orders o,
     LATERAL jsonb_array_elements(o.mes_work_refs) WITH ORDINALITY AS t(ref, seq),
     LATERAL (SELECT t.seq AS seq) rn
WHERE o.mes_work_refs IS NOT NULL
  AND ref->>'work_setup_id' IS NOT NULL
  AND ref->>'work_setup_id' != ''
ON CONFLICT (id) DO NOTHING;

-- 6. Drop old JSONB columns
ALTER TABLE quotes DROP COLUMN IF EXISTS mes_work_refs;
ALTER TABLE sales_orders DROP COLUMN IF EXISTS mes_work_refs;

COMMIT;
