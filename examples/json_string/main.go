package main

import (
	"fmt"
	"log"

	rules "github.com/tuannguyensn2001/go-json-rule-engine"
)

func main() {
	// Create a new rules engine
	engine := rules.NewEngine()

	// Example JSON string containing rules
	jsonRules := `[
		{
			"id": "discount-rule",
			"name": "10% Discount for orders over $100",
			"priority": 1,
			"conditions": {
				"operator": "and",
				"conditions": [
					{
						"fact": "total_amount",
						"operator": "greaterThan",
						"value": 100
					}
				]
			},
			"event": {
				"type": "apply_discount",
				"params": {
					"percentage": 10
				}
			}
		},
		{
			"id": "vip-discount",
			"name": "15% Discount for VIP customers",
			"priority": 2,
			"conditions": {
				"operator": "and",
				"conditions": [
					{
						"fact": "customer_type",
						"operator": "equal",
						"value": "VIP"
					}
				]
			},
			"event": {
				"type": "apply_discount",
				"params": {
					"percentage": 15
				}
			}
		}
	]`

	// Create a new Rule instance
	rule := rules.NewRules()

	// Load rules from JSON string
	err := rule.LoadRulesFromJSONString(jsonRules)
	if err != nil {
		log.Fatalf("Failed to load rules: %v", err)
	}

	// Create facts for evaluation
	facts := map[string]interface{}{
		"total_amount":  120.00,
		"customer_type": "VIP",
	}

	// Evaluate rules
	events, err := engine.Evaluate(rule, facts)
	if err != nil {
		log.Fatalf("Failed to evaluate rules: %v", err)
	}

	// Process the events
	for _, event := range events {
		fmt.Printf("Rule matched! Event type: %s, Parameters: %v\n", event.Type, event.Params)
	}
}
