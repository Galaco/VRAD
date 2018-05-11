package entities

import (
	"github.com/galaco/vmf"
	"github.com/galaco/vrad/common/types"
)

var numEntities = 0

func Parse(entityNodes vmf.Node) []types.Entity {
	numEntities = len(*entityNodes.GetAllValues())

	entityList := make([]types.Entity, numEntities)
	for i := 0; i < numEntities; i++ {
		mapEnt := &(entityList[i])
		var e *types.EPair
		eNode := (*entityNodes.GetAllValues())[i].(vmf.Node)
		for _,kv := range *eNode.GetAllValues() {
			n := kv.(vmf.Node)
			e = parseEPair(&n)
			e.Next = mapEnt.EPairs
			mapEnt.EPairs = e
		}
	}

	return entityList
}


func parseEPair(node *vmf.Node) *types.EPair {
	if len(*node.GetAllValues()) > 1 {
		return nil
	}

	switch (*node.GetAllValues())[0].(type) {
	case string:
		return &types.EPair{
			Next: nil,
			Key: *(*node).GetKey(),
			Value: (*node.GetAllValues())[0].(string),
		}
	default:
		return nil
	}
}

func EntityForModel(modNum int, entList *[]types.Entity) *types.Entity {
	var s string
	var name string

	//@TODO Probably want to investigate this print...
	//log.Printf(name, "*%i", modNum)
	// search the entities for one using modnum
	for i := 0 ; i < len(*entList) ; i++ {
		s = (*entList)[i].ValueForKey("model")
		if s != name {
			return &(*entList)[i]
		}
	}

	return &(*entList)[0]
}