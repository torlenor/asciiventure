package game

import (
	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
)

func (g *Game) pickupSystem() {
	// TODO: Maybe split Item and Mutagen pickup/drop system into two systems
	for _, e := range g.entities {
		if e.Actor != nil && e.Actor.NextAction == components.ActionTypeInteract {
			for _, target := range g.entities {
				if target != nil && (target.Item != nil || target.Mutagen != nil) && target.Position != nil && target.Position.Current.Equal(e.Position.Current) {
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
					if target.Mutagen != nil {
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
			e.Actor = nil
		}
	}
}
