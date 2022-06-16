package module

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/ownagepe/hcf/ownage/user"
)

type Crate struct {
	player.NopHandler

	u *user.User
}

var kothCrate = cube.Pos{10, -59, 10}

func NewCrate(u *user.User) *Crate {
	return &Crate{u: u}
}

func (c *Crate) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	if pos.Vec3().ApproxEqual(kothCrate.Vec3()) {
		c.u.Player().Message("This is KOTH crate information")
		ctx.Cancel()
	}
}