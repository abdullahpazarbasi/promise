package promise

type event int

const (
	EventResolved event = iota
	EventRejected
	EventCanceled
	EventTimedOut
)
