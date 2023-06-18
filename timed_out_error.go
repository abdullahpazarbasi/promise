package promise

type timedOutError string

func (e timedOutError) Error() string {
	return string(e)
}
