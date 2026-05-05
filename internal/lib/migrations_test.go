package lib

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrationsDir_Success(t *testing.T) {
	dir, err := MigrationsDir()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join("migrations"), filepath.Base(dir))

	info, err := os.Stat(dir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestMigrationsDir_NotFound(t *testing.T) {
	tmp := t.TempDir()
	t.Chdir(tmp)

	_, err := MigrationsDir()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not find project root containing go.mod")
}
