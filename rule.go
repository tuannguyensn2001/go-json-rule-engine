// Package types provides the core data structures for the JSON rules engine.
// It defines the types for rules, conditions, and events that are used throughout the engine.
package go_json_rules_engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
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

type ruleOption struct {
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

type Rule struct {
	opts []ruleOption
}

func NewRules() *Rule {
	return &Rule{
		opts: make([]ruleOption, 0),
	}
}

func (r *Rule) AddRule(opts ruleOption) {
	r.opts = append(r.opts, opts)
}

func (r *Rule) LoadRulesFromJSON(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read rules file: %w", err)
	}

	return r.LoadRulesFromJSONString(string(data))
}

func (r *Rule) LoadRulesFromJSONString(jsonStr string) error {
	var rules []ruleOption
	if err := json.Unmarshal([]byte(jsonStr), &rules); err != nil {
		return fmt.Errorf("failed to parse rules: %w", err)
	}

	r.opts = rules
	r.sortRulesByPriority()
	return nil
}

func (r *Rule) GetRules() []ruleOption {
	return r.opts
}

func (r *Rule) sortRulesByPriority() {
	sort.Slice(r.opts, func(i, j int) bool {
		return r.opts[i].Priority > r.opts[j].Priority
	})
}
