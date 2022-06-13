package module

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/ownagepe/hcf/ownage/claim"
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/vasar-network/vails/lang"
)

type Claim struct {
	player.NopHandler

	u *user.User
}

// NewCombat ...
func NewClaim(u *user.User) *Claim {
	return &Claim{u: u}
}

// HandleMove ...
func (c *Claim) HandleMove(ctx *event.Context, pos mgl64.Vec3, _, _ float64) {
	cl, _ := claim.Lookup(int(pos.X()), int(pos.Z()))
	if c.u.Claim() == nil {
		c.u.SetClaim(cl)
	}

	if c.u.Claim().Name != cl.Name {
		leave := "claim.leaving.deathban"
		enter := "claim.entering.deathban"
		if c.u.Claim().Faction && c.u.Claim().Name == c.u.Faction() {
			leave = "claim.leaving.friendly.deathban"
		}
		if c.u.Claim().Spawn {
			leave = "claim.leaving.special.nondeathban"
		}
		if c.u.Claim().Koth || c.u.Claim().Warzone || c.u.Claim().Road {
			leave = "claim.leaving.special.deathban"
		}

		if cl.Faction && cl.Name == c.u.Faction() {
			enter = "claim.leaving.friendly.deathban"
		}
		if cl.Spawn {
			enter = "claim.entering.special.nondeathban"
		}
		if cl.Koth || cl.Warzone || cl.Road {
			enter = "claim.entering.special.deathban"
		}

		c.u.Player().Message(lang.Translatef(c.u.Player().Locale(), enter, cl.Name))
		c.u.Player().Message(lang.Translatef(c.u.Player().Locale(), leave, c.u.Claim().Name))

		c.u.SetClaim(cl)
	}

}
