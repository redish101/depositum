package v1

type PaginationParams struct {
	Page     int `json:"page" query:"page"`           // 页码，从1开始
	PageSize int `json:"page_size" query:"page_size"` // 每页大小
}

// PaginationResponse 分页响应
type PaginationResponse[T any] struct {
	Data       []T   `json:"data"`        // 数据列表
	Page       int   `json:"page"`        // 当前页码
	PageSize   int   `json:"page_size"`   // 每页大小
	Total      int64 `json:"total"`       // 总记录数
	TotalPages int   `json:"total_pages"` // 总页数
	HasNext    bool  `json:"has_next"`    // 是否有下一页
	HasPrev    bool  `json:"has_prev"`    // 是否有上一页
}
