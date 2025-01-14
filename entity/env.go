package entity

type Env struct {
	VPS     map[string]*EnvCloudServer `json:"vps"`
	Bastion map[string]*EnvBastion     `json:"bastion"`
}

type EnvResourceToken struct {
	DigitaloceanAPI string `json:"digitalocean_api,omitempty"`
	VultrAPI        string `json:"vultr_api,omitempty"`
}

type EnvBastion struct {
	SShCertFile         string                                                       `json:"ssh_cert_file"`
	DeviceURL           map[CloudServiceEnvName]string                               `json:"device_url"`
	ResourceURL         map[CloudServiceTechStackName]map[CloudServiceEnvName]string `json:"resource_url"`
	ResourcePort        map[CloudServiceTechStackName]int                            `json:"resource_port"`
	ResourceExposedPort map[CloudServiceTechStackName]map[CloudServiceEnvName]int    `json:"resource_exposed_port"`
}

type EnvCloudServer struct {
	Domains                 []*EnvCloudServerDomain   `json:"domains"`
	SSHKey                  string                    `json:"ssh_key"`
	VmType                  EnvCloudServerVmBrand     `json:"vm_type"`
	VmTypeApiToken          string                    `json:"vm_type_api_token"`
	VmLabel                 string                    `json:"vm_label"`
	VmChooseSnapshotOverISO bool                      `json:"vm_choose_snapshot_over_iso,omitempty"`
	VmISO                   string                    `json:"vm_iso"`
	VmFirewall              string                    `json:"vm_firewall,omitempty"`
	VmBlockStoragesLabel    []string                  `json:"vm_block_storages_label,omitempty"`
	VmBackupPlan            *EnvCloudServerBackupPlan `json:"vm_backup_plan,omitempty"`
	VmResourcePlan          string                    `json:"vm_resource_plan"`
	VmRegion                string                    `json:"vm_region"`
}

type EnvCloudServerDomain struct {
	APIToken      string `json:"api_token"`
	Name          string `json:"name"`
	SubdomainName string `json:"subdomain_name,omitempty"`
	IsProxied     bool   `json:"proxied,omitempty"`
}

type EnvCloudServerVmBrand string

var (
	EnvCloudServerVmBrandDigitalOcean EnvCloudServerVmBrand = "digitalocean"
	EnvCloudServerVmBrandVultr        EnvCloudServerVmBrand = "vultr"
)

type EnvCloudServerBackupPlan struct {
	Plan    string `json:"plan"`
	Weekday string `json:"weekday"`
	Hour    int    `json:"hour"`
}
