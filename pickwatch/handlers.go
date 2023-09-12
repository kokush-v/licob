package pickwatch

import (
	"fmt"
	"sort"
	"time"

	"github.com/cam-per/licob"
	"github.com/gookit/event"
)

func (p *Pickwatch) onDrop(e event.Event) error {
	data, _ := e.Get("data").(*licob.DropData)
	if data == nil {
		return nil
	}

	if data.M.ChannelID != p.opt.ChannelID {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.drop = newDrop(data.M)
	p.drop.Currency = data.Currency

	go p.trigger(p.drop)

	fmt.Println("DROP   ", data.Currency)
	return nil
}

func (p *Pickwatch) onPick(e event.Event) error {
	data, _ := e.Get("data").(*licob.PickData)
	if data == nil {
		return nil
	}

	if data.M.ChannelID != p.opt.ChannelID {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.drop == nil {
		return nil
	}

	if data.M.Author == nil {
		return nil
	}

	if p.drop == nil {
		return nil
	}

	pick := newPick(p.drop, data.M)
	pick.Code = data.Code
	pick.Ok = false

	p.drop.Picks = append(p.drop.Picks, pick)

	go p.trigger(pick)

	fmt.Println("PICK   ", data.Code, data.M.Author)
	return nil
}

func (p *Pickwatch) onWin(e event.Event) error {
	data, _ := e.Get("data").(*licob.WinData)
	if data == nil {
		return nil
	}

	if data.M.ChannelID != p.opt.ChannelID {
		return nil
	}

	time.Sleep(2 * time.Second)

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.drop == nil {
		return nil
	}

	var winpick *Pick
	for i := len(p.drop.Picks) - 1; i >= 0; i-- {
		if p.drop.Picks[i].User.ID == data.WinnerID {
			winpick = p.drop.Picks[i]
			break
		}
	}

	if winpick == nil {
		return nil
	}
	winpick.Ok = true

	p.drop.Since = winpick.Since
	p.drop.Winner = winpick.User

	sort.SliceStable(p.drop.Picks, func(i, j int) bool {
		return p.drop.Picks[i].t.Unix() < p.drop.Picks[j].t.Unix()
	})

	win := EventWin(*p.drop)
	go p.trigger(&win)

	fmt.Println("WIN    ", winpick.Code, data.M.Author)
	return nil
}
