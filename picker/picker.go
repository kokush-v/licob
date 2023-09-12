package picker

import (
	"sync"

	"github.com/cam-per/licob"
	"github.com/gookit/event"
)

type Opt struct {
	ChannelID string
	Events    chan any
}

type listners struct {
	drop   event.Listener
	picked event.Listener
}

func (e *listners) listen() {
	event.Listen(licob.EventDrop, e.drop)
	event.Listen(licob.EventPicked, e.picked)
}

func (e *listners) close() {
	event.Std().RemoveListener(licob.EventDrop, e.drop)
	event.Std().RemoveListener(licob.EventPicked, e.picked)
}

type Picker struct {
	mu       sync.Mutex
	listners listners
	opt      *Opt
	drops    map[string]struct{}
	isOpen   bool
}

func NewPicker(opt *Opt) *Picker {
	result := &Picker{
		opt:   opt,
		drops: make(map[string]struct{}),
	}

	result.listners = listners{
		drop:   event.ListenerFunc(result.onDrop),
		picked: event.ListenerFunc(result.onPicked),
	}
	return result
}

func (p *Picker) trigger(e any) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isOpen {
		p.opt.Events <- e
	}
}

func (p *Picker) Listen() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isOpen = true

	p.listners.listen()
}

func (p *Picker) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isOpen = false

	p.listners.close()
}
