package api_vultr

import (
	"dev_scripts/adapters"
	"dev_scripts/entity"
)

func (vultr *VultrAPI) ListFirewallGroups() (
	[]entity.VultrFirewallGroup,
	entity.VultrAPIMeta,
	error,
) {
	response, err := adapters.CallApi[struct {
		FirewallGroups []entity.VultrFirewallGroup `json:"firewall_groups"`
		Meta           entity.VultrAPIMeta         `json:"meta"`
	}](vultr.getAPICallParams(adapters.HttpMethodGet, "/firewalls", struct{}{}))
	return response.FirewallGroups, response.Meta, err
}
