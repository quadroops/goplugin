package supervisor

// Payload used as main data when some plugin from some host indicated as error / cannot be reached
type Payload struct {
	Host   string
	Plugin string
}

// OnErrorHandler used as main type for handling plugin's error
type OnErrorHandler func(payload *Payload)

// Driver used as main interface to run supervisor activities
type Driver interface {
	Watch() <-chan *Payload
	OnError(event <-chan *Payload, handlers ...OnErrorHandler)
}
