package promise

type alreadyDoneError string

func (e alreadyDoneError) Error() string {
	return string(e)
}
