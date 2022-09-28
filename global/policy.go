package global

import "github.com/LaChimere/proxy-rule-classifier/rule"

var Policies = []rule.PolicyType{
	rule.JAPAN,
	rule.HONGKONG,
	rule.US,
	rule.AUSTRALIA,
	rule.DIRECT,
	rule.REJECT,
	rule.REJECT_TINYGIF,
}
