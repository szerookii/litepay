package env

import (
	"fmt"
	"os"
	"strings"
)

type Var struct {
	Key      string
	Required bool
	Default  string // Default value if not set
	// Validate is an optional additional check on the value.
	Validate func(v string) error
}

// Required declares a mandatory env var.
func Required(key string, validators ...func(string) error) Var {
	return Var{Key: key, Required: true, Validate: chain(validators)}
}

// Optional declares an optional env var with optional validators.
func Optional(key string, validators ...func(string) error) Var {
	return Var{Key: key, Required: false, Validate: chain(validators)}
}

// WithDefault declares an optional env var with a default fallback value.
func WithDefault(key string, defaultValue string, validators ...func(string) error) Var {
	return Var{Key: key, Required: false, Default: defaultValue, Validate: chain(validators)}
}

// RequiredWithDefault declares an env var that can use a default but validates if set.
func RequiredWithDefault(key string, defaultValue string, validators ...func(string) error) Var {
	return Var{Key: key, Required: false, Default: defaultValue, Validate: chain(validators)}
}

// Get returns os.Getenv(key). Sugar for use after Check passes.
func Get(key string) string { return os.Getenv(key) }

// GetOrDefault returns the env var value or a default if not set.
func GetOrDefault(key string, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// Check validates all declared vars and returns a combined error.
// Call this at startup — fail fast before any subsystem initializes.
// It also sets default values in the environment for vars that have defaults.
func Check(vars ...Var) error {
	var errs []string
	for _, v := range vars {
		val := os.Getenv(v.Key)
		
		// Use default if value is empty and default is set
		if val == "" && v.Default != "" {
			val = v.Default
			// Set the default in the environment so Get() works as expected
			os.Setenv(v.Key, val)
		}
		
		if val == "" {
			if v.Required {
				errs = append(errs, fmt.Sprintf("  %s: required but not set", v.Key))
			}
			continue
		}
		
		if v.Validate != nil {
			if err := v.Validate(val); err != nil {
				errs = append(errs, fmt.Sprintf("  %s: %s", v.Key, err.Error()))
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("environment configuration errors:\n%s", strings.Join(errs, "\n"))
	}
	return nil
}

// MinLen returns a validator that rejects values shorter than n characters.
func MinLen(n int) func(string) error {
	return func(v string) error {
		if len(v) < n {
			return fmt.Errorf("must be at least %d characters", n)
		}
		return nil
	}
}

// OneOf returns a validator that requires the value to be one of the allowed options.
func OneOf(options ...string) func(string) error {
	return func(v string) error {
		for _, o := range options {
			if v == o {
				return nil
			}
		}
		return fmt.Errorf("must be one of: %s", strings.Join(options, ", "))
	}
}

func chain(validators []func(string) error) func(string) error {
	if len(validators) == 0 {
		return nil
	}
	return func(v string) error {
		for _, fn := range validators {
			if err := fn(v); err != nil {
				return err
			}
		}
		return nil
	}
}
