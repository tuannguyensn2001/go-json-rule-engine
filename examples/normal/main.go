package main

import (
	"fmt"

	go_json_rules_engine "github.com/tuannguyensn2001/go-json-rule-engine"
)

func main() {
	// Create a new engine
	eng := go_json_rules_engine.NewEngine()

	// Create rules programmatically
	rules := go_json_rules_engine.NewRules()

	// Add premium customer rule
	rules.AddRule(go_json_rules_engine.RuleOption{
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
	})

	// Add basic customer rule
	rules.AddRule(go_json_rules_engine.RuleOption{
		ID:       "basic-customer",
		Name:     "Basic Customer Check",
		Priority: 5,
		Conditions: go_json_rules_engine.ConditionGroup{
			Operator: go_json_rules_engine.And,
			Conditions: []interface{}{
				go_json_rules_engine.Condition{
					Fact:     "age",
					Operator: go_json_rules_engine.GreaterThanInc,
					Value:    18,
				},
				go_json_rules_engine.Condition{
					Fact:     "membershipLevel",
					Operator: go_json_rules_engine.Equal,
					Value:    "basic",
				},
			},
		},
		Event: go_json_rules_engine.Event{
			Type: "basic-eligible",
			Params: map[string]interface{}{
				"message":  "Customer is eligible for basic status",
				"discount": 5,
			},
		},
	})

	// Test cases
	testCases := []map[string]interface{}{
		{
			"age":             25,
			"yearlyPurchases": 1200.0,
			"membershipLevel": "gold",
		},
		{
			"age":             19,
			"yearlyPurchases": 500.0,
			"membershipLevel": "basic",
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
				if msg, ok := event.Params["message"]; ok {
					fmt.Printf("Message: %s\n", msg)
				}
				if discount, ok := event.Params["discount"]; ok {
					fmt.Printf("Discount: %v%%\n", discount)
				}
			}
		} else {
			fmt.Printf("No rules triggered\n")
		}
	}
}
