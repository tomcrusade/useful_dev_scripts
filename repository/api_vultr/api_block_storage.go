package api_vultr

import (
	"dev_scripts/adapters"
	"dev_scripts/entity"
	"fmt"
)

func (vultr *VultrAPI) GetBlockStorages() ([]entity.VultrBlockStorage, entity.VultrAPIMeta, error) {
	response, err := adapters.CallApi[struct {
		Blocks []entity.VultrBlockStorage `json:"blocks"`
		Meta   entity.VultrAPIMeta        `json:"meta"`
	}](vultr.getAPICallParams(adapters.HttpMethodGet, "/blocks", struct{}{}))
	return response.Blocks, response.Meta, err
}

func (vultr *VultrAPI) AttachBlockStorage(blockId string, instanceId string) error {
	var param = struct {
		InstanceId string `json:"instance_id"`
		Live       bool   `json:"live"`
	}{
		InstanceId: instanceId,
		Live:       true,
	}
	_, err := adapters.CallApi[struct{}](
		vultr.getAPICallParams(
			adapters.HttpMethodPost,
			fmt.Sprintf("/blocks/%v/attach", blockId),
			param,
		),
	)
	return err
}

func (vultr *VultrAPI) DetachBlockStorage(blockId string) error {
	var param = struct {
		Live bool `json:"live"`
	}{Live: true}
	_, err := adapters.CallApi[struct{}](
		vultr.getAPICallParams(
			adapters.HttpMethodPost,
			fmt.Sprintf("/blocks/%v/detach", blockId),
			param,
		),
	)
	return err
}
