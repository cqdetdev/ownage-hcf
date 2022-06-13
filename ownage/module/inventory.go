package module

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/ownagepe/hcf/ownage/user"
	"time"
)

// Inventory is a module that adds all items required to access basic functions outside the lobby.
type Inventory struct {
	player.NopHandler

	u *user.User
}

// NewInventory ...
func NewInventory(u *user.User) *Inventory {
	return &Inventory{u: u}
}

// HandleItemUseOnBlock ...
func (i *Inventory) HandleItemUseOnBlock(ctx *event.Context, _ cube.Pos, _ cube.Face, _ mgl64.Vec3) {
}

// HandleItemUse ...
func (i *Inventory) HandleItemUse(*event.Context) {
}

// HandleItemConsume ...
func (i *Inventory) HandleItemConsume(_ *event.Context, h item.Stack) {
	if _, ok := h.Value("head"); ok {
		i.u.Player().AddEffect(effect.New(effect.Regeneration{}, 3, time.Second*9))
	}
}
