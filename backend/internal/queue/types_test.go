package queue

import (
	"encoding/json"
	"testing"
)

func TestGenerateReportPayload_RoundTrip(t *testing.T) {
	in := GenerateReportPayload{
		CompanyID:   "c1",
		InterviewID: "i1",
		JobID:       "j1",
		GeneratedBy: "u1",
		RecruiterID: "r1",
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out GenerateReportPayload
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out != in {
		t.Errorf("round-trip mismatch: got %+v want %+v", out, in)
	}
}

func TestJobTypeConstants(t *testing.T) {
	// Guard against accidental renames that would silently break routing.
	cases := map[string]string{
		TypeGenerateReport:  "report:generate",
		TypeAnalyzeCV:       "cv:analyze",
		TypeBatchTranscript: "transcript:batch",
	}
	for got, want := range cases {
		if got != want {
			t.Errorf("job type = %q, want %q", got, want)
		}
	}
}
