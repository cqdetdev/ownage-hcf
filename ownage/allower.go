package ownage

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/ownagepe/hcf/ownage/data"
	"github.com/vasar-network/vails/lang"
	"github.com/vasar-network/vails/role"
	"golang.org/x/text/language"
	"net"
	"strings"
)

// allower ensures that all players who join are whitelisted if whitelisting is enabled.
type allower struct {
	v *Ownage
}

// Allow ...
func (a *allower) Allow(_ net.Addr, identity login.IdentityData, client login.ClientData) (string, bool) {
	if a.v.config.Ownage.Whitelisted {
		locale, _ := language.Parse(strings.Replace(client.LanguageCode, "_", "-", 1))
		u, err := data.LoadOfflineUser(identity.DisplayName)
		if err != nil {
			return lang.Translatef(locale, "user.server.whitelist"), false
		}
		return lang.Translatef(locale, "user.server.whitelist"), u.Whitelisted || u.Roles.Contains(role.Trial{}, role.Operator{})
	}
	return "", true
}
