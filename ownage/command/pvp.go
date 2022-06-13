package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/ownagepe/hcf/ownage/data"
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/vasar-network/vails/lang"
	"github.com/vasar-network/vails/role"
)

type PvpEnable struct {
	Sub enable
}

// Run ...
func (p PvpEnable) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	u, _ := user.Lookup(s.(*player.Player))
	if !u.HasTimer() {
		o.Print(lang.Translatef(l, "command.pvp.enable.notimer"))
		return
	}
	u.ExpireTimer()
	data.SaveUser(u)
	o.Print(lang.Translatef(l, "command.pvp.enable.success"))
}

func (PvpEnable) Allow(s cmd.Source) bool {
	return allow(s, false, role.Default{})
}

type (
	enable string
)

func (enable) SubName() string {
	return "enable"
}
