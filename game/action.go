package game

import (
	"log"

	"github.com/torlenor/asciiventure/entity"
)

// actionType holds the type of the action to trigger.
type actionType int

// List of actionTypes.
const (
	actionTypeUnknown actionType = iota
	// actionTypeInteract can either be picking something up, consuming a mutation from the floor or going through a portal
	actionTypeInteract
	actionTypeDropItem
)

func (d actionType) String() string {
	return [...]string{"Unknown", "Interact", "Drop"}[d]
}

func (g *Game) cleanupEntities() {
	n := 0
	for _, e := range g.entities {
		if e != nil {
			g.entities[n] = e
			n++
		}
	}
	g.entities = g.entities[:n]
}

func (g *Game) performAction(at actionType) {
	switch at {
	case actionTypeInteract:
		for _, e := range g.entities {
			if e != nil && e.Position != nil && e.Position.Equal(*g.player.Position) && e != g.player {
				result := g.player.PickUpItem(e)
				for _, r := range result {
					switch r.Type {
					case entity.ActionResultItemPickedUp:
						g.nextStep = true
					case entity.ActionResultMessage:
						g.ui.AddLogEntry(r.StringValue)
					}
				}

				result = g.player.ConsumeMutation(e)
				for _, r := range result {
					switch r.Type {
					case entity.ActionResultMutationConsumed:
						g.nextStep = true
					case entity.ActionResultMessage:
						g.ui.AddLogEntry(r.StringValue)
					}
				}
			}
		}

		if g.currentGameMap.IsPortal(*g.player.Position) {
			g.selectGameMap(g.currentGamMapID + 1)
		}
	case actionTypeDropItem:
		if len(g.player.Inventory) > 0 {
			result := g.player.DropItem(g.player.Inventory[len(g.player.Inventory)-1])
			for _, r := range result {
				switch r.Type {
				case entity.ActionResultItemDropped:
					g.nextStep = true
				case entity.ActionResultMessage:
					g.ui.AddLogEntry(r.StringValue)
				}
			}
		}
	case actionTypeUnknown:
		log.Printf("Unknown action type provided")
	}
}
