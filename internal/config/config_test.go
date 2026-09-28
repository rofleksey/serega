package config

import "testing"

func TestConfigurationDefaultsAndValidation(t *testing.T) {
	const (
		databaseKey  = "DATABASE_URL"
		testDatabase = "test"
	)

	tests := []struct {
		name    string
		values  map[string]string
		invalid bool
	}{
		{name: "defaults", values: map[string]string{databaseKey: "postgres://local/test"}},
		{name: "missing database", invalid: true},
		{name: "invalid address", values: map[string]string{databaseKey: testDatabase, "SEREGA_HTTP_ADDR": "invalid"}, invalid: true},
		{name: "invalid cookies", values: map[string]string{databaseKey: testDatabase, "SEREGA_COOKIE_SECURE": "maybe"}, invalid: true},
		{name: "invalid proxy", values: map[string]string{databaseKey: testDatabase, "SEREGA_TRUSTED_PROXY_CIDRS": "localhost"}, invalid: true},
		{name: "invalid timeout", values: map[string]string{databaseKey: testDatabase, "SEREGA_SHUTDOWN_TIMEOUT": "0s"}, invalid: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := FromLookup(func(key string) (string, bool) { value, ok := test.values[key]; return value, ok })
			if (err != nil) != test.invalid {
				t.Fatalf("config error = %v", err)
			}

			if err == nil && !cfg.SecureCookies {
				t.Error("cookies must be secure by default")
			}
		})
	}
}
