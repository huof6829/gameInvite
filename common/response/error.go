package response

import "fmt"

type CodeError struct {
	code    int
	message string
}

var (
	ResourceExisted    = NewCodeError(1062, "resource is already existed.") /// unique key
	WalletNotConnected = NewCodeError(10001, "wallet is not connected.")
	WalletIncorrect    = NewCodeError(10002, "wallet is not correct.")

	InviteCodeIncorrect = NewCodeError(10003, "invite code is not correct.")
	InvitedAlready      = NewCodeError(10004, "you are already invited.")
	InvitedNever        = NewCodeError(10005, "you are never invited.")

	NotERC721  = NewCodeError(10101, "token is not ERC721.")
	NotERC1155 = NewCodeError(10102, "token is not ERC1155.")

	ParameterErr = NewCodeError(20001, "request parameter validate failed.")
	JWTErr       = NewCodeError(20002, "token generate failed.")
	DBErr        = NewCodeError(30001, "Database is busy now, please try again later.")
	RedisErr     = NewCodeError(30002, "Cache is busy now, please try again later.")
	ServerErr    = NewCodeError(50001, "Server is busy now, please try again later.")
)

func (e *CodeError) GetErrCode() int {
	return e.code
}

func (e *CodeError) GetErrMsg() string {
	return e.message
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.code, e.message)
}

func NewCodeError(code int, message string) *CodeError {
	return &CodeError{code: code, message: message}
}
