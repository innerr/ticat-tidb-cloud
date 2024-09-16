package tidb_cloud

import (
	"fmt"
)

// TODO: there are two json name styles

type Project struct {
	ID              uint64 `json:"id,string"`
	OrgID           uint64 `json:"orgId,string"`
	Name            string `json:"name"`
	ClusterCount    int64  `json:"clusterCount"`
	UserCount       int64  `json:"userCount"`
	CreateTimestamp int64  `json:"createTimestamp,string"`
}

// TODO: (serverless) the values of VpcPeering and Standard are the same ?!

type ConnectionStandard struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type ConnectionVpcPeering struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type ConnectionString struct {
	DefaultUser string               `json:"default_user"`
	Standard    ConnectionStandard   `json:"standard"`
	VpcPeering  ConnectionVpcPeering `json:"vpc_peering"`
}

type IPAccess struct {
	CIDR        string `json:"cidr"`
	Description string `json:"description"`
}

type ComponentTiDB struct {
	NodeSize       string `json:"node_size"`
	StorageSizeGib int32  `json:"storage_size_gib"`
	NodeQuantity   int32  `json:"node_quantity"`
}

type ComponentTiKV struct {
	NodeSize       string `json:"node_size"`
	StorageSizeGib int32  `json:"storage_size_gib"`
	NodeQuantity   int32  `json:"node_quantity"`
}

type ComponentTiFlash struct {
	NodeSize       string `json:"node_size"`
	StorageSizeGib int32  `json:"storage_size_gib"`
	NodeQuantity   int32  `json:"node_quantity"`
}

type Components struct {
	TiDB    *ComponentTiDB    `json:"tidb,omitempty"`
	TiKV    *ComponentTiKV    `json:"tikv,omitempty"`
	TiFlash *ComponentTiFlash `json:"tiflash,omitempty"`
}

type ClusterConfig struct {
	RootPassword string     `json:"root_password"`
	Port         int32      `json:"port"`
	Components   Components `json:"components"`
	IPAccessList []IPAccess `json:"ip_access_list"`
}

type ClusterStatus struct {
	TidbVersion   string `json:"tidb_version"`
	ClusterStatus string `json:"cluster_status"`
	// TODO: `json:"node_map"`
	ConnectionStrings ConnectionString `json:"connection_strings"`
}

type CreateClusterReq struct {
	Name          string        `json:"name"`
	ClusterType   string        `json:"cluster_type"`
	CloudProvider string        `json:"cloud_provider"`
	Region        string        `json:"region"`
	Config        ClusterConfig `json:"config"`
}

type CreateClusterResp struct {
	ClusterID uint64 `json:"id,string"`
	Message   string `json:"message"`
}

type GetAllProjectsResp struct {
	Items []Project `json:"items"`
	Total int64     `json:"total"`
}

type GetClusterResp struct {
	ID              uint64        `json:"id,string"`
	ProjectID       uint64        `json:"project_id,string"`
	Name            string        `json:"name"`
	ClusterType     string        `json:"cluster_type"`
	CloudProvider   string        `json:"cloud_provider"`
	Region          string        `json:"region"`
	Status          ClusterStatus `json:"status"`
	CreateTimestamp string        `json:"create_timestamp"`
	Config          ClusterConfig `json:"config"`
}

type NodeQuantityRange struct {
	Min  int `json:"min"`
	Step int `json:"step"`
}

func (n *NodeQuantityRange) String() string {
	return fmt.Sprintf("%v+N*%v node(s)", n.Min, n.Step)
}

type TiDBSpec struct {
	NodeSize          string            `json:"node_size"`
	NodeQuantityRange NodeQuantityRange `json:"node_quantity_range"`
}

func (t *TiDBSpec) String() string {
	return fmt.Sprintf("%v %v", t.NodeSize, t.NodeQuantityRange.String())
}

type StorageSizeGibRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

func (s *StorageSizeGibRange) String() string {
	return fmt.Sprintf("%v~%v GB", s.Min, s.Max)
}

type StorageSpec struct {
	NodeSize            string              `json:"node_size"`
	NodeQuantityRange   NodeQuantityRange   `json:"node_quantity_range"`
	StorageSizeGibRange StorageSizeGibRange `json:"storage_size_gib_range"`
}

func (t *StorageSpec) String() string {
	return fmt.Sprintf("%v %v, %v", t.NodeSize, t.NodeQuantityRange.String(), t.StorageSizeGibRange.String())
}

type Specification struct {
	ClusterType   string        `json:"cluster_type"`
	CloudProvider string        `json:"cloud_provider"`
	Region        string        `json:"region"`
	TiDB          []TiDBSpec    `json:"tidb"`
	TiKV          []StorageSpec `json:"tikv"`
	TiFlash       []StorageSpec `json:"tiflash"`
}

type GetSpecificationsResp struct {
	Items []Specification `json:"items"`
}
