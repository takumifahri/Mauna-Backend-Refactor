package validation

import "testing"

type sampleRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=12"`
}

func TestValidatePassesValidStruct(t *testing.T) {
	req := sampleRequest{
		Email:    "user@example.com",
		Password: "secret1",
	}

	if err := Validate(req); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
}

func TestValidateFailsInvalidStruct(t *testing.T) {
	req := sampleRequest{
		Email:    "not-an-email",
		Password: "123",
	}

	err := Validate(req)
	if err == nil {
		t.Fatal("Validate() error = nil, want error")
	}

	validationErr, ok := err.(ValidationError)
	if !ok {
		t.Fatalf("Validate() error type = %T, want ValidationError", err)
	}
	if len(validationErr.Fields) != 2 {
		t.Fatalf("Validate() field errors = %d, want 2", len(validationErr.Fields))
	}
}
