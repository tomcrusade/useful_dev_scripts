package cmdserver

import (
	"dev_scripts/entity"
	"github.com/spf13/cobra"
)

type CmdServer struct{ CobraCmd *cobra.Command }

func NewCmdServer(env *entity.Env) *CmdServer {
	cmdHandler := &CmdServer{}
	cmdHandler.CobraCmd = &cobra.Command{
		Use:   "server",
		Short: "Manage Cloud Resources (currently only vultr, and cloudflare)",
	}
	startServerCmd := NewCmdServerStart(env)
	cmdHandler.CobraCmd.AddCommand(startServerCmd.CobraCmd)
	stopServerCmd := NewCmdServerStop(env)
	cmdHandler.CobraCmd.AddCommand(stopServerCmd.CobraCmd)
	return cmdHandler
}
