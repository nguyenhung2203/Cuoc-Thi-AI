package validator

import (
	"testing"

	"backend/internal/dto/request"
)

type testStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=0,lte=150"`
}

func TestValidate_Valid(t *testing.T) {
	v := testStruct{
		Name:  "Test User",
		Email: "user@example.com",
		Age:   30,
	}
	if msgs := Validate(v); msgs != nil {
		t.Fatalf("expected no errors, got %v", msgs)
	}
}

func TestValidate_Invalid(t *testing.T) {
	v := testStruct{
		Name:  "",
		Email: "not-an-email",
		Age:   200,
	}
	msgs := Validate(v)
	if msgs == nil {
		t.Fatal("expected errors")
	}
	if len(msgs) < 2 {
		t.Fatalf("expected at least 2 errors, got %d: %v", len(msgs), msgs)
	}

	// Check specific messages
	foundRequired := false
	foundEmail := false
	for _, m := range msgs {
		if m == "Name: failed required validation" {
			foundRequired = true
		}
		if m == "Email: failed email validation" {
			foundEmail = true
		}
	}
	if !foundRequired {
		t.Error("missing Name:required message")
	}
	if !foundEmail {
		t.Error("missing Email:email message")
	}
}

func TestValidate_Nil(t *testing.T) {
	// Should not panic on nil
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("should not panic on nil")
		}
	}()
	_ = Validate(nil)
}

func TestValidate_Ptr(t *testing.T) {
	v := &testStruct{
		Name:  "Valid",
		Email: "a@b.com",
		Age:   25,
	}
	if msgs := Validate(v); msgs != nil {
		t.Fatalf("expected no errors, got %v", msgs)
	}
}

// K-AUTH-04/05: the register DTO's tags are the only guard once the OTP gate
// is disabled — an empty payload must produce all three field errors.
func TestValidate_RegisterRequest(t *testing.T) {
	msgs := Validate(&request.RegisterRequest{})
	if len(msgs) != 3 {
		t.Fatalf("empty register: got %d messages (%v), want 3", len(msgs), msgs)
	}

	ok := &request.RegisterRequest{
		Email:    "user@example.com",
		Password: "secret123",
		FullName: "Nguyen Van A",
	}
	if msgs := Validate(ok); msgs != nil {
		t.Fatalf("valid register: got %v", msgs)
	}
}
