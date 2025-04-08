package main

import (
	"fmt"
	"github.com/tuannguyensn2001/go-json-rule-engine"
	"math"
	"reflect"
	"strings"
	"time"
)

func main() {
	// Create a new engine
	e := go_json_rules_engine.NewEngine()

	// Register custom operators
	err := e.RegisterCustomOperator("divisibleBy", func(a, b interface{}) bool {
		// Convert both values to float64
		va := reflect.ValueOf(a)
		vb := reflect.ValueOf(b)
		if !e.IsNumeric(va) || !e.IsNumeric(vb) {
			return false
		}
		aFloat := e.ToFloat64(va)
		bFloat := e.ToFloat64(vb)
		if bFloat == 0 {
			return false
		}
		return math.Mod(aFloat, bFloat) == 0
	})
	if err != nil {
		fmt.Printf("Error registering divisibleBy operator: %v\n", err)
		return
	}

	err = e.RegisterCustomOperator("startsWith", func(a, b interface{}) bool {
		str, ok1 := a.(string)
		prefix, ok2 := b.(string)
		if !ok1 || !ok2 {
			return false
		}
		return strings.HasPrefix(str, prefix)
	})
	if err != nil {
		fmt.Printf("Error registering startsWith operator: %v\n", err)
		return
	}

	err = e.RegisterCustomOperator("olderThan", func(a, b interface{}) bool {
		// Convert age to float64
		va := reflect.ValueOf(a)
		if !e.IsNumeric(va) {
			return false
		}
		age := e.ToFloat64(va)

		// Convert years to float64
		vb := reflect.ValueOf(b)
		if !e.IsNumeric(vb) {
			return false
		}
		years := e.ToFloat64(vb)

		return age > years
	})
	if err != nil {
		fmt.Printf("Error registering olderThan operator: %v\n", err)
		return
	}

	err = e.RegisterCustomOperator("isWeekend", func(a, _ interface{}) bool {
		date, ok := a.(time.Time)
		if !ok {
			return false
		}
		return date.Weekday() == time.Saturday || date.Weekday() == time.Sunday
	})
	if err != nil {
		fmt.Printf("Error registering isWeekend operator: %v\n", err)
		return
	}

	// Create rules using custom operators
	rules := []go_json_rules_engine.Rule{
		{
			ID:       "age-rule",
			Name:     "Age Check Rule",
			Priority: 1,
			Conditions: go_json_rules_engine.ConditionGroup{
				Operator: go_json_rules_engine.And,
				Conditions: []interface{}{
					go_json_rules_engine.Condition{
						Fact:     "age",
						Operator: "divisibleBy",
						Value:    5,
					},
					go_json_rules_engine.Condition{
						Fact:     "age",
						Operator: "olderThan",
						Value:    18,
					},
				},
			},
			Event: go_json_rules_engine.Event{
				Type: "ageRuleMatched",
				Params: map[string]interface{}{
					"message": "Age requirements met!",
				},
			},
		},
		{
			ID:       "name-rule",
			Name:     "Name Check Rule",
			Priority: 2,
			Conditions: go_json_rules_engine.ConditionGroup{
				Operator: go_json_rules_engine.And,
				Conditions: []interface{}{
					go_json_rules_engine.Condition{
						Fact:     "name",
						Operator: "startsWith",
						Value:    "John",
					},
				},
			},
			Event: go_json_rules_engine.Event{
				Type: "nameRuleMatched",
				Params: map[string]interface{}{
					"message": "Name starts with John!",
				},
			},
		},
		{
			ID:       "weekend-rule",
			Name:     "Weekend Check Rule",
			Priority: 3,
			Conditions: go_json_rules_engine.ConditionGroup{
				Operator: go_json_rules_engine.And,
				Conditions: []interface{}{
					go_json_rules_engine.Condition{
						Fact:     "date",
						Operator: "isWeekend",
						Value:    nil, // Value is not used for this operator
					},
				},
			},
			Event: go_json_rules_engine.Event{
				Type: "weekendRuleMatched",
				Params: map[string]interface{}{
					"message": "It's the weekend!",
				},
			},
		},
	}

	// Add rules to the engine
	for _, rule := range rules {
		if err := e.AddRule(rule); err != nil {
			fmt.Printf("Error adding rule: %v\n", err)
			return
		}
	}

	// Test cases
	testCases := []map[string]interface{}{
		{
			"age":  25,                                          // divisible by 5 and > 18
			"name": "John Doe",                                  // starts with "John"
			"date": time.Date(2024, 3, 9, 0, 0, 0, 0, time.UTC), // Saturday
		},
		{
			"age":  16,                                          // not divisible by 5 and < 18
			"name": "John Smith",                                // starts with "John"
			"date": time.Date(2024, 3, 8, 0, 0, 0, 0, time.UTC), // Friday
		},
		{
			"age":  30,                                           // divisible by 5 and > 18
			"name": "Jane Doe",                                   // doesn't start with "John"
			"date": time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC), // Sunday
		},
	}

	// Evaluate each test case
	for i, facts := range testCases {
		events, err := e.Evaluate(facts)
		if err != nil {
			fmt.Printf("Error evaluating facts: %v\n", err)
			continue
		}

		fmt.Printf("\nTest Case %d:\n", i+1)
		fmt.Printf("Facts: %+v\n", facts)
		if len(events) > 0 {
			fmt.Printf("Matched Events: %+v\n", events)
		} else {
			fmt.Println("No rules matched")
		}
	}
}
