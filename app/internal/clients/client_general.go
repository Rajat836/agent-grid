package clients

import (
	"app/agent_grid/internal/agent_config"
	"context"
	"fmt"
	"net/http"

	networks "bitbucket.org/fyscal/be-commons/pkg/network"
)

type ClientGeneral struct {
	access *clientAccess
}

type ClientGeneralMethods interface {
	SendRequestByAgent(ctx context.Context, action *agent_config.Action, reqHeaders map[string]string, req any, queryParams map[string]any) (any, error)
}

func NewClientGeneral(access *clientAccess) ClientGeneralMethods {
	return &ClientGeneral{
		access: access,
	}
}

func (cu *ClientGeneral) SendRequestByAgent(ctx context.Context, action *agent_config.Action, reqHeaders map[string]string, req any, queryParams map[string]any) (any, error) {
	var (
		logger       = cu.access.logger.With(ctx)
		networkOps   = cu.access.networkOps
		apiConfig    = action.API
		headers      = make(http.Header)
		responseData any
	)

	headers.Add("Content-Type", "application/json")
	for k, v := range apiConfig.Headers {
		headers.Add(k, v)
	}

	queryParamsStr := make(map[string]string)
	for k, v := range queryParams {
		queryParamsStr[k] = fmt.Sprint(v)
	}

	payload := &networks.ApiPayload{
		Url:         fmt.Sprintf("%s/%s", apiConfig.Host, apiConfig.Endpoint),
		Method:      apiConfig.Method,
		Headers:     headers,
		Body:        req,
		JsonObject:  &responseData,
	}

	logger.Infof("sending request to %s with payload: %+v", payload.Url, payload.Body, payload.Headers, payload.Method)
	resp, reqErr := networkOps.SendRequest(ctx, payload)
	if reqErr != nil {
		logger.Errorf("error sending request: %v", reqErr)
		return nil, reqErr
	}

	logger.Infof("received response: %+v", string(resp.Data))

	return string(resp.Data), nil
}
