package utils

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Status  string `json:"status"` 
	Code    int    `json:"code"`   
	Message string `json:"message"`
}

func FormatResponse(status string, code int, message string, data interface{}) APIResponse {
	meta := Meta{
		Status:  status,
		Code:    code,
		Message: message,
	}
	return APIResponse{
		Meta: meta,
		Data: data,
	}
}

func FormatErrorResponse(status string, code int, message string) APIResponse {
	meta := Meta{
		Status:  status,
		Code:    code,
		Message: message,
	}
	return APIResponse{
		Meta: meta,
		Data: nil, // Data empty when error
	}
}

func SendSuccessResponse(c *gin.Context, message string, data interface{}) {
    response := FormatResponse("success", 200, message, data)
    c.JSON(200, response)
}

func SendCreatedResponse(c *gin.Context, message string, data interface{}) {
    response := FormatResponse("success", 201, message, data)
    c.JSON(201, response)
}

// HTTP error (4xx, 5xx)
func SendErrorResponse(c *gin.Context, code int, message string) {
    response := FormatErrorResponse("error", code, message)
    c.JSON(code, response)
}