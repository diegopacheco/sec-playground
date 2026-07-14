package models

type Book struct {
	ID        int64  `json:"id"`
	Title     string `json:"title" binding:"required"`
	Author    string `json:"author" binding:"required"`
	ISBN      string `json:"isbn" binding:"required"`
	Year      int    `json:"year" binding:"required"`
	CreatedAt string `json:"created_at"`
}

type Page struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
	Data       any   `json:"data"`
}
