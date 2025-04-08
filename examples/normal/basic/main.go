package main

import (
	"fmt"

	go_json_rules_engine "github.com/tuannguyensn2001/go-json-rule-engine"
)

func main() {
	// Create a new engine
	eng := go_json_rules_engine.NewEngine()

	// Define rules
	rulesContainer := go_json_rules_engine.NewRules()
	jsonStr := `[{
		"id": "age-check",
		"name": "Adult Check",
		"priority": 1,
		"conditions": {
			"operator": "and",
			"conditions": [{
				"fact": "age",
				"operator": "greaterThan",
				"value": 18
			}]
		},
		"event": {
			"type": "adult",
			"params": {
				"message": "User is an adult"
			}
		}
	}]`
	if err := rulesContainer.LoadRulesFromJSONString(jsonStr); err != nil {
		panic(err)
	}

	// Evaluate facts
	facts := map[string]interface{}{
		"age": 25,
	}

	events, err := eng.Evaluate(rulesContainer, facts)
	if err != nil {
		panic(err)
	}

	// Handle events
	for _, event := range events {
		fmt.Printf("Event: %s\n", event.Type)
		fmt.Printf("Message: %s\n", event.Params["message"])
	}
}
