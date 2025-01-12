package entity

type VultrISO struct {
	Id       string `json:"id"`
	Filename string `json:"filename"`
}

type VultrFirewallGroup struct {
	Id          string `json:"id"`
	Description string `json:"description"`
}

type VultrSSHKey struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type VultrReservedIp struct {
	Id         string `json:"id"`
	Subnet     string `json:"subnet"`
	Label      string `json:"label"`
	InstanceId string `json:"instance_id"`
}

type VultrBlockStorage struct {
	Id      string `json:"id"`
	Status  string `json:"status"`
	Cost    int    `json:"cost"`
	Label   string `json:"label"`
	MountId string `json:"mount_id"`
}

type VultrSnapshot struct {
	ID          string `json:"id"`
	DateCreated string `json:"date_created"`
	Description string `json:"description"`
	Size        int    `json:"size"`
	Status      string `json:"status"`
	OsID        int    `json:"os_id"`
	AppID       int    `json:"app_id"`
}

type VultrInstance struct {
	ID               string `json:"id"`
	OS               string `json:"os"`
	RAM              int    `json:"ram"`
	Disk             int    `json:"disk"`
	MainIP           string `json:"main_ip"`
	VCpuCount        int    `json:"vcpu_count"`
	Region           string `json:"region"`
	DefaultPassword  string `json:"default_password"`
	DateCreated      string `json:"date_created"`
	Status           string `json:"status"`
	PowerStatus      string `json:"power_status"`
	ServerStatus     string `json:"server_status"`
	AllowedBandwidth int    `json:"allowed_bandwidth"`
	NetmaskV4        string `json:"netmask_v4"`
	GatewayV4        string `json:"gateway_v4"`
	Hostname         string `json:"hostname"`
	Label            string `json:"label"`
	InternalIP       string `json:"internal_ip"`
	ImageID          string `json:"image_id"`
	Plan             string `json:"plan"`
}

type VultrKubeCluster struct {
	Id            string             `json:"id"`
	Label         string             `json:"label"`
	ClusterSubnet string             `json:"cluster_subnet"`
	ServiceSubnet string             `json:"service_subnet"`
	Ip            string             `json:"ip"`
	Endpoint      string             `json:"endpoint"`
	Region        string             `json:"region"`
	Status        string             `json:"status"`
	NodePools     VultrKubeNodePools `json:"node_pools"`
	DateCreated   string             `json:"date_created"`
}

type VultrKubeNodePools struct {
	Id           string `json:"id"`
	Label        string `json:"label"`
	Tag          string `json:"tag"`
	Plan         string `json:"plan"`
	Status       string `json:"status"`
	NodeQuantity int    `json:"node_quantity"`
	MinNodes     int    `json:"min_nodes"`
	MaxNodes     int    `json:"max_nodes"`
	AutoScaler   bool   `json:"auto_scaler"`
	DateCreated  string `json:"date_created"`
}

type VultrKubeNode struct {
	Id          string `json:"id"`
	Label       string `json:"label"`
	Status      string `json:"status"`
	DateCreated string `json:"date_created"`
}

type VultrAPIMeta struct {
	Links VultrAPIMetaLinks `json:"links"`
	Total int               `json:"total"`
}
type VultrAPIMetaLinks struct {
	Next string `json:"next"`
	Prev string `json:"prev"`
}

type VultrAPIRequestCreateNodePoolConfig struct {
	NodeQuantity int    `json:"node_quantity"`
	Label        string `json:"label"`
	Plan         string `json:"plan"`
	Tag          string `json:"tag"`
	AutoScaler   string `json:"auto_scaler"`
	MinNodes     string `json:"min_nodes"`
	MaxNodes     string `json:"max_nodes"`
}

type VultrAPIRequestCreateInstanceConfig struct {
	Region          string   `json:"region"`
	Plan            string   `json:"plan"`
	FirewallGroupID string   `json:"firewall_group_id,omitempty"`
	IsoID           string   `json:"iso_id,omitempty"`
	SnapshotID      string   `json:"snapshot_id,omitempty"`
	EnableIPV6      bool     `json:"enable_ipv6,omitempty"`
	SshKeyIDs       []string `json:"sshkey_id,omitempty"`
	Backups         string   `json:"backups,omitempty"`
	Label           string   `json:"label,omitempty"`
	ReservedIPV4    string   `json:"reserved_ipv4,omitempty"`
}
