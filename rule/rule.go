package rule

import (
	"fmt"
	"strings"

	"github.com/LaChimere/proxy-rule-classifier/pkg/errors"
)

// Rule refers to the proxy rule.
// The Type and Policy should not be empty.
// If Type is not FINAL, the Value should not be empty either.
type Rule struct {
	Type      RuleType
	Value     string
	Policy    PolicyType
	NoResolve bool
	Comment   string
}

// NewRule creates a Rule entity with the type, value and policy of a common proxy rule.
func NewRule(_type RuleType, value string, policy PolicyType) (*Rule, error) {
	if !validPolicy(policy) {
		return nil, errors.InvalidPolicy
	}

	return &Rule{
		Type:      _type,
		Value:     value,
		Policy:    policy,
		NoResolve: false,
	}, nil
}

// NewRuleFromString creates a Rule entity by the input string.
func NewRuleFromString(str string) (rule *Rule, err error) {
	fields, err := parseFields(str)
	if err != nil {
		return nil, err
	}

	fieldNum := len(fields)
	lastField, comment := parseComment(fields[fieldNum-1])

	var (
		_type  RuleType
		value  string
		policy PolicyType
	)

	// Special treatment for FINAL.
	if fields[0] == FINAL {
		policy = fields[1]
		if !validPolicy(policy) {
			return nil, errors.InvalidPolicy
		}
		return &Rule{Type: FINAL, Policy: policy, Comment: comment}, nil
	}

	_type, value = fields[0], fields[1]
	rule = &Rule{Type: _type, Value: value, Comment: comment}
	if fieldNum == 3 {
		policy = lastField
	} else {
		policy = fields[2]
		rule.NoResolve = true
	}
	rule.Policy = policy

	if !validPolicy(rule.Policy) {
		return nil, errors.InvalidPolicy
	}
	return rule, nil
}

// WithComment adds a comment to a Rule entity.
func (rule *Rule) WithComment(comment string) *Rule {
	rule.Comment = comment
	return rule
}

// String converts the Rule objects to the corresponding strings the same format as those in proxy config files.
func (rule *Rule) String() string {
	if rule.Type == FINAL {
		return fmt.Sprintf("%s,%s,dns-failed", rule.Type, rule.Policy)
	}

	s := fmt.Sprintf("%s,%s,%s", rule.Type, rule.Value, rule.Policy)
	if rule.NoResolve {
		s += ",no-resolve"
	}

	if rule.Comment != "" {
		s += fmt.Sprintf(" // %s", rule.Comment)
	}
	return s
}

// parseComment processes the last field of the split string to get the real last field and the comment.
func parseComment(lastField string) (field, comment string) {
	commentPosition := strings.Index(lastField, "//")
	if commentPosition == -1 {
		return lastField, ""
	}

	field = strings.TrimSpace(lastField[:commentPosition])
	comment = strings.TrimSpace(lastField[commentPosition+2:])
	return
}

// parseFields splits the input string by ",", and checks if the number of fields is valid.
func parseFields(str string) (fields []string, err error) {
	fields = strings.Split(str, ",")
	fieldNum := len(fields)

	if fieldNum < 3 {
		return nil, errors.InsufficientFields
	} else if fieldNum > 4 {
		return nil, errors.TooManyFields
	}

	for _, field := range fields {
		field = strings.TrimSpace(field)
	}

	return fields, nil
}

func validPolicy(policy PolicyType) bool {
	for _, p := range ValidPolicies {
		if policy == p {
			return true
		}
	}
	return false
}
