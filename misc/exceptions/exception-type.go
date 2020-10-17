package exceptions

const TYPE_NOT_FOUND = 1
const TYPE_BAD_REQUEST = 2
const TYPE_PRECONDITION_FAILED = 2
const TYPE_INTERNAL_ERROR = 2

type Exception struct {
	Message string
	Type    int8
	Cause   string
	Error   error
}
