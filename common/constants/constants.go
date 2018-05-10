package constants

const MAX_COORD_INTEGER = 16384
const MAX_MAP_FACES = 65536
const MAX_MAP_CLUSTERS = 65536
const MAX_MAP_TEXINFO = 12288
const MAX_MAP_LEAFS = 65536
const MAX_MAP_NODES = 65536

// We can have larger lightmaps on displacements
const MAX_DISP_LIGHTMAP_DIM_WITHOUT_BORDER	= 125
const MAX_DISP_LIGHTMAP_DIM_INCLUDING_BORDER =128


// This is the actual max.. (change if you change the brush lightmap dim or disp lightmap dim
const MAX_LIGHTMAP_DIM_WITHOUT_BORDER	= MAX_DISP_LIGHTMAP_DIM_WITHOUT_BORDER
const MAX_LIGHTMAP_DIM_INCLUDING_BORDER	= MAX_DISP_LIGHTMAP_DIM_INCLUDING_BORDER