-- Migration: Rename quote_status enum value ENVIADA → EMITIDA
-- Reason: "Emitida" is more accurate than "Enviada" as it's delivery-channel agnostic
-- (covers email, print, in-person delivery)

DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM pg_enum
    WHERE enumlabel = 'ENVIADA'
      AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'quote_status')
  ) THEN
    ALTER TYPE quote_status RENAME VALUE 'ENVIADA' TO 'EMITIDA';
  END IF;
END
$$;
