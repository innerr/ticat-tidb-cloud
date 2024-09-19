package tidb_cloud

import (
	"fmt"

	"github.com/innerr/ticat/pkg/core/model"
)

type LegacyProject struct {
	Id              uint64 `json:"id,string"`
	OrgId           uint64 `json:"orgId,string"`
	Name            string `json:"name"`
	ClusterCount    int64  `json:"clusterCount"`
	UserCount       int64  `json:"userCount"`
	CreateTimestamp int64  `json:"createTimestamp,string"`
}

type LegacyGetAllProjectsResp struct {
	Items []LegacyProject `json:"items"`
	Total int64           `json:"total"`
}

func LegacyGetAllProjects(host string, client *RestApiClient, cmd model.ParsedCmd) []LegacyProject {
	url := fmt.Sprintf("%s/api/v1beta/projects", host)
	var result LegacyGetAllProjectsResp
	client.DoGET(url, nil, &result, cmd)
	return result.Items
}

type LegacyNodeQuantityRange struct {
	Min  int `json:"min"`
	Step int `json:"step"`
}

func (n *LegacyNodeQuantityRange) String() string {
	return fmt.Sprintf("%v+N*%v node(s)", n.Min, n.Step)
}

type LegacyTiDBSpec struct {
	NodeSize                string                  `json:"node_size"`
	LegacyNodeQuantityRange LegacyNodeQuantityRange `json:"node_quantity_range"`
}

func (t *LegacyTiDBSpec) String() string {
	return fmt.Sprintf("%v %v", t.NodeSize, t.LegacyNodeQuantityRange.String())
}

type LegacyStorageSizeGibRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

func (s *LegacyStorageSizeGibRange) String() string {
	return fmt.Sprintf("%v~%v GB", s.Min, s.Max)
}

type LegacyStorageSpec struct {
	NodeSize            string                    `json:"node_size"`
	NodeQuantityRange   LegacyNodeQuantityRange   `json:"node_quantity_range"`
	StorageSizeGibRange LegacyStorageSizeGibRange `json:"storage_size_gib_range"`
}

func (t *LegacyStorageSpec) String() string {
	return fmt.Sprintf("%v %v, %v", t.NodeSize, t.NodeQuantityRange.String(), t.StorageSizeGibRange.String())
}

type LegacySpecification struct {
	ClusterType   string              `json:"cluster_type"`
	CloudProvider string              `json:"cloud_provider"`
	Region        string              `json:"region"`
	TiDB          []LegacyTiDBSpec    `json:"tidb"`
	TiKV          []LegacyStorageSpec `json:"tikv"`
	TiFlash       []LegacyStorageSpec `json:"tiflash"`
}

type LegacyGetSpecificationsResp struct {
	Items []LegacySpecification `json:"items"`
}

func LegacyGetSpecifications(host string, client *RestApiClient, cmd model.ParsedCmd) *LegacyGetSpecificationsResp {
	url := fmt.Sprintf("%s/api/v1beta/clusters/provider/regions", host)
	var result LegacyGetSpecificationsResp
	client.DoGET(url, nil, &result, cmd)
	return &result
}
