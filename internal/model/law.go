package model

type LawList struct {
	VersionID       string `gorm:"column:versionId;primaryKey"`
	Title           string `gorm:"column:title"`
	LawTypeID       int    `gorm:"column:lawTypeId"`
	LawType         string `gorm:"column:lawType"`
	PublishDate     string `gorm:"column:publishDate"`
	EffectDate      string `gorm:"column:effectDate"`
	DetailJSON      string `gorm:"column:detailJson"`
	EffectiveStatus int    `gorm:"column:effectiveStatus"`
	AuthorityID     int    `gorm:"column:authorityId"`
	AuthorityName   string `gorm:"column:authorityName"`
	ParseState      int    `gorm:"column:parse_state"`
}

func (LawList) TableName() string {
	return "laws_list"
}

type LawSummary struct {
	VersionID       string `json:"versionId"`
	Title           string `json:"title"`
	LawTypeID       int    `json:"lawTypeId"`
	LawType         string `json:"lawType"`
	PublishDate     string `json:"publishDate"`
	EffectDate      string `json:"effectDate"`
	EffectiveStatus int    `json:"effectiveStatus"`
	AuthorityName   string `json:"authorityName"`
}

type LawMeta struct {
	VersionID string
	Title     string
}
