package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables for testing
	PORT := "PORT"
	TRUSTED_PROXIES := "TRUSTED_PROXIES"
	port := "4000"
	expectedPort := 4000
	trustedProxies := "192.168.1.1, 10.0.0.1"
	expectedTrustedProxies := []string{"192.168.1.1", "10.0.0.1"}

	os.Setenv(PORT, port)
	os.Setenv(TRUSTED_PROXIES, trustedProxies)
	defer os.Unsetenv(PORT)
	defer os.Unsetenv(TRUSTED_PROXIES)

	cfg := LoadConfig()

	assert.Equal(t, expectedPort, cfg.PORT)
	assert.Equal(t, expectedTrustedProxies, cfg.TRUSTED_PROXIES)
}

func TestGetenv(t *testing.T) {
	t.Run("String default", func(t *testing.T) {
		val := getenv("NON_EXISTENT", "default")
		assert.Equal(t, "default", val)
	})

	t.Run("String env", func(t *testing.T) {
		os.Setenv("TEST_KEY", "test_val")
		defer os.Unsetenv("TEST_KEY")
		val := getenv("TEST_KEY", "default")
		assert.Equal(t, "test_val", val)
	})

	t.Run("Int env", func(t *testing.T) {
		os.Setenv("TEST_INT", "123")
		defer os.Unsetenv("TEST_INT")
		val := getenv("TEST_INT", 0)
		assert.Equal(t, 123, val)
	})

	t.Run("[]string env", func(t *testing.T) {
		os.Setenv("TEST_LIST", "a, b, , c")
		defer os.Unsetenv("TEST_LIST")
		val := getenv("TEST_LIST", []string{})
		assert.Equal(t, []string{"a", "b", "c"}, val)
	})
}
