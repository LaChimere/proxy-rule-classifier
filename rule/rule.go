package rule

import (
	"fmt"
	"strings"
)

// Rule refers to the proxy rule.
// The Type and Policy should not be empty.
// If Type is not FINAL, the Value should not be empty either.
type Rule struct {
	Type      RuleType
	Value     string
	Policy    string
	NoResolve bool
	Comment   string
}

// NewRule creates a Rule entity with the type, value and policy of a common proxy rule.
func NewRule(_type RuleType, value string, policy string) *Rule {
	return &Rule{
		Type:      _type,
		Value:     value,
		Policy:    policy,
		NoResolve: false,
	}
}

// NewRuleFromString creates a Rule entity by the input string.
func NewRuleFromString(str string) (rule *Rule, err error) {
	fields, err := parseFields(str)
	if err != nil {
		return nil, err
	}

	fieldNum := len(fields)
	lastField, comment := parseComment(fields[fieldNum-1])

	// Special treatment for FINAL.
	if fields[0] == FINAL {
		return &Rule{Type: FINAL, Policy: fields[1], Comment: comment}, nil
	}

	if fieldNum == 3 {
		rule = NewRule(fields[0], fields[1], lastField).WithComment(comment)
	} else {
		rule = NewRule(fields[0], fields[1], fields[2]).WithComment(comment)
		rule.NoResolve = true
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

// parseComment processes the last field of the splitted string to get the real last field and the comment.
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
		return nil, NewError("insufficient fields to resolve a Rule entity")
	} else if fieldNum > 4 {
		return nil, NewError("too many fields to resolve a Rule entity")
	}

	for _, field := range fields {
		field = strings.TrimSpace(field)
	}

	return fields, nil
}
