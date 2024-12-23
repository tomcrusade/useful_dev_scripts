package api_cloudflare

import (
	"dev_scripts/adapters"
	"dev_scripts/entity"
	"fmt"
)

func (api *CloudflareAPI) GetDNSRecordList(dnsZone string, searchQuery string) (
	[]entity.CloudflareDNSRecord,
	error,
) {
	result, err := CallCloudflareApi[[]entity.CloudflareDNSRecord](
		api.getAPICallParams(
			adapters.HttpMethodGet,
			fmt.Sprintf("/zones/%s/dns_records?%s", dnsZone, searchQuery),
			searchQuery,
		),
	)
	return result.Data, err
}

type CloudflareDNSRecordUpdateParams struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Type    string `json:"type"`
	Proxied bool   `json:"proxied,omitempty"`
	Comment string `json:"comment,omitempty"`
	TTL     int    `json:"ttl,omitempty"`
}

func (api *CloudflareAPI) UpdateDNSRecord(dnsZone string, dnsRecord string, data CloudflareDNSRecordUpdateParams) (
	entity.CloudflareDNSRecord,
	error,
) {
	result, err := CallCloudflareApi[entity.CloudflareDNSRecord](
		api.getAPICallParams(
			adapters.HttpMethodPut,
			fmt.Sprintf("/zones/%s/dns_records/%s", dnsZone, dnsRecord),
			data,
		),
	)
	return result.Data, err
}
