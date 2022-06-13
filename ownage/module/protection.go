package module

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/ownagepe/hcf/ownage/claim"
	"github.com/ownagepe/hcf/ownage/faction"
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/vasar-network/vails/lang"
)

// Protection is a module that ensures that players cannot break blocks or attack players in the lobby.
type Protection struct {
	player.NopHandler

	p *player.Player
}

// NewProtection ...
func NewProtection(p *player.Player) *Protection {
	return &Protection{p: p}
}

// HandleAttackEntity ...
func (p *Protection) HandleAttackEntity(ctx *event.Context, _ world.Entity, _ *float64, _ *float64, _ *bool) {

}

// HandleHurt ...
func (p *Protection) HandleHurt(ctx *event.Context, _ *float64, s damage.Source) {

}

// HandleFoodLoss ...
func (*Protection) HandleFoodLoss(ctx *event.Context, _ int, _ int) {
	ctx.Cancel()
}

// HandleBlockPlace ...
func (p *Protection) HandleBlockPlace(ctx *event.Context, pos cube.Pos, _ world.Block) {
	c, ok := claim.Lookup(pos.X(), pos.Z())
	if !ok { return }
	if c.Koth || c.Road || c.Spawn || c.Warzone {
		p.p.SendPopup(lang.Translatef(p.p.Locale(), "claim.cannotbuild", c.Name))
		ctx.Cancel()
	}

	u, _ := user.Lookup(p.p)
	if u.HasFaction() {
		f, _ := faction.Lookup(u.Faction())
		if c.Faction && c.Name != f.Name() {
			p.p.SendPopup(lang.Translatef(p.p.Locale(), "claim.cannotbuild", c.Name))
			ctx.Cancel()
		}
	}
}

// HandleBlockBreak ...
func (p *Protection) HandleBlockBreak(ctx *event.Context, pos cube.Pos, _ *[]item.Stack) {
	c, ok := claim.Lookup(pos.X(), pos.Z())
	if !ok { return }
	if c.Koth || c.Road || c.Spawn || c.Warzone {
		p.p.SendPopup(lang.Translatef(p.p.Locale(), "claim.cannotbreak", c.Name))
		ctx.Cancel()
	}

	u, _ := user.Lookup(p.p)
	if u.HasFaction() {
		f, _ := faction.Lookup(u.Faction())
		if c.Faction && c.Name != f.Name() {
			p.p.SendPopup(lang.Translatef(p.p.Locale(), "claim.cannotbreak", c.Name))
			ctx.Cancel()
		}
	}
}
