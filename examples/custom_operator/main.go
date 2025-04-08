package main

import (
	"fmt"

	go_json_rules_engine "github.com/tuannguyensn2001/go-json-rule-engine"
)

func main() {
	eng := go_json_rules_engine.NewEngine()

	// Register a custom operator for checking if a number is divisible by another number
	err := eng.RegisterCustomOperator("divisibleBy", func(a, b interface{}) bool {
		// Convert both values to float64 for comparison
		aFloat, aOk := a.(float64)
		bFloat, bOk := b.(float64)
		if !aOk || !bOk {
			return false
		}
		return aFloat != 0 && bFloat != 0 && int(aFloat)%int(bFloat) == 0
	})
	if err != nil {
		panic(err)
	}

	// Register a custom operator for checking if a string contains a substring
	err = eng.RegisterCustomOperator("containsSubstring", func(a, b interface{}) bool {
		aStr, aOk := a.(string)
		bStr, bOk := b.(string)
		if !aOk || !bOk {
			return false
		}
		return len(bStr) > 0 && len(aStr) >= len(bStr) && aStr != "" && bStr != ""
	})
	if err != nil {
		panic(err)
	}

	// Create rules container
	rules := go_json_rules_engine.NewRules()

	// Load rules from JSON string
	jsonStr := `[
		{
			"id": "divisible-check",
			"name": "Divisible Check",
			"priority": 1,
			"conditions": {
				"operator": "and",
				"conditions": [
					{
						"fact": "number",
						"operator": "divisibleBy",
						"value": 5
					}
				]
			},
			"event": {
				"type": "divisible",
				"params": {
					"message": "Number is divisible by 5"
				}
			}
		},
		{
			"id": "substring-check",
			"name": "Substring Check",
			"priority": 2,
			"conditions": {
				"operator": "and",
				"conditions": [
					{
						"fact": "text",
						"operator": "containsSubstring",
						"value": "hello"
					}
				]
			},
			"event": {
				"type": "contains",
				"params": {
					"message": "Text contains the substring"
				}
			}
		}
	]`

	if err := rules.LoadRulesFromJSONString(jsonStr); err != nil {
		panic(err)
	}

	// Test cases
	testCases := []map[string]interface{}{
		{
			"number": 10.0,
			"text":   "hello world",
		},
		{
			"number": 7.0,
			"text":   "goodbye world",
		},
		{
			"number": 15.0,
			"text":   "hello there",
		},
	}

	for i, facts := range testCases {
		fmt.Printf("\nTesting Case %d:\n", i+1)
		events, err := eng.Evaluate(rules.GetRules(), facts)
		if err != nil {
			panic(err)
		}

		if len(events) > 0 {
			for _, event := range events {
				fmt.Printf("Rule triggered: %s\n", event.Type)
				fmt.Printf("Message: %s\n", event.Params["message"])
			}
		} else {
			fmt.Printf("No rules triggered\n")
		}
	}
}
