package user

import (
	"github.com/df-mc/dragonfly/server/session"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// PlayerAuthInputHandler ...
type PlayerAuthInputHandler struct {
	u *User
}

// Handle ...
func (h PlayerAuthInputHandler) Handle(p packet.Packet, s *session.Session) error {
	return (session.PlayerAuthInputHandler{}).Handle(p, s)
}
