package repository

import (
	"fmt"
	"strings"
)

func NewSSHTunnel(originPort int, destinationPort int, originUrl string, destinationUrl string, customConfigPath string, certificateFile string) error {
	var err error
	screenName := fmt.Sprintf("fwd_stack_%d", destinationPort)
	screenCmd := GetScreen()
	if delScreenOutput, err := screenCmd.SetAsDelete(screenName).Run(); err != nil && strings.TrimSpace(delScreenOutput) != "No screen session found." {
		return err
	}
	if _, err = screenCmd.SetAsCreate(screenName).Run(); err != nil {
		return err
	}

	configPath := "~/.ssh/config"
	if customConfigPath != "" {
		configPath = customConfigPath
	}

	if certificateFile != "" {
		_, err = screenCmd.SetAsSendCommand(
			fmt.Sprintf(
				"ssh -NL %d:%s:%d %s -F %s -v -o CertificateFile=%s",
				destinationPort,
				destinationUrl,
				originPort,
				originUrl,
				configPath,
				certificateFile,
			),
			screenName,
		).Run()
	} else {
		_, err = screenCmd.SetAsSendCommand(
			fmt.Sprintf(
				"ssh -NL %d:%s:%d %s -F %s -v",
				destinationPort,
				destinationUrl,
				originPort,
				originUrl,
				configPath,
			),
			screenName,
		).Run()
	}

	if err != nil {
		return err
	}
	return nil
}
