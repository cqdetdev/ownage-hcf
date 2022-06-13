package user

import (
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"
	_ "unsafe"

	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	c "github.com/ownagepe/hcf/ownage/claim"
	ent "github.com/ownagepe/hcf/ownage/entity"
	it "github.com/ownagepe/hcf/ownage/item"
	"github.com/vasar-network/vails/lang"
	"github.com/vasar-network/vails/role"
	"github.com/vasar-network/vails/sets"
	"golang.org/x/exp/maps"
)

var (
	userMu    sync.Mutex
	users     = map[*player.Player]*User{}
	usersXUID = map[string]*User{}

	frozen = sets.New[string]()
)

// All returns a slice of all the users.
func All() []*User {
	userMu.Lock()
	defer userMu.Unlock()
	return maps.Values(users)
}

// Message will broadcast a message to every user using that user's locale.
func Message(key string, args ...any) {
	for _, u := range All() {
		u.Player().Message(lang.Translatef(u.Player().Locale(), key, args...))
	}
}

// Count returns the total user count.
func Count() int {
	userMu.Lock()
	defer userMu.Unlock()
	return len(users)
}

// Lookup looks up the user.User of a player.Player passed.
func Lookup(p *player.Player) (*User, bool) {
	userMu.Lock()
	defer userMu.Unlock()
	u, ok := users[p]
	return u, ok
}

// LookupXUID looks up the user.User of a XUID passed.
func LookupXUID(xuid string) (*User, bool) {
	userMu.Lock()
	defer userMu.Unlock()
	u, ok := usersXUID[xuid]
	return u, ok
}

// User is an extension of the Dragonfly player that adds a few extra features required by Vasar.
type User struct {
	p *player.Player

	hashedAddress string
	whitelisted   atomic.Bool
	address       net.Addr

	firstLogin time.Time
	joinTime   time.Time
	playTime   time.Duration

	displayName atomic.Value[string]

	lastMessageFrom atomic.Value[string]
	lastMessage     atomic.Value[time.Time]

	launchDelay         atomic.Value[time.Time]
	pearlCoolDown       atomic.Bool
	projectilesDisabled atomic.Bool

	frozen atomic.Bool

	tagMu         sync.Mutex
	tagExpiration time.Time
	attacker      *player.Player
	tagC          chan struct{}

	clickWatchersMu sync.Mutex
	clickWatchers   sets.Set[*User]
	watchingClick   *User

	rodMu sync.Mutex
	hook  *ent.FishingHook

	vanished atomic.Bool

	recentOpponent atomic.Value[string]

	settings atomic.Value[Settings]
	stats    atomic.Value[Stats]

	money int

	roles *Roles

	cooldowns []*Cooldown

	timer *Timer

	faction string
	claim   *c.Claim

	s *session.Session
}

// NewUser creates a new user from a Dragonfly player along with list of roles and settings.
func NewUser(p *player.Player, r *Roles, cooldowns []*Cooldown, timer *Timer, settings Settings, stats Stats, money int, firstLogin time.Time, playTime time.Duration, hashedAddress string, whitelisted bool, faction string) *User {
	u := &User{
		p:             p,
		address:       p.Addr(),
		whitelisted:   *atomic.NewBool(whitelisted),
		hashedAddress: hashedAddress,
		s:             player_session(p),

		joinTime:   time.Now(),
		firstLogin: firstLogin,
		playTime:   playTime,
 
		tagC: make(chan struct{}, 1),

		money: money,

		cooldowns: cooldowns,
		roles: r,

		clickWatchers: sets.New[*User](),

		settings: *atomic.NewValue(settings),
		stats:    *atomic.NewValue(stats),

		timer: timer,

		faction: faction,
	}
	u.displayName.Store(p.Name())

	if u.s != session.Nop {
		f := reflect.ValueOf(u.s).Elem().FieldByName("handlers")

		f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
		f.SetMapIndex(reflect.ValueOf(uint32(packet.IDPlayerAuthInput)), reflect.ValueOf(PlayerAuthInputHandler{u: u}))

		s := u.Settings()
		u.SetSettings(s)
	}
	u.roles.sortRoles()
	u.SetNameTagFromRole()

	userMu.Lock()
	users[p] = u
	usersXUID[p.XUID()] = u
	if frozen.Contains(p.XUID()) {
		u.p.SetImmobile()
		u.frozen.Toggle()
	}
	userMu.Unlock()
	return u
}

// DisplayName returns the display name of the user.
func (u *User) DisplayName() string {
	return u.displayName.Load()
}

// SetDisplayName sets the display name of the user.
func (u *User) SetDisplayName(name string) {
	u.displayName.Store(name)
}

// Vanished returns whether the user is vanished or not.
func (u *User) Vanished() bool { return u.vanished.Load() }

// ToggleVanish toggles the user's vanish state.
func (u *User) ToggleVanish() { u.vanished.Toggle() }

// Potions returns the amount of potions the user has.
func (u *User) Potions() (n int) {
	for _, i := range u.p.Inventory().Items() {
		if p, ok := i.Item().(it.VasarPotion); ok && p.Type == potion.StrongHealing() {
			n++
		}
	}
	return n
}

// RenewLastMessage renews the last time a message was sent from a player.
func (u *User) RenewLastMessage() {
	u.lastMessage.Store(time.Now())
}

// CanSendMessage returns true if the user can send a message.
func (u *User) CanSendMessage() bool {
	return u.Roles().Contains(role.Operator{}) || time.Since(u.lastMessage.Load()) > time.Second*2
}

// StartWatchingClicks starts watchingClick the user's clicks.
func (u *User) StartWatchingClicks(user *User) {
	u.clickWatchersMu.Lock()
	if u.watchingClick != nil {
		u.watchingClick.RemoveClickWatcher(u)
	}
	u.watchingClick = user
	u.clickWatchersMu.Unlock()
	user.AddClickWatcher(u)
}

// StopWatchingClicks stops watchingClick the user's clicks.
func (u *User) StopWatchingClicks() {
	u.clickWatchersMu.Lock()
	if u.watchingClick == nil {
		u.clickWatchersMu.Unlock()
		return
	}
	user := u.watchingClick
	u.watchingClick = nil
	u.clickWatchersMu.Unlock()
	user.RemoveClickWatcher(u)

}

// AddClickWatcher adds a user to the clicksWatchers set.
func (u *User) AddClickWatcher(user *User) {
	u.clickWatchersMu.Lock()
	defer u.clickWatchersMu.Unlock()
	u.clickWatchers.Add(user)
}

// RemoveClickWatcher removes a user from the clicksWatchers set.
func (u *User) RemoveClickWatcher(user *User) {
	u.clickWatchersMu.Lock()
	defer u.clickWatchersMu.Unlock()
	u.clickWatchers.Delete(user)
}

// ClickWatchers returns the users watchingClick the user.
func (u *User) ClickWatchers() (users []*User) {
	u.clickWatchersMu.Lock()
	for usr := range u.clickWatchers {
		users = append(users, usr)
	}
	u.clickWatchersMu.Unlock()
	return users
}

// WatchingClicks returns the user it is currently watchingClick.
func (u *User) WatchingClicks() *User {
	u.clickWatchersMu.Lock()
	defer u.clickWatchersMu.Unlock()
	return u.watchingClick
}

// SetRecentOpponent sets the opponent that the user last matched against.
func (u *User) SetRecentOpponent(opponent *User) {
	u.recentOpponent.Store(opponent.Player().XUID())
}

// ResetRecentOpponent resets the opponent that the user last matched against.
func (u *User) ResetRecentOpponent() {
	u.recentOpponent.Store("")
}

// RecentOpponent returns the opponent that the user last matched against.
func (u *User) RecentOpponent() (*User, bool) {
	return LookupXUID(u.recentOpponent.Load())
}

// Whitelisted returns true if the user is whitelisted.
func (u *User) Whitelisted() bool {
	return u.whitelisted.Load()
}

// Whitelist adds the user to the whitelist.
func (u *User) Whitelist() {
	u.whitelisted.Store(true)
}

// Unwhitelist removes the user from the whitelist.
func (u *User) Unwhitelist() {
	u.whitelisted.Store(false)
}

// FirstLogin returns the time the user first logged in.
func (u *User) FirstLogin() time.Time {
	return u.firstLogin
}

// JoinTime returns the time the user joined.
func (u *User) JoinTime() time.Time {
	return u.joinTime
}

// PlayTime returns the time the user has played.
func (u *User) PlayTime() time.Duration {
	return u.playTime + time.Since(u.joinTime)
}

// Roles returns the role manager of the user.
func (u *User) Roles() *Roles {
	return u.roles
}

// SetNameTagFromRole sets the name tag from the user's highest role.
func (u *User) SetNameTagFromRole() {
	highest := u.Roles().Highest()
	tag := highest.Tag(u.DisplayName())
	if _, ok := highest.(role.Plus); ok {
		tag = strings.ReplaceAll(tag, "§0", u.Settings().Advanced.VasarPlusColour)
	}
	u.Player().SetNameTag(tag)
}

// MultiplyParticles multiplies the hit particles for the user.
func (u *User) MultiplyParticles(e world.Entity, multiplier int) {
	for i := 0; i < multiplier; i++ {
		u.s.ViewEntityAction(e, entity.CriticalHitAction{})
	}
}

// Rotate rotates the user with the specified yaw and pitch deltas.
// TODO: Remove this once Dragonfly supports a way to do this properly.
func (u *User) Rotate(deltaYaw, deltaPitch float64) {
	currentYaw, currentPitch := u.p.Rotation()
	session_writePacket(u.s, &packet.MovePlayer{
		EntityRuntimeID: 1, // Always 1 on Dragonfly.
		Position:        vec64To32(u.p.Position().Add(mgl64.Vec3{0, 1.62})),
		Pitch:           float32(currentPitch + deltaPitch),
		Yaw:             float32(currentYaw + deltaYaw),
		HeadYaw:         float32(currentYaw + deltaYaw),
		Mode:            packet.MoveModeTeleport,
		OnGround:        u.p.OnGround(),
	})
}

// Launch launches the user in their direction vector.
func (u *User) Launch() {
	now := time.Now()
	if now.Before(u.launchDelay.Load()) {
		return
	}

	u.SendCustomParticle(8, 0, u.p.Position(), true) // Add a flame particle.
	u.SendCustomSound("mob.vex.hurt", 1, 0.5, true)

	motion := entity.DirectionVector(u.p).Mul(1.5)
	motion[1] = 0.85

	u.p.StopSprinting()
	u.p.SetVelocity(motion)

	u.launchDelay.Store(now.Add(time.Second * 2))
}

// SetLastMessageFrom sets the player passed as the last player who messaged the user.
func (u *User) SetLastMessageFrom(p *player.Player) {
	u.lastMessageFrom.Store(p.XUID())
}

// LastMessageFrom returns the last user that messaged the user.
func (u *User) LastMessageFrom() (*User, bool) {
	u, ok := LookupXUID(u.lastMessageFrom.Load())
	return u, ok
}

// Frozen returns the frozen state of the user.
func (u *User) Frozen() bool { return u.frozen.Load() }

// ToggleRod toggles a fishing hook. If the user is already using a hook, it will be removed, otherwise a new hook will
// be created.
func (u *User) ToggleRod() {
	u.rodMu.Lock()
	defer u.rodMu.Unlock()

	if u.hook == nil || u.hook.World() == nil {
		if !u.ProjectilesDisabled() {
			u.hook = ent.NewFishingHook(entity.EyePosition(u.p), entity.DirectionVector(u.p).Mul(1.3), u.p)

			w := u.p.World()
			w.AddEntity(u.hook)
		}
	} else {
		_ = u.hook.Close()
	}
}

// PearlCoolDown returns true if ender pearls currently are on cool down.
func (u *User) PearlCoolDown() bool {
	return u.pearlCoolDown.Load()
}

// TogglePearlCoolDown toggles the ender pearl cool down.
func (u *User) TogglePearlCoolDown() {
	if u.pearlCoolDown.Toggle() {
		u.ResetExperienceProgress()
	}
}

// DisableProjectiles disables the user's projectiles.
func (u *User) DisableProjectiles() {
	u.projectilesDisabled.Store(true)
}

// EnableProjectiles enables the user's projectiles.
func (u *User) EnableProjectiles() {
	u.projectilesDisabled.Store(false)
}

// ProjectilesDisabled returns true if the user's projectiles are disabled.
func (u *User) ProjectilesDisabled() bool {
	return u.projectilesDisabled.Load()
}

// Settings returns the settings of the user.
func (u *User) Settings() Settings {
	return u.settings.Load()
}

// SetSettings sets the settings of the user.
func (u *User) SetSettings(settings Settings) {
	u.settings.Store(settings)
}

// Stats returns the stats of the user.
func (u *User) Stats() Stats {
	return u.stats.Load()
}

// SetStats sets the stats of the user.
func (u *User) SetStats(stats Stats) {
	u.stats.Store(stats)
}

// Player ...
func (u *User) Player() *player.Player {
	return u.p
}

// Device returns the device of the user.
func (u *User) Device() string {
	return u.s.ClientData().DeviceModel
}

// DeviceID returns the device ID of the user.
func (u *User) DeviceID() string {
	return u.s.ClientData().DeviceID
}

// SelfSignedID returns the self-signed ID of the user.
func (u *User) SelfSignedID() string {
	return u.s.ClientData().SelfSignedID
}

// Address returns the address of the user.
func (u *User) Address() net.Addr {
	return u.address
}

// HashedAddress returns the hashed IP address of the user.
func (u *User) HashedAddress() string {
	return u.hashedAddress
}

func (u *User) Cooldowns() []*Cooldown {
	return u.cooldowns
}

func (u *User) Timer() *Timer {
	return u.timer
}

// Faction returns the faction name of the user
func (u *User) Faction() string {
	return u.faction
}

// HasFaction returns whether the user has a faction
func (u *User) HasFaction() bool {
	return u.faction != ""
}

func (u *User) Money() int {
	return u.money
}

func (u *User) AddMoney(m int) {
	u.money += m
}

func (u *User) TakeMoney(m int) {
	u.money -= m
}

func (u *User) SetMoney(m int) {
	u.money = m
}

// SetFaction sets the faction name of the user
func (u *User) SetFaction(f string) {
	u.faction = f
}

// Claim returns the claim name of the user at his/her position
func (u *User) Claim() *c.Claim {
	return u.claim
}

// SetClaim sets the claim name of the user at his/her position
func (u *User) SetClaim(c *c.Claim) {
	u.claim = c
}

// Close ...
func (u *User) Close() {
	u.tagMu.Lock()
	close(u.tagC)
	u.tagMu.Unlock()

	userMu.Lock()
	delete(users, u.p)
	delete(usersXUID, u.p.XUID())
	userMu.Unlock()
}

// viewers returns a list of all viewers of the Player.
func (u *User) viewers() []world.Viewer {
	viewers := u.p.World().Viewers(u.p.Position())
	for _, v := range viewers {
		if v == u.s {
			return viewers
		}
	}
	return append(viewers, u.s)
}

// vec64To32 converts a mgl64.Vec3 to a mgl32.Vec3.
func vec64To32(vec3 mgl64.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{float32(vec3[0]), float32(vec3[1]), float32(vec3[2])}
}

//go:linkname player_session github.com/df-mc/dragonfly/server/player.(*Player).session
//noinspection ALL
func player_session(*player.Player) *session.Session

//go:linkname session_writePacket github.com/df-mc/dragonfly/server/session.(*Session).writePacket
//noinspection ALL
func session_writePacket(*session.Session, packet.Packet)
