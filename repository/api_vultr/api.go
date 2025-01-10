package api_vultr

import (
	"dev_scripts/adapters"
	"dev_scripts/entity"
	"fmt"
)

type VultrAPI struct {
	env      *entity.EnvCloudServer
	tokenEnv *entity.EnvResourceToken
}

func (vultr *VultrAPI) getAPICallParams(method adapters.HttpMethod, urlPath string, requestParams interface{}) adapters.CallApiArgs {
	return adapters.CallApiArgs{
		FullPath:      fmt.Sprintf("https://api.vultr.com/v2%v", urlPath),
		Token:         vultr.tokenEnv.VultrAPI,
		Method:        method,
		RequestParams: requestParams,
	}
}

func NewVultrAPI(env *entity.EnvCloudServer, tokenEnv *entity.EnvResourceToken) *VultrAPI {
	return &VultrAPI{env, tokenEnv}
}
