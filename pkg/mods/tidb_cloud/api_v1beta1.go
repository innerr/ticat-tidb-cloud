package tidb_cloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/innerr/ticat/pkg/core/model"
)

type V1Beta1Region struct {
	DisplayName   string `json:"displayName"`
	Name          string `json:"name"`
	CloudProvider string `json:"cloudProvider"`
	RegionId      string `json:"regionId"`
}

type V1Beta1Endpoints struct {
	Public bool `json:"public"`
}

type V1Beta1SpendingLimit struct {
	Monthly int32 `json:"monthly"`
}

type V1Beta1CreateClusterReq struct {
	DisplayName   string                `json:"displayName"`
	Region        V1Beta1Region         `json:"region"`
	RootPassword  string                `json:"rootPassword"`
	Endpoints     V1Beta1Endpoints      `json:"endpoint"`
	Labels        map[string]string     `json:"labels"`
	SpendingLimit *V1Beta1SpendingLimit `json:"spendingLimit,omitempty"`
}

type V1Beta1EncryptionConfig struct {
	EnhancedEncryptionEnabled bool `json:"enhancedEncryptionEnabled"`
}

type V1Beta1Cluster struct {
	Name             string                  `json:"name"`
	DisplayName      string                  `json:"displayName"`
	ClusterId        string                  `json:"clusterId",string`
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
	projectId uint64,
	name string,
	rootPwd string,
	cloudProvider string,
	cloudRegion string,
	limitMonthlyUSCents int,
	cmd model.ParsedCmd) *V1Beta1Cluster {

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
			"tidb.cloud/project": fmt.Sprintf("%v", projectId),
		},
	}
	if limitMonthlyUSCents > 0 {
		payload.SpendingLimit = &V1Beta1SpendingLimit{
			Monthly: int32(limitMonthlyUSCents),
		}
	}
	url := fmt.Sprintf("%s/v1beta1/clusters", host)
	var result V1Beta1Cluster
	client.DoPOST(url, payload, &result, cmd)
	return &result
}

type V1Beta1ChangeRootPasswordReq struct {
	ClusterId string `json:"clusterId"`
	Password  string `json:"password"`
}

func V1Beta1ChangeRootPassword(
	host string,
	client *RestApiClient,
	clusterId string,
	rootPwd string,
	cmd model.ParsedCmd) {

	payload := V1Beta1ChangeRootPasswordReq{
		ClusterId: clusterId,
		Password:  rootPwd,
	}
	url := fmt.Sprintf("%s/v1beta1/clusters/%s/password", host, clusterId)
	client.DoPOST(url, payload, nil, cmd)
}
