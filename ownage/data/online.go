package data

import (
	"encoding/hex"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/upper/db/v4"
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/vasar-network/vails"
	"github.com/vasar-network/vails/role"
	"golang.org/x/crypto/sha3"
)

// LoadUser loads a *user.User from the database using a *player.Player. If the user does not exist, it will be created.
func LoadUser(p *player.Player) (*user.User, error) {
	result := sess.Collection("users").Find(db.Or(db.Cond{"name": strings.ToLower(p.Name())}, db.Cond{"xuid": p.XUID()}))
	addr, _ := netip.ParseAddrPort(p.Addr().String())

	s := sha3.New256()
	s.Write(addr.Addr().AsSlice())
	s.Write([]byte(salt))
	address := hex.EncodeToString(s.Sum(nil))

	if ok, _ := result.Exists(); !ok {
		return user.NewUser(p,
			user.NewRoles([]vails.Role{role.Default{}}, map[vails.Role]time.Time{}),
			[]*user.Cooldown{},
			user.DefaultTimer(),
			user.DefaultSettings(),
			user.DefaultStats(),
			0,
			time.Now(),
			0,
			address,
			false,
			"",
		), nil
	}

	var data userData
	if err := result.One(&data); err != nil {
		return nil, fmt.Errorf("load user: %v", err)
	}

	var roles []vails.Role
	expirations := make(map[vails.Role]time.Time)
	for _, dat := range data.Roles {
		if dat.Expires && time.Now().After(dat.Expiration) {
			continue
		}
		r, ok := role.ByName(dat.Name)
		if !ok {
			return nil, fmt.Errorf("load user: role %s does not exist", dat.Name)
		}
		roles = append(roles, r)
		if dat.Expires {
			expirations[r] = dat.Expiration
		}
	}
	var cooldowns []*user.Cooldown
	for _, c := range data.Cooldowns {
		cooldowns = append(cooldowns, user.NewCooldown(c.Name, c.Length, c.Last))
	}
	return user.NewUser(p, user.NewRoles(roles, expirations), cooldowns, &user.Timer{Has: data.Timer.Has, Expires: data.Timer.Expires}, data.Settings, data.Stats, data.Money, data.FirstLogin, data.PlayTime, address, data.Whitelisted, data.Faction), nil
}

// SaveUser saves a *user.User to the database. If an error occurs, it will be returned to the second return value.
func SaveUser(u *user.User) error {
	var roles []roleData
	for _, r := range u.Roles().All() {
		data := roleData{Name: r.Name()}
		if e, ok := u.Roles().Expiration(r); ok {
			data.Expiration, data.Expires = e, true
		}
		roles = append(roles, data)
	}

	var cooldowns []cooldownData
	for _, c := range u.Cooldowns() {
		data := cooldownData{
			Name: c.Name,
			Length: c.Length,
			Last: c.Last,
		}
		cooldowns = append(cooldowns, data)
	}

	p := u.Player()
	users := sess.Collection("users")
	data := userData{
		XUID:         p.XUID(),
		DisplayName:  p.Name(),
		Name:         strings.ToLower(p.Name()),
		DeviceID:     u.DeviceID(),
		SelfSignedID: u.SelfSignedID(),
		Address:      u.HashedAddress(),
		Whitelisted:  u.Whitelisted(),

		FirstLogin: u.FirstLogin(),
		PlayTime:   u.PlayTime(),

		Settings: u.Settings(),
		Stats: u.Stats(),

		Timer: timerData{
			Has: u.Timer().Has,
			Expires: u.Timer().Expires,
		},

		Cooldowns: cooldowns,
		Roles: roles,

		Faction: u.Faction(),
	}

	entry := users.Find(db.Or(db.Cond{"name": strings.ToLower(p.Name())}, db.Cond{"xuid": p.XUID()}))
	if ok, _ := entry.Exists(); ok {
		return entry.Update(data)
	}
	_, err := users.Insert(data)
	return err
}
