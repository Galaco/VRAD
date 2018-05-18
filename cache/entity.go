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