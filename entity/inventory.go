package entity

import "fmt"

// Inventory holds the inventory of an entity.
type Inventory struct {
	MaxSlots int
	Items    []*Entity
}

// Add an item to the inventory.
// Returns an error if it is not possible.
func (i *Inventory) Add(item *Entity) error {
	if len(i.Items) >= i.MaxSlots {
		return fmt.Errorf("Inventory full")
	}
	i.Items = append(i.Items, item)
	return nil
}

// RemoveByID removes an item by id.
func (i *Inventory) RemoveByID(id int) {
	if len(i.Items) > id {
		i.Items[id] = nil

		n := 0
		for _, e := range i.Items {
			if e != nil {
				i.Items[n] = e
				n++
			}
		}
		i.Items = i.Items[:n]
	}
}

// PopOneByName returns one item by name and removes it from inventory.
// Returns nil if it is not found.
func (i *Inventory) PopOneByName(name string) *Entity {
	for j, item := range i.Items {
		if item.Name == name {
			i.RemoveByID(j)
			return item
		}
	}
	return nil
}

// PopOneByID returns one item by id and removes it from inventory.
// Returns nil if it is not found.
func (i *Inventory) PopOneByID(id int) *Entity {
	if len(i.Items) > id {
		item := i.Items[id]
		i.RemoveByID(id)
		return item
	}
	return nil
}

// GetInventoryList returns a list of all available inventory entries
func (i *Inventory) GetInventoryList() (list []string) {
	for _, item := range i.Items {
		list = append(list, item.Name)
	}
	return
}
