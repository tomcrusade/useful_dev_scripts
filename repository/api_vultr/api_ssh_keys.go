package api_vultr

import (
	"dev_scripts/adapters"
	"dev_scripts/entity"
)

func (vultr *VultrAPI) GetSSHKeys() ([]entity.VultrSSHKey, entity.VultrAPIMeta, error) {
	response, err := adapters.CallApi[struct {
		SSHKeys []entity.VultrSSHKey `json:"ssh_keys"`
		Meta    entity.VultrAPIMeta  `json:"meta"`
	}](vultr.getAPICallParams(adapters.HttpMethodGet, "/ssh-keys", struct{}{}))
	return response.SSHKeys, response.Meta, err
}
