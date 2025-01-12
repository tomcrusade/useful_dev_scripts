package repository

import (
	"dev_scripts/adapters"
	"fmt"
	"strings"
)

func SecretToolStore(attributeKey string, attributeVal string, label string, storedValue string) error {
	var commandArgs []string
	if label != "" {
		commandArgs = append(commandArgs, fmt.Sprintf("--label %s", label))
	}
	if storedValue == "" {
		return fmt.Errorf("no value stored to new key value")
	}
	if attributeKey == "" || attributeVal == "" {
		return fmt.Errorf("attributeKey or attributeVal is required")
	}
	commandArgs = append(commandArgs, fmt.Sprintf("%s %s", attributeKey, attributeVal))

	_, err := adapters.NewOSCmdBuilder(
		fmt.Sprintf("secret-tool store %s", strings.Join(commandArgs, " ")),
		[]string{},
	).RunWithInput(storedValue)
	if err != nil {
		return fmt.Errorf("failed to store secret to GNOME keyring because: %v", err.Error())
	}
	return nil
}

func SecretToolGet(attributeKey string, attributeVal string, label string) (string, error) {
	var commandArgs []string
	if label != "" {
		commandArgs = append(commandArgs, fmt.Sprintf("--label %s", label))
	}
	if attributeKey == "" || attributeVal == "" {
		return "", fmt.Errorf("attributeKey or attributeVal is required")
	}
	commandArgs = append(commandArgs, fmt.Sprintf("%s %s", attributeKey, attributeVal))

	result, err := adapters.NewOSCmdBuilder(
		fmt.Sprintf("secret-tool lookup %s", strings.Join(commandArgs, " ")),
		[]string{},
	).Run()
	if err != nil {
		return "", fmt.Errorf(
			"failed to get secret to GNOME keyring because: %v",
			err.Error(),
		)
	}
	return result, nil
}
