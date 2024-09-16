package tidb_cloud

import (
	"github.com/innerr/ticat-tidb-cloud/pkg/mods/tidb_cloud"
	"github.com/innerr/ticat/pkg/ticat"
)

func Integrate(ticat *ticat.TiCat) {
	ticat.AddIntegratedModVersion("tidb-cloud 1.0")
	ticat.AddInitRepo("ticat-mods/tidb.cloud")
	ticat.AddInitRepo("ticat-mods/tidb.bench")
	tidb_cloud.RegisterCmds(ticat.Cmds)
}
