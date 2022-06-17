package partner

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/vasar-network/vails/lang"
	"golang.org/x/text/language"
)

// VasarPearl is an edited item for ender pearls.
type SwitchStick struct {
	PartnerItem

	Locale language.Tag
}

var count map[string]countTarget = map[string]countTarget{}

func (s SwitchStick) Run(user *user.User, on *user.User) {
	if on == nil { return }
	if _, ok := count[user.Player().Name()]; !ok {
		count[user.Player().Name()] = countTarget{
			target: on.Player().Name(),
			times: 1,
		}
		return
	} else {
		o := count[user.Player().Name()]
		o.times++
		count[user.Player().Name()] = o
	}

	if count[user.Player().Name()].times > 2 {
		on.Player().Move(mgl64.Vec3{}, on.Player().Data().Yaw / 2, 0)
		delete(count, user.Player().Name())
	}
}

func (s SwitchStick) Name() string {
	return lang.Translatef(s.Locale, "pi.switch_stick.name")
}

func (s SwitchStick) Meta() string {
	return "Switch Stick"
}

func (s SwitchStick) Description() string {
	return lang.Translatef(s.Locale, "pi.switch_stick.description")
}

// Cooldown ...
func (SwitchStick) Cooldown() time.Duration {
	return time.Second * 45
}

// MaxCount ...
func (SwitchStick) MaxCount() int {
	return 64
}

// EncodeItem ...
func (SwitchStick) EncodeItem() (name string, meta int16) {
	return "minecraft:stick", 0
}
