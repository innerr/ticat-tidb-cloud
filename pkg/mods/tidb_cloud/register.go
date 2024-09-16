package tidb_cloud

import (
	"github.com/innerr/ticat/pkg/core/model"
)

func RegisterCmds(cmds *model.CmdTree) {
	addAuthArgs := func(cmd *model.Cmd) {
		cmd.AddArg("public-key", "", "public", "pub").
			AddArg2Env(EnvKeyPubKey, "public-key").
			AddArg("private-key", "", "private", "pri").
			AddArg2Env(EnvKeyPriKey, "private-key").
			AddArg("api-address", "", "api-addr").
			AddArg2Env(EnvKeyApiAddr, "api-address").
			AddEnvOp(EnvKeyPubKey, model.EnvOpTypeRead).
			AddEnvOp(EnvKeyPriKey, model.EnvOpTypeRead).
			AddEnvOp(EnvKeyApiAddr, model.EnvOpTypeRead)
	}

	addClusterIdArgs := func(cmd *model.Cmd) {
		cmd.AddArg("cluster-id", "", "cluster", "id").
			AddArg2Env(EnvKeyClusterId, "cluster-id").
			AddArg("project-id", "", "project", "pro", "prj").
			AddArg2Env(EnvKeyProject, "project-id").
			AddEnvOp(EnvKeyProject, model.EnvOpTypeRead).
			AddEnvOp(EnvKeyClusterId, model.EnvOpTypeRead)
	}

	tc := cmds.AddSub(
		"tidb-cloud", "tc").RegEmptyCmd(
		"tidb-cloud toolbox").Owner()

	tc.AddSub(
		"projects", "prjs", "prj", "pro", "p").RegPowerCmd(ProjectsList,
		"list projects")

	tc.AddSub(
		"specs", "spec").RegPowerCmd(SpecsList,
		"list specs").
		AddArg("type", "", "t").
		AddArg("region", "", "r").
		AddArg("provider", "", "p")

	conf := tc.AddSub(
		"config", "conf", "cfg").RegEmptyCmd(
		"config tidb cloud").Owner()

	confPrj := conf.AddSub("project", "prj", "pro").
		RegPowerCmd(ProjectSelectTheOnlyOne,
			"assume there is only one project and select it").
		AddEnvOp(EnvKeyProject, model.EnvOpTypeWrite)
	addAuthArgs(confPrj)

	confPrjByName := confPrj.AddSub("by-name", "name").
		RegPowerCmd(ProjectSelectByName,
			"select the only match project with search string").
		AddArg("find-str", "", "name", "str").
		AddEnvOp(EnvKeyProject, model.EnvOpTypeWrite)
	addAuthArgs(confPrjByName)

	clusters := tc.AddSub(
		"clusters", "cluster", "c").RegEmptyCmd(
		"tidb cloud cluster service toolbox").Owner()

	clusterGet := clusters.AddSub("get").
		RegPowerCmd(ClusterGet,
			"get info of a cluster by id")
	addAuthArgs(clusterGet)
	addClusterIdArgs(clusterGet)

	clusterDelete := clusters.AddSub("delete", "remove", "del", "rm").
		RegPowerCmd(ClusterDelete,
			"delete a cluster")
	addAuthArgs(clusterDelete)
	addClusterIdArgs(clusterDelete)

	clusterWait := clusters.AddSub("wait-status", "wait").
		RegPowerCmd(ClusterWait,
			"wait for cluster to become the specify status (default: AVAILABLE)").
		AddArg("interval-secs", "2", "interval").
		AddArg("timeout-secs", "600", "timeout").
		AddArg("status", "AVAILABLE", "s").
		AddEnvOp("mysql.host", model.EnvOpTypeWrite).
		AddEnvOp("mysql.port", model.EnvOpTypeWrite).
		AddEnvOp("mysql.user", model.EnvOpTypeWrite)
	addAuthArgs(clusterWait)
	addClusterIdArgs(clusterWait)

	serverless := tc.AddSub(
		"serverless", "sl", "s").RegEmptyCmd(
		"tidb cloud serverless cluster toolbox").Owner()
	serverlessClusters := serverless.AddSub(
		"clusters", "cluster", "c").RegEmptyCmd(
		"tidb cloud serverless cluster service toolbox").Owner()

	serverlessCreate := serverlessClusters.AddSub("create", "new", "c").
		RegPowerCmd(ClusterServerlessCreate,
			"create a serverless cluster").
		AddArg("cluster-name", "", "name").
		AddArg2Env(EnvKeyClusterName, "cluster-name").
		AddArg("project-id", "", "project", "pro", "prj").
		AddArg2Env(EnvKeyProject, "project-id").
		AddArg("root-password", "", "root-pwd", "password", "pwd").
		AddArg2Env(EnvKeyRootPwd, "root-password").
		AddArg("cloud-provider", "AWS", "provider").
		AddArg2Env(EnvKeyCloudProvider, "cloud-provider").
		AddArg("cloud-region", "us-east-1", "region").
		AddArg2Env(EnvKeyCloudRegion, "cloud-region").
		AddEnvOp(EnvKeyClusterName, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyRootPwd, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyCloudProvider, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyCloudRegion, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyProject, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyClusterId, model.EnvOpTypeWrite).
		AddEnvOp("mysql.pwd", model.EnvOpTypeWrite)
	addAuthArgs(serverlessCreate)

	dedicated := tc.AddSub(
		"dedicated", "dd", "d").RegEmptyCmd(
		"tidb cloud dedicated cluster toolbox").Owner()

	dedicatedClusters := dedicated.AddSub(
		"clusters", "cluster", "c").RegEmptyCmd(
		"tidb cloud dedicated cluster service toolbox").Owner()

	dedicatedCreate := dedicatedClusters.AddSub("create", "new", "c").
		RegPowerCmd(ClusterDedicatedCreate,
			"create a dedicated cluster").
		AddArg("cluster-name", "", "name").
		AddArg2Env(EnvKeyClusterName, "cluster-name").
		AddArg("project-id", "", "project", "pro", "prj").
		AddArg2Env(EnvKeyProject, "project-id").
		AddArg("root-password", "", "root-pwd", "password", "pwd").
		AddArg2Env(EnvKeyRootPwd, "root-password").
		AddArg("cloud-provider", "AWS", "provider").
		AddArg2Env(EnvKeyCloudProvider, "cloud-provider").
		AddArg("cloud-region", "us-east-1", "region").
		AddArg2Env(EnvKeyCloudRegion, "cloud-region").
		AddArg("tidb-node-size", "4C16G", "tidb-size", "tidb", "db").
		AddArg("tidb-node-count", "1", "tidb-count", "tidb-cnt", "db-cnt").
		AddArg("tikv-node-size", "4C16G", "tikv-size", "tikv", "kv").
		AddArg("tikv-node-count", "3", "tikv-count", "tikv-cnt", "kv-cnt").
		AddArg("tikv-storage-gb", "10", "tikv-storage", "tikv-gb", "tikv-stg", "kv-stg", "stg-gb", "stg", "gb").
		AddArg("cidr", "172.16.0.0/21").
		AddEnvOp(EnvKeyClusterName, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyRootPwd, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyCloudProvider, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyCloudRegion, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyProject, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyClusterId, model.EnvOpTypeWrite).
		AddEnvOp("mysql.pwd", model.EnvOpTypeWrite)
	addAuthArgs(dedicatedCreate)
}

const (
	EnvKeyPubKey        = "tidb-cloud.auth.public-key"
	EnvKeyPriKey        = "tidb-cloud.auth.private-key"
	EnvKeyApiAddr       = "tidb-cloud.address.api"
	EnvKeyClusterName   = "tidb-cloud.cluster.name"
	EnvKeyRootPwd       = "tidb-cloud.cluster.root-pwd"
	EnvKeyCloudProvider = "tidb-cloud.provider"
	EnvKeyCloudRegion   = "tidb-cloud.provider.region"
	EnvKeyProject       = "tidb-cloud.project"
	EnvKeyClusterId     = "tidb-cloud.cluster.id"
)
