package api_vultr

import (
	"dev_scripts/adapters"
	"dev_scripts/entity"
)

func (vultr *VultrAPI) GetReservedIps() ([]entity.VultrReservedIp, entity.VultrAPIMeta, error) {
	response, err := adapters.CallApi[struct {
		ReservedIps []entity.VultrReservedIp `json:"reserved_ips"`
		Meta        entity.VultrAPIMeta      `json:"meta"`
	}](vultr.getAPICallParams(adapters.HttpMethodGet, "/reserved-ips", struct{}{}))
	return response.ReservedIps, response.Meta, err
}
