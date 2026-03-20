package validator

import (
	"strings"
	"testing"
)

// TestInput is a test struct for validation
type TestInput struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Price int64  `json:"price" validate:"required,min=1"`
}

func TestValidator_Validate_ValidInput(t *testing.T) {
	v := New()
	input := &TestInput{
		Name:  "John Doe",
		Email: "john@example.com",
		Price: 1000,
	}

	errors := v.Validate(input)

	if errors != nil {
		t.Errorf("Expected nil, got %v", errors)
	}
}

func TestValidator_Validate_EmptyFields(t *testing.T) {
	v := New()
	input := &TestInput{
		Name:  "",
		Email: "",
		Price: 0,
	}

	errors := v.Validate(input)

	if errors == nil {
		t.Fatal("Expected errors, got nil")
	}

	if len(errors) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(errors))
	}

	expectedKeys := []string{"name", "email", "price"}
	for _, key := range expectedKeys {
		if _, exists := errors[key]; !exists {
			t.Errorf("Expected error key '%s', got %v", key, errors)
		}
	}
}

func TestValidator_Validate_InvalidEmailFormat(t *testing.T) {
	v := New()
	input := &TestInput{
		Name:  "John Doe",
		Email: "invalid-email",
		Price: 1000,
	}

	errors := v.Validate(input)

	if errors == nil {
		t.Fatal("Expected errors, got nil")
	}

	emailError, exists := errors["email"]
	if !exists {
		t.Fatalf("Expected 'email' error key, got %v", errors)
	}

	if !strings.Contains(emailError, "format email tidak valid") {
		t.Errorf("Expected email error to contain 'format email tidak valid', got '%s'", emailError)
	}
}

func TestValidator_Validate_NameTooShort(t *testing.T) {
	v := New()
	input := &TestInput{
		Name:  "J",
		Email: "john@example.com",
		Price: 1000,
	}

	errors := v.Validate(input)

	if errors == nil {
		t.Fatal("Expected errors, got nil")
	}

	nameError, exists := errors["name"]
	if !exists {
		t.Fatalf("Expected 'name' error key, got %v", errors)
	}

	if !strings.Contains(nameError, "minimal 2") {
		t.Errorf("Expected name error to contain 'minimal 2', got '%s'", nameError)
	}
}
