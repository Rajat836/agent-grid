package clients

import (
	"app/agent_grid/internal/config"
	"fmt"
)

type AgentName string

const (
	AgentNameOllama AgentName = "Ollama"
)

type ClientAgent struct {
	Access       *clientAccess
	ClientOllama ClientAgentMethods
}

type ClientAgentMethods interface {
	GenerateResponse(agentName AgentName, agentCfg *config.ModelConfig, prompt string) (*ResponseAgentModel, error)
}

func NewClientAgent(access *clientAccess) *ClientAgent {
	return &ClientAgent{
		Access:       access,
		ClientOllama: NewOllamaAgentClient(access),
	}
}

func (c *ClientAgent) GenerateResponse(agentName AgentName, agentCfg *config.ModelConfig, prompt string) (*ResponseAgentModel, error) {
	switch agentName {
	case AgentNameOllama:
		return c.ClientOllama.GenerateResponse(agentName, agentCfg, prompt)
	}
	return nil, fmt.Errorf("unsupported agent name: %s", agentName)
}
