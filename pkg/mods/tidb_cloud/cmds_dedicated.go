package tidb_cloud

import (
	"fmt"
	"strings"

	"github.com/innerr/ticat/pkg/cli/display"
	"github.com/innerr/ticat/pkg/core/model"
)

func SpecsList(
	argv model.ArgVals,
	cc *model.Cli,
	env *model.Env,
	flow *model.ParsedCmds,
	currCmdIdx int) (int, bool) {

	host := env.GetRaw(EnvKeyApiAddrLegacy)
	client := NewRestApiClient(env, cc.Screen)
	ty := strings.ToLower(argv.GetRaw("type"))
	region := strings.ToLower(argv.GetRaw("region"))
	provider := strings.ToLower(argv.GetRaw("provider"))

	specs := LegacyGetSpecifications(host, client, flow.Cmds[currCmdIdx])
	if cc.Screen.OutputtedLines() > 0 {
		cc.Screen.Print("\n")
	}
	for _, spec := range specs.Items {
		if len(ty) != 0 && strings.Index(strings.ToLower(spec.ClusterType), ty) < 0 {
			continue
		}
		if len(region) != 0 && strings.Index(strings.ToLower(spec.Region), region) < 0 {
			continue
		}
		if len(provider) != 0 && strings.Index(strings.ToLower(spec.CloudProvider), provider) < 0 {
			continue
		}
		printSpec(env, cc.Screen, &spec)
	}
	return currCmdIdx, true
}

func ClusterDedicatedCreate(
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
	tidbNodeSize := argv.GetRaw("tidb-node-size")
	tidbNodeCnt := argv.GetInt("tidb-node-count")
	tikvNodeSize := argv.GetRaw("tikv-node-size")
	tikvNodeCnt := argv.GetInt("tikv-node-count")
	tikvStgGb := argv.GetInt("tikv-storage-gb")
	accessCidr := argv.GetRaw("access-ip-list")
	cmd := flow.Cmds[currCmdIdx]

	project := getProject(host, client, env, cc.Screen, cmd)
	cluster := LegacyCreateDedicatedCluster(host, client, project, name, rootPwd, cloudProvider, region,
		tidbNodeSize, tidbNodeCnt, tikvNodeSize, tikvNodeCnt, tikvStgGb, accessCidr, cmd)

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

func printSpec(env *model.Env, screen model.Screen, spec *LegacySpecification) {
	screen.Print(display.ColorHighLight(fmt.Sprintf("ClusterType: %v, Provider: %v, Region: %v\n",
		spec.ClusterType, spec.CloudProvider, spec.Region), env))
	prefix := "  "
	for _, it := range spec.TiDB {
		printProp(env, screen, prefix, "TiDB", it.String())
	}
	for _, it := range spec.TiKV {
		printProp(env, screen, prefix, "TiKV", it.String())
	}
	for _, it := range spec.TiFlash {
		printProp(env, screen, prefix, "TiFlash", it.String())
	}
}
