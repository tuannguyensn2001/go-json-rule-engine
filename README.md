# Go JSON Rules Engine

A flexible and powerful rules engine for Go that allows you to define business rules using JSON configuration.

[![Go Report Card](https://goreportcard.com/badge/github.com/tuannguyensn2001/go-json-rule-engine)](https://goreportcard.com/report/github.com/tuannguyensn2001/go-json-rule-engine)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/tuannguyensn2001/go-json-rule-engine.svg)](https://pkg.go.dev/github.com/tuannguyensn2001/go-json-rule-engine)

## Features

- Define complex business rules using JSON
- Support for nested logical conditions (AND/OR)
- Multiple comparison operators
- Priority-based rule evaluation
- Event-driven results with custom parameters
- Easy integration with existing Go applications

## Installation

```bash
go get github.com/tuannguyensn2001/go-json-rule-engine
```

## Quick Start

Here's a simple example showing how to use the rules engine:

```go
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
            "id": "premium-customer",
            "name": "Premium Customer Rule",
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
                        "fact": "yearlyPurchases",
                        "operator": "greaterThan",
                        "value": 1000.0
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
    ]`
    
    if err := rules.LoadRulesFromJSONString(jsonStr); err != nil {
        panic(err)
    }

    // Define facts to evaluate
    facts := map[string]interface{}{
        "age":             25,
        "yearlyPurchases": 1200.0,
    }

    // Evaluate rules
    events, err := eng.Evaluate(rules, facts)
    if err != nil {
        panic(err)
    }

    // Handle results
    for _, event := range events {
        fmt.Printf("Rule triggered: %s\n", event.Type)
        fmt.Printf("Message: %v\n", event.Params["message"])
        fmt.Printf("Discount: %v%%\n", event.Params["discount"])
    }
}
```

## Rule Structure

Rules are defined in JSON with the following structure:

```json
{
    "id": "rule-id",
    "name": "Rule Name",
    "priority": 10,
    "conditions": {
        "operator": "and|or",
        "conditions": [
            {
                "fact": "factName",
                "operator": "operatorType",
                "value": "compareValue"
            }
        ]
    },
    "event": {
        "type": "eventType",
        "params": {
            "key": "value"
        }
    }
}
```

## Supported Operators

The engine supports these comparison operators:

- `equal` - Exact match
- `notEqual` - Not equal to
- `greaterThan` - Greater than
- `lessThan` - Less than
- `greaterThanInclusive` - Greater than or equal to
- `lessThanInclusive` - Less than or equal to
- `in` - Value exists in array
- `notIn` - Value does not exist in array
- `regex` - Matches regular expression
- `isNull` - Value is null
- `isNotNull` - Value is not null

## Advanced Features

### Priority-based Evaluation

Rules are evaluated in order of priority (higher numbers first). This allows you to control the execution order of your rules:

```json
[
    {
        "id": "high-priority",
        "priority": 100,
        "conditions": { ... }
    },
    {
        "id": "low-priority",
        "priority": 1,
        "conditions": { ... }
    }
]
```

### Nested Conditions

You can create complex rule conditions by nesting AND/OR operators:

```json
{
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
    }
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 