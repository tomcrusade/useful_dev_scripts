package api_digitalocean

import (
	"dev_scripts/adapters"
	"dev_scripts/entity"
	"fmt"
)

type DigitaloceanAPI struct {
	env      *entity.EnvCloudServer
	tokenEnv *entity.EnvResourceToken
}

func (do *DigitaloceanAPI) getAPICallParams(method adapters.HttpMethod, urlPath string, requestParams interface{}) adapters.CallApiArgs {
	return adapters.CallApiArgs{
		FullPath:      fmt.Sprintf("https://api.digitalocean.com%v", urlPath),
		Token:         do.tokenEnv.DigitaloceanAPI,
		Method:        method,
		RequestParams: requestParams,
	}
}

func NewDigitaloceanAPI(env *entity.EnvCloudServer, tokenEnv *entity.EnvResourceToken) *DigitaloceanAPI {
	return &DigitaloceanAPI{env, tokenEnv}
}
