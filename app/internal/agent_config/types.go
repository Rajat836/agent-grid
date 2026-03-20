package agent_config

import networks "bitbucket.org/fyscal/be-commons/pkg/network"

type FilterParam struct {
	Key         string `json:"key"`
	Type        string `json:"type"` // string, number, bool, date
	Description string `json:"description"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type APIConfig struct {
	Host     string            `json:"host"`
	Method   networks.Method   `json:"method"`
	Endpoint string            `json:"endpoint"`
	Headers  map[string]string `json:"headers"`
}

type Action struct {
	Name         string   `json:"name"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	UserExamples []string `json:"user_examples"`
	ResponseJSON string   `json:"response_json"`
	Notes        string   `json:"notes,omitempty"`

	// 🔥 NEW
	Filters    []FilterParam `json:"filters,omitempty"`
	Pagination bool          `json:"pagination_supported"`

	API APIConfig `json:"api"`
}

type DecisionConfig struct {
	SystemInstruction string   `json:"system_instruction"`
	Actions           []Action `json:"actions"`
}
