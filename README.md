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

## Usage Examples

### 1. Basic Rule Evaluation

```go
package main

import (
    "fmt"
    "github.com/tuannguyensn2001/go-json-rule-engine/pkg/engine"
)

func main() {
    // Create a new engine
    eng := engine.NewEngine()

    // Define a simple rule: If age > 18, then adult
    rule := `{
        "id": "age-check",
        "name": "Adult Check",
        "priority": 1,
        "conditions": {
            "operator": "and",
            "conditions": [
                {
                    "fact": "age",
                    "operator": "greaterThan",
                    "value": 18
                }
            ]
        },
        "event": {
            "type": "adult",
            "params": {
                "message": "User is an adult"
            }
        }
    }`

    // Load the rule
    if err := eng.LoadRulesFromJSONString(rule); err != nil {
        panic(err)
    }

    // Evaluate facts
    facts := map[string]interface{}{
        "age": 25,
    }

    events, err := eng.Evaluate(facts)
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
    "github.com/tuannguyensn2001/go-json-rule-engine/pkg/engine"
)

func main() {
    eng := engine.NewEngine()

    // Define a complex rule for customer eligibility
    rule := `{
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
    }`

    if err := eng.LoadRulesFromJSONString(rule); err != nil {
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
    }

    for i, facts := range testCases {
        fmt.Printf("\nTesting Case %d:\n", i+1)
        events, err := eng.Evaluate(facts)
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

### 3. Loading Rules from File

```go
package main

import (
    "fmt"
    "github.com/tuannguyensn2001/go-json-rule-engine/pkg/engine"
)

func main() {
    eng := engine.NewEngine()

    // Load rules from a JSON file
    if err := eng.LoadRulesFromJSON("rules.json"); err != nil {
        panic(err)
    }

    // Save rules to a file
    if err := eng.SaveRulesToJSON("rules_backup.json"); err != nil {
        panic(err)
    }
}
```

### 4. Custom Operators

```go
package main

import (
    "fmt"
    "github.com/tuannguyensn2001/go-json-rule-engine/pkg/engine"
)

func main() {
    eng := engine.NewEngine()

    // Register a custom operator
    err := eng.RegisterCustomOperator("divisibleBy", func(a, b interface{}) bool {
        // Implementation of divisibleBy operator
        return true
    })
    if err != nil {
        panic(err)
    }

    // Use the custom operator in a rule
    rule := `{
        "id": "divisible-check",
        "name": "Divisible Check",
        "priority": 1,
        "conditions": {
            "operator": "and",
            "conditions": [
                {
                    "fact": "number",
                    "operator": "divisibleBy",
                    "value": 5
                }
            ]
        },
        "event": {
            "type": "divisible",
            "params": {
                "message": "Number is divisible by 5"
            }
        }
    }`

    if err := eng.LoadRulesFromJSONString(rule); err != nil {
        panic(err)
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