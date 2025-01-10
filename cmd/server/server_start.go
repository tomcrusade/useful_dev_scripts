package cmdserver

import (
	"dev_scripts/entity"
	usecase2 "dev_scripts/usecase"
	"fmt"
	"github.com/spf13/cobra"
)

type Start struct {
	CobraCmd *cobra.Command
	env      *entity.Env
}

func NewCmdServerStart(env *entity.Env) *Start {
	cmdHandler := &Start{env: env}
	cmdHandler.CobraCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts the server",
		Args:  cmdHandler.Args,
		RunE:  cmdHandler.RunE,
	}
	return cmdHandler
}

func (handler *Start) Args(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}
	if vps, ok := handler.env.VPS[args[0]]; !ok {
		var knownServerNames []string
		for serverName := range handler.env.VPS {
			knownServerNames = append(knownServerNames, serverName)
		}
		return fmt.Errorf(
			"server named %s not registered in env. Available servers: %s",
			args[0],
			knownServerNames,
		)
	} else if vps.VmBrand != entity.EnvCloudServerVmBrandDigitalOcean && vps.VmBrand != entity.EnvCloudServerVmBrandVultr {
		return fmt.Errorf("server brand %s is not known", vps.VmBrand)
	}
	return nil
}

func (handler *Start) RunE(_ *cobra.Command, args []string) error {
	vpsConfig := handler.env.VPS[args[0]]

	mainIP := ""
	if vpsConfig.VmBrand == entity.EnvCloudServerVmBrandVultr {
		vultrUseCase := usecase2.NewCloudSvcVultr(vpsConfig, handler.env.Tokens)
		chosenCloudVM, err := vultrUseCase.StartInstance()
		if err != nil {
			return err
		}
		mainIP = chosenCloudVM.MainIP
	} else {
		doUseCase := usecase2.NewCloudSvcDigitalocean(vpsConfig, handler.env.Tokens)
		chosenDroplet, err := doUseCase.StartInstance()
		if err != nil {
			return err
		}

		for _, v4Network := range chosenDroplet.Networks.V4 {
			if v4Network.Type == "public" {
				mainIP = v4Network.IPAddress
			}
		}
	}

	if mainIP != "" {
		cloudflareUseCase := usecase2.NewCloudSvcCloudflare(vpsConfig, handler.env.Tokens)
		if err := cloudflareUseCase.UpdateDNS(
			vpsConfig.DomainName,
			vpsConfig.SubdomainName,
			mainIP,
		); err != nil {
			return err
		}
	}

	return nil
}
