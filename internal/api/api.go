package api

// const hardBodyLimit = 1024 * 1024 * 64

// Response is a response for api
type Response struct {
	Data interface{} `json:"data"`
}

// Responses is a response for api
type Responses struct {
	Data []interface{} `json:"data"`
}
