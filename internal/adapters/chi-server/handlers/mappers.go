package handlers

func MapErrorIntoErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Message: err.Error(),
		},
	}
}
