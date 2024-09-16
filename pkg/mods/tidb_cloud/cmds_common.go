package tidb_cloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/innerr/ticat/pkg/cli/display"
	"github.com/innerr/ticat/pkg/core/model"
)

func ProjectSelectTheOnlyOne(
	argv model.ArgVals,
	cc *model.Cli,
	env *model.Env,
	flow *model.ParsedCmds,
	currCmdIdx int) (int, bool) {

	host := env.GetRaw(EnvKeyApiAddr)
	client := NewRestApiClient(env, cc.Screen)
	cmd := flow.Cmds[currCmdIdx]

	projects := LegacyGetAllProjects(host, client, flow.Cmds[currCmdIdx])
	if len(projects) == 0 {
		msg := "no project found"
		panic(model.NewCmdError(cmd, msg))
	}
	if len(projects) > 1 {
		msg := "too many projects found, can't auto decide use which one"
		panic(model.NewCmdError(cmd, msg))
	}
	project := projects[0]
	env.GetLayer(model.EnvLayerSession).SetUint64(EnvKeyProject, project.ID)

	if cc.Screen.OutputtedLines() > 0 {
		cc.Screen.Print("\n")
	}
	printProject(env, cc.Screen, &project)
	return currCmdIdx, true
}

func ProjectSelectByName(
	argv model.ArgVals,
	cc *model.Cli,
	env *model.Env,
	flow *model.ParsedCmds,
	currCmdIdx int) (int, bool) {

	host := env.GetRaw(EnvKeyApiAddr)
	client := NewRestApiClient(env, cc.Screen)
	str := argv.GetRaw("find-str")
	cmd := flow.Cmds[currCmdIdx]

	all := LegacyGetAllProjects(host, client, flow.Cmds[currCmdIdx])
	var projects []LegacyProject
	for _, it := range all {
		if strings.Index(it.Name, str) >= 0 {
			projects = append(projects, it)
		}
	}
	if len(projects) == 0 {
		msg := fmt.Sprintf("no matched project with name `%s`", str)
		panic(model.NewCmdError(cmd, msg))
	}
	if len(projects) > 1 {
		msg := fmt.Sprintf("too many projects found with name `%s`, can't decide which one", str)
		panic(model.NewCmdError(cmd, msg))
	}
	project := projects[0]
	env.GetLayer(model.EnvLayerSession).SetUint64(EnvKeyProject, project.ID)

	if cc.Screen.OutputtedLines() > 0 {
		cc.Screen.Print("\n")
	}
	printProject(env, cc.Screen, &project)
	return currCmdIdx, true
}

func ProjectsList(
	argv model.ArgVals,
	cc *model.Cli,
	env *model.Env,
	flow *model.ParsedCmds,
	currCmdIdx int) (int, bool) {

	host := env.GetRaw(EnvKeyApiAddr)
	client := NewRestApiClient(env, cc.Screen)

	projects := LegacyGetAllProjects(host, client, flow.Cmds[currCmdIdx])
	if cc.Screen.OutputtedLines() > 0 {
		cc.Screen.Print("\n")
	}
	for _, project := range projects {
		printProject(env, cc.Screen, &project)
	}
	return currCmdIdx, true
}

func ClustersList(
	argv model.ArgVals,
	cc *model.Cli,
	env *model.Env,
	flow *model.ParsedCmds,
	currCmdIdx int) (int, bool) {

	cc.Screen.Print("TODO\n")
	return currCmdIdx, true
}

func ClusterDelete(
	argv model.ArgVals,
	cc *model.Cli,
	env *model.Env,
	flow *model.ParsedCmds,
	currCmdIdx int) (int, bool) {

	host := env.GetRaw(EnvKeyApiAddr)
	client := NewRestApiClient(env, cc.Screen)
	id := env.GetUint64(EnvKeyClusterId)
	cmd := flow.Cmds[currCmdIdx]

	project := getProject(host, client, env, cc.Screen, cmd)
	LegacyDeleteClusterByID(host, client, project, id, cmd)

	if cc.Screen.OutputtedLines() > 0 {
		cc.Screen.Print("\n")
	}
	cc.Screen.Print(fmt.Sprintf("cluster deleted, ID: %v\n", id))
	return currCmdIdx, true
}

func ClusterWait(
	argv model.ArgVals,
	cc *model.Cli,
	env *model.Env,
	flow *model.ParsedCmds,
	currCmdIdx int) (int, bool) {

	host := env.GetRaw(EnvKeyApiAddr)
	client := NewRestApiClient(env, cc.Screen)
	id := env.GetUint64(EnvKeyClusterId)
	interval := argv.GetInt("interval-secs")
	timeout := argv.GetInt("timeout-secs")
	status := argv.GetRaw("status")
	cmd := flow.Cmds[currCmdIdx]

	project := getProject(host, client, env, cc.Screen, cmd)
	cluster, ok := waitClusterStatus(host, client, project, id, interval, timeout, status, env, cc.Screen, cmd)
	if !ok {
		panic(model.NewCmdError(flow.Cmds[currCmdIdx], fmt.Sprintf("wait for cluster to be %s, timeout", status)))
	}
	if cc.Screen.OutputtedLines() > 0 {
		cc.Screen.Print("\n")
	}
	envSession := env.GetLayer(model.EnvLayerSession)
	envSession.Set("mysql.host", cluster.Status.ConnectionStrings.Standard.Host)
	envSession.SetInt("mysql.port", cluster.Status.ConnectionStrings.Standard.Port)
	envSession.Set("mysql.user", cluster.Status.ConnectionStrings.DefaultUser)
	printCluster(env, cc.Screen, cluster)
	return currCmdIdx, true
}

func ClusterGet(
	argv model.ArgVals,
	cc *model.Cli,
	env *model.Env,
	flow *model.ParsedCmds,
	currCmdIdx int) (int, bool) {

	host := env.GetRaw(EnvKeyApiAddr)
	client := NewRestApiClient(env, cc.Screen)
	id := env.GetUint64(EnvKeyClusterId)
	cmd := flow.Cmds[currCmdIdx]

	project := getProject(host, client, env, cc.Screen, cmd)
	cluster := LegacyGetClusterByID(host, client, project, id, cmd)

	if cc.Screen.OutputtedLines() > 0 {
		cc.Screen.Print("\n")
	}
	printCluster(env, cc.Screen, cluster)
	return currCmdIdx, true
}

func waitClusterStatus(
	host string,
	client *RestApiClient,
	project uint64,
	clusterId uint64,
	intervalSecs int,
	timeoutSecs int,
	status string,
	env *model.Env,
	screen model.Screen,
	cmd model.ParsedCmd) (*LegacyGetClusterResp, bool) {

	deadline := time.Now().Add(time.Duration(timeoutSecs) * time.Second)
	ticker := time.NewTicker(time.Duration(intervalSecs) * time.Second)
	defer ticker.Stop()
	var cluster *LegacyGetClusterResp
	for range ticker.C {
		cluster = LegacyGetClusterByID(host, client, project, clusterId, cmd)
		if cluster.Status.ClusterStatus == status {
			return cluster, true
		}
		if time.Now().After(deadline) {
			break
		}
		screen.Print(cluster.Status.ClusterStatus + "\n")
	}
	return cluster, false
}

func getProject(host string, client *RestApiClient, env *model.Env, screen model.Screen, cmd model.ParsedCmd) uint64 {
	project := env.GetRaw(EnvKeyProject)
	if len(project) == 0 {
		projects := LegacyGetAllProjects(host, client, cmd)
		if len(projects) == 0 {
			msg := "no project found"
			panic(model.NewCmdError(cmd, msg))
		}
		if len(projects) > 1 {
			msg := "too many projects found, can't auto decide use which one"
			panic(model.NewCmdError(cmd, msg))
		}
		return projects[0].ID
	}
	return env.GetUint64(EnvKeyProject)
}

func printCluster(env *model.Env, screen model.Screen, cluster *LegacyGetClusterResp) {
	prefix := "  "
	screen.Print(display.ColorHighLight(fmt.Sprintf("ID: %v\n", cluster.ID), env))
	printProp(env, screen, prefix, "Name", cluster.Name)
	printProp(env, screen, prefix, "ProjectID", cluster.ProjectID)
	printProp(env, screen, prefix, "ClusterType", cluster.ClusterType)
	printProp(env, screen, prefix, "CloudProvider", cluster.CloudProvider)
	printProp(env, screen, prefix, "Region", cluster.Region)
	printProp(env, screen, prefix, "Status.TiDBVersion", cluster.Status.TidbVersion)
	printProp(env, screen, prefix, "Status.ClusterStatus", cluster.Status.ClusterStatus)
	printProp(env, screen, prefix, "CreateTimestamp", cluster.CreateTimestamp)
	// printProp(env, screen, prefix, "Config.RootPassword", cluster.Config.RootPassword)
	printProp(env, screen, prefix, "Config.Port", cluster.Config.Port)
	printProp(env, screen, prefix, "Config.IPAccessList", cluster.Config.IPAccessList)
	printProp(env, screen, prefix, "Config.Components.TiDB.NodeSize", cluster.Config.Components.TiDB.NodeSize)
	printProp(env, screen, prefix, "Config.Components.TiDB.StorageSizeGib", cluster.Config.Components.TiDB.StorageSizeGib)
	printProp(env, screen, prefix, "Config.Components.TiDB.NodeQuantity", cluster.Config.Components.TiDB.NodeQuantity)
	printProp(env, screen, prefix, "Config.Components.TiKV.NodeSize", cluster.Config.Components.TiKV.NodeSize)
	printProp(env, screen, prefix, "Config.Components.TiKV.StorageSizeGib", cluster.Config.Components.TiKV.StorageSizeGib)
	printProp(env, screen, prefix, "Config.Components.TiKV.NodeQuantity", cluster.Config.Components.TiKV.NodeQuantity)
	if cluster.Config.Components.TiFlash != nil {
		printProp(env, screen, prefix, "Config.Components.TiFlash.NodeSize", cluster.Config.Components.TiFlash.NodeSize)
		printProp(env, screen, prefix, "Config.Components.TiFlash.StorageSizeGib", cluster.Config.Components.TiFlash.StorageSizeGib)
		printProp(env, screen, prefix, "Config.Components.TiFlash.NodeQuantity", cluster.Config.Components.TiFlash.NodeQuantity)
	}
	printProp(env, screen, prefix, "Status.ConnectionStrings.DefaultUser", cluster.Status.ConnectionStrings.DefaultUser)
	printProp(env, screen, prefix, "Status.ConnectionStrings.Standard.Host", cluster.Status.ConnectionStrings.Standard.Host)
	printProp(env, screen, prefix, "Status.ConnectionStrings.Standard.Port", cluster.Status.ConnectionStrings.Standard.Port)
	printProp(env, screen, prefix, "Status.ConnectionStrings.VpcPeering.Host", cluster.Status.ConnectionStrings.VpcPeering.Host)
	printProp(env, screen, prefix, "Status.ConnectionStrings.VpcPeering.Port", cluster.Status.ConnectionStrings.VpcPeering.Port)
}

func printProject(env *model.Env, screen model.Screen, project *LegacyProject) {
	prefix := "  "
	screen.Print(display.ColorHighLight(fmt.Sprintf("ID: %v\n", project.ID), env))
	printProp(env, screen, prefix, "Name", project.Name)
	printProp(env, screen, prefix, "OrgID", project.OrgID)
	printProp(env, screen, prefix, "ClusterCount", project.ClusterCount)
	printProp(env, screen, prefix, "UserCount", project.UserCount)
	printProp(env, screen, prefix, "CreateTimestamp", project.CreateTimestamp)
}

func printProp(env *model.Env, screen model.Screen, prefix string, name string, value interface{}) {
	prefix = prefix + display.ColorProp("-", env) + " "
	sep := display.ColorProp(":", env)
	valStr := ""
	if value != nil {
		valStr = fmt.Sprintf("%v\n", value)
	}
	if len(valStr) != 0 {
		valStr = " " + valStr
	}
	screen.Print(prefix + display.ColorArg(name, env) + sep + valStr)
}
