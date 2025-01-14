package api_digitalocean

import (
	"dev_scripts/adapters"
	"dev_scripts/entity"
	"fmt"
)

type DigitaloceanAPI struct {
	env *entity.EnvCloudServer
}

func (do *DigitaloceanAPI) getAPICallParams(method adapters.HttpMethod, urlPath string, requestParams interface{}) adapters.CallApiArgs {
	return adapters.CallApiArgs{
		FullPath:      fmt.Sprintf("https://api.digitalocean.com%v", urlPath),
		Token:         do.env.VmTypeApiToken,
		Method:        method,
		RequestParams: requestParams,
	}
}

func NewDigitaloceanAPI(env *entity.EnvCloudServer) *DigitaloceanAPI {
	return &DigitaloceanAPI{env}
}
