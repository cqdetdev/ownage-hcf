package ownage

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
)

type InventoryHandler struct{}

func (InventoryHandler) HandleTake(ctx *event.Context, slot int, _ item.Stack)  {}
func (InventoryHandler) HandlePlace(ctx *event.Context, slot int, _ item.Stack) {}
func (InventoryHandler) HandleDrop(ctx *event.Context, slot int, _ item.Stack)  {}