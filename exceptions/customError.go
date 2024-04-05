package exceptions

const (
	// param errors
	CodeParamInvalid = 1001
	CodeParamBlank   = 1002

	//user triggered errors
	CodeEmailAlreadyExists      = 2001
	CodeVerificationCodeInvalid = 2002
	CodeUserNotFound            = 2003

	// service triggered errors
	CodeSendEmailFailed = 3001
	CodeUploadFailed    = 3002
	CodeDatabaseError   = 3003
)

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *CustomError) Error() string {
	return c.Message
}
