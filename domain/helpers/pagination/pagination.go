package pagination

type Response struct {
	TotalData int64  `json:"total_data"`
	TotalPage uint16 `json:"total_page"`
	Page      uint16 `json:"page"`
	Limit     uint16 `json:"limit"`
}

func NewPaginationResponse(totalData int64, page, limit uint64) Response {
	if limit == 0 {
		var totalPage uint16
		if totalData > 0 {
			totalPage = 1
		}

		return Response{
			TotalData: totalData,
			TotalPage: totalPage,
			Page:      uint16(page),
		}
	}

	totalPage := int(totalData) / int(limit)
	if int(totalData)%int(limit) != 0 {
		totalPage++
	}

	return Response{
		TotalData: totalData,
		TotalPage: uint16(totalPage),
		Page:      uint16(page),
		Limit:     uint16(limit),
	}
}

func PaginateInMemory[T any](items []T, offset, limit uint64) []T {
	if offset >= uint64(len(items)) {
		return []T{}
	}
	end := offset + limit
	if end > uint64(len(items)) {
		end = uint64(len(items))
	}
	return items[offset:end]
}
