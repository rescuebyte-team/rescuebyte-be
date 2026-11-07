package pagination_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rescuebyte/rescuebyte-be/domain/helpers/pagination"
)

func TestNewPaginationResponse_ZeroLimitMeansUnpaged(t *testing.T) {
	assert.Equal(t, pagination.Response{TotalData: 25, TotalPage: 1, Page: 1}, pagination.NewPaginationResponse(25, 1, 0))
	assert.Equal(t, pagination.Response{TotalData: 0, TotalPage: 0, Page: 1}, pagination.NewPaginationResponse(0, 1, 0))
}

func TestNewPaginationResponse_Paged(t *testing.T) {
	assert.Equal(t, pagination.Response{TotalData: 25, TotalPage: 3, Page: 2, Limit: 10}, pagination.NewPaginationResponse(25, 2, 10))
}

func TestNormalize(t *testing.T) {
	var limit, page, offset uint64
	pagination.Normalize(&limit, &page, &offset, "product")

	assert.Equal(t, uint64(10), limit)
	assert.Equal(t, uint64(1), page)
	assert.Equal(t, uint64(0), offset)

	limit = 20
	page = 3
	pagination.Normalize(&limit, &page, &offset, "product")
	assert.Equal(t, uint64(20), limit)
	assert.Equal(t, uint64(3), page)
	assert.Equal(t, uint64(40), offset)
}
