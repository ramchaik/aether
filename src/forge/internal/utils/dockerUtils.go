package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

func RunDockerCommand(command string, args ...string) error {
	argsSlice := []string{command}
	argsSlice = append(argsSlice, args...)

	// Execute the command.
	cmd := exec.Command(argsSlice[0], argsSlice[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to execute docker command: %w\nOutput:\n%s", err, out.String())
	}

	return nil
}

func CloneRepo(repoURL, tempDir string) error {
	return RunDockerCommand("clone", repoURL, tempDir)
}

func BuildAndRunDockerContainer(tempDir string) error {
	return RunDockerCommand("build", ".", "-t", "secure-build-env")
}
