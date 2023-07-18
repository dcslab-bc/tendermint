package commands

import (
	"testing"

	cfg "github.com/Finschia/ostracon/config"
	"github.com/stretchr/testify/require"
)

func TestRollbackStateCmd(t *testing.T) {
	config = cfg.TestConfig()
	err := RollbackStateCmd.RunE(RollbackStateCmd, nil)
	require.Error(t, err)
}
