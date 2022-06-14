package user

import "time"

// These are a list of constants of constant partner item cooldowns
const (
	PartnerItem = "partner_item"
)

type Cooldown struct {
	Name string
	expiration time.Time
}

func (cd *Cooldown) Expired() bool                  { return cd.expiration.Before(time.Now()) }
func (cd *Cooldown) Expiration() time.Time          { return cd.expiration }
func (cd *Cooldown) UntilExpiration() time.Duration { return time.Until(cd.expiration) }
func (cd *Cooldown) SetCooldown(d time.Duration)    { cd.expiration = time.Now().Add(d) }

func NewCooldown(name string, length time.Duration) *Cooldown {
	return &Cooldown{Name: name, expiration: time.Now().Add(length)}
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