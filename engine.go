package go_json_rules_engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"sort"
	"sync"
)

type Engine struct {
	rules           []Rule
	customOperators map[Operator]CustomOperatorFunc
	mu              sync.RWMutex
}

type CustomOperatorFunc func(a, b interface{}) bool

func NewEngine() *Engine {
	return &Engine{
		rules:           make([]Rule, 0),
		customOperators: make(map[Operator]CustomOperatorFunc),
	}
}

func (e *Engine) RegisterCustomOperator(op Operator, fn CustomOperatorFunc) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.customOperators[op]; exists {
		return fmt.Errorf("operator %s is already registered", op)
	}

	e.customOperators[op] = fn
	return nil
}

func (e *Engine) UnregisterCustomOperator(op Operator) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.customOperators, op)
}

func (e *Engine) AddRule(rule Rule) error {
	e.rules = append(e.rules, rule)
	e.sortRulesByPriority()
	return nil
}

func (e *Engine) LoadRulesFromJSON(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read rules file: %w", err)
	}

	return e.LoadRulesFromJSONString(string(data))
}

func (e *Engine) LoadRulesFromJSONString(jsonStr string) error {
	var rules []Rule
	if err := json.Unmarshal([]byte(jsonStr), &rules); err != nil {
		return fmt.Errorf("failed to parse rules: %w", err)
	}

	e.rules = rules
	e.sortRulesByPriority()
	return nil
}

func (e *Engine) SaveRulesToJSON(filename string) error {
	data, err := json.MarshalIndent(e.rules, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal rules: %w", err)
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write rules file: %w", err)
	}

	return nil
}

func (e *Engine) sortRulesByPriority() {
	sort.Slice(e.rules, func(i, j int) bool {
		return e.rules[i].Priority > e.rules[j].Priority
	})
}

func (e *Engine) Evaluate(facts map[string]interface{}) ([]Event, error) {
	var events []Event

	for _, rule := range e.rules {
		if e.evaluateConditionGroup(rule.Conditions, facts) {
			events = append(events, rule.Event)
		}
	}

	return events, nil
}

func (e *Engine) evaluateConditionGroup(group ConditionGroup, facts map[string]interface{}) bool {
	if len(group.Conditions) == 0 {
		return true
	}

	switch group.Operator {
	case And:
		for _, condition := range group.Conditions {
			switch cond := condition.(type) {
			case Condition:
				if !e.evaluateCondition(cond, facts) {
					return false
				}
			case ConditionGroup:
				if !e.evaluateConditionGroup(cond, facts) {
					return false
				}
			}
		}
		return true

	case Or:
		for _, condition := range group.Conditions {
			switch cond := condition.(type) {
			case Condition:
				if e.evaluateCondition(cond, facts) {
					return true
				}
			case ConditionGroup:
				if e.evaluateConditionGroup(cond, facts) {
					return true
				}
			}
		}
		return false

	default:
		return false
	}
}

func (e *Engine) evaluateCondition(condition Condition, facts map[string]interface{}) bool {
	factValue, exists := facts[condition.Fact]
	if !exists {
		return false
	}

	// Check for custom operator first
	e.mu.RLock()
	customFn, isCustom := e.customOperators[condition.Operator]
	e.mu.RUnlock()

	if isCustom {
		return customFn(factValue, condition.Value)
	}

	// Handle built-in operators
	switch condition.Operator {
	case Equal:
		return e.compareEqual(factValue, condition.Value)
	case NotEqual:
		return !e.compareEqual(factValue, condition.Value)
	case GreaterThan:
		return e.compareGreaterThan(factValue, condition.Value)
	case LessThan:
		return e.compareLessThan(factValue, condition.Value)
	case GreaterThanInc:
		return e.compareGreaterThanOrEqual(factValue, condition.Value)
	case LessThanInc:
		return e.compareLessThanOrEqual(factValue, condition.Value)
	case In:
		return e.evaluateIn(factValue, condition.Value)
	case NotIn:
		return !e.evaluateIn(factValue, condition.Value)
	case Regex:
		return e.evaluateRegex(factValue, condition.Value)
	case IsNull:
		return factValue == nil
	case IsNotNull:
		return factValue != nil
	default:
		return false
	}
}

func (e *Engine) compareEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	// Handle numeric types
	if e.IsNumeric(va) && e.IsNumeric(vb) {
		return e.compareNumeric(va, vb) == 0
	}

	// Handle string types
	if va.Kind() == reflect.String && vb.Kind() == reflect.String {
		return va.String() == vb.String()
	}

	// Handle boolean types
	if va.Kind() == reflect.Bool && vb.Kind() == reflect.Bool {
		return va.Bool() == vb.Bool()
	}

	// Use DeepEqual for other types
	return reflect.DeepEqual(a, b)
}

func (e *Engine) compareGreaterThan(a, b interface{}) bool {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	// Only handle numeric types
	if e.IsNumeric(va) && e.IsNumeric(vb) {
		return e.compareNumeric(va, vb) > 0
	}

	return false
}

func (e *Engine) compareLessThan(a, b interface{}) bool {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	// Only handle numeric types
	if e.IsNumeric(va) && e.IsNumeric(vb) {
		return e.compareNumeric(va, vb) < 0
	}

	return false
}

func (e *Engine) compareGreaterThanOrEqual(a, b interface{}) bool {
	return e.compareGreaterThan(a, b) || e.compareEqual(a, b)
}

func (e *Engine) compareLessThanOrEqual(a, b interface{}) bool {
	return e.compareLessThan(a, b) || e.compareEqual(a, b)
}

func (e *Engine) evaluateIn(a, b interface{}) bool {
	slice, ok := b.([]interface{})
	if !ok {
		return false
	}

	for _, item := range slice {
		if e.compareEqual(a, item) {
			return true
		}
	}
	return false
}

func (e *Engine) evaluateRegex(a, pattern interface{}) bool {
	str, ok := a.(string)
	if !ok {
		return false
	}

	patternStr, ok := pattern.(string)
	if !ok {
		return false
	}

	matched, err := regexp.MatchString(patternStr, str)
	if err != nil {
		return false
	}
	return matched
}

// IsNumeric checks if a value is any numeric type
func (e *Engine) IsNumeric(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// ToFloat64 converts any numeric type to float64
func (e *Engine) ToFloat64(v reflect.Value) float64 {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return v.Float()
	default:
		return 0
	}
}

func (e *Engine) compareNumeric(a, b reflect.Value) int {
	// Convert both values to float64 for comparison
	aFloat := e.ToFloat64(a)
	bFloat := e.ToFloat64(b)

	if aFloat < bFloat {
		return -1
	}
	if aFloat > bFloat {
		return 1
	}
	return 0
}
