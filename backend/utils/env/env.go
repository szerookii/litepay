package env

import (
	"fmt"
	"os"
	"strings"
)

type Var struct {
	Key      string
	Required bool
	Default  string
	Validate func(v string) error
}

func Required(key string, validators ...func(string) error) Var {
	return Var{Key: key, Required: true, Validate: chain(validators)}
}

func Optional(key string, validators ...func(string) error) Var {
	return Var{Key: key, Required: false, Validate: chain(validators)}
}

func WithDefault(key string, defaultValue string, validators ...func(string) error) Var {
	return Var{Key: key, Required: false, Default: defaultValue, Validate: chain(validators)}
}

func RequiredWithDefault(key string, defaultValue string, validators ...func(string) error) Var {
	return Var{Key: key, Required: false, Default: defaultValue, Validate: chain(validators)}
}

func Get(key string) string { return os.Getenv(key) }

func GetOrDefault(key string, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func Check(vars ...Var) error {
	var errs []string
	for _, v := range vars {
		val := os.Getenv(v.Key)

		if val == "" && v.Default != "" {
			val = v.Default
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

func MinLen(n int) func(string) error {
	return func(v string) error {
		if len(v) < n {
			return fmt.Errorf("must be at least %d characters", n)
		}
		return nil
	}
}

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
