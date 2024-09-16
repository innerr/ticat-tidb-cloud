package main

import (
	"os"

	"github.com/innerr/ticat/pkg/ticat"

	"github.com/innerr/ticat-tidb-cloud/pkg/integrate/tidb_cloud"
)

func main() {
	ticat := ticat.NewTiCat()
	tidb_cloud.Integrate(ticat)
	ticat.RunCli(os.Args[1:]...)
}
