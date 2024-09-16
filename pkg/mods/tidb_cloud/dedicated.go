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

	host := env.GetRaw(EnvKeyApiAddr)
	client := NewRestApiClient(env, cc.Screen)
	ty := strings.ToLower(argv.GetRaw("type"))
	region := strings.ToLower(argv.GetRaw("region"))
	provider := strings.ToLower(argv.GetRaw("provider"))

	specs := getSpecifications(host, client, flow.Cmds[currCmdIdx])
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

	host := env.GetRaw(EnvKeyApiAddr)
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
	cidr := argv.GetRaw("cidr")
	cmd := flow.Cmds[currCmdIdx]

	project := getProject(host, client, env, cc.Screen, cmd)
	cluster := createDedicatedCluster(host, client, project, name, rootPwd, cloudProvider, region,
		tidbNodeSize, tidbNodeCnt, tikvNodeSize, tikvNodeCnt, tikvStgGb, cidr, cmd)

	if cc.Screen.OutputtedLines() > 0 {
		cc.Screen.Print("\n")
	}
	sep := display.ColorProp(":", env)
	cc.Screen.Print(fmt.Sprintf("%s%s %v\n", display.ColorArg("ID", env), sep, cluster.ClusterID))
	if len(cluster.Message) != 0 {
		cc.Screen.Print(display.ColorExplain(cluster.Message, env))
	}
	env.GetLayer(model.EnvLayerSession).SetUint64(EnvKeyClusterId, cluster.ClusterID)
	env.GetLayer(model.EnvLayerSession).Set("mysql.pwd", rootPwd)
	return currCmdIdx, true
}

func createDedicatedCluster(
	host string,
	client *RestApiClient,
	projectID uint64,
	name string,
	rootPwd string,
	cloudProvider string,
	cloudRegion string,
	tidbNodeSize string,
	tidbNodeCnt int,
	tikvNodeSize string,
	tikvNodeCnt int,
	tikvStgGb int,
	cidr string,
	cmd model.ParsedCmd) *CreateClusterResp {

	payload := CreateClusterReq{
		Name:          name,
		ClusterType:   "DEDICATED",
		CloudProvider: cloudProvider,
		Region:        cloudRegion,
		Config: ClusterConfig{
			RootPassword: rootPwd,
			Port:         4000,
			Components: Components{
				TiDB: &ComponentTiDB{
					NodeSize:     tidbNodeSize,
					NodeQuantity: int32(tidbNodeCnt),
				},
				TiKV: &ComponentTiKV{
					NodeSize:       tikvNodeSize,
					NodeQuantity:   int32(tikvNodeCnt),
					StorageSizeGib: int32(tikvStgGb),
				},
			},
			IPAccessList: []IPAccess{
				{
					CIDR:        cidr,
					Description: cidr,
				},
			},
		},
	}

	url := fmt.Sprintf("%s/api/v1beta/projects/%d/clusters", host, projectID)
	var result CreateClusterResp
	client.DoPOST(url, payload, &result, cmd)
	return &result
}

func getSpecifications(host string, client *RestApiClient, cmd model.ParsedCmd) *GetSpecificationsResp {
	url := fmt.Sprintf("%s/api/v1beta/clusters/provider/regions", host)
	var result GetSpecificationsResp
	client.DoGET(url, nil, &result, cmd)
	return &result
}

func printSpec(env *model.Env, screen model.Screen, spec *Specification) {
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
