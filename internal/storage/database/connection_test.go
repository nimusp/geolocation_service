package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildConnString(t *testing.T) {
	t.Setenv(envHost, "db_host")
	t.Setenv(envPort, "8888")
	t.Setenv(envUser, "db_user")
	t.Setenv(envPassword, "db_passwd")

	want := "sslmode=disable host=db_host port=8888 user=db_user password=db_passwd"

	got := buildConnString()

	assert.Equal(t, want, got)
}

func Test_readMigrations(t *testing.T) {
	want := []string{"hello", "world"}
	dirName := t.TempDir()
	firstFileName := filepath.Join(dirName, "first")
	secondFileName := filepath.Join(dirName, "second")

	f1, err := os.OpenFile(firstFileName, os.O_CREATE|os.O_WRONLY, 0666)
	assert.NoError(t, err)
	_, err = f1.WriteString(want[0])
	assert.NoError(t, err)
	err = f1.Close()
	assert.NoError(t, err)

	f2, err := os.OpenFile(secondFileName, os.O_CREATE|os.O_WRONLY, 0666)
	assert.NoError(t, err)
	_, err = f2.WriteString(want[1])
	assert.NoError(t, err)
	err = f2.Close()
	assert.NoError(t, err)

	got, err := readMigrations(dirName)
	assert.NoError(t, err)

	assert.Equal(t, want, got)
}
