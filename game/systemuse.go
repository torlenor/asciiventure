package game

import (
	"github.com/torlenor/asciiventure/components"
	"github.com/torlenor/asciiventure/entity"
)

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
			e.Actor = nil
		}
	}
}
