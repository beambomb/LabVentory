package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"labventroy/backend/internal/apperror"
)

// JSON writes a JSON-encoded API Gateway v2 response with the given HTTP status code.
func JSON(status int, payload any) events.APIGatewayV2HTTPResponse {
	body, err := json.Marshal(payload)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"code":"INTERNAL_ERROR","error":"failed to encode response"}`,
		}
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}
}

// Error maps an error (AppError or generic error) into a standardized JSON error response.
func Error(err error) events.APIGatewayV2HTTPResponse {
	appErr := apperror.FromError(err)
	return JSON(appErr.HTTPStatus, appErr)
}
