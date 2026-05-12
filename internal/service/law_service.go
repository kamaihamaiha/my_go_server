package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"LawHelperServer/internal/model"
	"LawHelperServer/internal/repository"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
	previewLimit    = 20
)

var (
	ErrTypeNotFound = errors.New("type not found")
	ErrLawNotFound  = errors.New("law not found")
)

type LawService struct {
	typeRepo      *repository.TypeRepository
	lawRepo       *repository.LawRepository
	parsedLawRepo *repository.ParsedLawRepository
}

type TypePreview struct {
	TypeID   int                `json:"typeId"`
	TypeName string             `json:"typeName"`
	ParentID *int               `json:"parentId"`
	Total    int64              `json:"total"`
	Items    []model.LawSummary `json:"items"`
}

type PaginatedLawList struct {
	Type       TypeInfo           `json:"type"`
	Page       int                `json:"page"`
	PageSize   int                `json:"pageSize"`
	Total      int64              `json:"total"`
	TotalPages int                `json:"totalPages"`
	Items      []model.LawSummary `json:"items"`
}

type TypeInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID *int   `json:"parentId"`
}

type ParsedLawDetail struct {
	VersionID string           `json:"versionId"`
	Title     string           `json:"title"`
	Available bool             `json:"available"`
	Content   *json.RawMessage `json:"content"`
}

func NewLawService(typeRepo *repository.TypeRepository, lawRepo *repository.LawRepository, parsedLawRepo *repository.ParsedLawRepository) *LawService {
	return &LawService{
		typeRepo:      typeRepo,
		lawRepo:       lawRepo,
		parsedLawRepo: parsedLawRepo,
	}
}

func (s *LawService) ListTypePreviews(ctx context.Context) ([]TypePreview, error) {
	lawTypes, err := s.typeRepo.ListConcreteTypesWithLawCount(ctx)
	if err != nil {
		return nil, err
	}

	previews := make([]TypePreview, 0, len(lawTypes))
	for _, lawType := range lawTypes {
		items, err := s.lawRepo.ListByType(ctx, lawType.ID, 0, previewLimit)
		if err != nil {
			return nil, err
		}

		previews = append(previews, TypePreview{
			TypeID:   lawType.ID,
			TypeName: lawType.Name,
			ParentID: lawType.ParentID,
			Total:    lawType.LawCount,
			Items:    items,
		})
	}

	return previews, nil
}

func (s *LawService) ListLawsByType(ctx context.Context, typeID, page, pageSize int) (*PaginatedLawList, error) {
	lawType, err := s.typeRepo.GetByID(ctx, typeID)
	if err != nil {
		return nil, err
	}
	if lawType == nil {
		return nil, ErrTypeNotFound
	}

	page, pageSize = normalizePagination(page, pageSize)

	total, err := s.lawRepo.CountByType(ctx, typeID)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	items, err := s.lawRepo.ListByType(ctx, typeID, offset, pageSize)
	if err != nil {
		return nil, err
	}

	return &PaginatedLawList{
		Type: TypeInfo{
			ID:       lawType.ID,
			Name:     lawType.Name,
			ParentID: lawType.ParentID,
		},
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages(total, pageSize),
		Items:      items,
	}, nil
}

func (s *LawService) GetParsedLaw(ctx context.Context, versionID string) (*ParsedLawDetail, error) {
	versionID = strings.TrimSpace(versionID)

	lawMeta, err := s.lawRepo.GetMetaByVersionID(ctx, versionID)
	if err != nil {
		return nil, err
	}
	if lawMeta == nil {
		return nil, ErrLawNotFound
	}

	raw, err := s.parsedLawRepo.GetByVersionID(ctx, versionID, lawMeta.LawTypeID)
	if err != nil {
		if errors.Is(err, repository.ErrParsedLawNotFound) {
			return &ParsedLawDetail{
				VersionID: lawMeta.VersionID,
				Title:     lawMeta.Title,
				Available: false,
				Content:   nil,
			}, nil
		}
		return nil, err
	}

	return &ParsedLawDetail{
		VersionID: lawMeta.VersionID,
		Title:     lawMeta.Title,
		Available: true,
		Content:   &raw,
	}, nil
}

func normalizePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = defaultPage
	}

	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return page, pageSize
}

func totalPages(total int64, pageSize int) int {
	if total == 0 {
		return 0
	}

	return int((total + int64(pageSize) - 1) / int64(pageSize))
}
