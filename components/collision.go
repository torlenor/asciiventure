package components

import "github.com/torlenor/asciiventure/entities"

type Collision struct {
	V                  bool
	DestroyOnCollision bool
}

type CollisionManager map[entities.Entity]Collision
