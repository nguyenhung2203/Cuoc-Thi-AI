DO $$
DECLARE
    cid uuid := '34032bf9-3254-48af-b5fb-e81361ac1633';
    uid uuid;
    i int;
BEGIN
    SELECT id INTO uid FROM users WHERE email = 'luulai2006@gmail.com';
    FOR i IN 1..10 LOOP
        INSERT INTO jobs (id, company_id, title, description, status, created_by, created_at, updated_at)
        VALUES (
            gen_random_uuid(),
            cid,
            'Job Title ' || i || ' (Generated)',
            'This is a generated job description for job ' || i || '. It has at least 10 characters so validation passes.',
            'open',
            uid,
            NOW(),
            NOW()
        );
    END LOOP;
END $$;
