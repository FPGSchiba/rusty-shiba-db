package util

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type Pagination struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
