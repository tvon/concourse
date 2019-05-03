BEGIN;
  ALTER TABLE next_build_inputs
    ADD COLUMN "resolve_error" text;
COMMIT;
