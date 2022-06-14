package module

import (
	"time"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/ownagepe/hcf/ownage/data"
	pi "github.com/ownagepe/hcf/ownage/item/partner"
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/vasar-network/vails/lang"
)

// PartnerItem is a module that handles partner item actions.
type PartnerItem struct {
	player.NopHandler

	u *user.User
}

// NewPartnerItem creates a new partner item module.
func NewPartnerItem(u *user.User) *PartnerItem {
	return &PartnerItem{u: u}
}

// HandleItemUse ...
func (m *PartnerItem) HandleItemUse(ctx *event.Context) {
	item, off := m.u.Player().HeldItems()
	it := item.Item()
	if pi, ok := it.(pi.PartnerItem); ok {
		if c, has := m.u.Cooldown(pi.Meta()); has {
			if c.Expired() {
				m.u.RemoveCooldown(pi.Meta())
			} else {
				m.u.Player().Message(lang.Translatef(m.u.Player().Locale(), "pi.cooldown.item", int(c.UntilExpiration().Seconds())))
				ctx.Cancel()
				return
			}
		}
		if c, has := m.u.Cooldown(user.PartnerItem); has {
			if c.Expired() {
				m.u.RemoveCooldown(user.PartnerItem)
			} else {
				m.u.Player().Message(lang.Translatef(m.u.Player().Locale(), "pi.cooldown", int(c.UntilExpiration().Seconds())))
				ctx.Cancel()
				return
			}
		}
		pi.Run(m.u, nil)
		m.u.Player().SetHeldItems(item.Grow(-1), off)
		m.u.AddCooldown(user.NewCooldown(pi.Meta(), pi.Cooldown()))
		m.u.AddCooldown(user.NewCooldown(user.PartnerItem, time.Second * 15))
		go data.SaveUser(m.u)
	}
}
