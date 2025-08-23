package types

type EID int

type BaseEntity struct {
	Id EID
}

type GetQueryParams struct {
	PageSize   *int    `query:"page-size"`
	PageNumber *int    `query:"page"`
	Embed      *string `query:"embed"`
}

type PayloadValidator interface {
	Validate() error
}

type GetAllDTO[T any] struct {
	Data          T   `json:"data"`
	PageNumber    int `json:"pageNumber"`
	PageSize      int `json:"pageSize"`
	TotalRowCount int `json:"totalRowCount"`
}
