package partner

import (
	"time"

	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/vasar-network/vails/lang"
	"golang.org/x/text/language"
)

// VasarPearl is an edited item for ender pearls.
type StrengthPowder struct {
	PartnerItem

	Locale language.Tag
}

func (s StrengthPowder) Run(user *user.User, on *user.User) {
	user.Player().AddEffect(effect.New(effect.Strength{}, 2, time.Second*7))
	user.Player().Message(text.Colourf(lang.Translate(s.Locale, "pi.strength_powder.use")))
}

func (s StrengthPowder) Name() string {
	return lang.Translatef(s.Locale, "pi.strength_powder.name")
}

func (s StrengthPowder) Meta() string {
	return "Strength Powder"
}

func (s StrengthPowder) Description() string {
	return lang.Translatef(s.Locale, "pi.strength_powder.description")
}

// Cooldown ...
func (StrengthPowder) Cooldown() time.Duration {
	return time.Second * 30
}

// MaxCount ...
func (StrengthPowder) MaxCount() int {
	return 64
}

// EncodeItem ...
func (StrengthPowder) EncodeItem() (name string, meta int16) {
	return "minecraft:blaze_powder", 0
}
