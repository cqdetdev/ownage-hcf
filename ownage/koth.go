package ownage

import (
	"time"

	"github.com/ownagepe/hcf/ownage/module"
)

func (v *Ownage) startKOTH() {
	t := time.NewTicker(time.Hour)
	defer t.Stop()
	for {
		select {
		case <-v.c:
			return
		case <-t.C:
			koth := module.Rotate()
			koth.Start()
		}
	}
}