package core

type Event struct {
	Data     interface{}
	Listener chan interface{}
}

func InitEvent() *Event {
	return &Event{
		Listener: make(chan interface{}),
	}
}

func (e *Event) Publish(data interface{}) {
	e.Data = data
	go func() {
		e.Listener <- e.Data
	}()
}

func (e *Event) Subscribe() <-chan interface{} {
	return e.Listener
}

func (e *Event) Unsubscribe() {
	close(e.Listener)
}
