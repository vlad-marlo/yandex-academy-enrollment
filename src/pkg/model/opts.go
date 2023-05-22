package model

type PaginationOpts interface {
	Limit() int
	Offset() int
}
