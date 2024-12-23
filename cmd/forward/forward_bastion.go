package cmdforward

import (
	"dev_scripts/entity"
	"dev_scripts/usecase"
	"fmt"
	"github.com/spf13/cobra"
)

type Bastion struct {
	CobraCmd *cobra.Command
	env      *entity.Env
}

func NewCmdForwardBastion(env *entity.Env) *Bastion {
	cmdHandler := &Bastion{env: env}
	cmdHandler.CobraCmd = &cobra.Command{
		Use:   "bastion",
		Short: "Connect resources through bastion",
		Args:  cmdHandler.Args,
		RunE:  cmdHandler.RunE,
	}
	return cmdHandler
}

func (handler *Bastion) Args(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(3)(cmd, args); err != nil {
		return err
	}

	var config *entity.EnvBastion
	var ok bool

	if config, ok = handler.env.Bastion[args[0]]; !ok {
		return fmt.Errorf("bastion config group %s not found", args[0])
	}

	if _, ok = config.DeviceURL[entity.CloudServiceEnvName(args[2])]; !ok {
		return fmt.Errorf("env %s not registered inside device_url configuration", args[2])
	}

	if _, ok = config.ResourceURL[entity.CloudServiceTechStackName(args[1])][entity.CloudServiceEnvName(args[2])]; !ok {
		return fmt.Errorf(
			"env %s stack %s not registered resource_url configuration",
			args[2],
			args[1],
		)
	}

	if _, ok = config.ResourcePort[entity.CloudServiceTechStackName(args[1])]; !ok {
		return fmt.Errorf(
			"stack %s not registered inside resource_port configuration",
			args[1],
		)
	}

	if _, ok = config.ResourceExposedPort[entity.CloudServiceTechStackName(args[1])][entity.CloudServiceEnvName(args[2])]; !ok {
		return fmt.Errorf(
			"env %s stack %s not registered resource_exposed_port configuration",
			args[2],
			args[1],
		)
	}

	return nil
}

func (handler *Bastion) RunE(_ *cobra.Command, args []string) error {
	bastion := usecase.NewCloudBastion(handler.env.Bastion[args[0]])
	if err := bastion.PortForward(
		entity.CloudServiceTechStackName(args[1]),
		entity.CloudServiceEnvName(args[2]),
	); err != nil {
		return err
	}
	return nil
}
