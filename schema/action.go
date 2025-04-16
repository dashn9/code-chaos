package schema

type ExpectedResult struct {
	Condition     string `yaml:"condition"`
	ExpectedValue string `yaml:"expected_value"`
}

type Action struct {
	ID                 int              `yaml:"id"`
	Variables          []string         `yaml:"variables"`
	ActionType         string           `yaml:"action_type"`
	ActionData         []string         `yaml:"action_data"`
	ColumnName         string           `yaml:"column_name"`
	BehaviourOnFail    string           `yaml:"behaviour_on_fail"`
	BehaviourOnSuccess string           `yaml:"behaviour_on_success"`
	ExpectedResults    []ExpectedResult `yaml:"expected_results"`
}

type Test struct {
	ID      int      `yaml:"id"`
	Actions []Action `yaml:"actions"`
}
