package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

var ErrParsedLawNotFound = errors.New("parsed law not found")

type ParsedLawRepository struct {
	baseDir string
}

func NewParsedLawRepository(baseDir string) *ParsedLawRepository {
	return &ParsedLawRepository{baseDir: baseDir}
}

func (r *ParsedLawRepository) GetByVersionID(_ context.Context, versionID string, lawTypeID int) (json.RawMessage, error) {
	for _, filePath := range r.candidatePaths(versionID, lawTypeID) {
		data, err := os.ReadFile(filePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, err
		}

		data = bytes.TrimSpace(data)
		if len(data) == 0 {
			continue
		}

		if !json.Valid(data) {
			return nil, fmt.Errorf("parsed law file %s is not valid json", filePath)
		}

		return json.RawMessage(data), nil
	}

	return nil, ErrParsedLawNotFound
}

func (r *ParsedLawRepository) candidatePaths(versionID string, lawTypeID int) []string {
	paths := make([]string, 0, 2)
	if lawTypeID > 0 {
		typeDir := "type_" + strconv.Itoa(lawTypeID)
		paths = append(paths, filepath.Join(r.baseDir, "laws_by_type", typeDir, versionID+".json"))
	}

	paths = append(paths, filepath.Join(r.baseDir, versionID+".json"))
	return paths
}
