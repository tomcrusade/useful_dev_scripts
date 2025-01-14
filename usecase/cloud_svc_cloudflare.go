package usecase

import (
	"dev_scripts/entity"
	"dev_scripts/repository/api_cloudflare"
	"fmt"
)

type CloudSvcCloudflare struct {
	env            *entity.EnvCloudServerDomain
	cloudflareRepo *api_cloudflare.CloudflareAPI
}

func NewCloudSvcCloudflare(env *entity.EnvCloudServerDomain) *CloudSvcCloudflare {
	cloudflareRepo := api_cloudflare.NewCloudflareAPI(env)
	return &CloudSvcCloudflare{env, cloudflareRepo}
}

// --

func (uc *CloudSvcCloudflare) UpdateDNS(newIPAddress string) error {
	searchName := uc.env.Name
	if uc.env.SubdomainName != "" {
		searchName = uc.env.SubdomainName + "." + uc.env.Name
	}
	zoneList, err := uc.cloudflareRepo.GetZoneList(fmt.Sprintf("name=%s", uc.env.Name))
	if err != nil {
		return fmt.Errorf("failed to get dns zone list. Error: %v", err)
	}
	if len(zoneList) <= 0 {
		return fmt.Errorf("no zone found for domain %s", uc.env.Name)
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
			Proxied: uc.env.IsProxied,
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
