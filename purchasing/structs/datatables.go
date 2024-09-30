package structs

// Corrected PaginatedResponse struct
type PaginatedResponse[T any] struct {
	TotalRecords    int64  `json:"recordsTotal"`
	FilteredRecords int64  `json:"recordsFiltered"`
	Data            []T    `json:"data"` // This should be []T, not [][]T
	Code            int    `json:"code"`
	Message         string `json:"message"`
}
