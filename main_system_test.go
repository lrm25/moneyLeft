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

func execProgram(t *testing.T, fileLocation string) string {

	wd, err := os.Getwd()
	require.NoError(t, err)

	fullExePath := filepath.Join(wd, "moneyLeft.exe")

	cmd := exec.Command(fullExePath, "-"+FLAG_FILE_LOCATION, fileLocation)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()
	require.NoError(t, err)

	stdoutBytes := stdout.Bytes()
	return string(stdoutBytes)
}

func getJsonFileName(t *testing.T) string {
	testJsonFile := t.Name() + "_" + time.Now().Format("20060102150405") + ".json"
	return filepath.Join(os.TempDir(), testJsonFile)
}

func Test_main_fileNotFound(t *testing.T) {

	stdoutString := execProgram(t, "nowhere")
	t.Log(stdoutString)

	require.Contains(t, stdoutString, "The system cannot find the file specified")
}

func Test_main_moneyPerMonthJsonFormat(t *testing.T) {

	testJsonFullPath := getJsonFileName(t)
	jsonBytes := []byte("{\"totalMoney\": \"500\"}")
	err := os.WriteFile(testJsonFullPath, jsonBytes, 0644)
	require.NoError(t, err)

	stdoutString := execProgram(t, testJsonFullPath)
	t.Log(stdoutString)

	require.Contains(t, stdoutString, "json: cannot unmarshal")
}

func Test_main_moneyPerMonthLessThanZero(t *testing.T) {

	testJsonFullPath := getJsonFileName(t)
	// omit required parameters
	jsonBytes := []byte("{\"totalMoney\": 500}")
	err := os.WriteFile(testJsonFullPath, jsonBytes, 0644)
	require.NoError(t, err)

	stdoutString := execProgram(t, testJsonFullPath)
	t.Log(stdoutString)

	require.Equal(t, ERROR_MONEYPERMONTH, stdoutString, "Errors should match")
}

func Test_main_valid(t *testing.T) {

	testJsonFullPath := getJsonFileName(t)
	// omit required parameters
	jsonBytes := []byte("{\"totalMoney\": 500, \"moneyPerMonth\": 100}")
	err := os.WriteFile(testJsonFullPath, jsonBytes, 0644)
	require.NoError(t, err)

	stdoutString := execProgram(t, testJsonFullPath)
	require.Contains(t, stdoutString, "5.00")
}
