package configs

import (
	"flag"
	"github.com/stretchr/testify/require"
	//"github.com/vladislaoramos/gophermart/configs"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	// Setting environment variables
	err := os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "http://localhost:8081")
	require.NoError(t, err)

	err = os.Setenv("RUN_ADDRESS", "0.0.0.0:8080")
	require.NoError(t, err)

	err = os.Setenv("DATABASE_URI", "postgres://user:pass@localhost/gophermart")
	require.NoError(t, err)

	// Generating the config
	cfg := NewConfig()

	// Checking the correctness of parsing environment variables
	require.Equal(t, "http://localhost:8081", cfg.App.AccrualSystemAddress)
	require.Equal(t, "0.0.0.0:8080", cfg.Client.Address)
	require.Equal(t, "postgres://user:pass@localhost/gophermart", cfg.Database.URI)

	// Checking the correctness of the default values
	require.Equal(t, "gophermart", cfg.App.Name)
	require.Equal(t, "0.0.1", cfg.App.Version)
	require.Equal(t, "debug", cfg.Logger.Level)

	// Checking the correctness of flag parsing
	err = flag.Set("a", "0.0.0.0:8888")
	require.NoError(t, err)

	err = flag.Set("d", "postgres://user:pass@localhost/otherdb")
	require.NoError(t, err)

	err = flag.Set("r", "http://localhost:8082")
	require.NoError(t, err)

	flag.Parse()

	flags := &Config{
		App: App{
			AccrualSystemAddress: flag.Lookup("r").Value.String(),
		},
		Database: Database{
			URI: flag.Lookup("d").Value.String(),
		},
		Client: Client{
			Address: flag.Lookup("a").Value.String(),
		},
	}

	cfg.updateConfig(flags)

	require.Equal(t, "http://localhost:8082", cfg.App.AccrualSystemAddress)
	require.Equal(t, "0.0.0.0:8888", cfg.Client.Address)
	require.Equal(t, "postgres://user:pass@localhost/otherdb", cfg.Database.URI)
}
