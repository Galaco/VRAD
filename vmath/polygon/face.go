package polygon

import "github.com/galaco/bsp/primitives/face"

func ValidDispFace(f *face.Face) bool {
	if f == nil {
		return false
	}
	if f.DispInfo == -1 {
		return false
	}
	if f.NumEdges != 4 {
		return false
	}
	return true
}
