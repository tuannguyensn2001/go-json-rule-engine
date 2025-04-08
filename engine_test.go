package go_json_rules_engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	assert.NotNil(t, engine, "Engine should not be nil")
	assert.Empty(t, engine.rules, "Rules should be empty")
	assert.Empty(t, engine.customOperators, "Custom operators should be empty")
}

func TestAddRule(t *testing.T) {
	engine := NewEngine()
	rule := Rule{
		Priority: 1,
		Conditions: ConditionGroup{
			Operator: And,
			Conditions: []interface{}{
				Condition{
					Fact:     "age",
					Operator: GreaterThan,
					Value:    18,
				},
			},
		},
		Event: Event{
			Type: "adult",
		},
	}

	err := engine.AddRule(rule)
	require.NoError(t, err, "Adding rule should not return error")
	assert.Len(t, engine.rules, 1, "Should have one rule")
	assert.Equal(t, rule, engine.rules[0], "Added rule should match")
}

func TestRegisterCustomOperator(t *testing.T) {
	engine := NewEngine()

	// Test successful registration
	err := engine.RegisterCustomOperator("is_even", func(a, b interface{}) bool {
		num, ok := a.(int)
		if !ok {
			return false
		}
		return num%2 == 0
	})
	require.NoError(t, err, "Registering custom operator should not return error")
	assert.Len(t, engine.customOperators, 1, "Should have one custom operator")

	// Test duplicate registration
	err = engine.RegisterCustomOperator("is_even", func(a, b interface{}) bool { return false })
	assert.Error(t, err, "Duplicate registration should return error")
}

func TestUnregisterCustomOperator(t *testing.T) {
	engine := NewEngine()

	// Register operator first
	err := engine.RegisterCustomOperator("is_even", func(a, b interface{}) bool { return false })
	require.NoError(t, err)

	// Test unregistering
	engine.UnregisterCustomOperator("is_even")
	assert.Len(t, engine.customOperators, 0, "Custom operators should be empty after unregistering")
}

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name     string
		rules    []Rule
		facts    map[string]interface{}
		expected []Event
	}{
		{
			name: "Simple age check - match",
			rules: []Rule{
				{
					Priority: 1,
					Conditions: ConditionGroup{
						Operator: And,
						Conditions: []interface{}{
							Condition{
								Fact:     "age",
								Operator: GreaterThan,
								Value:    18,
							},
						},
					},
					Event: Event{
						Type: "adult",
					},
				},
			},
			facts: map[string]interface{}{
				"age": 20,
			},
			expected: []Event{
				{
					Type: "adult",
				},
			},
		},
		{
			name: "Multiple conditions - match",
			rules: []Rule{
				{
					Priority: 1,
					Conditions: ConditionGroup{
						Operator: And,
						Conditions: []interface{}{
							Condition{
								Fact:     "age",
								Operator: GreaterThan,
								Value:    18,
							},
							Condition{
								Fact:     "country",
								Operator: Equal,
								Value:    "US",
							},
						},
					},
					Event: Event{
						Type: "us_adult",
					},
				},
			},
			facts: map[string]interface{}{
				"age":     20,
				"country": "US",
			},
			expected: []Event{
				{
					Type: "us_adult",
				},
			},
		},
		{
			name: "No matching conditions",
			rules: []Rule{
				{
					Priority: 1,
					Conditions: ConditionGroup{
						Operator: And,
						Conditions: []interface{}{
							Condition{
								Fact:     "age",
								Operator: GreaterThan,
								Value:    18,
							},
						},
					},
					Event: Event{
						Type: "adult",
					},
				},
			},
			facts: map[string]interface{}{
				"age": 15,
			},
			expected: nil,
		},
		{
			name: "Nested condition groups",
			rules: []Rule{
				{
					Priority: 1,
					Conditions: ConditionGroup{
						Operator: And,
						Conditions: []interface{}{
							ConditionGroup{
								Operator: Or,
								Conditions: []interface{}{
									Condition{
										Fact:     "age",
										Operator: GreaterThan,
										Value:    18,
									},
									Condition{
										Fact:     "country",
										Operator: Equal,
										Value:    "US",
									},
								},
							},
						},
					},
					Event: Event{
						Type: "match",
					},
				},
			},
			facts: map[string]interface{}{
				"age":     15,
				"country": "US",
			},
			expected: []Event{
				{
					Type: "match",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewEngine()
			for _, rule := range tt.rules {
				err := engine.AddRule(rule)
				require.NoError(t, err)
			}

			events, err := engine.Evaluate(tt.facts)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, events)
		})
	}
}

func TestLoadAndSaveRules(t *testing.T) {
	engine := NewEngine()

	// Test loading rules from JSON string
	jsonStr := `[
		{
			"priority": 1,
			"conditions": {
				"operator": "and",
				"conditions": [
					{
						"fact": "age",
						"operator": "gt",
						"value": 18
					}
				]
			},
			"event": {
				"type": "adult"
			}
		}
	]`

	err := engine.LoadRulesFromJSONString(jsonStr)
	require.NoError(t, err, "Loading rules from JSON string should not return error")
	assert.Len(t, engine.rules, 1, "Should have one rule after loading")

	// Test saving rules to JSON
	err = engine.SaveRulesToJSON("test_rules.json")
	require.NoError(t, err, "Saving rules to JSON should not return error")

	// Test loading rules from file
	engine = NewEngine()
	err = engine.LoadRulesFromJSON("test_rules.json")
	require.NoError(t, err, "Loading rules from file should not return error")
	assert.Len(t, engine.rules, 1, "Should have one rule after loading from file")
}

func TestCompareOperators(t *testing.T) {
	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		operator Operator
		expected bool
	}{
		{
			name:     "Equal - numbers",
			a:        42,
			b:        42,
			operator: Equal,
			expected: true,
		},
		{
			name:     "Equal - strings",
			a:        "hello",
			b:        "hello",
			operator: Equal,
			expected: true,
		},
		{
			name:     "NotEqual - numbers",
			a:        42,
			b:        43,
			operator: NotEqual,
			expected: true,
		},
		{
			name:     "GreaterThan - numbers",
			a:        43,
			b:        42,
			operator: GreaterThan,
			expected: true,
		},
		{
			name:     "LessThan - numbers",
			a:        41,
			b:        42,
			operator: LessThan,
			expected: true,
		},
		{
			name:     "In - string in slice",
			a:        "apple",
			b:        []interface{}{"banana", "apple", "orange"},
			operator: In,
			expected: true,
		},
		{
			name:     "NotIn - string not in slice",
			a:        "grape",
			b:        []interface{}{"banana", "apple", "orange"},
			operator: NotIn,
			expected: true,
		},
		{
			name:     "Regex - matching pattern",
			a:        "hello123",
			b:        "^hello\\d+$",
			operator: Regex,
			expected: true,
		},
		{
			name:     "IsNull - nil value",
			a:        nil,
			b:        nil,
			operator: IsNull,
			expected: true,
		},
		{
			name:     "IsNotNull - non-nil value",
			a:        "not null",
			b:        nil,
			operator: IsNotNull,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := Rule{
				Priority: 1,
				Conditions: ConditionGroup{
					Operator: And,
					Conditions: []interface{}{
						Condition{
							Fact:     "test",
							Operator: tt.operator,
							Value:    tt.b,
						},
					},
				},
				Event: Event{
					Type: "test",
				},
			}

			engine := NewEngine()
			err := engine.AddRule(rule)
			require.NoError(t, err)

			events, err := engine.Evaluate(map[string]interface{}{
				"test": tt.a,
			})
			require.NoError(t, err)

			if tt.expected {
				assert.Len(t, events, 1, "Expected one matching event")
			} else {
				assert.Empty(t, events, "Expected no matching events")
			}
		})
	}
}

func TestCustomOperatorIntegration(t *testing.T) {
	engine := NewEngine()

	// Register a custom operator that checks if a number is even
	err := engine.RegisterCustomOperator("is_even", func(a, b interface{}) bool {
		num, ok := a.(int)
		if !ok {
			return false
		}
		return num%2 == 0
	})
	require.NoError(t, err)

	rule := Rule{
		Priority: 1,
		Conditions: ConditionGroup{
			Operator: And,
			Conditions: []interface{}{
				Condition{
					Fact:     "number",
					Operator: "is_even",
					Value:    nil,
				},
			},
		},
		Event: Event{
			Type: "even_number",
		},
	}

	err = engine.AddRule(rule)
	require.NoError(t, err)

	// Test with even number
	events, err := engine.Evaluate(map[string]interface{}{
		"number": 4,
	})
	require.NoError(t, err)
	assert.Len(t, events, 1, "Expected one event for even number")
	assert.Equal(t, "even_number", events[0].Type)

	// Test with odd number
	events, err = engine.Evaluate(map[string]interface{}{
		"number": 3,
	})
	require.NoError(t, err)
	assert.Empty(t, events, "Expected no events for odd number")
}

func TestCompareFunctions(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		compare  func(a, b interface{}) bool
		expected bool
	}{
		{
			name:     "compareEqual - numbers equal",
			a:        42,
			b:        42,
			compare:  engine.compareEqual,
			expected: true,
		},
		{
			name:     "compareEqual - numbers not equal",
			a:        42,
			b:        43,
			compare:  engine.compareEqual,
			expected: false,
		},
		{
			name:     "compareEqual - strings equal",
			a:        "hello",
			b:        "hello",
			compare:  engine.compareEqual,
			expected: true,
		},
		{
			name:     "compareEqual - strings not equal",
			a:        "hello",
			b:        "world",
			compare:  engine.compareEqual,
			expected: false,
		},
		{
			name:     "compareEqual - nil values",
			a:        nil,
			b:        nil,
			compare:  engine.compareEqual,
			expected: true,
		},
		{
			name:     "compareEqual - one nil value",
			a:        nil,
			b:        42,
			compare:  engine.compareEqual,
			expected: false,
		},
		{
			name:     "compareGreaterThan - numbers greater",
			a:        43,
			b:        42,
			compare:  engine.compareGreaterThan,
			expected: true,
		},
		{
			name:     "compareGreaterThan - numbers equal",
			a:        42,
			b:        42,
			compare:  engine.compareGreaterThan,
			expected: false,
		},
		{
			name:     "compareGreaterThan - numbers less",
			a:        41,
			b:        42,
			compare:  engine.compareGreaterThan,
			expected: false,
		},
		{
			name:     "compareLessThan - numbers less",
			a:        41,
			b:        42,
			compare:  engine.compareLessThan,
			expected: true,
		},
		{
			name:     "compareLessThan - numbers equal",
			a:        42,
			b:        42,
			compare:  engine.compareLessThan,
			expected: false,
		},
		{
			name:     "compareLessThan - numbers greater",
			a:        43,
			b:        42,
			compare:  engine.compareLessThan,
			expected: false,
		},
		{
			name:     "compareGreaterThanOrEqual - numbers greater",
			a:        43,
			b:        42,
			compare:  engine.compareGreaterThanOrEqual,
			expected: true,
		},
		{
			name:     "compareGreaterThanOrEqual - numbers equal",
			a:        42,
			b:        42,
			compare:  engine.compareGreaterThanOrEqual,
			expected: true,
		},
		{
			name:     "compareGreaterThanOrEqual - numbers less",
			a:        41,
			b:        42,
			compare:  engine.compareGreaterThanOrEqual,
			expected: false,
		},
		{
			name:     "compareLessThanOrEqual - numbers less",
			a:        41,
			b:        42,
			compare:  engine.compareLessThanOrEqual,
			expected: true,
		},
		{
			name:     "compareLessThanOrEqual - numbers equal",
			a:        42,
			b:        42,
			compare:  engine.compareLessThanOrEqual,
			expected: true,
		},
		{
			name:     "compareLessThanOrEqual - numbers greater",
			a:        43,
			b:        42,
			compare:  engine.compareLessThanOrEqual,
			expected: false,
		},
		{
			name:     "evaluateIn - string in slice",
			a:        "apple",
			b:        []interface{}{"banana", "apple", "orange"},
			compare:  engine.evaluateIn,
			expected: true,
		},
		{
			name:     "evaluateIn - string not in slice",
			a:        "grape",
			b:        []interface{}{"banana", "apple", "orange"},
			compare:  engine.evaluateIn,
			expected: false,
		},
		{
			name:     "evaluateIn - number in slice",
			a:        42,
			b:        []interface{}{40, 41, 42, 43},
			compare:  engine.evaluateIn,
			expected: true,
		},
		{
			name:     "evaluateIn - number not in slice",
			a:        44,
			b:        []interface{}{40, 41, 42, 43},
			compare:  engine.evaluateIn,
			expected: false,
		},
		{
			name:     "evaluateRegex - matching pattern",
			a:        "hello123",
			b:        "^hello\\d+$",
			compare:  engine.evaluateRegex,
			expected: true,
		},
		{
			name:     "evaluateRegex - non-matching pattern",
			a:        "hello123",
			b:        "^world\\d+$",
			compare:  engine.evaluateRegex,
			expected: false,
		},
		{
			name:     "evaluateRegex - invalid pattern",
			a:        "hello123",
			b:        "[",
			compare:  engine.evaluateRegex,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.compare(tt.a, tt.b)
			assert.Equal(t, tt.expected, result, "Comparison result mismatch")
		})
	}
}
