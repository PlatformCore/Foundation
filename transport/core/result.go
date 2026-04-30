package core

type Result struct {
	OK       bool
	Code     string
	Message  string
	Data     any
	Error    *Error
	Metadata Metadata
}

func Success(data any) Result { return Result{OK: true, Code: "OK", Data: data, Metadata: Metadata{}} }
func Failure(err error) Result {
	e := ToError(err)
	return Result{OK: false, Code: e.Code, Message: e.Message, Error: e, Metadata: Metadata{}}
}
