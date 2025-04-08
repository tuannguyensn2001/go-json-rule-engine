package main

import (
	"fmt"
	"github.com/tuannguyensn2001/go-json-rule-engine"
)

func main() {
	// Create a new engine
	eng := go_json_rules_engine.NewEngine()

	// Create a rule programmatically
	rule := go_json_rules_engine.Rule{
		ID:       "customer-eligibility",
		Name:     "Premium Customer Eligibility",
		Priority: 10,
		Conditions: go_json_rules_engine.ConditionGroup{
			Operator: go_json_rules_engine.And,
			Conditions: []interface{}{
				go_json_rules_engine.Condition{
					Fact:     "age",
					Operator: go_json_rules_engine.GreaterThanInc,
					Value:    21,
				},
				go_json_rules_engine.ConditionGroup{
					Operator: go_json_rules_engine.Or,
					Conditions: []interface{}{
						go_json_rules_engine.Condition{
							Fact:     "yearlyPurchases",
							Operator: go_json_rules_engine.GreaterThan,
							Value:    1000.0,
						},
						go_json_rules_engine.Condition{
							Fact:     "membershipLevel",
							Operator: go_json_rules_engine.In,
							Value:    []interface{}{"gold", "platinum"},
						},
					},
				},
			},
		},
		Event: go_json_rules_engine.Event{
			Type: "premium-eligible",
			Params: map[string]interface{}{
				"message":  "Customer is eligible for premium status",
				"discount": 20,
			},
		},
	}

	// Add rule to engine
	eng.AddRule(rule)

	// Evaluate facts
	facts := map[string]interface{}{
		"age":             25,
		"yearlyPurchases": 1200.0,
		"membershipLevel": "gold",
	}

	events, err := eng.Evaluate(facts)
	if err != nil {
		panic(err)
	}

	// Handle events
	for _, event := range events {
		fmt.Printf("Rule triggered: %s\n", event.Type)
		if msg, ok := event.Params["message"]; ok {
			fmt.Printf("Message: %s\n", msg)
		}
		if discount, ok := event.Params["discount"]; ok {
			fmt.Printf("Discount: %v%%\n", discount)
		}
	}
}
