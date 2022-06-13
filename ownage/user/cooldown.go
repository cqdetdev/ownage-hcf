package user

import "time"

// These are a list of constants of constant partner item cooldowns
const (
	PartnerItem = "partner_item"
)

type Cooldown struct {
	Name string
	Length time.Duration
	Last time.Time
}

func (c *Cooldown) Expired() bool {
	diff := time.Until(c.Last)
	return diff > c.Length
}

func (c *Cooldown) TimeLeft() time.Duration {
	return c.Length - time.Since(c.Last)
}

func NewCooldown(name string, length time.Duration, last time.Time) *Cooldown {
	return &Cooldown{Name: name, Length: length, Last: last}
}

func (u *User) Cooldown(n string) (*Cooldown, bool) {
	for _, c := range u.cooldowns {
		if c.Name == n {
			return c, true
		}
	}
	return nil, false
}

func (u *User) AddCooldown(cd *Cooldown) {
	u.cooldowns = append(u.cooldowns, cd)
}

func (u *User) RemoveCooldown(n string) bool {
	for i, c := range u.cooldowns {
		if c.Name == n {
			u.cooldowns = append(u.cooldowns[:i], u.cooldowns[i+1:]...)
			return true
		}
	}
	return false
}

func (u *User) ClearCooldowns() {
	u.cooldowns = []*Cooldown{}
}