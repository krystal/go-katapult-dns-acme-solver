package solver

import (
	"log"

	"github.com/krystal/go-katapult-acme-dns/dns"
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

func (s *Solver) Delete(zoneName string, recordName string) error {
	zone, err := s.dns.DNSZone(zoneName)
	if err != nil {
		s.logger.Printf("failed to find DNS zone with name %s (%s)", zoneName, err)
		return err
	}

	s.logger.Printf("got zone %s with id %s\n", zone.Name, zone.ID)

	err = s.deleteAllRecordsByName(zone, recordName)
	if err != nil {
		s.logger.Printf("failed to delete DNS all records with name %s (%s)", recordName, err)
		return err
	}

	return nil
}

func (s *Solver) deleteAllRecordsByName(zone *dns.DNSZone, recordName string) error {
	records, err := s.dns.DNSRecordsByName(zone, recordName, "TXT")
	if err != nil {
		return err
	}

	if len(records) == 0 {
		s.logger.Printf("no records found matching record %s in zone %s (%s)\n", recordName, zone.Name, zone.ID)
	} else {
		for _, r := range records {
			s.logger.Printf("deleting record %s (%s)\n", r.FullName, r.ID)
			err = s.dns.DeleteRecord(r)
			if err != nil {
				s.logger.Printf("failed to delete record %s (%s)\n", r.ID, err)
			}
		}
	}

	return nil
}
