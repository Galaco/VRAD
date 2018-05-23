package trace

import (
	"github.com/galaco/vrad/vmath/ssemath"
	"github.com/galaco/vrad/vmath/ssemath/simd"
	"log"
)

func TestLineDoesHitSky(start *ssemath.FourVectors, stop ssemath.FourVectors,
	fractionVisible *simd.Flt4x, canRecurse bool, staticPropToSkip int, doDebug bool) {
		log.Panicln("TestLineDoesHitSky: NOT IMPLEMENTED")
/*	FourRays myrays;
	myrays.origin = start;
	myrays.direction = stop;
	myrays.direction -= myrays.origin;
	fltx4 len = myrays.direction.length();
	myrays.direction *= ReciprocalSIMD( len );
	RayTracingResult rt_result;
	CCoverageCountTexture coverageCallback;

	g_RtEnv.Trace4Rays(myrays, Four_Zeros, len, &rt_result, TRACE_ID_STATICPROP | static_prop_to_skip, g_bTextureShadows? &coverageCallback : 0);

	if ( bDoDebug )
	{
		WriteTrace( "trace.txt", myrays, rt_result );
	}

	float aOcclusion[4];
	for ( int i = 0; i < 4; i++ )
	{
		aOcclusion[i] = 0.0f;
		if ( ( rt_result.HitIds[i] != -1 ) &&
		     ( rt_result.HitDistance.m128_f32[i] < len.m128_f32[i] ) )
		{
			int id = g_RtEnv.OptimizedTriangleList[rt_result.HitIds[i]].m_Data.m_IntersectData.m_nTriangleID;
			if ( !( id & TRACE_ID_SKY ) )
				aOcclusion[i] = 1.0f;
		}
	}
	fltx4 occlusion = LoadUnalignedSIMD( aOcclusion );
	if (g_bTextureShadows)
		occlusion = MaxSIMD ( occlusion, coverageCallback.GetCoverage() );

	bool fullyOccluded = ( TestSignSIMD( CmpGeSIMD( occlusion, Four_Ones ) ) == 0xF );

	// if we hit sky, and we're not in a sky camera's area, try clipping into the 3D sky boxes
	if ( (! fullyOccluded) && canRecurse && (! g_bNoSkyRecurse ) )
	{
		FourVectors dir = stop;
		dir -= start;
		dir.VectorNormalize();

		int leafIndex = -1;
		leafIndex = PointLeafnum( start.Vec( 0 ) );
		if ( leafIndex >= 0 )
		{
			int area = dleafs[leafIndex].area;
			if (area >= 0 && area < numareas)
			{
				if (area_sky_cameras[area] < 0)
				{
					int cam;
					for (cam = 0; cam < num_sky_cameras; ++cam)
					{
						FourVectors skystart, skytrans, skystop;
						skystart.DuplicateVector( sky_cameras[cam].origin );
						skystop = start;
						skystop *= sky_cameras[cam].world_to_sky;
						skystart += skystop;

						skystop = dir;
						skystop *= MAX_TRACE_LENGTH;
						skystop += skystart;
						TestLine_DoesHitSky ( skystart, skystop, pFractionVisible, false, static_prop_to_skip, bDoDebug );
						occlusion = AddSIMD ( occlusion, Four_Ones );
						occlusion = SubSIMD ( occlusion, *pFractionVisible );
					}
				}
			}
		}
	}

	occlusion = MaxSIMD( occlusion, Four_Zeros );
	occlusion = MinSIMD( occlusion, Four_Ones );
	*pFractionVisible = SubSIMD( Four_Ones, occlusion );
*/
}
