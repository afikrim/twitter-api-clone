package pkg_pagination

type CursorPagination struct {
	Current int64 `json:"current"`
	Next    int64 `json:"next"`
}

func NewCursorPagination(total int64, limit int, offset int) *CursorPagination {
	next := int64(limit + offset)
	if next > total {
		next = -1
	}
	return &CursorPagination{
		Current: int64(offset),
		Next:    next,
	}
}
