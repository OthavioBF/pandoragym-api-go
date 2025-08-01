package services

import (
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
)

type MailService struct {
	queries *pgstore.Queries
}

func NewMailService(queries *pgstore.Queries) *MailService {
	return &MailService{
		queries: queries,
	}
}
