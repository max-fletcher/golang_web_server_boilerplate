package responses

type Response struct {
	// The part `json:"status"` is called a struct tag
	// It tells fo than when converting this struct to/from JSON, use this name
	// Without this, we will get { "Status": "ok" } instead of { "status": "ok" }
	// This tag has many variations e.g
	// Data interface{} `json:"data,omitempty"` - Field is excluded from JSON
	// Data interface{} `json:"data,omitempty"` - Only include this field if it has a value
	// UserID int `json:"user_id"` - COnvert UserID to user_id when converting to/from JSON(what we are using here)
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}
