//go:build system
// +build system

package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func execProgram(t *testing.T) *bytes.Buffer {

	wd, err := os.Getwd()
	require.NoError(t, err)

	fullExePath := filepath.Join(wd, "moneyLeft.exe")

	cmd := exec.Command(fullExePath, "-"+FLAG_FILE_LOCATION, "nowhere")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err = cmd.Run()
	require.NoError(t, err)
	return &stdout
}

func Test_main_fileNotFound(t *testing.T) {

	stdout := execProgram(t)
	stdoutBytes := stdout.Bytes()
	stdoutString := string(stdoutBytes)
	t.Log(stdoutString)

	require.Contains(t, stdoutString, "The system cannot find the file specified")
}
