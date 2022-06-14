package data

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/ownagepe/hcf/ownage/user"
	"github.com/upper/db/v4"
	"github.com/vasar-network/vails"
	"github.com/vasar-network/vails/role"
)

// User is a structure containing the data of an offline user. It also contains useful functions that can be used
// externally to modify offline user data, such as roles.
type User struct {
	// xuid is the xuid of the user.
	xuid string
	// displayName is the display name of the user.
	displayName string
	// name is the name of the user.
	name string
	// deviceID is the device ID of the user.
	deviceID string
	// selfSignedID is the self-signed ID of the user.
	selfSignedID string
	// address is the hashed IP address of the user.
	address string
	// firstLogin is the time the user first logged in.
	firstLogin time.Time
	// playTime is the duration the user has played for on the server.
	playTime time.Duration

	// Money is the amount of money of the User has.
	Money int
	// Roles is the roles manager of the User.
	Roles *user.Roles
	// Cooldowns is a list of cooldown of the User.
	Cooldowns []*user.Cooldown
	// Settings contains the settings of the User.
	Settings user.Settings
	// Stats contains the stats of the user.
	Stats user.Stats
	// Timer contains timer data of the user.
	Timer *user.Timer
	// Whitelisted is true if the user is whitelisted.
	Whitelisted bool

	// Faction is the name of the user's faction.
	Faction string
}

// NewOfflineUser creates a new offline user with the provided data.
func NewOfflineUser(name string) User {
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(rand.Intn(10))
	}
	return User{
		displayName: strings.ToLower(name),
		name:        strings.ToLower(name),
		Roles:       user.NewRoles([]vails.Role{role.Default{}}, map[vails.Role]time.Time{}),
		Timer: 		 user.DefaultTimer(),
		Settings:    user.DefaultSettings(),
		Stats:       user.DefaultStats(),
	}
}

// SearchOfflineUsers searches for offline users using the given conditions.
func SearchOfflineUsers(cond ...any) ([]User, error) {
	var data []userData
	err := sess.Collection("users").Find(cond...).All(&data)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(data))
	for _, d := range data {
		u, _ := parseData(d)
		users = append(users, u)
	}
	return users, nil
}

// OrderedOfflineUsers searches and orders offline users using a query and limit.
func OrderedOfflineUsers(query string, limit int) ([]User, error) {
	var data []userData
	err := sess.Collection("users").Find().Limit(limit).OrderBy(query).All(&data)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(data))
	for _, d := range data {
		u, _ := parseData(d)
		users = append(users, u)
	}
	return users, nil
}

// LoadOfflineUser loads an offline User from the database by checking XUID and Name. If the user does not exist, an error will be
// returned to the second return value.
func LoadOfflineUser(id string) (User, error) {
	result := sess.Collection("users").Find(db.Or(db.Cond{"name": strings.ToLower(id)}, db.Cond{"xuid": id}))
	if ok, _ := result.Exists(); !ok {
		return User{}, fmt.Errorf("load user: user does not exist (%s)", id)
	}
	var data userData
	if err := result.One(&data); err != nil {
		return User{}, fmt.Errorf("load user: %v", err)
	}
	return parseData(data)
}

// SaveOfflineUser saves an offline User to the database. If an error occurs, it will be returned to the second return
// value.
func SaveOfflineUser(u User) error {
	var roles []roleData
	for _, r := range u.Roles.All() {
		data := roleData{Name: r.Name()}
		if e, ok := u.Roles.Expiration(r); ok {
			data.Expiration, data.Expires = e, true
		}
		roles = append(roles, data)
	}

	var cooldowns []cooldownData
	for _, c := range u.Cooldowns {
		data := cooldownData{
			Name: c.Name,
			Expires: c.Expiration(),
		}
		cooldowns = append(cooldowns, data)
	}

	users := sess.Collection("users")
	data := userData{
		XUID:         u.XUID(),
		Name:         u.Name(),
		DisplayName:  u.DisplayName(),
		DeviceID:     u.DeviceID(),
		SelfSignedID: u.SelfSignedID(),

		FirstLogin: u.FirstLogin(),
		PlayTime:   u.PlayTime(),

		Whitelisted: u.Whitelisted,
		Settings:    u.Settings,
		Stats:    u.Stats,

		Money: u.Money,
		Roles: roles,
		Cooldowns: cooldowns,

		Timer: timerData{
			Has: u.Timer.Has,
			Expires: u.Timer.Expires,
		},

		Faction: u.Faction,
	}

	cond := db.Cond{"xuid": u.XUID()}
	if len(u.XUID()) == 0 {
		cond = db.Cond{"name": u.Name()}
	}
	entry := users.Find(cond)
	if ok, _ := entry.Exists(); ok {
		return entry.Update(data)
	}
	_, err := users.Insert(data)
	return err
}

// parseData parses userData into an offline User.
func parseData(data userData) (User, error) {
	roles := make([]vails.Role, 0, len(data.Roles))
	expirations := make(map[vails.Role]time.Time)
	for _, dat := range data.Roles {
		r, ok := role.ByName(dat.Name)
		if !ok {
			return User{}, fmt.Errorf("load user: role %s does not exist", dat.Name)
		}
		roles = append(roles, r)
		if dat.Expires {
			expirations[r] = dat.Expiration
		}
	}

	var cooldowns []*user.Cooldown
	for _, c := range data.Cooldowns {
		cooldowns = append(cooldowns, user.NewCooldown(c.Name, time.Now().Sub(c.Expires)))
	}

	return User{
		xuid:         data.XUID,
		displayName:  data.DisplayName,
		name:         data.Name,
		deviceID:     data.DeviceID,
		selfSignedID: data.SelfSignedID,
		address:      data.Address,
		firstLogin:   data.FirstLogin,
		playTime:     data.PlayTime,

		Money: data.Money,
		Roles:       user.NewRoles(roles, expirations),
		Cooldowns: 	 cooldowns,
		Timer: &user.Timer{Has: data.Timer.Has, Expires: data.Timer.Expires},
		Whitelisted: data.Whitelisted,
		Settings:    data.Settings,
		Stats:       data.Stats,
	}, nil
}

// XUID returns the XUID of the offline user.
func (u User) XUID() string {
	return u.xuid
}

// DisplayName returns the display name of the offline user.
func (u User) DisplayName() string {
	return u.displayName
}

// Name returns the name of the offline user.
func (u User) Name() string {
	return u.name
}

// DeviceID returns the device ID of the offline user.
func (u User) DeviceID() string {
	return u.deviceID
}

// SelfSignedID returns the self-signed id of the offline user.
func (u User) SelfSignedID() string {
	return u.selfSignedID
}

// Address returns the hashed and salted ip address of the offline user.
func (u User) Address() string {
	return u.address
}

// FirstLogin returns the time the user first logged in.
func (u User) FirstLogin() time.Time {
	return u.firstLogin
}

// PlayTime returns the duration of time the user has played.
func (u User) PlayTime() time.Duration {
	return u.playTime
}
