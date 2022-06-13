package user

import "time"

// Timer represents a users PVP timer
type Timer struct {
	// Has returns whether the user is on PVP Timer
	Has bool
	// Expires returns the expiration time of the PVP Timer 
	Expires time.Time
}

func DefaultTimer() *Timer {
	return &Timer{
		Has: true,
		Expires: time.Now().Add(time.Hour),
	}
}

func (t *Timer) TimeLeft() time.Duration {
	return time.Until(t.Expires)
}

func (t *Timer) Expired() bool {
	return time.Now().Before(t.Expires)
}

func (u *User) ExpireTimer() {
	u.timer.Has = false
	u.timer.Expires = time.Now()
}

func (u *User) HasTimer() bool {
	end := time.Unix(u.timer.Expires.Unix(), 0)
	return time.Now().Before(end) || u.timer.Has
}