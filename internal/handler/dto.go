package handler

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Result string `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}
