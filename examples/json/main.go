package main

import (
	"fmt"
	"github.com/tuannguyensn2001/go-json-rule-engine"
)

func main() {
	// Create a new engine
	eng := go_json_rules_engine.NewEngine()

	// Load rules from JSON string
	jsonStr := `[
        {
            "id": "vip-customer",
            "name": "VIP Customer Rule",
            "priority": 20,
            "conditions": {
                "operator": "and",
                "conditions": [
                    {
                        "fact": "membershipLevel",
                        "operator": "equal",
                        "value": "platinum"
                    },
                    {
                        "fact": "yearsAsMember",
                        "operator": "greaterThan",
                        "value": 5
                    }
                ]
            },
            "event": {
                "type": "vip-status",
                "params": {
                    "message": "Customer is a VIP member",
                    "benefits": ["priority support", "exclusive offers"]
                }
            }
        },
        {
            "id": "new-customer",
            "name": "New Customer Welcome",
            "priority": 5,
            "conditions": {
                "operator": "and",
                "conditions": [
                    {
                        "fact": "yearsAsMember",
                        "operator": "lessThan",
                        "value": 1
                    },
                    {
                        "fact": "firstPurchase",
                        "operator": "equal",
                        "value": true
                    }
                ]
            },
            "event": {
                "type": "welcome-offer",
                "params": {
                    "message": "Welcome to our service!",
                    "offer": "10% off your next purchase"
                }
            }
        }
    ]`

	if err := eng.LoadRulesFromJSONString(jsonStr); err != nil {
		panic(err)
	}

	// Evaluate facts for VIP customer
	fmt.Println("Testing VIP customer scenario:")
	vipFacts := map[string]interface{}{
		"membershipLevel": "platinum",
		"yearsAsMember":   6,
		"firstPurchase":   false,
	}

	events, err := eng.Evaluate(vipFacts)
	if err != nil {
		panic(err)
	}

	for _, event := range events {
		fmt.Printf("Rule triggered: %s\n", event.Type)
		if msg, ok := event.Params["message"]; ok {
			fmt.Printf("Message: %s\n", msg)
		}
		if benefits, ok := event.Params["benefits"]; ok {
			fmt.Printf("Benefits: %v\n", benefits)
		}
	}

	// Evaluate facts for new customer
	fmt.Println("\nTesting new customer scenario:")
	newCustomerFacts := map[string]interface{}{
		"membershipLevel": "basic",
		"yearsAsMember":   0,
		"firstPurchase":   true,
	}

	events, err = eng.Evaluate(newCustomerFacts)
	if err != nil {
		panic(err)
	}

	for _, event := range events {
		fmt.Printf("Rule triggered: %s\n", event.Type)
		if msg, ok := event.Params["message"]; ok {
			fmt.Printf("Message: %s\n", msg)
		}
		if offer, ok := event.Params["offer"]; ok {
			fmt.Printf("Offer: %s\n", offer)
		}
	}
}
