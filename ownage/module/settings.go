package module

import (
	"math"

	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/ownagepe/hcf/ownage/user"
)

// Settings is a module which handles player settings.
type Settings struct {
	player.NopHandler

	u *user.User
}

// NewSettings ...
func NewSettings(u *user.User) *Settings {
	return &Settings{u: u}
}

// HandleJoin ...
func (s *Settings) HandleJoin() {
	// cape, _ := cape.ByName(s.u.Settings().Advanced.Cape)
	// skin := s.u.Player().Skin()
	// skin.Cape = cape.Cape()
	// s.u.Player().SetSkin(skin)
}

// HandleSkinChange ...
func (s *Settings) HandleSkinChange(ctx *event.Context, skin *skin.Skin) {
	// cape, _ := cape.ByName(s.u.Settings().Advanced.Cape)
	// (*skin).Cape = cape.Cape()
}

// HandleMove ...
func (s *Settings) HandleMove(_ *event.Context, pos mgl64.Vec3, newYaw, _ float64) {
	p := s.u.Player()
	if !s.u.Settings().Gameplay.ToggleSprint || p.Sprinting() {
		return
	}
	delta := pos.Sub(p.Position())
	if mgl64.FloatEqual(delta[0], 0) && mgl64.FloatEqual(delta[2], 0) {
		return
	}
	diff := (mgl64.RadToDeg(math.Atan2(delta[2], delta[0])) - 90) - newYaw
	if diff < 0 {
		diff += 360
	}
	if diff <= 65 && diff >= -65 {
		p.StartSprinting()
	}
}

// HandleDeath ...
func (s *Settings) HandleDeath(src damage.Source) {
}

// HandleAttackEntity ...
func (s *Settings) HandleAttackEntity(_ *event.Context, e world.Entity, _, _ *float64, _ *bool) {
	if e, ok := e.(*player.Player); ok && !e.AttackImmune() {
		s.u.MultiplyParticles(e, s.u.Settings().Advanced.ParticleMultiplier)
	}
}
