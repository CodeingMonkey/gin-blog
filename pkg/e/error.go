package e

type CustomizeError struct {
	Msg string
	ErrorCode int
}

func (c CustomizeError) Error() string {
	return c.Msg
}

func (c CustomizeError) Code() int  {
	return c.ErrorCode
}
