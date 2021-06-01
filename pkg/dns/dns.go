package dns

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
)

type DomainManager interface {
	AddCNAMERecord(ctx context.Context, subdomain string) error
}

type Service struct {
	client      *cloudflare.API
	email       string
	cnameTarget string
}

func NewService(client *cloudflare.API, email, cnameTarget string) *Service {
	return &Service{
		client:      client,
		email:       email,
		cnameTarget: cnameTarget,
	}
}

func (s *Service) AddCNAMERecord(ctx context.Context, subdomain string) error {
	id, err := s.client.ZoneIDByName(s.email)
	if err != nil {
		return err
	}

	// todo enable proxy
	proxied := true
	_, err = s.client.CreateDNSRecord(ctx, id, cloudflare.DNSRecord{
		Name:    subdomain,
		Type:    "CNAME",
		Content: s.cnameTarget,
		TTL:     1,
		Proxied: &proxied,
	})

	return err
}
