package dns

import (
	"encoding/json"
	"fmt"
)

// This is a DNS client which can create DNS records for the Katapult API.
// At present, it supports pretty much nothing though and solely exists
// to allow for ACME challenge requests to be handled.
type Client struct {
	APIToken string
	Host     string
}

// Return an array of all DNS records for a given zone where the filter matches
// the record.
func (c *Client) DNSRecords(zoneName string, filter func(DNSRecord) bool) ([]*DNSRecord, error) {
	body, err := c.apiRequest("GET", "core/v1/dns_zones/_/records", map[string]string{
		"dns_zone[name]": zoneName,
	}, "")
	if err != nil {
		return nil, err
	}

	response := &DNSRecordsResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	records := make([]*DNSRecord, 0)
	for _, record := range response.DNSRecords {
		if filter(record) {
			records = append(records, &record)
		}
	}

	return records, nil
}

// Create a new TXT record in the given zone with the given name and content.
func (c *Client) CreateTXTRecord(zone string, recordName string, content string) (*DNSRecord, error) {
	requestBody := &DNSRecordCreateRequest{
		DNSZone: &DNSZoneLookup{Name: zone},
		Properties: &DNSRecordProperties{
			Name: recordName,
			TTL:  60,
			Type: "TXT",
			Content: &DNSRecordPropertiesContent{
				TXT: &DNSRecordPropertiesContentTXT{
					Content: content,
				},
			},
		},
	}

	requestJson, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request json: %w", err)
	}

	body, err := c.apiRequest("POST", "core/v1/dns_zones/_/records", map[string]string{}, string(requestJson))
	if err != nil {
		return nil, err
	}

	newRecordResponse := &DNSRecordResponse{}
	err = json.Unmarshal(body, newRecordResponse)
	if err != nil {
		return nil, err
	}

	return &newRecordResponse.DNSRecord, nil
}

// Delete the given record
func (c *Client) DeleteRecord(record *DNSRecord) error {
	_, err := c.apiRequest("DELETE", "core/v1/dns_records/_", map[string]string{
		"dns_record[id]": record.ID,
	}, "")
	if err != nil {
		return err
	}

	return nil
}
