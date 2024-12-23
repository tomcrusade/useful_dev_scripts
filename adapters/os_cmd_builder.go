package adapters

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type OSCmdBuilder struct {
	command string
	args    []string
}

func NewOSCmdBuilder(command string, params []string) OSCmdBuilder {
	return OSCmdBuilder{
		command: command,
		args:    params,
	}
}

func (cmd OSCmdBuilder) Run() (string, error) {
	combinedCmd := cmd.command
	if len(cmd.args) > 0 {
		combinedCmd += " " + strings.Join(cmd.args, " ")
	}
	executor := exec.Command("bash", "-lc", combinedCmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	executor.Stdout = &out
	executor.Stderr = &stderr
	err := executor.Run()
	if err != nil {
		return out.String(), fmt.Errorf(
			"failed to execute \"%s\" because error: \"%s\", output: \"%s\"",
			combinedCmd,
			fmt.Sprint(err)+"  > "+stderr.String(),
			out.String(),
		)
	}
	return out.String(), nil
}

func (cmd OSCmdBuilder) RunWithInput(input string) (string, error) {
	combinedCmd := cmd.command
	if len(cmd.args) > 0 {
		combinedCmd += " " + strings.Join(cmd.args, " ")
	}
	executor := exec.Command("bash", "-c", combinedCmd)

	stdin, err := executor.StdinPipe()
	if err != nil {
		return "", fmt.Errorf(
			"failed to create input pipe for command: \"%s\" because error: \"%s\"",
			cmd.ToString(),
			fmt.Sprint(err),
		)
	}

	go func() {
		defer stdin.Close()
		parsedInput := strings.Replace(input, "\\", "\\\\", -1)
		io.WriteString(stdin, fmt.Sprintf("%s\n", parsedInput))
	}()

	output, err := executor.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(
			"failed to execute \"%s\" because error: \"%s\", output: \"%s\"",
			cmd.ToString(),
			fmt.Sprint(err),
			string(output),
		)
	}

	return string(output), nil
}

func (cmd OSCmdBuilder) ToString() string {
	return cmd.command + " " + strings.Join(cmd.args, " ")
}
