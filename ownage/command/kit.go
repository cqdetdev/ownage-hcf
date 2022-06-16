package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/ownagepe/hcf/ownage/kit"
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/vasar-network/vails/role"
)

type Kit struct {
	Kit kits
}

func (c Kit) Run(s cmd.Source, o *cmd.Output) {
	// TODO: cooldowns
	p := s.(*player.Player)
	u, _ := user.Lookup(p)
	clearEffects(u)
	switch string(c.Kit) {
	case "diamond":
		kit.Apply(kit.Diamond{}, p)
		kit.ApplyEffects(kit.Diamond{}, p)
		u.SetKit(kit.Diamond{})
	case "bard":
		kit.Apply(kit.Bard{}, p)
		kit.ApplyEffects(kit.Bard{}, p)
		u.SetKit(kit.Bard{})
	}
}

func (Kit) Allow(s cmd.Source) bool {
	return allow(s, false, role.Default{})
}

func clearEffects(u *user.User) {
	if u.Kit() != nil {
		kit.RemoveEffects(u.Kit(), u.Player())
	}
}

type kits string

// Type ...
func (kits) Type() string {
	return "kits"
}

// Options ...
func (kits) Options(cmd.Source) []string {
	return []string{
		"diamond",
		"bard",
	}
}