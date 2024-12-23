package entity

type Env struct {
	VPS     map[string]*EnvCloudServer `json:"vps"`
	Bastion map[string]*EnvBastion     `json:"bastion"`
}

type EnvBastion struct {
	SShCertFile         string                                                       `json:"ssh_cert_file"`
	DeviceURL           map[CloudServiceEnvName]string                               `json:"device_url"`
	ResourceURL         map[CloudServiceTechStackName]map[CloudServiceEnvName]string `json:"resource_url"`
	ResourcePort        map[CloudServiceTechStackName]int                            `json:"resource_port"`
	ResourceExposedPort map[CloudServiceTechStackName]map[CloudServiceEnvName]int    `json:"resource_exposed_port"`
}

type EnvCloudServer struct {
	TokenCloudflare      string   `json:"cloudflare_token,omitempty"`
	TokenVultr           string   `json:"vultr_token"`
	DomainName           string   `json:"domain_name"`
	SubdomainName        string   `json:"subdomain_name,omitempty"`
	SSHKeyLabel          string   `json:"ssh_key_label"`
	VmLabel              string   `json:"vm_label"`
	VmISoFilename        string   `json:"vm_iso_filename"`
	VmFirewallLabel      string   `json:"vm_firewall_label,omitempty"`
	VmBlockStoragesLabel []string `json:"vm_block_storages_label,omitempty"`
	VmVultrPlan          string   `json:"vm_vultr_plan"`
	VmRegion             string   `json:"vm_region"`
}
