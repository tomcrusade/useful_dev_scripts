package usecase

import (
	"dev_scripts/adapters"
	"dev_scripts/entity"
	"dev_scripts/repository"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CloudBastion struct {
	env *entity.EnvBastion
}

func NewCloudBastion(env *entity.EnvBastion) *CloudBastion {
	return &CloudBastion{env}
}

// --

func (t *CloudBastion) PortForward(stackName entity.CloudServiceTechStackName, svcEnvName entity.CloudServiceEnvName) error {
	destinationPort, err := t.setupPortForward(stackName, svcEnvName)
	if err != nil {
		fmt.Println()
		fmt.Println("Perform gcloud login")
		if err := t.googleCloudLogin(t.env.SShCertFile, 7200); err != nil {
			fmt.Println(t.env)
			return err
		}
		destinationPort, err = t.setupPortForward(stackName, svcEnvName)
		if err != nil {
			return fmt.Errorf("cannot port forward even after google cloud login because: %v", err)
		}
	}
	if stackName == entity.CloudServiceTechStackNameVault {
		fmt.Println("Login onto vault")
		if err := repository.VaultLogin(destinationPort); err != nil {
			return err
		}
	}
	return nil
}

func (t *CloudBastion) setupPortForward(stackName entity.CloudServiceTechStackName, svcEnvName entity.CloudServiceEnvName) (
	int,
	error,
) {
	var err error

	exposedPort := t.env.ResourceExposedPort[stackName][svcEnvName]

	fmt.Print("checking")
	if err = t.verifyUntilPortIsActive(strconv.Itoa(exposedPort), 5, 1); err == nil {
		fmt.Printf("Port %d already active \n", exposedPort)
		return exposedPort, nil
	}

	err = repository.NewSSHTunnel(
		t.env.ResourcePort[stackName],
		exposedPort,
		t.env.DeviceURL[svcEnvName],
		t.env.ResourceURL[stackName][svcEnvName],
		"~/.ssh/config",
		t.env.SShCertFile,
	)
	if err != nil {
		return exposedPort, err
	}

	fmt.Print("verifying")
	if err = t.verifyUntilPortIsActive(strconv.Itoa(exposedPort), 6, 1); err != nil {
		return exposedPort, err
	}
	fmt.Println()

	return exposedPort, nil
}

func (t *CloudBastion) googleCloudLogin(keyFile string, ttl int) error {
	_, err := adapters.NewOSCmdBuilder(
		fmt.Sprintf(
			"gcloud compute os-login ssh-keys add --key-file %s --ttl %d",
			keyFile,
			ttl,
		), []string{},
	).Run()
	if err != nil {
		return fmt.Errorf("failed to login google because: %v", err.Error())
	}
	return nil
}

func (t *CloudBastion) verifyUntilPortIsActive(port string, maxRetries int, retries int) error {
	if retries > maxRetries {
		return fmt.Errorf(
			"\n unable to verify port %s to be open after %d seconds/times",
			port,
			retries-1,
		)
	}
	result, err := adapters.NewOSCmdBuilder(
		fmt.Sprintf("nc -z 127.0.0.1 %s && echo \"I\" || echo \"H\"", port),
		[]string{},
	).Run()
	fmt.Printf(" ...%d", retries)

	if err != nil || strings.Trim(result, "\n\r") != "I" {
		time.Sleep(1 * time.Second)
		return t.verifyUntilPortIsActive(port, maxRetries, retries+1)
	}
	return nil
}
