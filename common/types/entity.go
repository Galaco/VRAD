package types

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/bsp/primitives/portal"
)

type Entity struct {
	Origin mgl32.Vec3
	FirstBrush int
	NumBrushes int
	EPairs *EPair

	// only valid for func_areaportals
	AreaPortalNum int
	PortalAreas [2]int
	PortalsLeadingIntoAreas [2]*portal.Portal	// portals leading into portalareas
}

func (ent *Entity) ValueForKey(key string) string {
	e := ent.EPairs
	for e != nil {
		if e.Key == key {
			return e.Value
		}

		e = e.Next
	}

	return ""
}



type EPair struct {
	Next *EPair
	Key string
	Value string
}