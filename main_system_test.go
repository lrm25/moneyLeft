//go:build system
// +build system

package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func execProgram(t *testing.T, fileLocation string) *bytes.Buffer {

	wd, err := os.Getwd()
	require.NoError(t, err)

	fullExePath := filepath.Join(wd, "moneyLeft.exe")

	cmd := exec.Command(fullExePath, "-"+FLAG_FILE_LOCATION, fileLocation)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err = cmd.Run()
	require.NoError(t, err)
	return &stdout
}

func Test_main_fileNotFound(t *testing.T) {

	stdout := execProgram(t, "nowhere")
	stdoutBytes := stdout.Bytes()
	stdoutString := string(stdoutBytes)
	t.Log(stdoutString)

	require.Contains(t, stdoutString, "The system cannot find the file specified")
}

func Test_main_moneyPerMonthLessThanZero(t *testing.T) {

	testJsonFile := "moneyLeft_test_" + time.Now().Format("20060102150405") + ".json"
	testJsonFullPath := filepath.Join(os.TempDir(), testJsonFile)
	// omit required parameters
	jsonBytes := []byte("{\"totalMoney\": 500}")
	err := os.WriteFile(testJsonFullPath, jsonBytes, 0644)
	require.NoError(t, err)

	stdout := execProgram(t, testJsonFullPath)
	stdoutBytes := stdout.Bytes()
	stdoutString := string(stdoutBytes)
	t.Log(stdoutString)

	require.Equal(t, ERROR_MONEYPERMONTH, stdoutString, "Errors should match")
}
