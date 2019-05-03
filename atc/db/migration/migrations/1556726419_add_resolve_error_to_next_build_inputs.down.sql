BEGIN;
  ALTER TABLE next_build_inputs
    DROP COLUMN resolve_error;
COMMIT;
