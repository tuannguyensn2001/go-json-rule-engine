# Go JSON Rules Engine

A powerful and flexible rules engine for Go that allows you to define, evaluate, and manage business rules using JSON.

[![Go Report Card](https://goreportcard.com/badge/github.com/tuannguyensn2001/go-json-rule-engine)](https://goreportcard.com/report/github.com/tuannguyensn2001/go-json-rule-engine)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/tuannguyensn2001/go-json-rule-engine.svg)](https://pkg.go.dev/github.com/tuannguyensn2001/go-json-rule-engine)

## Features

- Define rules using JSON
- Support for complex conditions (AND/OR logic)
- Rich set of operators (equal, notEqual, greaterThan, lessThan, in, contains, regex, etc.)
- Rule priorities
- JSON persistence (save/load rules)
- Type-safe evaluation
- Extensible architecture

## Installation

```bash
go get github.com/tuannguyensn2001/go-json-rule-engine
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/tuannguyensn2001/go-json-rule-engine/pkg/engine"
    "github.com/tuannguyensn2001/go-json-rule-engine/pkg/types"
)

func main() {
    // Create a new engine
    eng := engine.NewEngine()

    // Create a rule
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
        fmt.Printf("Message: %s\n", event.Params["message"])
        fmt.Printf("Discount: %v%%\n", event.Params["discount"])
    }
}
```

## Rule JSON Format

```json
{
    "id": "rule-id",
    "name": "Rule Name",
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
}
```

## Supported Operators

- `equal` / `notEqual`
- `greaterThan` / `lessThan`
- `greaterThanInclusive` / `lessThanInclusive`
- `in` / `notIn`
- `contains` / `notContains`
- `regex`
- `isNull` / `isNotNull`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 