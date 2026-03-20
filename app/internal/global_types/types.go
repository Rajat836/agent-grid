package globaltypes

import "time"

type RequestOntologySummary struct {
	Prompt string `json:"prompt"`
}

type ResponseOntologySummary struct {
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	EntityCount int       `json:"entity_count"`
	EdgeCount   int       `json:"edge_count"`
	Timestamp   time.Time `json:"timestamp"`
}
