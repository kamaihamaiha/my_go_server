package model

type Type struct {
	ID       int    `gorm:"column:id;primaryKey"`
	Name     string `gorm:"column:name"`
	ParentID *int   `gorm:"column:parent_id"`
}

func (Type) TableName() string {
	return "types"
}

type TypeWithCount struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID *int   `json:"parentId"`
	LawCount int64  `json:"lawCount"`
}
