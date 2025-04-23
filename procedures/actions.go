package procedures

import (
	"fmt"

	"github.com/Ishogbon/code-chaos/schema"
)

type Actions struct {
	actions map[int]schema.Action
}

var ActionStore = &Actions{
	actions: make(map[int]schema.Action),
}

func (s *Actions) Add(action schema.Action) {
	// Check if an action with the given id already exists
	if _, exists := s.actions[action.ID]; exists {
		panic(fmt.Sprintf("attempt to overwrite existing action with ID %d", action.ID))
	}

	// Store the new action
	s.actions[action.ID] = action
}

func (s *Actions) Get(id int) (*schema.Action, error) {
	action, exists := s.actions[id]
	if !exists {
		return nil, fmt.Errorf("action with ID %d not found", id)
	}
	return &action, nil
}

func StoreAction(action schema.Action) {
	ActionStore.Add(action)
}

func StoreActions(actions []schema.Action) {
	for _, action := range actions {
		StoreAction(action)
	}
}

func GetAction(id int) (*schema.Action, error) {
	return ActionStore.Get(id)
}
