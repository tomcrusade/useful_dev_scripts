package entity

type CloudflareDNSRecord struct {
	ID      string `json:"id,omitempty"`
	Content string `json:"content,omitempty"`
	Name    string `json:"name,omitempty"`
}

type CloudflareZone struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
