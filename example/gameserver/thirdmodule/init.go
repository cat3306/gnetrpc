package thirdmodule

import (
	"github.com/cat3306/gnetrpc/util"
	"time"
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
}
