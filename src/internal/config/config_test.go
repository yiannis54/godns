package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfiguration(t *testing.T) {
	t.Run("should error without port", func(t *testing.T) {
		t.Setenv("SRV_PORT", "")
		_, err := LoadConfiguration()
		require.Error(t, err)
	})

	t.Run("should error with invalid port", func(t *testing.T) {
		t.Setenv("SRV_PORT", "abc")
		_, err := LoadConfiguration()
		require.Error(t, err)
	})

	t.Run("should succeed with valid port", func(t *testing.T) {
		t.Setenv("SRV_PORT", "1020")
		_, err := LoadConfiguration()
		require.NoError(t, err)
	})
}
