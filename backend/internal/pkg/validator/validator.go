package validator

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	instance *validator.Validate
	once     sync.Once
)

// get returns the singleton validator, initialising it on first call.
func get() *validator.Validate {
	once.Do(func() {
		instance = validator.New()
	})
	return instance
}

// Validate runs struct-level validation and returns human-readable error messages.
// Each message follows the pattern "<Field>: failed <tag> validation".
// Returns nil when v is valid.
func Validate(v any) []string {
	err := get().Struct(v)
	if err == nil {
		return nil
	}

	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return []string{err.Error()}
	}

	msgs := make([]string, 0, len(ve))
	for _, fe := range ve {
		msgs = append(msgs, fmt.Sprintf("%s: failed %s validation", fe.Field(), fe.Tag()))
	}
	return msgs
}
