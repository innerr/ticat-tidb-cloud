package tidb_cloud

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/icholy/digest"

	"github.com/innerr/ticat/pkg/cli/display"
	"github.com/innerr/ticat/pkg/core/model"
)

type RestApiClient struct {
	publicKey   string
	privateKey  string
	env         *model.Env
	screen      model.Screen
	screenCusor int
	restClient  *resty.Client
}

func NewRestApiClient(env *model.Env, screen model.Screen) *RestApiClient {
	client := &RestApiClient{
		publicKey:   env.GetRaw(EnvKeyPubKey),
		privateKey:  env.GetRaw(EnvKeyPriKey),
		env:         env,
		screen:      screen,
		screenCusor: screen.OutputtedLines(),
	}
	client.restClient = resty.New()
	client.restClient.SetTransport(&digest.Transport{
		Username: client.publicKey,
		Password: client.privateKey,
	})
	return client
}

func (client *RestApiClient) DoGET(url string, payload, output interface{}, cmd model.ParsedCmd) *resty.Response {
	return client.DoRequest(resty.MethodGet, url, payload, output, cmd)
}

func (client *RestApiClient) DoPOST(url string, payload, output interface{}, cmd model.ParsedCmd) *resty.Response {
	return client.DoRequest(resty.MethodPost, url, payload, output, cmd)
}

func (client *RestApiClient) DoDELETE(url string, payload, output interface{}, cmd model.ParsedCmd) *resty.Response {
	return client.DoRequest(resty.MethodDelete, url, payload, output, cmd)
}

func (client *RestApiClient) DoPATCH(url string, payload, output interface{}, cmd model.ParsedCmd) *resty.Response {
	return client.DoRequest(resty.MethodPatch, url, payload, output, cmd)
}

func (client *RestApiClient) DoRequest(method, url string, payload, output interface{}, cmd model.ParsedCmd) *resty.Response {
	request := client.restClient.R()
	lines := client.screen.OutputtedLines()
	outputted := false

	printf := func(msg string) {
		if !outputted && lines != client.screenCusor {
			client.screen.Print("\n")
			outputted = true
		}
		msg = display.ColorExplain(msg, client.env)
		client.screen.Print(msg)
	}

	_, _ = json.Marshal(payload)
	printf(fmt.Sprintf(">req: %s %s\n", method, url))
	if payload != nil {
		request.SetBody(payload)
		printf(fmt.Sprintf("payload: %+v\n", payload))
	}

	resp, err := request.Execute(method, url)
	if err != nil {
		panic(model.WrapCmdError(cmd, err))
	}

	printf(fmt.Sprintf("resp: %v %s\n", resp.StatusCode(), resp))
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("status %d, resp %s", resp.StatusCode(), resp)
		panic(model.WrapCmdError(cmd, err))
	}

	if output != nil {
		err = json.Unmarshal(resp.Body(), output)
		if err != nil {
			panic(model.WrapCmdError(cmd, err))
		}
	}

	client.screenCusor = client.screen.OutputtedLines()
	return resp
}
