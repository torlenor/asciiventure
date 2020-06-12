package game

import (
	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
)

func (g *Game) pickupSystem() {
	for _, e := range g.entities {
		if e.Actor != nil && e.Actor.NextAction == components.ActionTypeInteract {
			for _, target := range g.entities {
				if target != nil && (target.Item != nil || target.Mutation != nil) && target.Position != nil && target.Position.Equal(*e.Position) {
					if target.Item != nil {
						result := e.PickUpItem(target)
						for _, r := range result {
							switch r.Type {
							case entity.ActionResultItemPickedUp:
							case entity.ActionResultMessage:
								g.ui.AddLogEntry(r.StringValue)
							}
						}
					}
					if target.Mutation != nil {
						result := e.ConsumeMutation(target)
						for _, r := range result {
							switch r.Type {
							case entity.ActionResultMutationConsumed:
							case entity.ActionResultMessage:
								g.ui.AddLogEntry(r.StringValue)
							}
						}
					}
				}
			}
		}
	}
}

func (g *Game) useSystem() {
	for _, e := range g.entities {
		if e.Actor != nil && e.Actor.NextAction == components.ActionTypeUseItem {
			if item := e.Inventory.PopOneByID(e.Actor.IntValue); item != nil {
				result := e.UseItem(item)
				for _, r := range result {
					switch r.Type {
					case entity.ActionResultItemUsed:
					case entity.ActionResultMessage:
						g.ui.AddLogEntry(r.StringValue)
					}
				}
			}
		}
	}
}
