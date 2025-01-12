package usecase

import (
	"dev_scripts/entity"
	"dev_scripts/repository"
	"encoding/json"
	"errors"
	"fmt"
)

func LoadEnvFromFile() (*entity.Env, error) {
	var err error
	var rawConfig string

	rawConfig, err = repository.SecretToolGet(
		"tomcrusade_aio_commands",
		"config",
		"",
	)
	if rawConfig == "" || err != nil {
		fmt.Printf(err.Error())
		newConfig, err := json.Marshal(entity.Env{})
		if err = repository.SecretToolStore(
			"tomcrusade_aio_commands",
			"config",
			"aio_commands_config",
			string(newConfig),
		); err != nil {
			return nil, errors.New(
				fmt.Sprintf(
					"unable to save new secrets tomcrusade_aio_commands config with value %s because: %s",
					string(newConfig),
					err,
				),
			)
		}
		return nil, errors.New("secrets tomcrusade_aio_commands config is empty, creating it.... (edit the secret using ubuntu \"Passwords and Keys\", then run this command again)")
	}

	cfg := &entity.Env{}
	if err = json.Unmarshal([]byte(rawConfig), &cfg); err != nil {
		return nil, errors.New(
			fmt.Sprintf(
				"unable to parse config due to wrong format. details: %s",
				err,
			),
		)
	}

	return cfg, nil
}
