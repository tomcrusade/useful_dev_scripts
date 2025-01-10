package api_digitalocean

import (
	"dev_scripts/adapters"
	"fmt"
)

type DOAttachFirewallAPIResponse struct {
	ID        string `json:"id,omitempty"`
	Message   string `json:"message,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func (do *DigitaloceanAPI) AttachDropletToFirewall(firewallID string, dropletIDs []int) error {
	response, err := adapters.CallApi[DOAttachFirewallAPIResponse](
		do.getAPICallParams(
			adapters.HttpMethodPost, fmt.Sprintf("/v2/firewalls/%s/droplets", firewallID),
			struct {
				DropletIDs []int `json:"droplet_ids"`
			}{dropletIDs},
		),
	)

	if err != nil {
		return err
	}

	if response.ID != "" || response.Message != "" || response.RequestID != "" {
		return fmt.Errorf(response.Message)
	}

	return nil
}
