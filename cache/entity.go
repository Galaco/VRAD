package cache

import (
	"github.com/galaco/vrad/common/types"
)

var entityList []types.Entity

func SetEntityList(ents []types.Entity) {
	entityList = ents
}

func GetEntity(index int) *types.Entity {
	return &(entityList[index])
}

func GetAllEntities() *[]types.Entity {
	return &entityList
}
/*
  ==================
  FindTargetEntity
  ==================
*/
func FindTargetEntity (target string) *types.Entity {
	for i := 0 ; i < len(entityList) ; i++ {
		n := entityList[i].ValueForKey("targetname")
		if n == target {
			return &entityList[i]
		}
	}

	return nil
}