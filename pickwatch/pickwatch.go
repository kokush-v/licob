package pickwatch

import (
	"sync"

	"github.com/cam-per/licob"
	"github.com/gookit/event"
)

type Opt struct {
	Events    chan any
	ChannelID string
}

type listners struct {
	drop event.Listener
	pick event.Listener
	win  event.Listener
}

func (e *listners) listen() {
	event.Listen(licob.EventDrop, e.drop)
	event.Listen(licob.EventPick, e.pick)
	event.Listen(licob.EventWin, e.win)
}

func (e *listners) close() {
	event.Std().RemoveListener(licob.EventDrop, e.drop)
	event.Std().RemoveListener(licob.EventPick, e.pick)
	event.Std().RemoveListener(licob.EventWin, e.win)
}

type Pickwatch struct {
	mu       sync.Mutex
	listners listners
	opt      *Opt
	drop     *Drop
	isOpen   bool
}

func NewPickwatch(opt *Opt) *Pickwatch {
	result := &Pickwatch{
		opt: opt,
	}

	result.listners = listners{
		drop: event.ListenerFunc(result.onDrop),
		pick: event.ListenerFunc(result.onPick),
		win:  event.ListenerFunc(result.onWin),
	}
	return result
}

func (p *Pickwatch) trigger(e any) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isOpen {
		p.opt.Events <- e
	}
}

func (p *Pickwatch) Listen() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isOpen = true
	p.listners.listen()
}

func (p *Pickwatch) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isOpen = false
	p.listners.close()
}
