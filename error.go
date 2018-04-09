package log

type CriticalError struct {
	s      string
	Bytes  []int
	Errors []error
}

func (e CriticalError) Error() string {
	return e.s
}
