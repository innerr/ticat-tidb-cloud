package tidb_cloud

import (
	"github.com/innerr/ticat/pkg/core/model"
)

func IntegrationReplaceCmds(
	argv model.ArgVals,
	cc *model.Cli,
	env *model.Env,
	flow *model.ParsedCmds,
	currCmdIdx int) (int, bool) {

	old := cc.Cmds.GetSub("bench", "result")
	if old == nil {
		panic(model.NewCmdError(flow.Cmds[currCmdIdx],
			"failed to replace cmd `bench.result`: not existed"))
	}
	old.ReplaceCmdWithPowerCmd(BenchResult)
	return currCmdIdx, true
}

func BenchResult(
	argv model.ArgVals,
	cc *model.Cli,
	env *model.Env,
	flow *model.ParsedCmds,
	currCmdIdx int) (int, bool) {

	cc.Screen.Print("TODO: bench result\n")
	return currCmdIdx, true
}
