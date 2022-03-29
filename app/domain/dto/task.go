package dto

type GetTasksDTO struct {
	Limit   uint32 `json:"limit"`
	Page    uint32 `json:"page"`
	Keyword string `json:"keyword"`
	Sort    string `json:"sort"`
}
