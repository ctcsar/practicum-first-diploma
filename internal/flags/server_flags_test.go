package flags

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServerURL(t *testing.T) {
	f := &serverFlags{}
	t.Run("GetServerUrlFromEnv", func(t *testing.T) {
		os.Setenv("ADDRESS", "localhost:9010")
		// Check that the flags are set correctly
		assert.Equal(t, "localhost:9010", f.GetServerURL())
	})
	t.Run("GetEmptyServerURL", func(t *testing.T) {
		os.Unsetenv("ADDRESS")
		f.url = ""
		// Check that the flags are set correctly
		assert.Equal(t, "", f.GetServerURL())
	})

	t.Run("GetCustomServerURL", func(t *testing.T) {
		f.url = "localhost:9000"
		// Check that the flags are set correctly
		assert.Equal(t, "localhost:9000", f.GetServerURL())
	})

}

func TestGetAccural(t *testing.T) {
	f := &serverFlags{}
	t.Run("GetAccuralFromEnv", func(t *testing.T) {
		os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "localhost:9010")
		// Check that the flags are set correctly
		assert.Equal(t, "localhost:9010", f.GetAccural())
	})
	t.Run("GetEmptyAccural", func(t *testing.T) {
		os.Unsetenv("ACCRUAL_SYSTEM_ADDRESS")
		f.accural = ""
		// Check that the flags are set correctly
		assert.Equal(t, "", f.GetAccural())
	})

	t.Run("GetCustomAccural", func(t *testing.T) {
		f.accural = "localhost:9000"
		// Check that the flags are set correctly
		assert.Equal(t, "localhost:9000", f.GetAccural())
	})

}

func TestGetDatabasePath(t *testing.T) {
	f := &serverFlags{}
	t.Run("GetDatabasePathFromEnv", func(t *testing.T) {
		os.Setenv("DATABASE_DSN", "postgres://metrics:password@localhost:5432/metrics?sslmode=disable")
		// Check that the flags are set correctly
		assert.Equal(t, "postgres://metrics:password@localhost:5432/metrics?sslmode=disable", f.GetDatabasePath())
	})
	t.Run("GetEmptyDatabasePath", func(t *testing.T) {
		os.Unsetenv("DATABASE_DSN")
		f.databaseDSN = ""
		// Check that the flags are set correctly
		assert.Equal(t, "", f.GetDatabasePath())
	})

	t.Run("GetCustomDatabasePath", func(t *testing.T) {
		f.databaseDSN = "postgres://metrics:password@localhost:5432/metrics?sslmode=disable"
		// Check that the flags are set correctly
		assert.Equal(t, "postgres://metrics:password@localhost:5432/metrics?sslmode=disable", f.GetDatabasePath())
	})

}
