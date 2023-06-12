package promise

type unexpectedCaseError string

func (e unexpectedCaseError) Error() string {
	return string(e)
}
