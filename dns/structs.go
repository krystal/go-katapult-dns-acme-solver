package dns

// DNS Zones

type DNSZoneResponse struct {
	DNSZone DNSZone `json:"dns_zone"`
}

type DNSZone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DNS Records

type DNSRecordsResponse struct {
	DNSRecords []DNSRecord `json:"dns_records"`
}

type DNSRecordResponse struct {
	DNSRecord DNSRecord `json:"dns_record"`
}

type DNSRecord struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	TTL      int    `json:"ttl"`
	Type     string `json:"type"`
	Content  string `json:"content"`
}

// DNS Record Creation

type DNSRecordCreateRequest struct {
	DNSZone    *DNSZoneLookup       `json:"dns_zone"`
	Properties *DNSRecordProperties `json:"properties"`
}

type DNSZoneLookup struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type DNSRecordProperties struct {
	Name    string                      `json:"name"`
	Type    string                      `json:"type"`
	TTL     int                         `json:"ttl"`
	Content *DNSRecordPropertiesContent `json:"content"`
}

type DNSRecordPropertiesContent struct {
	TXT *DNSRecordPropertiesContentTXT `json:"TXT"`
}

type DNSRecordPropertiesContentTXT struct {
	Content string `json:"content"`
}
