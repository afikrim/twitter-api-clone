package pkg_pagination

type CursorPagination struct {
	Current int64 `json:"current"`
	Next    int64 `json:"next"`
}
