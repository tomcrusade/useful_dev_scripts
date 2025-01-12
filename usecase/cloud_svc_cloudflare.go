package usecase

import (
	"dev_scripts/entity"
	"dev_scripts/repository/api_cloudflare"
	"fmt"
)

type CloudSvcCloudflare struct {
	env            *entity.EnvCloudServer
	tokenEnv       *entity.EnvResourceToken
	cloudflareRepo *api_cloudflare.CloudflareAPI
}

func NewCloudSvcCloudflare(env *entity.EnvCloudServer, tokenEnv *entity.EnvResourceToken) *CloudSvcCloudflare {
	cloudflareRepo := api_cloudflare.NewCloudflareAPI(env, tokenEnv)
	return &CloudSvcCloudflare{env, tokenEnv, cloudflareRepo}
}

// --

func (uc *CloudSvcCloudflare) UpdateDNS(domainName string, subdomainName string, newIPAddress string) error {
	searchName := domainName
	if subdomainName != "" {
		searchName = subdomainName + "." + domainName
	}
	zoneList, err := uc.cloudflareRepo.GetZoneList(fmt.Sprintf("name=%s", domainName))
	if err != nil {
		return fmt.Errorf("failed to get dns zone list. Error: %v", err)
	}
	if len(zoneList) <= 0 {
		return fmt.Errorf("no zone found for domain %s", domainName)
	}
	fmt.Printf(
		"Zone found : %s \n",
		entity.ConvertToJSON[[]entity.CloudflareZone](zoneList),
	)

	dnsRecordList, err := uc.cloudflareRepo.GetDNSRecordList(
		zoneList[0].ID,
		fmt.Sprintf("name=%s&type=A", searchName),
	)
	if err != nil {
		return fmt.Errorf("cannot find DNS record because: %v", err)
	}
	if len(dnsRecordList) <= 0 {
		return fmt.Errorf("no DNS record found")
	}
	fmt.Printf(
		"DNS record found : %s \n",
		entity.ConvertToJSON[[]entity.CloudflareDNSRecord](dnsRecordList),
	)

	updateDNSResult, err := uc.cloudflareRepo.UpdateDNSRecord(
		zoneList[0].ID,
		dnsRecordList[0].ID,
		api_cloudflare.CloudflareDNSRecordUpdateParams{
			Name:    searchName,
			Content: newIPAddress,
			Type:    "A",
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update DNS record because: %v", err)
	}
	fmt.Printf(
		"DNS record updated to : %s \n",
		entity.ConvertToJSON[entity.CloudflareDNSRecord](updateDNSResult),
	)
	return nil
}
