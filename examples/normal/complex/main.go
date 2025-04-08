package main

import (
	"fmt"

	go_json_rules_engine "github.com/tuannguyensn2001/go-json-rule-engine"
)

func main() {
	eng := go_json_rules_engine.NewEngine()

	// Replace the rules definition with:
	rulesContainer := go_json_rules_engine.NewRules()
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
	if err := rulesContainer.LoadRulesFromJSONString(jsonStr); err != nil {
		panic(err)
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
		events, err := eng.Evaluate(rulesContainer, facts)
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
