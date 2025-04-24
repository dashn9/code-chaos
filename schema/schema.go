package schema

// Global represents the global configuration
type Global struct {
	Variables []string `yaml:"variables"`
	Broker    []Broker `yaml:"broker"`
}

// Broker represents a broker connection configuration
type Broker struct {
	Connection    string `yaml:"connection"`
	ConnectionURL string `yaml:"connection_url"`
}

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
	ID           int    `yaml:"id"`
	Type         string `yaml:"type"`
	Exchange     string `yaml:"exchange"`
	Queue        string `yaml:"queue"`
	RoutingKey   string `yaml:"routing_key"`
	Message      string `yaml:"message"`
	ConnectionID string `yaml:"connection_id"`
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
	ConnectionID       string           `yaml:"connection_id,omitempty"`
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

// Generate represents a generation configuration
type Generate struct {
	ID        int      `yaml:"id"`
	Count     int      `yaml:"count"`
	Variables []string `yaml:"variables"`
	Action    int      `yaml:"action"`
}

// TestFile represents the complete test file structure
type TestFile struct {
	Globals   Global           `yaml:"globals"`
	Results   []ExpectedResult `yaml:"results"`
	Resources []string         `yaml:"resources"`
	Actions   []Action         `yaml:"actions"`
	Generates []Generate       `yaml:"generates"`
	Tests     []Test           `yaml:"tests"`
}
