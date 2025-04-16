package schema

type ExpectedResult struct {
	Condition     string
	ExpectedValue string
}

type Action struct {
	ID                 int
	Variables          []string
	ActionType         string
	ActionData         []string
	ColumnName         string
	BehaviourOnFail    string
	BehaviourOnSuccess string
	ExpectedResults    []ExpectedResult
}

type Test struct {
	ID      int
	Actions []Action
}
