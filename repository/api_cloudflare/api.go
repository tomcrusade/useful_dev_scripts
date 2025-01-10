package api_cloudflare

import (
	"dev_scripts/adapters"
	"dev_scripts/entity"
	"encoding/json"
	"errors"
	"fmt"
)

type CloudflareAPI struct {
	env      *entity.EnvCloudServer
	tokenEnv *entity.EnvResourceToken
}

func (api *CloudflareAPI) getAPICallParams(method adapters.HttpMethod, urlPath string, requestParams interface{}) adapters.CallApiArgs {
	return adapters.CallApiArgs{
		FullPath:      fmt.Sprintf("https://api.cloudflare.com/client/v4%v", urlPath),
		Token:         api.tokenEnv.CloudflareAPI,
		Method:        method,
		RequestParams: requestParams,
	}
}

func NewCloudflareAPI(env *entity.EnvCloudServer, tokenEnv *entity.EnvResourceToken) *CloudflareAPI {
	return &CloudflareAPI{env, tokenEnv}
}

type CloudflareAPIResponse[data interface{}] struct {
	Data     data
	Messages []struct {
		Code    string
		Message string
	}
}

type cloudflareResponse[data interface{}] struct {
	Result data `json:"result,omitempty"`
	Errors []struct {
		Code    string
		Message string
	} `json:"errors,omitempty"`
	Messages []struct {
		Code    string
		Message string
	} `json:"messages,omitempty"`
	Success bool `json:"success,omitempty"`
}

func CallCloudflareApi[ResponseMap interface{}](params adapters.CallApiArgs) (
	CloudflareAPIResponse[ResponseMap],
	error,
) {
	response, err := adapters.CallApi[cloudflareResponse[ResponseMap]](params)
	if len(response.Errors) > 0 {
		var errorMessages []byte
		errorMessages, err = json.Marshal(response.Errors)
		if err == nil {
			err = errors.New(string(errorMessages))
		}
	}
	return CloudflareAPIResponse[ResponseMap]{
		Data:     response.Result,
		Messages: response.Messages,
	}, err
}
