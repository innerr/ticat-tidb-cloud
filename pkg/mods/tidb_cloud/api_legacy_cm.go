package tidb_cloud

import (
	"fmt"

	"github.com/innerr/ticat/pkg/core/model"
)

// TODO: there are two json name styles
// TODO: (serverless) the values of VpcPeering and Standard are the same ?!

type LegacyConnectionStandard struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type LegacyConnectionVpcPeering struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type LegacyConnectionString struct {
	DefaultUser string                     `json:"default_user"`
	Standard    LegacyConnectionStandard   `json:"standard"`
	VpcPeering  LegacyConnectionVpcPeering `json:"vpc_peering"`
}

type LegacyIPAccess struct {
	CIDR        string `json:"cidr"`
	Description string `json:"description"`
}

type LegacyComponent struct {
	NodeSize       string `json:"node_size"`
	StorageSizeGib int32  `json:"storage_size_gib"`
	NodeQuantity   int32  `json:"node_quantity"`
}

type LegacyComponents struct {
	TiDB    *LegacyComponent `json:"tidb,omitempty"`
	TiKV    *LegacyComponent `json:"tikv,omitempty"`
	TiFlash *LegacyComponent `json:"tiflash,omitempty"`
}

type LegacyClusterConfig struct {
	RootPassword string           `json:"root_password"`
	Port         int32            `json:"port"`
	Components   LegacyComponents `json:"components"`
	IPAccessList []LegacyIPAccess `json:"ip_access_list"`
}

type LegacyClusterStatus struct {
	TidbVersion   string `json:"tidb_version"`
	ClusterStatus string `json:"cluster_status"`
	// TODO: missed `json:"node_map"` here
	ConnectionStrings LegacyConnectionString `json:"connection_strings"`
}

type LegacyCreateClusterReq struct {
	Name          string              `json:"name"`
	ClusterType   string              `json:"cluster_type"`
	CloudProvider string              `json:"cloud_provider"`
	Region        string              `json:"region"`
	Config        LegacyClusterConfig `json:"config"`
}

type LegacyCreateClusterResp struct {
	ClusterID uint64 `json:"id,string"`
	Message   string `json:"message"`
}

func LegacyCreateDevCluster(
	host string,
	client *RestApiClient,
	projectID uint64,
	name string,
	rootPwd string,
	cloudProvider string,
	cloudRegion string, cmd model.ParsedCmd) *LegacyCreateClusterResp {

	payload := LegacyCreateClusterReq{
		Name:          name,
		ClusterType:   "DEVELOPER",
		CloudProvider: cloudProvider,
		Region:        cloudRegion,
		Config: LegacyClusterConfig{
			RootPassword: rootPwd,
			IPAccessList: []LegacyIPAccess{
				{
					CIDR:        "0.0.0.0/0",
					Description: "allow access from anywhere",
				},
			},
		},
	}

	url := fmt.Sprintf("%s/api/v1beta/projects/%d/clusters", host, projectID)
	var result LegacyCreateClusterResp
	client.DoPOST(url, payload, &result, cmd)
	return &result
}

func LegacyCreateDedicatedCluster(
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
	cmd model.ParsedCmd) *LegacyCreateClusterResp {

	payload := LegacyCreateClusterReq{
		Name:          name,
		ClusterType:   "DEDICATED",
		CloudProvider: cloudProvider,
		Region:        cloudRegion,
		Config: LegacyClusterConfig{
			RootPassword: rootPwd,
			Port:         4000,
			Components: LegacyComponents{
				TiDB: &LegacyComponent{
					NodeSize:     tidbNodeSize,
					NodeQuantity: int32(tidbNodeCnt),
				},
				TiKV: &LegacyComponent{
					NodeSize:       tikvNodeSize,
					NodeQuantity:   int32(tikvNodeCnt),
					StorageSizeGib: int32(tikvStgGb),
				},
			},
			IPAccessList: []LegacyIPAccess{
				{
					CIDR:        cidr,
					Description: cidr,
				},
			},
		},
	}

	url := fmt.Sprintf("%s/api/v1beta/projects/%d/clusters", host, projectID)
	var result LegacyCreateClusterResp
	client.DoPOST(url, payload, &result, cmd)
	return &result
}

type LegacyGetClusterResp struct {
	ID              uint64              `json:"id,string"`
	ProjectID       uint64              `json:"project_id,string"`
	Name            string              `json:"name"`
	ClusterType     string              `json:"cluster_type"`
	CloudProvider   string              `json:"cloud_provider"`
	Region          string              `json:"region"`
	Status          LegacyClusterStatus `json:"status"`
	CreateTimestamp string              `json:"create_timestamp"`
	Config          LegacyClusterConfig `json:"config"`
}

func LegacyGetClusterByID(host string, client *RestApiClient, project, cluster uint64, cmd model.ParsedCmd) *LegacyGetClusterResp {
	url := fmt.Sprintf("%s/api/v1beta/projects/%d/clusters/%d", host, project, cluster)
	var result LegacyGetClusterResp
	client.DoGET(url, nil, &result, cmd)
	return &result
}

func LegacyDeleteClusterByID(host string, client *RestApiClient, projectID, clusterID uint64, cmd model.ParsedCmd) {
	url := fmt.Sprintf("%s/api/v1beta/projects/%d/clusters/%d", host, projectID, clusterID)
	client.DoDELETE(url, nil, nil, cmd)
}
