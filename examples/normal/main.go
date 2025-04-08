package main

import (
	"fmt"

	go_json_rules_engine "github.com/tuannguyensn2001/go-json-rule-engine"
)

func main() {
	// Create a new engine
	eng := go_json_rules_engine.NewEngine()

	// Create rules using JSON
	rules := go_json_rules_engine.NewRules()
	jsonStr := `[
		{
			"id": "customer-eligibility",
			"name": "Premium Customer Eligibility",
			"priority": 10,
			"conditions": {
				"operator": "and",
				"conditions": [
					{
						"fact": "age",
						"operator": "greaterThanInclusive",
						"value": 21
					},
					{
						"operator": "or",
						"conditions": [
							{
								"fact": "yearlyPurchases",
								"operator": "greaterThan",
								"value": 1000.0
							},
							{
								"fact": "membershipLevel",
								"operator": "in",
								"value": ["gold", "platinum"]
							}
						]
					}
				]
			},
			"event": {
				"type": "premium-eligible",
				"params": {
					"message": "Customer is eligible for premium status",
					"discount": 20
				}
			}
		},
		{
			"id": "basic-customer",
			"name": "Basic Customer Check", 
			"priority": 5,
			"conditions": {
				"operator": "and",
				"conditions": [
					{
						"fact": "age",
						"operator": "greaterThanInclusive",
						"value": 18
					},
					{
						"fact": "membershipLevel",
						"operator": "equal",
						"value": "basic"
					}
				]
			},
			"event": {
				"type": "basic-eligible",
				"params": {
					"message": "Customer is eligible for basic status",
					"discount": 5
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
		events, err := eng.Evaluate(rules, facts)
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
