package service

import "testing"

func TestNormalizeCandidateLevel(t *testing.T) {
	cases := map[string]string{"Intern": "fresher", "entry-level": "fresher", "JR": "junior", "Middle": "mid", "mid-level": "mid", "unknown": "junior"}
	for input, want := range cases {
		if got := NormalizeCandidateLevel(input); got != want {
			t.Errorf("%q: got %q want %q", input, got, want)
		}
	}
}
