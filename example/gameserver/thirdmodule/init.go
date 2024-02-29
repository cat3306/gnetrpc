package thirdmodule

import (
	"time"

	"github.com/cat3306/gnetrpc/util"
)

func Init() {
	util.PanicRepeatRun(InitDb, util.PanicRepeatRunArgs{
		Sleep: time.Second,
		Try:   3,
	})
	util.PanicRepeatRun(InitCache, util.PanicRepeatRunArgs{
		Sleep: time.Second,
		Try:   3,
	})
	util.PanicRepeatRun(InitLocalCache, util.PanicRepeatRunArgs{
		Sleep: time.Second,
		Try:   3,
	})
}
