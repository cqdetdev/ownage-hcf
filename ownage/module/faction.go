package module

import (
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/ownagepe/hcf/ownage/data"
	"github.com/ownagepe/hcf/ownage/faction"
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/vasar-network/vails/lang"
)

type Faction struct {
	player.NopHandler

	u *user.User
}

// NewCombat ...
func NewFaction(u *user.User) *Faction {
	return &Faction{u: u}
}

func (f *Faction) HandleDeath(src damage.Source) {
	if f.u.HasFaction() {
		fac, _ := faction.Lookup(f.u.Faction())
		fac.DecDtr()
		fac.Broadcast(lang.Translatef(f.u.Player().Locale(), "faction.member.death", f.u.Player().Name(), fac.Dtr()))
		data.SaveFaction(fac)
	}
}