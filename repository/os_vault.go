package repository

import (
	"dev_scripts/adapters"
	"fmt"
)

func VaultLogin(port int) error {
	_, err := adapters.NewOSCmdBuilder(
		fmt.Sprintf(
			"export VAULT_ADDR=https://localhost:%d; export VAULT_SKIP_VERIFY=true;vault login -method=oidc",
			port,
		), []string{},
	).Run()
	if err != nil {
		return fmt.Errorf("failed to login to vault because: %v", err.Error())
	}
	return nil
}
