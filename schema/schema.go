package schema

// Global represents the global configuration
type Global struct {
	Variables []string `yaml:"variables"`
	Brokers   []Broker `yaml:"brokers"`
}

// Broker represents a broker connection configuration
type Broker struct {
	ConnectionID  string `yaml:"connection"`
	ConnectionURL string `yaml:"connection_url"`
	Type          string `yaml:"type"`
}

// ExpectedResult represents a condition and its expected value
type ExpectedResult struct {
	ID        int      `yaml:"id"`
	Condition string   `yaml:"condition"`
	Checks    []string `yaml:"checks,omitempty"`
}

// ResponseDataType represents the data configuration for an action
type ResponseDataType struct {
	ID   string `yaml:"id"`
	Type string `yaml:"type"`
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

// Action represents a generic action that can be either an endpoint, broker, or database action
type Action struct {
	ID               int              `yaml:"id"`
	Type             string           `yaml:"type"`
	MethodType       string           `yaml:"method_type,omitempty"`
	URL              string           `yaml:"url,omitempty"`
	SleepBefore      int              `yaml:"sleep_before,omitempty"`
	Exchange         string           `yaml:"exchange,omitempty"`
	Queue            string           `yaml:"queue,omitempty"`
	RoutingKey       string           `yaml:"routing_key,omitempty"`
	Message          string           `yaml:"message,omitempty"`
	ConnectionID     string           `yaml:"connection_id,omitempty"`
	Query            string           `yaml:"query,omitempty"`
	ResponseDataType ResponseDataType `yaml:"data,omitempty"`
}

// Test represents a test case with variables and actions
type Test struct {
	ID                 int      `yaml:"id"`
	Variables          []string `yaml:"variables"`
	Actions            []int    `yaml:"actions"`
	BehaviourOnFail    string   `yaml:"behaviour_on_fail"`
	BehaviourOnSuccess string   `yaml:"behaviour_on_success"`
	ExpectedResults    []int    `yaml:"expected_results"`
}

// Generate represents a generation configuration
type Generate struct {
	ID              int      `yaml:"id"`
	Count           int      `yaml:"count"`
	Variables       []string `yaml:"variables"`
	Actions         []int    `yaml:"actions"`
	ExpectedResults []int    `yaml:"expected_results"`
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
