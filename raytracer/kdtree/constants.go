package kdtree

// Both the "quick" and regular kd tree building algorithms here use the "surface area heuristic":
// the relative probability of hitting the "left" subvolume (Vl) from a split is equal to that
// subvolume's surface area divided by its parent's surface area (Vp) : P(Vl | V)=SA(Vl)/SA(Vp).
// The same holds for the right subvolume, Vp. Nl is the number of triangles in the left volume,
// and Nr in the right volume. if Ct is the cost of traversing one tree node, and Ci is the cost of
// intersection with the primitive, than the cost of splitting is estimated as:
//
//    Ct+Ci*((SA(Vl)/SA(V))*Nl+(SA(Vr)/SA(V)*Nr)).
// and the cost of not splitting is
//    Ci*N
//
//  This both provides a metric to minimize when computing how and where to split, and also a
//  termination criterion.
//
// the "quick" method just splits down the middle, while the slow method splits at the best
// discontinuity of the cost formula. The quick method splits along the longest axis ; the
// regular algorithm tries all 3 to find which one results in the minimum cost
//
// both methods use the additional optimization of "growing" empty nodes - if the split results in
// one side being devoid of triangles, the empty side is "grown" as much as possible.
//

const COST_OF_TRAVERSAL = 75								// approximate #operations
const COST_OF_INTERSECTION = 167							// approximate #operations

const MAX_TREE_DEPTH = 21
