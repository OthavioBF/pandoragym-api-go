package services

import (
	"context"
	"mime/multipart"

	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
)

type FileService struct {
	queries *pgstore.Queries
}

func NewFileService(queries *pgstore.Queries) *FileService {
	return &FileService{
		queries: queries,
	}
}

func (s *FileService) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	return "", nil
}

func (s *FileService) DeleteFile(ctx context.Context, fileId string) error {
	return nil
}
