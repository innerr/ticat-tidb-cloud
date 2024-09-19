package tidb_cloud

import (
	"fmt"

	"github.com/innerr/ticat/pkg/cli/display"
	"github.com/innerr/ticat/pkg/core/model"
)

func V1Beta1ClusterServerlessCreate(
	argv model.ArgVals,
	cc *model.Cli,
	env *model.Env,
	flow *model.ParsedCmds,
	currCmdIdx int) (int, bool) {

	host := env.GetRaw(EnvKeyApiAddrV1Beta1)
	client := NewRestApiClient(env, cc.Screen)
	name := env.GetRaw(EnvKeyClusterName)
	rootPwd := env.GetRaw(EnvKeyRootPwd)
	cloudProvider := env.GetRaw(EnvKeyCloudProvider)
	region := env.GetRaw(EnvKeyCloudRegion)
	cmd := flow.Cmds[currCmdIdx]

	project := getProject(host, client, env, cc.Screen, cmd)
	cluster := V1Beta1CreateDevCluster(host, client, project, name, rootPwd, cloudProvider, region, cmd)

	if cc.Screen.OutputtedLines() > 0 {
		cc.Screen.Print("\n")
	}
	v1beta1PrintCluster(env, cc.Screen, cluster)
	//V1Beta1ChangeRootPassword(host, client, cluster.ClusterId, rootPwd, cmd)
	env.GetLayer(model.EnvLayerSession).Set(EnvKeyClusterId, cluster.ClusterId)
	env.GetLayer(model.EnvLayerSession).Set("mysql.pwd", rootPwd)
	return currCmdIdx, true
}

func LegecyClusterServerlessCreate(
	argv model.ArgVals,
	cc *model.Cli,
	env *model.Env,
	flow *model.ParsedCmds,
	currCmdIdx int) (int, bool) {

	host := env.GetRaw(EnvKeyApiAddrLegacy)
	client := NewRestApiClient(env, cc.Screen)
	name := env.GetRaw(EnvKeyClusterName)
	rootPwd := env.GetRaw(EnvKeyRootPwd)
	cloudProvider := env.GetRaw(EnvKeyCloudProvider)
	region := env.GetRaw(EnvKeyCloudRegion)
	cmd := flow.Cmds[currCmdIdx]

	project := getProject(host, client, env, cc.Screen, cmd)
	cluster := LegacyCreateDevCluster(host, client, project, name, rootPwd, cloudProvider, region, cmd)

	if cc.Screen.OutputtedLines() > 0 {
		cc.Screen.Print("\n")
	}
	sep := display.ColorProp(":", env)
	cc.Screen.Print(fmt.Sprintf("%s%s %v\n", display.ColorArg("Id", env), sep, cluster.ClusterId))
	if len(cluster.Message) != 0 {
		cc.Screen.Print(display.ColorExplain(cluster.Message, env))
	}
	env.GetLayer(model.EnvLayerSession).SetUint64(EnvKeyClusterId, cluster.ClusterId)
	env.GetLayer(model.EnvLayerSession).Set("mysql.pwd", rootPwd)
	return currCmdIdx, true
}
