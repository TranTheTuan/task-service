package models

import (
	pbTypes "github.com/TranTheTuan/pbtypes/build/go/core"
)

type Pagination struct {
	Limit      uint32 `json:"limit"`
	Page       uint32 `json:"page"`
	Total      uint32 `json:"total"`
	TotalPages uint32 `json:"total_pages"`
}

func (p *Pagination) ToProto() *pbTypes.Pagination {
	return &pbTypes.Pagination{
		Limit:      p.Limit,
		Page:       p.Page,
		Total:      p.Total,
		TotalPages: p.TotalPages,
	}
}
