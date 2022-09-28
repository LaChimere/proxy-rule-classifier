package rule

type RuleType = string

const (
	DOMAIN           RuleType = "DOMAIN"
	DOMAIN_SUFFIX    RuleType = "DOMAIN-SUFFIX"
	DOMAIN_KEYWORD   RuleType = "DOMAIN-KEYWORD"
	IP_CIDR          RuleType = "IP-CIDR"
	IP_CIDR6         RuleType = "IP-CIDR6"
	GEOIP            RuleType = "GEOIP"
	USER_AGENT       RuleType = "USER_AGENT"
	URL_REGEX        RuleType = "URL-REGEX"
	SCRIPT           RuleType = "SCRIPT"
	PROTOCOL         RuleType = "PROTOCOL"
	DEST_PORT        RuleType = "DEST-PORT"
	DOMAIN_SET       RuleType = "DOMAIN-SET"
	SUBNET           RuleType = "SUBNET"
	CELLULAR_CARRIER RuleType = "CELLULAR-CARRIER"
	CELLULAR_RADIO   RuleType = "CELLULAR-RADIO"
	IP_ASN           RuleType = "IP-ASN"
	FINAL            RuleType = "FINAL"
)
