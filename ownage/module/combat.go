package module

import (
	"math/rand"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/ownagepe/hcf/ownage/data"
	ent "github.com/ownagepe/hcf/ownage/entity"
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/vasar-network/vails/lang"
	"golang.org/x/exp/slices"
)

// Combat is a handler that is used for combat related things, such as the combat tag.
type Combat struct {
	player.NopHandler

	u *user.User
}

// bannedCommands is a list of commands disallowed in combat.
var bannedCommands = []string{"/f home"}

// NewCombat ...
func NewCombat(u *user.User) *Combat {
	return &Combat{u: u}
}

// HandleCommandExecution ...
func (c *Combat) HandleCommandExecution(ctx *event.Context, cmd cmd.Command, _ []string) {
	if c.u.Tagged() && slices.Contains(bannedCommands, cmd.Name()) {
		c.u.Player().Message(lang.Translatef(c.u.Player().Locale(), "user.feature.disabled"))
		ctx.Cancel()
	}
}

// HandleHurt ...
func (c *Combat) HandleHurt(ctx *event.Context, _ *float64, s damage.Source) {
	if ctx.Cancelled() {
		// Was cancelled at some point, so just ignore this.
		return
	}

	var attacker *player.Player
	if a, ok := s.(damage.SourceEntityAttack); ok {
		if p, ok := a.Attacker.(*player.Player); ok {
			attacker = p
		}
	} else if t, ok := s.(damage.SourceProjectile); ok {
		if p, ok := t.Owner.(*player.Player); ok {
			attacker = p
		}
	}
	if attacker == nil {
		// No attacker, so we don't need to do anything.
		return
	}
}

var fall = []string{
	"death.message.fallDamage1",
	"death.message.fallDamage2",
	"death.message.fallDamage3",
}

var attack = []string{
	"death.message.attack1",
	"death.message.attack2",
	"death.message.attack3",
	"death.message.attack4",
} 

// HandleDeath ...
func (c *Combat) HandleDeath(src damage.Source) {
	if c.u.Tagged() {
		c.u.RemoveTag()
	}

	pots := c.u.Potions()
	pos := c.u.Player().Position()
	lightning := ent.NewLightning(pos)
	viewers := make([]*player.Player, 0, 16)
	for _, e := range c.u.Player().World().EntitiesWithin(cube.Box(-50, -25, -50, 50, 25, 50).Translate(pos), nil) {
		if p, ok := e.(*player.Player); ok {
			u, ok := user.Lookup(p)
			if ok {
				if set := u.Settings(); set.Visual.Lightning && (u != c.u || !set.Gameplay.InstantRespawn) {
					u.SendSound(pos, sound.Thunder{})
					u.SendSound(pos, sound.Explosion{})
					p.ShowEntity(lightning)
					viewers = append(viewers, p)
				}
			}
		}
	}

	time.AfterFunc(time.Millisecond*250, func() {
		for _, v := range viewers {
			v.HideEntity(lightning)
		}
	})

	if _, ok := src.(damage.SourceFall); ok {
		msg := fall[rand.Intn(len(fall))]
		c.u.Message(msg, c.u.Player().Name(), pots)
	}

	if _, ok := src.(damage.SourceEntityAttack); ok {
		msg := attack[rand.Intn(len(attack))]
		if c.u.Attacker() == nil { return }
		au, _ := user.Lookup(c.u.Attacker())
		c.u.Message(msg, c.u.Player().Name(), pots, au.Player().Name(), au.Potions())
	}

	if d := c.u.Attacker(); d != nil {
		if da, ok := user.Lookup(d); ok {
			stats := c.u.Stats()
			stats.Deaths++
			stats.KillStreak = 0
			c.u.SetStats(stats)
			go data.SaveUser(c.u)

			stats = da.Stats()
			stats.Kills++
			stats.KillStreak++
			if stats.KillStreak > stats.BestKillStreak {
				stats.BestKillStreak = stats.KillStreak
			}
			da.SetStats(stats)
			go data.SaveUser(da)
		}
	}
}
