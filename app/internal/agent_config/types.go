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
	Host     ServiceHost       `json:"host"`
	Method   networks.Method   `json:"method"`
	Endpoint string            `json:"endpoint"`
	Headers  map[string]string `json:"headers"`
}

type Param struct {
	Key         string
	Type        string
	Required    bool
	Description string
	Example     string
}

type Action struct {
	Name         ActionName
	Title        string
	Description  string
	UserExamples []string

	PathParams  []Param
	QueryParams []Param
	BodyParams  []Param

	Pagination bool

	ResponseJSON string
	IsActive     bool
	API          APIConfig
}
type DecisionConfig struct {
	SystemInstruction string   `json:"system_instruction"`
	Actions           []Action `json:"actions"`
}
