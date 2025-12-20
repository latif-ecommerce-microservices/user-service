package httputil

import (
	"context"
	"net/url"
	"strconv"
)

func ParseQueryToFilter(ctx context.Context, query url.Values) BasePaginationFilter {
	filter := BasePaginationFilter{
		Page: 1,
		Size: 10,
	}

	if page, err := strconv.Atoi(query.Get("page")); err == nil && page > 0 {
		filter.Page = int64(page)
	} else {
		filter.Page = 1
	}

	if size, err := strconv.Atoi(query.Get("size")); err == nil && size > 0 {
		filter.Size = int64(size)
	} else {
		filter.Size = 10
	}

	filter.Search = query.Get("search")
	filter.Sort = query.Get("sort")

	return filter
}
