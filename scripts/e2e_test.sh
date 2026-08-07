#!/usr/bin/env bash
# End-to-end smoke test hitting the local test backend on :8090.
set -e
BASE=http://localhost:8090/api/v1
jqget() { python -c "import sys,json;d=json.load(sys.stdin);print(eval('d'+sys.argv[1]))" "$1"; }

TOKEN=$(curl -s -X POST $BASE/auth/login -H "Content-Type: application/json" \
  -d '{"email":"recruiter@test.com","password":"Test1234!"}' | jqget "['data']['access_token']")
CID=$(curl -s $BASE/auth/me -H "Authorization: Bearer $TOKEN" | jqget "['data']['companies'][0]['id']")
AUTH="-H Authorization:Bearer_$TOKEN"
echo "COMPANY=$CID"

echo "=== CREATE JOB ==="
JOB=$(curl -s -X POST $BASE/companies/$CID/jobs -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"title":"Backend Engineer","description":"We need a strong backend engineer with Go experience and distributed systems knowledge.","status":"open"}')
echo "$JOB" | head -c 300; echo
JOB_ID=$(echo "$JOB" | jqget "['data']['id']")
echo "JOB_ID=$JOB_ID"

echo "=== CREATE CANDIDATE ==="
CEMAIL="cand+$(date +%s)@test.com"
CAND=$(curl -s -X POST $BASE/companies/$CID/candidates -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d "{\"full_name\":\"Nguyen Van Test\",\"email\":\"$CEMAIL\",\"job_id\":\"$JOB_ID\"}")
echo "$CAND" | head -c 300; echo
CAND_ID=$(echo "$CAND" | jqget "['data']['id']")
echo "CAND_ID=$CAND_ID"

echo "=== CREATE INTERVIEW ==="
IV=$(curl -s -X POST $BASE/companies/$CID/interviews -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d "{\"job_id\":\"$JOB_ID\",\"candidate_id\":\"$CAND_ID\",\"scheduled_at\":\"2026-08-01T09:00:00Z\",\"mode\":\"real\",\"send_invite\":false}")
echo "$IV" | head -c 400; echo
IV_ID=$(echo "$IV" | jqget "['data']['interview_id']")
echo "IV_ID=$IV_ID"

echo "=== SAVE NOTES ==="
curl -s -X PUT $BASE/companies/$CID/interviews/$IV_ID/notes -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"notes":"Ung vien co kinh nghiem tot, can hoi them ve he thong phan tan."}' | head -c 300; echo

echo "=== GET INTERVIEW (verify notes persisted) ==="
curl -s $BASE/companies/$CID/interviews/$IV_ID -H "Authorization: Bearer $TOKEN" | python -c "import sys,json;d=json.load(sys.stdin)['data'];print('recruiter_notes=',repr(d.get('recruiter_notes')))"

echo "=== SEND REMINDER ==="
curl -s -X POST $BASE/companies/$CID/interviews/$IV_ID/send-reminder -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{}' | head -c 300; echo

echo "DONE_IDS $CID $JOB_ID $CAND_ID $IV_ID"
