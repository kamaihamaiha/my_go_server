package model

import "time"

// CommonLawType 常用法律类型定义
type CommonLawType struct {
	ID             int       `gorm:"column:id;primaryKey;autoIncrement"`
	UUID           string    `gorm:"column:uuid;uniqueIndex;not null"`
	TypeID         int       `gorm:"column:type_id;not null;index"`     // 关联 types 表
	LawType        string    `gorm:"column:law_type;uniqueIndex;not null"`
	LawTypeDisplay string    `gorm:"column:law_type_display;not null"`
	Icon           string    `gorm:"column:icon;not null"`
	Keywords       string    `gorm:"column:keywords;not null"` // JSON 数组字符串
	SortOrder      int       `gorm:"column:sort_order;not null;default:0"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (CommonLawType) TableName() string {
	return "common_law_type"
}

// CommonLaw 法律与常用类型的映射关系
type CommonLaw struct {
	ID              int       `gorm:"column:id;primaryKey;autoIncrement"`
	CommonLawTypeID int       `gorm:"column:common_law_type_id;not null;index;uniqueIndex:idx_type_law"`
	LawID           string    `gorm:"column:law_id;not null;index;uniqueIndex:idx_type_law"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (CommonLaw) TableName() string {
	return "common_laws"
}
