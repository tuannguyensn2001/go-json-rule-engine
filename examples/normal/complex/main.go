package main

import (
	"fmt"

	go_json_rules_engine "github.com/tuannguyensn2001/go-json-rule-engine"
)

func main() {
	eng := go_json_rules_engine.NewEngine()

	// Define a complex rule for customer eligibility
	rules := []go_json_rules_engine.RuleOption{
		{
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
		},
		{
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
		},
	}

	// Test with different customer profiles
	testCases := []map[string]interface{}{
		{
			"age":             25,
			"yearlyPurchases": 1200.0,
			"membershipLevel": "silver",
		},
		{
			"age":             20,
			"yearlyPurchases": 800.0,
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
		events, err := eng.Evaluate(rules, facts)
		if err != nil {
			panic(err)
		}

		if len(events) > 0 {
			fmt.Printf("Customer is eligible for %s status!\n", events[0].Type)
			fmt.Printf("Discount: %v%%\n", events[0].Params["discount"])
			fmt.Printf("Message: %s\n", events[0].Params["message"])
		} else {
			fmt.Printf("Customer is not eligible for any status\n")
		}
	}
}
