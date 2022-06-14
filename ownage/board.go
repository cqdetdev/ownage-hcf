package ownage

import (
	"fmt"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/ownagepe/hcf/ownage/item/partner"
	"github.com/ownagepe/hcf/ownage/module"
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/vasar-network/vails/lang"
)

func (v *Ownage) startBoards() {
	t := time.NewTicker(time.Second * 1)
	defer t.Stop()
	for {
		select {
		case <-v.c:
				return
		case <-t.C:
			for _, u := range user.All() {
				p := u.Player()
				b := scoreboard.New(lang.Translatef(p.Locale(), "scoreboard.title"))
				var msg []string
				if v.SOTW() {
					left := time.Until(time.Unix(v.config.MapInfo.SOTWEnd, 0))
					f := time.Unix(0, 0).UTC().Add(time.Duration(left)).Format("15:04:05")
					msg = append(msg, fmt.Sprintf("<aqua>SOTW Timer</aqua><white>: %s</white>", f))
				}
				if u.HasTimer() {
					left := time.Until(time.Unix(u.Timer().Expires.Unix(), 0))
					f := time.Unix(0, 0).UTC().Add(time.Duration(left)).Format("15:04:05")
					msg = append(msg, fmt.Sprintf("<aqua>PVP Timer</aqua><white>: %s</white>", f))
				}
				for _, cd := range u.Cooldowns() {
					if cd.Expired() || cd.TimeLeft().Milliseconds() < 0 { continue }
					if cd.Name == "partner_item" {
						f := time.Unix(0, 0).UTC().Add(time.Duration(cd.TimeLeft())).Format("5")
						msg = append(msg, fmt.Sprintf("<dark-purple>Partner Item</dark-purple><white>: %ss</white>", f))
					}

					if pi, ok := partner.ItemByMeta(cd.Name); ok {
						f := time.Unix(0, 0).UTC().Add(time.Duration(cd.TimeLeft())).Format("5")
						msg = append(msg, fmt.Sprintf("<dark-purple>%s</dark-purple><white>: %ss</white>", pi.Meta(), f))
					}
				}
				if module.Current().Started() {
					f := time.Unix(0, 0).UTC().Add(time.Duration(module.Current().TimeLeft())).Format("4:05")
					msg = append(msg, fmt.Sprintf("<red>%s</red><white>: %s</white>", module.Current().Name(), f))
				}
				_, _ = b.WriteString(lang.Translatef(p.Locale(), "scoreboard", " " + strings.Join(msg, "\n ")))
				b.RemovePadding()
				p.SendScoreboard(b)
			}

		}
	}
}