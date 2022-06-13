package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/player"
)

type Kill struct {
	Sub kill
}

func (c Kill) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)
	p.Hurt(25, damage.SourceFall{})
}

func (Kill) Allow(s cmd.Source) bool {
	return s.Name() == "xCqzzz"
}

type (
	kill string
)

func (kill) SubName() string {
	return "kill"
}