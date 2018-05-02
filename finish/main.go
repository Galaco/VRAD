package finish

import (
	"github.com/galaco/vrad/cmd"
	"log"
)

func Main(args *cmd.Args, transfered interface{}) (interface{}, error) {
	log.Println("Ready to Finish")

	if args.Verbose == true {
		//PrintBSPFileSizes()
	}

	log.Printf("Writing %s\n", args.Filename)
	//writer := bsp.Writer{}
	// Maybe want to compile the bsp here too?
	// writer.SetBsp(transfered.(bsp.Bsp))

	if args.DumpPatches == true {
		for iStyle := 0; iStyle < 4; iStyle++ {
			for iBump := 0; iBump < 4; iBump++ {
				//g_pFileSystem->Close( pFileSamples[iStyle][iBump] );
			}
		}
	}

	//CloseDispLuxels();

	//StaticPropMgr()->Shutdown();

	//double end = Plat_FloatTime();

	//char str[512];
	//GetHourMinuteSecondsString( (int)( end - g_flStartTime ), str, sizeof( str ) );
	//Msg( "%s elapsed\n", str );

	//ReleasePakFileLumps();

	return nil,nil
}
