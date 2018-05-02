package common

import (
	"github.com/galaco/bsp"
	"log"
)

// Is this awful?
// YES.
// Why is it here?
// Because separating out the swathe of global uses of the bsp is difficult.
// At least this way we can minimise and control the globals

var global_BSP *bsp.Bsp

func GLOBALSET_BSP(file *bsp.Bsp) {
	if global_BSP == nil {
		global_BSP = file
	} else {
		log.Fatal("Must not re-set global after")
	}
}

func GLOBALGET_BSP() *bsp.Bsp {
	return global_BSP
}