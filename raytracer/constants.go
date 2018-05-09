package raytracer



const RTE_FLAGS_FAST_TREE_GENERATION = 1
const RTE_FLAGS_DONT_STORE_TRIANGLE_COLORS = 2				// saves memory if not needed
const RTE_FLAGS_DONT_STORE_TRIANGLE_MATERIALS = 4

const TRACE_ID_SKY        = 0x01000000  // sky face ray blocker
const TRACE_ID_OPAQUE     = 0x02000000  // everyday light blocking face
const TRACE_ID_STATICPROP = 0x04000000  // static prop - lower bits are prop ID

const KDNODE_STATE_XSPLIT = 0								// this node is an x split
const KDNODE_STATE_YSPLIT = 1								// this node is a ysplit
const KDNODE_STATE_ZSPLIT = 2								// this node is a zsplit
const KDNODE_STATE_LEAF   = 3								// this node is a leaf

const NEVER_SPLIT = 0