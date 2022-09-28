package rule

import (
	"fmt"
	"log"
	"testing"
)

func TestNewRuleFromString(t *testing.T) {
	ruleStrings := []string{
		//"RULE-SET,https://s.trojanflare.com/Surge/Basic/Apple-proxy.list,DIRECT",
		//"DOMAIN-SUFFIX,disney-plus.net,Japan // Added for: prod-ripcut-delivery.disney-plus.net:443",
		//"DOMAIN-SUFFIX,co.uk,Japan // Added for: https://example.co.uk",
		//"IP-CIDR,91.108.8.0/22,Japan,no-resolve",
		"RULE-SET,https://s.trojanflare.com/Surge/Basic/common-ad-keyword.list,REJECT-TINYGIF",
		"FINAL,DIRECT,dns-failed",
	}

	for _, str := range ruleStrings {
		rule, err := NewRuleFromString(str)
		if err != nil {
			log.Println(err.Error())
		}

		fmt.Println(rule)
	}
}
