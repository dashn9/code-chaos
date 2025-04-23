package schema

// ExpectedResult represents a condition and its expected value
type ExpectedResult struct {
	ID            int    `yaml:"id"`
	Condition     string `yaml:"condition"`
	ExpectedValue string `yaml:"expected_value"`
}

// EndpointAction represents an HTTP endpoint action
type EndpointAction struct {
	ID         int    `yaml:"id"`
	Type       string `yaml:"type"`
	MethodType string `yaml:"method_type"`
	URL        string `yaml:"url"`
	Payload    string `yaml:"payload"`
}

// RmqBrokerAction represents a RabbitMQ broker action
type RmqBrokerAction struct {
	ID         int    `yaml:"id"`
	Type       string `yaml:"type"`
	Exchange   string `yaml:"exchange"`
	Queue      string `yaml:"queue"`
	RoutingKey string `yaml:"routing_key"`
	Message    string `yaml:"message"`
}

// Action represents a generic action that can be either an endpoint or broker action
type Action struct {
	ID                 int              `yaml:"id"`
	Type               string           `yaml:"type"`
	MethodType         string           `yaml:"method_type,omitempty"`
	URL                string           `yaml:"url,omitempty"`
	Payload            string           `yaml:"payload,omitempty"`
	Exchange           string           `yaml:"exchange,omitempty"`
	Queue              string           `yaml:"queue,omitempty"`
	RoutingKey         string           `yaml:"routing_key,omitempty"`
	Message            string           `yaml:"message,omitempty"`
	BehaviourOnFail    string           `yaml:"behaviour_on_fail,omitempty"`
	BehaviourOnSuccess string           `yaml:"behaviour_on_success,omitempty"`
	ExpectedResults    []ExpectedResult `yaml:"expected_results,omitempty"`
}

// Test represents a test case with variables and actions
type Test struct {
	ID                 int              `yaml:"id"`
	Variables          []string         `yaml:"variables"`
	Action             int              `yaml:"action"`
	BehaviourOnFail    string           `yaml:"behaviour_on_fail"`
	BehaviourOnSuccess string           `yaml:"behaviour_on_success"`
	ExpectedResults    []ExpectedResult `yaml:"expected_results"`
}

// TestFile represents the complete test file structure
type TestFile struct {
	Globals []Global         `yaml:"globals"`
	Results []ExpectedResult `yaml:"results"`
	Actions []Action         `yaml:"actions"`
	Tests   []Test           `yaml:"tests"`
}

// Global represents global variables and settings
type Global struct {
	Variables []string `yaml:"variables"`
}
