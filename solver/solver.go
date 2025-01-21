package solver

import (
	"log"

	"github.com/krystal/go-katapult-dns-acme-solver/dns"
)

type Solver struct {
	dns    dns.Client
	logger *log.Logger
}

func NewSolver(apiToken string, logger *log.Logger) *Solver {
	dnsClient := dns.Client{APIToken: apiToken}
	return &Solver{
		dns:    dnsClient,
		logger: logger,
	}
}

func NewSolverWithHost(host string, apiToken string, logger *log.Logger) *Solver {
	dnsClient := dns.Client{Host: host, APIToken: apiToken}
	return &Solver{
		dns:    dnsClient,
		logger: logger,
	}
}

func (s *Solver) Set(zoneName string, recordName string, key string) error {
	zone, err := s.dns.DNSZone(zoneName)
	if err != nil {
		s.logger.Printf("failed to find DNS zone with name %s (%s)", zoneName, err)
		return err
	}

	s.logger.Printf("got zone %s with id %s\n", zone.Name, zone.ID)

	newRecord, err := s.dns.CreateTXTRecord(zone, recordName, key)
	if err != nil {
		s.logger.Printf("failed to create TXT record %s (%s)", recordName, err)
		return err
	}

	s.logger.Printf("created new record %s (%s)", newRecord.FullName, newRecord.ID)

	return nil
}

func (s *Solver) CleanupAll(zoneName string, recordName string) error {
	zone, err := s.dns.DNSZone(zoneName)
	if err != nil {
		s.logger.Printf("failed to find DNS zone with name %s (%s)", zoneName, err)
		return err
	}

	s.logger.Printf("got zone %s with id %s\n", zone.Name, zone.ID)

	records, err := s.dns.DNSRecords(zone, func(record dns.DNSRecord) bool {
		return record.FullName == recordName && record.Type == "TXT"
	})
	if err != nil {
		return err
	}

	if len(records) == 0 {
		s.logger.Printf("no txt records found matching record %s in zone %s (%s)\n", recordName, zone.Name, zone.ID)
		return nil
	} else {
		return s.deleteRecords(records)
	}
}

func (s *Solver) Cleanup(zoneName string, recordName string, value string) error {
	zone, err := s.dns.DNSZone(zoneName)
	if err != nil {
		s.logger.Printf("failed to find DNS zone with name %s (%s)", zoneName, err)
		return err
	}

	s.logger.Printf("got zone %s with id %s\n", zone.Name, zone.ID)

	records, err := s.dns.DNSRecords(zone, func(record dns.DNSRecord) bool {
		return record.FullName == recordName && record.Type == "TXT" && record.Content == value
	})
	if err != nil {
		return err
	}

	if len(records) == 0 {
		s.logger.Printf("no txt records found matching record %s in zone %s (%s) with value %s\n", recordName, zone.Name, zone.ID, value)
		return nil
	} else {
		return s.deleteRecords(records)
	}
}

func (s *Solver) deleteRecords(records []*dns.DNSRecord) error {
	for _, r := range records {
		s.logger.Printf("deleting record %s (%s)\n", r.FullName, r.ID)
		err := s.dns.DeleteRecord(r)
		if err != nil {
			s.logger.Printf("failed to delete record %s (%s)\n", r.ID, err)
		}
	}

	return nil
}
