package entity

type DigitaloceanSnapshot struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	CreatedAt     string   `json:"created_at"`
	Regions       []string `json:"regions"`
	ResourceID    string   `json:"resource_id"`
	ResourceType  string   `json:"resource_type"`
	MinDiskSize   int      `json:"min_disk_size"`
	SizeGigabytes float64  `json:"size_gigabytes"`
	Tags          []string `json:"tags"`
}

type DigitaloceanRegion struct {
	Name      string   `json:"name"`
	Slug      string   `json:"slug"`
	Features  []string `json:"features"`
	Available bool     `json:"available"`
	Sizes     []string `json:"sizes"`
}

type DigitaloceanDropletSize struct {
	Slug         string                 `json:"slug"`
	Memory       int                    `json:"memory"`
	VCPUs        int                    `json:"vcpus"`
	Disk         int                    `json:"disk"`
	Transfer     float64                `json:"transfer"`
	PriceMonthly float64                `json:"price_monthly"`
	PriceHourly  float64                `json:"price_hourly"`
	Regions      []string               `json:"regions"`
	Available    bool                   `json:"available"`
	Description  string                 `json:"description"`
	DiskInfo     []DigitaloceanDiskInfo `json:"disk_info"`
}

type DigitaloceanDiskInfo struct {
	Type string                   `json:"type"`
	Size DigitaloceanDiskInfoSize `json:"size"`
}

type DigitaloceanDiskInfoSize struct {
	Amount int `json:"amount"`
	Unit   int `json:"unit"`
}

type DigitaloceanImage struct {
	ID            int      `json:"id"`
	Name          string   `json:"name,omitempty"`
	Distribution  string   `json:"distribution"`
	Slug          string   `json:"slug"`
	Public        bool     `json:"public"`
	Regions       []string `json:"regions"`
	CreatedAt     string   `json:"created_at"`
	Type          string   `json:"type"`
	MinDiskSize   int      `json:"min_disk_size"`
	SizeGigabytes float64  `json:"size_gigabytes"`
	Description   string   `json:"description,omitempty"`
	Tags          []string `json:"tags"`
	Status        string   `json:"status"`
	ErrorMessage  string   `json:"error_message,omitempty"`
}

type DigitaloceanDroplet struct {
	ID               int                         `json:"id"`
	Name             string                      `json:"name"`
	Memory           int                         `json:"memory"`
	VCPUs            int                         `json:"vcpus"`
	Disk             int                         `json:"disk"`
	DiskInfo         DigitaloceanDiskInfo        `json:"disk_info"`
	Locked           bool                        `json:"locked"`
	Status           string                      `json:"status"`
	Kernel           *interface{}                `json:"kernel,omitempty"`
	CreatedAt        string                      `json:"created_at"`
	Features         []string                    `json:"features"`
	BackupIDs        []int                       `json:"backup_ids"`
	NextBackupWindow *interface{}                `json:"next_backup_window,omitempty"`
	SnapshotIDs      []int                       `json:"snapshot_ids"`
	Image            DigitaloceanImage           `json:"image"`
	VolumeIDs        []int                       `json:"volume_ids"`
	Size             DigitaloceanDropletSize     `json:"size"`
	SizeSlug         string                      `json:"size_slug"`
	Networks         DigitaloceanDropletNetworks `json:"networks"`
	Region           DigitaloceanRegion          `json:"region"`
	Tags             []string                    `json:"tags"`
}

type DigitaloceanDropletNetworks struct {
	V4 []*DigitaloceanDropletNetwork `json:"v4"`
	V6 []*DigitaloceanDropletNetwork `json:"v6"`
}

type DigitaloceanDropletNetwork struct {
	IPAddress string `json:"ip_address,omitempty"`
	Netmask   int    `json:"netmask,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Type      string `json:"type,omitempty"`
}

type DigitaloceanDropletActionable struct {
	ID   int    `json:"id"`
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type DigitaloceanDropletAction struct {
	ID           int                `json:"id"`
	Status       string             `json:"status"`
	Type         string             `json:"type"`
	StartedAt    string             `json:"started_at"`
	CompletedAt  string             `json:"completed_at"`
	ResourceID   int                `json:"resource_id"`
	ResourceType string             `json:"resource_type"`
	Region       DigitaloceanRegion `json:"region"`
	RegionSlug   string             `json:"region_slug"`
}

type DigitaloceanAPIResultCursor struct {
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Pref  string `json:"pref,omitempty"`
	Next  string `json:"next,omitempty"`
}

type DigitaloceanAPIResultMeta struct {
	Total int `json:"total"`
}
