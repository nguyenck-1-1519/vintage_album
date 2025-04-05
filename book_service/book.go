package book_service

type Book struct {
	ID     int     `json:"id"`
	Title  string  `json:"title" binding:"required"`      // Title là bắt buộc
	Author string  `json:"author" binding:"required"`     // Author là bắt buộc
	Price  float64 `json:"price" binding:"required,gt=0"` // Price là bắt buộc và phải lớn hơn 0
	Stock  int     `json:"stock" binding:"gte=0"`         // Stock phải lớn hơn hoặc bằng 0
}
