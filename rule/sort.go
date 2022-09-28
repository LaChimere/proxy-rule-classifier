package rule

import "sort"

type ruleSlice []*Rule

func (x ruleSlice) Len() int {
	return len(x)
}

func (x ruleSlice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x ruleSlice) Less(i, j int) bool {
	return x[i].String() < x[j].String()
}

func Sort(x []*Rule) {
	sort.Sort(ruleSlice(x))
}
