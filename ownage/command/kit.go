package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/ownagepe/hcf/ownage/kit"
	"github.com/vasar-network/vails/role"
)

type Kit struct {
	Kit kits
}

func (c Kit) Run(s cmd.Source, o *cmd.Output) {
	// TODO: cooldowns
	p := s.(*player.Player)
	switch string(c.Kit) {
	case "diamond":
		kit.Apply(kit.Diamond{}, p)
	case "bard":
		kit.Apply(kit.Bard{}, p)
	}
}

func (Kit) Allow(s cmd.Source) bool {
	return allow(s, false, role.Default{})
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