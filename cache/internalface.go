package cache

import "github.com/galaco/bsp/primitives/face"

var targetFaces []face.Face

func GetTargetFaces() *[]face.Face {
	return &targetFaces
}

func SetTargetFaces(faces *[]face.Face) {
	targetFaces = *faces
}