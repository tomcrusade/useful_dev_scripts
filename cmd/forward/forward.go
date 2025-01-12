package cmdforward

import (
	"dev_scripts/entity"
	"github.com/spf13/cobra"
)

type CmdForward struct{ CobraCmd *cobra.Command }

func NewCmdForward(env *entity.Env) *CmdForward {
	cmdHandler := &CmdForward{}
	cmdHandler.CobraCmd = &cobra.Command{
		Use:   "forward",
		Short: "Port Forward to existing services",
	}
	bastionCmd := NewCmdForwardBastion(env)
	cmdHandler.CobraCmd.AddCommand(bastionCmd.CobraCmd)
	return cmdHandler
}
