package ownage

import (
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/vasar-network/vails/lang"
	"github.com/vasar-network/vails/role"
	"strings"
	"time"
)

// startBroadcasts starts sending a new broadcast every five minutes.
func (v *Ownage) startBroadcasts() {
	broadcasts := []string{
		"vasar.broadcast.discord",
		"vasar.broadcast.store",
		"vasar.broadcast.emojis",
		"vasar.broadcast.settings",
		"vasar.broadcast.duels",
		"vasar.broadcast.feedback",
		"vasar.broadcast.report",
	}

	var cursor int
	t := time.NewTicker(time.Minute * 5)
	defer t.Stop()
	for {
		select {
		case <-v.c:
			return
		case <-t.C:
			message := broadcasts[cursor]
			for _, u := range user.All() {
				p := u.Player()
				p.Message(lang.Translatef(p.Locale(), "vasar.broadcast.notice", lang.Translate(p.Locale(), message)))
			}

			if cursor++; cursor == len(broadcasts) {
				cursor = 0
			}
		}
	}
}

// startPlayerBroadcasts starts sending a new player broadcast every five minutes.
func (v *Ownage) startPlayerBroadcasts() {
	t := time.NewTicker(time.Minute * 10)
	defer t.Stop()
	for {
		select {
		case <-v.c:
			return
		case <-t.C:
			users := user.All()
			var plus []string
			for _, u := range users {
				if u.Roles().Contains(role.Plus{}) {
					plus = append(plus, u.Player().Name())
				}
			}

			for _, u := range users {
				p := u.Player()
				p.Message(lang.Translatef(p.Locale(), "vasar.broadcast.plus", len(plus), strings.Join(plus, ", ")))
			}
		}
	}
}
