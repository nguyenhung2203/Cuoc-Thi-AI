-- Wipe corrupted batch2 rows before re-import (UTF-8).
BEGIN;

DELETE FROM jobs
WHERE id >= 'dddddddd-0000-4000-8000-000000000201'
  AND id <= 'dddddddd-0000-4000-8000-000000000228';

DELETE FROM company_members
WHERE company_id >= 'aaaaaaaa-0000-0000-0000-000000000021'
  AND company_id <= 'aaaaaaaa-0000-0000-0000-000000000034';

DELETE FROM companies
WHERE id >= 'aaaaaaaa-0000-0000-0000-000000000021'
  AND id <= 'aaaaaaaa-0000-0000-0000-000000000034';

COMMIT;
