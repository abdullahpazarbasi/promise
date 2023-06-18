package promise

type canceledError string

func (e canceledError) Error() string {
	return string(e)
}
