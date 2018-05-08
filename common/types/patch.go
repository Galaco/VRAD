package types


type Patch struct {

}

/**

	winding_t	*winding;
	Vector		mins, maxs, face_mins, face_maxs;

	Vector		origin;				// adjusted off face by face normal

	dplane_t	*plane;				// plane (corrected for facing)

	unsigned short		m_IterationKey;	// Used to prevent touching the same patch multiple times in the same query.
										// See IncrementPatchIterationKey().

	// these are packed into one dword
	unsigned int normalMajorAxis : 2;	// the major axis of base face normal
	unsigned int sky : 1;
	unsigned int needsBumpmap : 1;
	unsigned int pad : 28;

	Vector		normal;				// adjusted for phong shading

	float		planeDist;			// Fixes up patch planes for brush models with an origin brush

	float		chop;				// smallest acceptable width of patch face
	float		luxscale;			// average luxels per world coord
	float		scale[2];			// Scaling of texture in s & t

	bumplights_t totallight;		// accumulated by radiosity
									// does NOT include light
									// accounted for by direct lighting
	Vector		baselight;			// emissivity only
	float		basearea;			// surface per area per baselight instance

	Vector		directlight;		// direct light value
	float		area;

	Vector		reflectivity;		// Average RGB of texture, modified by material type.

	Vector		samplelight;
	float		samplearea;		// for averaging direct light
	int			faceNumber;
	int			clusterNumber;

	int			parent;			// patch index of parent
	int			child1;			// patch index for children
	int			child2;

	int			ndxNext;					// next patch index in face
	int			ndxNextParent;				// next parent patch index in face
	int			ndxNextClusterChild;		// next terminal child index in cluster
//	struct		patch_s		*next;					// next in face
//	struct		patch_s		*nextparent;		    // next in face
//	struct		patch_s		*nextclusterchild;		// next terminal child in cluster

	int			numtransfers;
	transfer_t	*transfers;

	short		indices[3];				// displacement use these for subdivision
 */