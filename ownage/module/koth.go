package module

import (
	"math/rand"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/vasar-network/vails/lang"
)

var koths []*KOTH
var current *KOTH

func KOTHS() []*KOTH {
	return koths
}

func Current() *KOTH {
	return current
}

func Rotate() *KOTH {
	current = random()
	return current
}

func register(k *KOTH) {
	koths = append(koths, k)
}

func random() *KOTH {
	return koths[rand.Intn(len(koths))]
}


func init() {
	register(&KOTH{
		name: "Kingdom",
		captureArea: cube.Box(float64(-30), -256, float64(-30), float64(-20), 256, float64(-20)),
		duration: time.Minute * 5,
	})

	current = random()
	current.Start()
}

type kothSource interface{}

type kothSourceCapture struct {
	winner *player.Player
}

type KOTH struct {
	name            string
	captureArea     cube.BBox
	duration        time.Duration
	h				*Koth				
	capturing       *player.Player
	shouldCaptureAt time.Time
	started         bool
}

func NewKOTH(name string, captureArea cube.BBox, duration time.Duration) *KOTH {
	return &KOTH{
		name:        name,
		captureArea: captureArea,
		duration:    duration,
	}
}

func (k *KOTH) Name() string            { return k.name }
func (k *KOTH) Started() bool           { return k.started }
func (k *KOTH) CaptureArea() cube.BBox  { return k.captureArea }
func (k *KOTH) TimeLeft() time.Duration {
	if _, ok := k.Capturing(); !ok {
		return k.duration
	}
	return time.Until(k.shouldCaptureAt)
}
func (k *KOTH) Duration() time.Duration { return k.duration }
func (k *KOTH) handler() *Koth        { return k.h }
func (k *KOTH) Capturing() (*player.Player, bool) {
	return k.capturing, k.capturing != nil
}

func (k *KOTH) Start() {
	if !k.started {
		ctx := event.C()
		k.handler().HandleStart(ctx)
		if !ctx.Cancelled() {
			k.started = true
		}
	}
}
func (k *KOTH) Stop() {
	if k.started {
		ctx := event.C()
		k.handler().HandleStop(ctx, k)
		if !ctx.Cancelled() {
			k.started = false
		}
	}

}
func (k *KOTH) StartCapturing(p *player.Player) {
	if k.started {
		if k.capturing != p {
			ctx := event.C()
			k.handler().HandleStartCapturing(ctx, p)
			if !ctx.Cancelled() {
				k.capturing = p
				k.shouldCaptureAt = time.Now().Add(k.duration)
				time.AfterFunc(k.duration, k.captureFunc(p))
			}
		}
	}
}
func (k *KOTH) StopCapturing(p *player.Player) {
	if k.started {
		if k.capturing == p {
			ctx := event.C()
			k.handler().HandleStopCapturing(ctx, p)
			if !ctx.Cancelled() {
				k.capturing = nil
				k.shouldCaptureAt = time.Now().Add(43830 * time.Minute)
			}
		}
	}
}
func (k *KOTH) captureFunc(p *player.Player) func() {
	return func() {
		if k.capturing != nil && k.capturing == p {
			if k.shouldCaptureAt.Before(time.Now()) || k.shouldCaptureAt.Equal(time.Now()) {
				ctx := event.C()
				k.h.HandleCapture(ctx, p)
				if !ctx.Cancelled() {
					k.Stop()
				}
			}
		}
	}
}

type Koth struct {
	player.NopHandler

	u *user.User
}

// NewModeration creates a new Koth module.
func NewKoth(u *user.User) *Koth {
	return &Koth{u: u}
}

func (k *Koth) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	if current.started {
		if current.captureArea.Vec3WithinXZ(newPos) {
			current.StartCapturing(k.u.Player())
		} else {
			current.StopCapturing(k.u.Player())
		}
	}
}


func (*Koth) HandleStartCapturing(ctx *event.Context, p *player.Player) {
	for _, u := range user.All() {
		u.Player().Message(lang.Translatef(u.Player().Locale(), "koth.control", u.Player().Name(), current.name))
	}
}
func (*Koth) HandleStopCapturing(ctx *event.Context, p *player.Player)  {
	for _, u := range user.All() {
		u.Player().Message(lang.Translatef(u.Player().Locale(), "koth.lostcontrol", u.Player().Name(), current.name))
	}
}
func (*Koth) HandleCapture(ctx *event.Context, p *player.Player)        {
	for _, u := range user.All() {
		u.Player().Message(lang.Translatef(u.Player().Locale(), "koth.captured", u.Player().Name(), current.name))
	}
}
func (k *Koth) HandleStart(ctx *event.Context) {
	for _, u := range user.All() {
		u.Player().Message(lang.Translatef(u.Player().Locale(), "koth.started", current.name))
	}
}
func (k *Koth) HandleStop(ctx *event.Context, koth *KOTH) {
	for _, u := range user.All() {
		u.Player().Message(lang.Translatef(u.Player().Locale(), "koth.ended", koth.name))
	}
}


func actuallyMoved(old, new mgl64.Vec3) bool {
	return !mgl64.FloatEqual(old.X(), new.X()) && !mgl64.FloatEqual(old.Z(), new.Z())
}