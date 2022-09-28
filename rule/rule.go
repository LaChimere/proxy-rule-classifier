package rule

import (
	"fmt"
	"strings"
)

type Rule struct {
	Type      RuleType
	Value     string
	Policy    string
	NoResolve bool
	Comment   string
}

func NewRule(_type RuleType, value string, policy string) *Rule {
	return &Rule{
		Type:      _type,
		Value:     value,
		Policy:    policy,
		NoResolve: false,
	}
}

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

func (rule *Rule) WithComment(comment string) *Rule {
	rule.Comment = comment
	return rule
}

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

func parseComment(lastField string) (field, comment string) {
	commentPosition := strings.Index(lastField, "//")
	if commentPosition == -1 {
		return lastField, ""
	}

	field = strings.TrimSpace(lastField[:commentPosition])
	comment = strings.TrimSpace(lastField[commentPosition+2:])
	return
}

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
