package dto

type PaginatedResponse struct {
	Data       any        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	CurrentPage int32 `json:"current_page"`
	PageSize    int32 `json:"page_size"`
	HasMore     bool  `json:"has_more"`
}
