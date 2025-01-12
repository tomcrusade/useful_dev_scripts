package api_vultr

import (
	"dev_scripts/adapters"
	"dev_scripts/entity"
	"fmt"
)

func (vultr *VultrAPI) CreateInstance(param entity.VultrAPIRequestCreateInstanceConfig) (
	entity.VultrInstance,
	entity.VultrAPIMeta,
	error,
) {
	response, err := adapters.CallApi[struct {
		Instance entity.VultrInstance `json:"instance"`
		Meta     entity.VultrAPIMeta  `json:"meta"`
	}](vultr.getAPICallParams(adapters.HttpMethodPost, "/instances", param))
	return response.Instance, response.Meta, err
}

func (vultr *VultrAPI) ListInstances() ([]entity.VultrInstance, entity.VultrAPIMeta, error) {
	response, err := adapters.CallApi[struct {
		Instances []entity.VultrInstance `json:"instances"`
		Meta      entity.VultrAPIMeta    `json:"meta"`
	}](vultr.getAPICallParams(adapters.HttpMethodGet, "/instances", struct{}{}))
	return response.Instances, response.Meta, err
}

func (vultr *VultrAPI) GetInstance(instanceId string) (*entity.VultrInstance, error) {
	response, err := adapters.CallApi[struct {
		Instance entity.VultrInstance `json:"instance"`
	}](
		vultr.getAPICallParams(
			adapters.HttpMethodGet,
			fmt.Sprintf("/instances/%v", instanceId),
			struct{}{},
		),
	)
	return &response.Instance, err
}

func (vultr *VultrAPI) HaltInstances(instanceIds []string) error {
	var param = struct {
		InstanceIds []string `json:"instance_ids"`
	}{InstanceIds: instanceIds}
	_, err := adapters.CallApi[struct{}](
		vultr.getAPICallParams(
			adapters.HttpMethodPost,
			"/instances/halt",
			param,
		),
	)
	return err
}

func (vultr *VultrAPI) RemoveInstance(instanceId string) error {
	_, err := adapters.CallApi[struct{}](
		vultr.getAPICallParams(
			adapters.HttpMethodDelete,
			fmt.Sprintf("/instances/%v", instanceId),
			struct{}{},
		),
	)
	return err
}
