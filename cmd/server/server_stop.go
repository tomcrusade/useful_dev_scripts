package cmdserver

import (
	"dev_scripts/entity"
	"dev_scripts/usecase"
	"fmt"
	"github.com/spf13/cobra"
)

type Stop struct {
	CobraCmd *cobra.Command
	env      *entity.Env
}

func NewCmdServerStop(env *entity.Env) *Stop {
	cmdHandler := &Stop{env: env}
	cmdHandler.CobraCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stops (and snapshot) the server",
		Args:  cmdHandler.Args,
		RunE:  cmdHandler.RunE,
	}
	return cmdHandler
}

func (handler *Stop) Args(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}
	if _, ok := handler.env.VPS[args[0]]; !ok {
		var knownServerNames []string
		for serverName := range handler.env.VPS {
			knownServerNames = append(knownServerNames, serverName)
		}
		return fmt.Errorf(
			"server named %s not registered in env. Available servers: %s",
			args[0],
			knownServerNames,
		)
	}
	return nil
}

func (handler *Stop) RunE(_ *cobra.Command, args []string) error {
	serverConfig := handler.env.VPS[args[0]]
	vultrUseCase := usecase.NewCloudSvcVultr(serverConfig)

	return vultrUseCase.StopInstance()
}
