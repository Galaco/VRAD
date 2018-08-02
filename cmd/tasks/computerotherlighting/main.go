package computerotherlighting

import (
	"github.com/galaco/vrad/cmd"
	"github.com/galaco/vrad/compute/leafambient"
)

func Main(args *cmd.Args, transfered interface{}) (interface{}, error) {
	// Compute lighting for the bsp file
	if args.NoDetailLighting == false {
		//ComputeDetailPropLighting( THREADINDEX_MAIN )
	}

	leafambient.ComputePerLeafAmbientLighting()

	// bake the static props high quality vertex lighting into the bsp
	if args.Fast == false && args.StaticPropLighting == true {
		//StaticPropMgr()->ComputeLighting( THREADINDEX_MAIN );
	}

	return nil,nil
}