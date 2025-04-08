# Go JSON Rules Engine

A powerful and flexible rules engine for Go that allows you to define, evaluate, and manage business rules using JSON.

[![Go Report Card](https://goreportcard.com/badge/github.com/tuannguyensn2001/go-json-rule-engine)](https://goreportcard.com/report/github.com/tuannguyensn2001/go-json-rule-engine)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/tuannguyensn2001/go-json-rule-engine.svg)](https://pkg.go.dev/github.com/tuannguyensn2001/go-json-rule-engine)

## Features

- Define rules using JSON
- Support for complex conditions (AND/OR logic)
- Rich set of operators (equal, notEqual, greaterThan, lessThan, in, contains, regex, etc.)
- Custom operator support
- Type-safe evaluation
- Thread-safe operations
- Extensible architecture

## Installation

```bash
go get github.com/tuannguyensn2001/go-json-rule-engine
```

## Usage Examples

### 1. Basic Rule Evaluation

```go
package main

import (
    "fmt"
    "github.com/tuannguyensn2001/go-json-rule-engine"
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
```

### 2. Complex Rules with Multiple Conditions

```go
package main

import (
    "fmt"
    "github.com/tuannguyensn2001/go-json-rule-engine"
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
    }

    for i, facts := range testCases {
        fmt.Printf("\nTesting Case %d:\n", i+1)
        events, err := eng.Evaluate(rules, facts)
        if err != nil {
            panic(err)
        }

        if len(events) > 0 {
            fmt.Printf("Customer is eligible!\n")
            fmt.Printf("Discount: %v%%\n", events[0].Params["discount"])
        } else {
            fmt.Printf("Customer is not eligible\n")
        }
    }
}
```

### 3. Custom Operators

```go
package main

import (
    "fmt"
    "github.com/tuannguyensn2001/go-json-rule-engine"
)

func main() {
    eng := go_json_rules_engine.NewEngine()

    // Register a custom operator
    err := eng.RegisterCustomOperator("divisibleBy", func(a, b interface{}) bool {
        // Implementation of divisibleBy operator
        return true
    })
    if err != nil {
        panic(err)
    }

    // Use the custom operator in a rule
    rules := []go_json_rules_engine.RuleOption{
        {
            ID:       "divisible-check",
            Name:     "Divisible Check",
            Priority: 1,
            Conditions: go_json_rules_engine.ConditionGroup{
                Operator: go_json_rules_engine.And,
                Conditions: []interface{}{
                    go_json_rules_engine.Condition{
                        Fact:     "number",
                        Operator: "divisibleBy",
                        Value:    5,
                    },
                },
            },
            Event: go_json_rules_engine.Event{
                Type: "divisible",
                Params: map[string]interface{}{
                    "message": "Number is divisible by 5",
                },
            },
        },
    }

    // Evaluate with test data
    facts := map[string]interface{}{
        "number": 10,
    }

    events, err := eng.Evaluate(rules, facts)
    if err != nil {
        panic(err)
    }

    for _, event := range events {
        fmt.Printf("Event: %s\n", event.Type)
        fmt.Printf("Message: %s\n", event.Params["message"])
    }
}
```

## Supported Operators

- `equal` / `notEqual`
- `greaterThan` / `lessThan`
- `greaterThanInclusive` / `lessThanInclusive`
- `in` / `notIn`
- `regex`
- `isNull` / `isNotNull`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 