package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"LawHelperServer/internal/model"
)

const lawListOrder = `
CASE WHEN effectDate IS NULL OR TRIM(effectDate) = '' THEN 1 ELSE 0 END ASC,
effectDate DESC,
CASE WHEN publishDate IS NULL OR TRIM(publishDate) = '' THEN 1 ELSE 0 END ASC,
publishDate DESC,
versionId DESC
`

type LawRepository struct {
	db *gorm.DB
}

func NewLawRepository(db *gorm.DB) *LawRepository {
	return &LawRepository{db: db}
}

func (r *LawRepository) ListByType(ctx context.Context, typeID, offset, limit int) ([]model.LawSummary, error) {
	var laws []model.LawSummary

	err := r.db.WithContext(ctx).
		Model(&model.LawList{}).
		Select("versionId, title, lawTypeId, lawType, publishDate, effectDate, effectiveStatus, authorityName").
		Where("lawTypeId = ?", typeID).
		Order(lawListOrder).
		Offset(offset).
		Limit(limit).
		Find(&laws).Error
	if err != nil {
		return nil, err
	}

	return laws, nil
}

func (r *LawRepository) CountByType(ctx context.Context, typeID int) (int64, error) {
	var total int64

	err := r.db.WithContext(ctx).
		Model(&model.LawList{}).
		Where("lawTypeId = ?", typeID).
		Count(&total).Error
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *LawRepository) GetMetaByVersionID(ctx context.Context, versionID string) (*model.LawMeta, error) {
	var law model.LawMeta

	err := r.db.WithContext(ctx).
		Model(&model.LawList{}).
		Select("versionId, title, lawTypeId").
		Where("versionId = ?", versionID).
		Take(&law).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &law, nil
}
