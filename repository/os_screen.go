package repository

import "dev_scripts/adapters"

func GetScreen() ScreenExecutable {
	return ScreenExecutable{
		builder: func(args []string) adapters.OSCmdBuilder {
			return adapters.NewOSCmdBuilder("screen", args)
		},
	}
}

type ScreenExecutable struct {
	builder func(args []string) adapters.OSCmdBuilder
}

func (cmd ScreenExecutable) SetAsCreate(screenName string) adapters.OSCmdBuilder {
	return cmd.builder([]string{"-dmS", screenName})
}

func (cmd ScreenExecutable) SetAsDelete(screenName string) adapters.OSCmdBuilder {
	return cmd.builder([]string{"-S", screenName, "-X", "quit"})
}

func (cmd ScreenExecutable) SetAsSendCommand(command string, onScreenName string) adapters.OSCmdBuilder {
	return cmd.builder([]string{"-S", onScreenName, "-p", "0", "-X", "stuff", "'" + command + " ^M'"})
}
