package joint

// Pager 定义了分页信息结构体。
type Pager struct {
	Page      int   `json:"page"`       // 当前页码
	PageSize  int   `json:"page_size"`  // 每页数据量
	TotalRows int64 `json:"total_rows"` // 总数据量
}

// PageResp 定义了带分页信息的页面响应结构体。
type PageResp struct {
	List  interface{} `json:"list"`  // 列表数据
	Pager Pager       `json:"pager"` // 分页信息
}

// PageRespFrom 根据传入的列表数据、当前页码、每页数据量和总数据量创建一个 PageResp 实例。
func PageRespFrom(list interface{}, page int, pageSize int, totalRows int64) *PageResp {
	return &PageResp{
		List: list,
		Pager: Pager{
			Page:      page,
			PageSize:  pageSize,
			TotalRows: totalRows,
		},
	}
}
