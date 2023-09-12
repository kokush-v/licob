package picker

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cam-per/licob"
	"github.com/gookit/event"
)

func (p *Picker) onDrop(e event.Event) error {
	data, _ := e.Get("data").(*licob.DropData)
	if data == nil {
		return nil
	}

	if data.M.ChannelID != p.opt.ChannelID {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.drops[data.M.ID] = struct{}{}

	resp, err := http.Get(data.CaptchaURL)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	captcha, err := NewCaptchaDecoder(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	err = captcha.Decode()
	if err != nil {
		captcha.SaveDebug(data.M.ID)
		log.Println(err)
		return err
	}

	captcha.Save(data.M.ID)

	go p.trigger(&EventCaptcha{
		DropID: data.M.ID,
		Codes:  [2]string(captcha.Codes()),
	})
	return nil
}

func (p *Picker) onPicked(e event.Event) error {
	data, _ := e.Get("data").(*licob.PickedData)
	if data == nil {
		return nil
	}

	if data.M.ChannelID != p.opt.ChannelID {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.drops[data.M.ID]; !ok {
		return nil
	}

	delete(p.drops, data.M.ID)

	go p.trigger(&EventPicked{
		DropID: data.M.ID,
	})

	fmt.Println("PICKED ", data.M.ID)
	return nil
}
