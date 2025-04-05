package base_response

type BaseResponse struct {
	Data    any            `json:"data,omitempty"`
	Error   *error         `json:"error,omitempty"`
	Message string         `json:"message" binding:"required"`
	Page    PaginationData `json:"page,omitzero"`
	Status  Status         `json:"status" binding:"required"`
}

type PaginationData struct {
	TotalItems  int `json:"total_items" binding:"required,gte=0"`
	CurrentPage int `json:"current_page" binding:"required,gte=0"`
	PageSize    int `json:"page_size" binding:"required,gte=0"`
	TotalPages  int `json:"total_pages" binding:"required,gte=0"`
}

type Status string

const (
	StatusOK    = "OK"
	StatusError = "Error"
)
