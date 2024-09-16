package tidb_cloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/innerr/ticat/pkg/core/model"
)

type V1Beta1IPAccess struct {
	CIDR        string `json:"cidr"`
	Description string `json:"description"`
}

type V1Beta1ClusterConfig struct {
	RootPassword string            `json:"rootPassword"`
	Port         int32             `json:"port"`
	Components   LegacyComponents  `json:"components"`
	IPAccessList []V1Beta1IPAccess `json:"ipAccessList"`
}

type V1Beta1Region struct {
	DisplayName   string `json:"displayName"`
	Name          string `json:"name"`
	CloudProvider string `json:"cloudProvider"`
	RegionId      string `json:"regionId"`
}

type V1Beta1Endpoints struct {
	Public bool `json:"public"`
}

type V1Beta1CreateClusterReq struct {
	DisplayName  string            `json:"displayName"`
	Region       V1Beta1Region     `json:"region"`
	RootPassword string            `json:"rootPassword"`
	Endpoints    V1Beta1Endpoints  `json:"endpoint"`
	Labels       map[string]string `json:"labels"`
}

type V1Beta1EncryptionConfig struct {
	EnhancedEncryptionEnabled bool `json:"enhancedEncryptionEnabled"`
}

type V1Beta1CreateClusterResp struct {
	Name             string                  `json:"name"`
	DisplayName      string                  `json:"displayName"`
	ClusterID        string                  `json:"clusterId",string`
	Region           V1Beta1Region           `json:"region"`
	Version          string                  `json:"version"`
	CreatedBy        string                  `json:"createdBy"`
	UserPrefix       string                  `json:"userPrefix"`
	State            string                  `json:"state"`
	Labels           map[string]string       `json:"labels"`
	Annotations      map[string]string       `json:"annotations"`
	CreateTime       time.Time               `json:"createTime",string`
	UpdateTime       time.Time               `json:"updateTime",string`
	EncryptionConfig V1Beta1EncryptionConfig `json:"encryptionConfig"`
}

func V1Beta1CreateDevCluster(
	host string,
	client *RestApiClient,
	projectID uint64,
	name string,
	rootPwd string,
	cloudProvider string,
	cloudRegion string, cmd model.ParsedCmd) *V1Beta1CreateClusterResp {

	cloudProvider = strings.ToLower(cloudProvider)

	payload := V1Beta1CreateClusterReq{
		DisplayName:  name,
		RootPassword: rootPwd,
		Region: V1Beta1Region{
			Name:          "regions/" + cloudProvider + "-" + cloudRegion,
			CloudProvider: cloudProvider,
		},
		Endpoints: V1Beta1Endpoints{
			Public: true,
		},
		Labels: map[string]string{
			"tidb.cloud/project": fmt.Sprintf("%v", projectID),
		},
	}
	url := fmt.Sprintf("%s/v1beta1/clusters", host)
	var result V1Beta1CreateClusterResp
	client.DoPOST(url, payload, &result, cmd)
	return &result
}
