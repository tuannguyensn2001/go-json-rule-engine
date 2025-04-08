package main

import (
	"fmt"

	"github.com/tuannguyensn2001/go-json-rule-engine/pkg/engine"
	"github.com/tuannguyensn2001/go-json-rule-engine/pkg/types"
)

func main() {
	// Create a new engine
	eng := engine.NewEngine()

	// Create a rule programmatically
	rule := types.Rule{
		ID:       "customer-eligibility",
		Name:     "Premium Customer Eligibility",
		Priority: 10,
		Conditions: types.ConditionGroup{
			Operator: types.And,
			Conditions: []interface{}{
				types.Condition{
					Fact:     "age",
					Operator: types.GreaterThanInc,
					Value:    21,
				},
				types.ConditionGroup{
					Operator: types.Or,
					Conditions: []interface{}{
						types.Condition{
							Fact:     "yearlyPurchases",
							Operator: types.GreaterThan,
							Value:    1000.0,
						},
						types.Condition{
							Fact:     "membershipLevel",
							Operator: types.In,
							Value:    []interface{}{"gold", "platinum"},
						},
					},
				},
			},
		},
		Event: types.Event{
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
