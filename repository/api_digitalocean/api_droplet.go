package api_digitalocean

import (
	"dev_scripts/adapters"
	"dev_scripts/entity"
	"fmt"
)

type DOCreateDropletRequest struct {
	Name         string                              `json:"name"`
	Region       string                              `json:"region"`
	Size         string                              `json:"size"`
	Image        string                              `json:"image"`
	SSHKeys      []string                            `json:"ssh_keys"`
	Backups      bool                                `json:"backups"`
	BackupPolicy *DOCreateDropletRequestBackupPolicy `json:"backup_policy,omitempty"`
	IPv6         bool                                `json:"ipv6"`
	Monitoring   bool                                `json:"monitoring"`
	Tags         []string                            `json:"tags"`
	UserData     string                              `json:"user_data"`
	VPCUUID      string                              `json:"vpc_uuid"`
}

type DOCreateDropletRequestBackupPolicy struct {
	Plan    string `json:"plan"`
	Weekday string `json:"weekday"`
	Hour    int    `json:"hour"`
}

type DOCreateDropletAPIResponse struct {
	Droplets []*entity.DigitaloceanDroplet `json:"droplets,omitempty"`
	Droplet  *entity.DigitaloceanDroplet   `json:"droplet,omitempty"`
	Links    struct {
		Actions []*entity.DigitaloceanDropletActionable `json:"actions"`
	} `json:"links"`
}

func (do *DigitaloceanAPI) CreateDroplets(param DOCreateDropletRequest) (
	*entity.DigitaloceanDroplet,
	[]*entity.DigitaloceanDropletActionable,
	error,
) {
	response, err := adapters.CallApi[DOCreateDropletAPIResponse](
		do.getAPICallParams(adapters.HttpMethodPost, "/v2/droplets", param),
	)

	if response.Droplet != nil {
		return response.Droplet, response.Links.Actions, err
	}

	return response.Droplets[0], response.Links.Actions, err
}

type DOListDropletAPIRequest struct {
	TagName string `json:"tag_name,omitempty"`
	Name    string `json:"name,omitempty"`
}

type DOListDropletAPIResponse struct {
	Droplets []*entity.DigitaloceanDroplet       `json:"droplets,omitempty"`
	Links    *entity.DigitaloceanAPIResultCursor `json:"links,omitempty"`
	Meta     *entity.DigitaloceanAPIResultMeta   `json:"meta,omitempty"`
}

func (do *DigitaloceanAPI) ListDroplets(filter DOListDropletAPIRequest) (
	[]*entity.DigitaloceanDroplet,
	*entity.DigitaloceanAPIResultCursor,
	*entity.DigitaloceanAPIResultMeta,
	error,
) {
	response, err := adapters.CallApi[DOListDropletAPIResponse](
		do.getAPICallParams(
			adapters.HttpMethodGet, "/v2/droplets",
			filter,
		),
	)
	return response.Droplets, response.Links, response.Meta, err
}

type getDropletAPIResponse struct {
	Droplet *entity.DigitaloceanDroplet `json:"droplet,omitempty"`
}

func (do *DigitaloceanAPI) GetDroplet(id string) (*entity.DigitaloceanDroplet, error) {
	response, err := adapters.CallApi[getDropletAPIResponse](
		do.getAPICallParams(
			adapters.HttpMethodGet,
			fmt.Sprintf("/v2/droplets/%v", id),
			struct{}{},
		),
	)
	return response.Droplet, err
}

func (do *DigitaloceanAPI) RemoveDroplet(id string) error {
	_, err := adapters.CallApi[struct{}](
		do.getAPICallParams(
			adapters.HttpMethodDelete,
			fmt.Sprintf("/v2/droplets/%v", id),
			struct{}{},
		),
	)
	return err
}

type getDropletActionAPIResponse struct {
	Action *entity.DigitaloceanDropletAction `json:"action"`
}

func (do *DigitaloceanAPI) GetDropletAction(id int, actionID int) (
	*entity.DigitaloceanDropletAction,
	error,
) {
	response, err := adapters.CallApi[getDropletActionAPIResponse](
		do.getAPICallParams(
			adapters.HttpMethodGet,
			fmt.Sprintf("/v2/droplets/%d/actions/%d", id, actionID),
			struct{}{},
		),
	)
	return response.Action, err
}