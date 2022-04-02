package efibootmgr

import (
	"fmt"
	"os/exec"
	"strings"
)

func Execute(flags []Flag, options map[Option]string, args ...string) error {
	_, err := ExecuteWithOutput(flags, options, args...)

	return err
}

func ExecuteWithOutput(flags []Flag, options map[Option]string, args ...string) (string, error) {
	cmdParts := []string{"efibootmgr"}
	for _, flag := range flags {
		cmdParts = append(cmdParts, fmt.Sprintf("--%s", string(flag)))
	}
	for option, value := range options {
		cmdParts = append(cmdParts, fmt.Sprintf("--%s", string(option)), value)
	}
	cmdParts = append(cmdParts, args...)
	fmt.Printf("* %s\n", strings.Join(cmdParts, " "))

	//nolint:gosec
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	d, err := cmd.CombinedOutput()

	out := string(d)
	lines := strings.Split(out, "\n")
	lines = lines[:len(lines)-1]
	for _, line := range lines {
		fmt.Printf("+ %s\n", line)
	}

	return strings.Join(lines, "\n"), err
}
