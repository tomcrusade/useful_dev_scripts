package api_vultr

import (
	"dev_scripts/adapters"
	"dev_scripts/entity"
)

func (vultr *VultrAPI) GetISOs() ([]entity.VultrISO, entity.VultrAPIMeta, error) {
	response, err := adapters.CallApi[struct {
		SSHKeys []entity.VultrISO   `json:"isos"`
		Meta    entity.VultrAPIMeta `json:"meta"`
	}](vultr.getAPICallParams(adapters.HttpMethodGet, "/iso", struct{}{}))
	return response.SSHKeys, response.Meta, err
}
