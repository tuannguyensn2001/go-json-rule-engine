// Package types provides the core data structures for the JSON rules engine.
// It defines the types for rules, conditions, and events that are used throughout the engine.
package go_json_rules_engine

import (
	"encoding/json"
	"errors"
)

// Operator represents the comparison operators that can be used in conditions.
type Operator string

const (
	Equal          Operator = "equal"
	NotEqual       Operator = "notEqual"
	GreaterThan    Operator = "greaterThan"
	LessThan       Operator = "lessThan"
	GreaterThanInc Operator = "greaterThanInclusive"
	LessThanInc    Operator = "lessThanInclusive"
	In             Operator = "in"
	NotIn          Operator = "notIn"
	Regex          Operator = "regex"
	IsNull         Operator = "isNull"
	IsNotNull      Operator = "isNotNull"
)

type LogicalOperator string

const (
	And LogicalOperator = "and"
	Or  LogicalOperator = "or"
)

type Condition struct {
	Fact     string      `json:"fact"`
	Operator Operator    `json:"operator"`
	Value    interface{} `json:"value"`
}

type ConditionGroup struct {
	// Operator is the logical operator to use when combining conditions
	Operator LogicalOperator `json:"operator"`
	// Conditions is a slice of conditions or condition groups
	Conditions []interface{} `json:"conditions"`
}

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

type Rule struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Priority   int            `json:"priority"`
	Conditions ConditionGroup `json:"conditions"`
	Event      Event          `json:"event"`
}

// Event represents what should happen when a rule's conditions are met.
type Event struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params,omitempty"`
}
