package cache

const KDNODE_STATE_XSPLIT = 0								// this node is an x split
const KDNODE_STATE_YSPLIT = 1								// this node is a ysplit
const KDNODE_STATE_ZSPLIT = 2								// this node is a zsplit
const KDNODE_STATE_LEAF = 3									// this node is a leaf

// this is the cache intensive data structure. "Tricks" are used to fit it into 8 bytes:
//
// A) the right child is always stored after the left child, which means we only need one
// pointer
// B) The type of node (KDNODE_xx) is stored in the lower 2 bits of the pointer.
// C) for leaf nodes, we store the number of triangles in the leaf in the same place as the floating
//    point splitting parameter is stored in a non-leaf node
type OptimisedKDNode struct {
	Children int
	SplittingPlaneValue float32
}

func (node *OptimisedKDNode) NodeType() int{
	return node.Children & 3
}

func (node *OptimisedKDNode) TriangleIndexStart() int{
	//assert node.NodeType==KDNODE_STATE_LEAF
	return node.Children >> 2
}

func (node *OptimisedKDNode) LeftChild() int{
	//assert node.NoteType!==KDNODE_STATE_LEAF
	return node.Children >> 2
}

func (node *OptimisedKDNode) RightChild() int{
	return node.LeftChild() + 1
}


func (node *OptimisedKDNode) NumberOfTrianglesInLeaf() int{
	//assert node.NodeType==KDNODE_STATE_LEAF
	// @TODO verify this is correct
	return int(node.SplittingPlaneValue)
}


func (node *OptimisedKDNode) SetNumberOfTrianglesInLeaf(n int) {
	//assert node.NodeType==KDNODE_STATE_LEAF
	 node.SplittingPlaneValue = float32(n)
}