package promise

type proprietyError string

func (e proprietyError) Error() string {
	return string(e)
}
