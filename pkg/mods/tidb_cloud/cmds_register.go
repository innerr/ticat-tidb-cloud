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
			AddEnvOp(EnvKeyPubKey, model.EnvOpTypeRead).
			AddEnvOp(EnvKeyPriKey, model.EnvOpTypeRead)
	}

	addApiAddrArgLegacy := func(cmd *model.Cmd) {
		cmd.AddArg("api-address", "", "api-addr").
			AddArg2Env(EnvKeyApiAddrLegacy, "api-address").
			AddEnvOp(EnvKeyApiAddrLegacy, model.EnvOpTypeRead)
	}

	addApiAddrArgV1Beta1 := func(cmd *model.Cmd) {
		cmd.AddArg("api-address", "", "api-addr").
			AddArg2Env(EnvKeyApiAddrV1Beta1, "api-address").
			AddEnvOp(EnvKeyApiAddrLegacy, model.EnvOpTypeRead)
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
	addApiAddrArgLegacy(confPrj)

	confPrjByName := confPrj.AddSub("by-name", "name").
		RegPowerCmd(ProjectSelectByName,
			"select the only match project with search string").
		AddArg("find-str", "", "name", "str").
		AddEnvOp(EnvKeyProject, model.EnvOpTypeWrite)
	addAuthArgs(confPrjByName)
	addApiAddrArgLegacy(confPrjByName)

	clusters := tc.AddSub(
		"clusters", "cluster", "c").RegEmptyCmd(
		"tidb cloud cluster service toolbox").Owner()

	clusterGet := clusters.AddSub("get").
		RegPowerCmd(ClusterGet,
			"get info of a cluster by id")
	addAuthArgs(clusterGet)
	addClusterIdArgs(clusterGet)
	addApiAddrArgLegacy(clusterGet)

	clusterDelete := clusters.AddSub("delete", "remove", "del", "rm").
		RegPowerCmd(ClusterDelete,
			"delete a cluster")
	addAuthArgs(clusterDelete)
	addClusterIdArgs(clusterDelete)
	addApiAddrArgLegacy(clusterDelete)

	clusterWait := clusters.AddSub("wait-status", "wait").
		RegPowerCmd(ClusterWait,
			"wait for cluster to become the specify status (default: AVAILABLE)").
		AddArg("interval-secs", "3", "interval", "interv").
		AddArg("timeout-secs", "3600", "timeout").
		AddArg("status", "AVAILABLE", "s").
		AddEnvOp("mysql.host", model.EnvOpTypeWrite).
		AddEnvOp("mysql.port", model.EnvOpTypeWrite).
		AddEnvOp("mysql.user", model.EnvOpTypeWrite)
	addAuthArgs(clusterWait)
	addClusterIdArgs(clusterWait)
	addApiAddrArgLegacy(clusterWait)

	serverless := tc.AddSub(
		"serverless", "sl", "s").RegEmptyCmd(
		"tidb cloud serverless cluster toolbox").Owner()
	serverlessClusters := serverless.AddSub(
		"clusters", "cluster", "c").RegEmptyCmd(
		"tidb cloud serverless cluster service toolbox").Owner()

	serverlessCreate := serverlessClusters.AddSub("create", "new", "c").
		RegPowerCmd(V1Beta1ClusterServerlessCreate,
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
		AddArg("spending-limit-monthly-us-cents", "-1", "spending-limit", "spend-limit", "limit").
		AddArg2Env(EnvKeySpendLimit, "spending-limit-monthly-us-cents").
		AddEnvOp(EnvKeyClusterName, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyRootPwd, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyCloudProvider, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyCloudRegion, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyProject, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyClusterId, model.EnvOpTypeWrite).
		AddEnvOp("mysql.pwd", model.EnvOpTypeWrite)
	addAuthArgs(serverlessCreate)
	addApiAddrArgV1Beta1(serverlessCreate)

	serverlessCreateLegacy := serverlessCreate.AddSub("legacy", "l").
		RegPowerCmd(V1Beta1ClusterServerlessCreate,
			"use legacy api to create a serverless cluster").
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
	addAuthArgs(serverlessCreateLegacy)
	addApiAddrArgLegacy(serverlessCreateLegacy)

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
		AddArg("cloud-region", "us-west-2", "region").
		AddArg2Env(EnvKeyCloudRegion, "cloud-region").
		AddArg("tidb-node-size", "2C4G", "tidb-size", "tidb", "db").
		AddArg("tidb-node-count", "1", "tidb-count", "tidb-cnt", "db-cnt").
		AddArg("tikv-node-size", "2C4G", "tikv-size", "tikv", "kv").
		AddArg("tikv-node-count", "3", "tikv-count", "tikv-cnt", "kv-cnt").
		AddArg("tikv-storage-gb", "10", "tikv-storage", "tikv-gb", "tikv-stg", "kv-stg", "stg-gb", "stg", "gb").
		AddArg("access-ip-list", "0.0.0.0/0", "access-cidr", "access-ip", "access").
		AddEnvOp(EnvKeyClusterName, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyRootPwd, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyCloudProvider, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyCloudRegion, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyProject, model.EnvOpTypeRead).
		AddEnvOp(EnvKeyClusterId, model.EnvOpTypeWrite).
		AddEnvOp("mysql.pwd", model.EnvOpTypeWrite)
	addAuthArgs(dedicatedCreate)
	addApiAddrArgLegacy(dedicatedCreate)

	cmds.GetOrAddSub("init", "tidb-cloud", "replace").RegPowerCmd(
		IntegrationReplaceCmds, "tidb-cloud integration: replace some cmds")
}

/*
bench.add-tags
bench.record.tags
bench.record

bench.result
bench.result.clear

--

bench.result.list
bench.result.show

bench.record.config
bench.record.jitter
bench.record.usage

bench.result.tag.add
bench.result.tag.all
bench.result.tag.remove

bench.result.last
bench.result.last.aggregate
bench.result.select-result
bench.result.select-result.last

bench.result.info
bench.result.monitor-links
bench.result.tiup-yaml
bench.result.dashboard
*/

const (
	EnvKeyPubKey         = "tidb-cloud.auth.public-key"
	EnvKeyPriKey         = "tidb-cloud.auth.private-key"
	EnvKeyApiAddrLegacy  = "tidb-cloud.address.api.legacy"
	EnvKeyApiAddrV1Beta1 = "tidb-cloud.address.api.v1beta1"
	EnvKeyClusterName    = "tidb-cloud.cluster.name"
	EnvKeyRootPwd        = "tidb-cloud.cluster.root-pwd"
	EnvKeyCloudProvider  = "tidb-cloud.provider"
	EnvKeyCloudRegion    = "tidb-cloud.provider.region"
	EnvKeyProject        = "tidb-cloud.project"
	EnvKeyClusterId      = "tidb-cloud.cluster.id"
	EnvKeySpendLimit     = "tidb-cloud.cluster.serverless.spend-limit"
)
