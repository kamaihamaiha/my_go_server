package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrParsedLawNotFound = errors.New("parsed law not found")

type ParsedLawRepository struct {
	baseDir string
}

func NewParsedLawRepository(baseDir string) *ParsedLawRepository {
	return &ParsedLawRepository{baseDir: baseDir}
}

func (r *ParsedLawRepository) GetByVersionID(_ context.Context, versionID string) (json.RawMessage, error) {
	filePath := filepath.Join(r.baseDir, versionID+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrParsedLawNotFound
		}
		return nil, err
	}

	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil, ErrParsedLawNotFound
	}

	if !json.Valid(data) {
		return nil, fmt.Errorf("parsed law file %s is not valid json", filePath)
	}

	return json.RawMessage(data), nil
}
