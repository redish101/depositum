package v1

type PaginationParams struct {
	Page     int `json:"page" query:"page"`         // 页码，从1开始
	PageSize int `json:"pageSize" query:"pageSize"` // 每页大小
}

// PaginationResponse 分页响应
type PaginationResponse[T any] struct {
	Data       []T   `json:"data"`        // 数据列表
	Page       int   `json:"page"`        // 当前页码
	PageSize   int   `json:"pageSize"`    // 每页大小
	Total      int64 `json:"total"`       // 总记录数
	TotalPages int   `json:"totalPages"` // 总页数
	HasNext    bool  `json:"hasNext"`    // 是否有下一页
	HasPrev    bool  `json:"hasPrev"`    // 是否有上一页
}
