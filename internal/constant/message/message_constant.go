package message

const (
	Success = "Successfuly"

	InternalUserAuthNotFound = "user authentication data not found in context"

	ClientInvalidEmailOrPassword = "Make sure you have provide valid email or password"
	ClientUserAlreadyExist       = "Email already been used, please use another email"
	ClientUnauthenticated        = "Unauthenticated, please try login again"
	ClientPermissionDenied       = "Permission denied for accessing this resource"

	InternalGracefulError = "Something wrong happened. Please try again"
)
