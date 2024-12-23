package api_cloudflare

import (
	"dev_scripts/adapters"
	"dev_scripts/entity"
)

func (api *CloudflareAPI) GetZoneList(searchQuery string) ([]entity.CloudflareZone, error) {
	result, err := CallCloudflareApi[[]entity.CloudflareZone](
		api.getAPICallParams(adapters.HttpMethodGet, "/zones", searchQuery),
	)
	return result.Data, err
}
