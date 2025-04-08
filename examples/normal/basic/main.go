package main

import (
	"fmt"

	go_json_rules_engine "github.com/tuannguyensn2001/go-json-rule-engine"
)

func main() {
	// Create a new engine
	eng := go_json_rules_engine.NewEngine()

	// Define rules
	rules := []go_json_rules_engine.RuleOption{
		{
			ID:       "age-check",
			Name:     "Adult Check",
			Priority: 1,
			Conditions: go_json_rules_engine.ConditionGroup{
				Operator: go_json_rules_engine.And,
				Conditions: []interface{}{
					go_json_rules_engine.Condition{
						Fact:     "age",
						Operator: go_json_rules_engine.GreaterThan,
						Value:    18,
					},
				},
			},
			Event: go_json_rules_engine.Event{
				Type: "adult",
				Params: map[string]interface{}{
					"message": "User is an adult",
				},
			},
		},
	}

	// Evaluate facts
	facts := map[string]interface{}{
		"age": 25,
	}

	events, err := eng.Evaluate(rules, facts)
	if err != nil {
		panic(err)
	}

	// Handle events
	for _, event := range events {
		fmt.Printf("Event: %s\n", event.Type)
		fmt.Printf("Message: %s\n", event.Params["message"])
	}
}
