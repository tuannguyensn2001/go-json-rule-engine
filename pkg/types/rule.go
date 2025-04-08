// Package types provides the core data structures for the JSON rules engine.
// It defines the types for rules, conditions, and events that are used throughout the engine.
package types

import (
	"encoding/json"
	"errors"
)

// Operator represents the comparison operators that can be used in conditions.
type Operator string

const (
	// Equal checks if two values are equal
	Equal Operator = "equal"
	// NotEqual checks if two values are not equal
	NotEqual Operator = "notEqual"
	// GreaterThan checks if the first value is greater than the second
	GreaterThan Operator = "greaterThan"
	// LessThan checks if the first value is less than the second
	LessThan Operator = "lessThan"
	// GreaterThanInc checks if the first value is greater than or equal to the second
	GreaterThanInc Operator = "greaterThanInclusive"
	// LessThanInc checks if the first value is less than or equal to the second
	LessThanInc Operator = "lessThanInclusive"
	// In checks if a value exists in a collection
	In Operator = "in"
	// NotIn checks if a value does not exist in a collection
	NotIn Operator = "notIn"
	// Contains checks if a string contains another string
	Contains Operator = "contains"
	// NotContains checks if a string does not contain another string
	NotContains Operator = "notContains"
	// Regex checks if a string matches a regular expression
	Regex Operator = "regex"
	// IsNull checks if a value is null
	IsNull Operator = "isNull"
	// IsNotNull checks if a value is not null
	IsNotNull Operator = "isNotNull"
)

// LogicalOperator represents the logical operators that can be used to combine conditions.
type LogicalOperator string

const (
	// And represents a logical AND operation
	And LogicalOperator = "and"
	// Or represents a logical OR operation
	Or LogicalOperator = "or"
)

// Condition represents a single condition in a rule.
// It consists of a fact to evaluate, an operator to use, and a value to compare against.
type Condition struct {
	// Fact is the name of the fact to evaluate
	Fact string `json:"fact"`
	// Operator is the comparison operator to use
	Operator Operator `json:"operator"`
	// Value is the value to compare against
	Value interface{} `json:"value"`
}

// ConditionGroup represents a group of conditions with a logical operator.
// It can contain both individual conditions and nested condition groups.
type ConditionGroup struct {
	// Operator is the logical operator to use when combining conditions
	Operator LogicalOperator `json:"operator"`
	// Conditions is a slice of conditions or condition groups
	Conditions []interface{} `json:"conditions"`
}

// UnmarshalJSON implements custom JSON unmarshaling for ConditionGroup.
// It handles the conversion of JSON conditions into either Condition or ConditionGroup types.
func (cg *ConditionGroup) UnmarshalJSON(data []byte) error {
	type Alias ConditionGroup
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(cg),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Validate and convert conditions
	for i, condition := range cg.Conditions {
		condData, err := json.Marshal(condition)
		if err != nil {
			return err
		}

		// Try to unmarshal as Condition
		var singleCond Condition
		if err := json.Unmarshal(condData, &singleCond); err == nil {
			cg.Conditions[i] = singleCond
			continue
		}

		// Try to unmarshal as ConditionGroup
		var groupCond ConditionGroup
		if err := json.Unmarshal(condData, &groupCond); err == nil {
			cg.Conditions[i] = groupCond
			continue
		}

		return errors.New("invalid condition format")
	}

	return nil
}

// Rule represents a business rule in the engine.
// It consists of an ID, name, priority, conditions to evaluate, and an event to trigger.
type Rule struct {
	// ID is the unique identifier for the rule
	ID string `json:"id"`
	// Name is the human-readable name of the rule
	Name string `json:"name"`
	// Priority determines the order in which rules are evaluated (higher numbers first)
	Priority int `json:"priority"`
	// Conditions is the root condition group to evaluate
	Conditions ConditionGroup `json:"conditions"`
	// Event is the event to trigger when the rule's conditions are met
	Event Event `json:"event"`
}

// Event represents what should happen when a rule's conditions are met.
type Event struct {
	// Type is the type of event to trigger
	Type string `json:"type"`
	// Params contains additional parameters for the event
	Params map[string]interface{} `json:"params,omitempty"`
}
