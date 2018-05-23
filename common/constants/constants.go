package constants


// Following values should be +16384, -16384, +15/16, -15/16
// NOTE THAT IF THIS GOES ANY BIGGER THEN DISK NODES/LEAVES CANNOT USE SHORTS TO STORE THE BOUNDS
const MAX_COORD_INTEGER = 16384
const MIN_COORD_INTEGER	= -MAX_COORD_INTEGER
const MAX_COORD_FRACTION = 1.0-(1.0/16.0)
const MIN_COORD_FRACTION = -1.0+(1.0/16.0)

const MAX_COORD_FLOAT = 16384.0
const MIN_COORD_FLOAT = -MAX_COORD_FLOAT

// Width of the coord system, which is TOO BIG to send as a client/server coordinate value
const COORD_EXTENT = 2 * MAX_COORD_INTEGER

// Maximum traceable distance ( assumes cubic world and trace from one corner to opposite )
// COORD_EXTENT * sqrt(3)
const MAX_TRACE_LENGTH = 1.732050807569 * COORD_EXTENT

// We can have larger lightmaps on displacements
const MAX_DISP_LIGHTMAP_DIM_WITHOUT_BORDER	= 125
const MAX_DISP_LIGHTMAP_DIM_INCLUDING_BORDER =128


// This is the actual max.. (change if you change the brush lightmap dim or disp lightmap dim
const MAX_LIGHTMAP_DIM_WITHOUT_BORDER	= MAX_DISP_LIGHTMAP_DIM_WITHOUT_BORDER
const MAX_LIGHTMAP_DIM_INCLUDING_BORDER	= MAX_DISP_LIGHTMAP_DIM_INCLUDING_BORDER


const CONSTRUCTS_INVALID_INDEX = -1
const MAX_POINTS_ON_WINDING = 64
const NUM_BUMP_VECTS = 3

const ANGLE_UP = -1
const ANGLE_DOWN = -2

const PITCH = 0
const YAW = 1
const ROLL = 2


const TEMPCONST_NUM_THREADS = 1