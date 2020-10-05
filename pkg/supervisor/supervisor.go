package supervisor

// Runner is main struct to store supervisor's driver
type Runner struct {
	driver        Driver
	errorHandlers []OnErrorHandler
	errorEvent    <-chan *Payload
}

// New used to create new instance of supervisor
func New(driver Driver, handlers ...OnErrorHandler) *Runner {
	return &Runner{driver: driver, errorHandlers: handlers}
}

// Start used to start supervisor's process
func (r *Runner) Start() *Runner {
	r.errorEvent = r.driver.Watch()
	return r
}

// Handle used to listen to error events, by default we put main process to listening events
// inside a goroutine
func (r *Runner) Handle() {
	go func() {
		for p := range r.errorEvent {
			if len(r.errorHandlers) >= 1 && p != nil {
				r.driver.OnError(p, r.errorHandlers...)
			}
		}
	}()
}
