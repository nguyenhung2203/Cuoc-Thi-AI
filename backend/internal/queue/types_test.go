package queue

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSendOTPEmailPayload_ExcludesRawOTP(t *testing.T) {
	in := SendOTPEmailPayload{To: "user@example.com", Purpose: "reset", GenerationID: "gen-1", TemplateType: "otp"}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(b), `"otp":`) {
		t.Fatalf("payload contains raw OTP field: %s", b)
	}
	var out SendOTPEmailPayload
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.To != in.To || out.GenerationID != in.GenerationID || out.Purpose != in.Purpose {
		t.Fatalf("payload mismatch: got %+v want %+v", out, in)
	}
}

func TestGenerateReportPayload_RoundTrip(t *testing.T) {
	in := GenerateReportPayload{CompanyID: "c1", InterviewID: "i1", JobID: "j1", GeneratedBy: "u1", RecruiterID: "r1"}
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
	cases := map[string]string{TypeGenerateReport: "report:generate", TypeAnalyzeCV: "cv:analyze", TypeBatchTranscript: "transcript:batch", TypeSendOTPEmail: "email:send-otp"}
	for got, want := range cases {
		if got != want {
			t.Errorf("job type = %q, want %q", got, want)
		}
	}
}
