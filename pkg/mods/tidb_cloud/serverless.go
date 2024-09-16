package tidb_cloud

import (
	"fmt"

	"github.com/innerr/ticat/pkg/cli/display"
	"github.com/innerr/ticat/pkg/core/model"
)

func ClusterServerlessCreate(
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
	cmd := flow.Cmds[currCmdIdx]

	project := getProject(host, client, env, cc.Screen, cmd)
	cluster := createDevCluster(host, client, project, name, rootPwd, cloudProvider, region, cmd)

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

func createDevCluster(host string, client *RestApiClient, projectID uint64, name string, rootPwd string,
	cloudProvider string, cloudRegion string, cmd model.ParsedCmd) *CreateClusterResp {

	payload := CreateClusterReq{
		Name:          name,
		ClusterType:   "DEVELOPER",
		CloudProvider: "AWS",
		Region:        cloudRegion,
		Config: ClusterConfig{
			RootPassword: rootPwd,
			IPAccessList: []IPAccess{
				{
					CIDR:        "0.0.0.0/0",
					Description: "allow access from anywhere",
				},
			},
		},
	}

	url := fmt.Sprintf("%s/api/v1beta/projects/%d/clusters", host, projectID)
	var result CreateClusterResp
	client.DoPOST(url, payload, &result, cmd)
	return &result
}
