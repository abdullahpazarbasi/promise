package promise

type event int

// Promise events
const (
	EventResolved event = iota
	EventRejected
	EventCanceled
	EventTimedOut
	EventEliminated
)
