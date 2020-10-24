package exceptions

const TYPE_NOT_FOUND = 1
const TYPE_BAD_REQUEST = 2
const TYPE_PRECONDITION_FAILED = 3
const TYPE_INTERNAL_ERROR = 4

type Exception struct {
	Message string
	Type    int8
	Cause   string
	Error   error
}

type ExceptionInterface interface {
	GetStatusCode() int
}

func (e Exception) GetStatusCode() int {
	switch e.Type {
	case TYPE_NOT_FOUND:
		return 404
	case TYPE_BAD_REQUEST:
		return 400
	case TYPE_PRECONDITION_FAILED:
		return 412
	}
	return 500
}
