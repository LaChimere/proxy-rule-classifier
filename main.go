package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"testing"
)

var (
	rulesDirPath string
	ruleFiles    []os.DirEntry

	existedRules    = make(map[string]bool)
	classifiedRules = make(map[string][]string)
)

const OUTPUT_FILENAME = "output"

func init() {
	flag.StringVar(&rulesDirPath, "rule", "", "Rule directory")
	testing.Init()
	flag.Parse()

	if rulesDirPath == "" {
		log.Fatalln("Empty rule directory path")
	}

	var err error
	ruleFiles, err = os.ReadDir(rulesDirPath)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	for _, file := range ruleFiles {
		filePath := fmt.Sprintf("%s/%s", rulesDirPath, file.Name())
		if err := readRules(filePath); err != nil {
			log.Fatalln(err.Error())
		}
	}

	classifyRules()
	if err := outputClassifiedRules(); err != nil {
		log.Fatalln(err.Error())
	}
}

func readRules(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, _, err := reader.ReadLine()
	for err == nil {
		rule := string(line)
		if rule != "" && rule != "\n" {
			existedRules[rule] = true
		}

		line, _, err = reader.ReadLine()
	}

	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func parseRuleComment(fullRule string) (rule, comment string) {
	commentPosition := strings.Index(fullRule, ` // `)
	if commentPosition == -1 {
		return fullRule, ""
	}

	return fullRule[:commentPosition], fullRule[commentPosition:]
}

func parseRuleFields(rule string) []string {
	return strings.Split(rule, ",")
}

func parseRulePolicy(rule string) string {
	fields := parseRuleFields(rule)
	return fields[len(fields)-1]
}

func classifyRules() {
	for rule := range existedRules {
		rule, comment := parseRuleComment(rule)
		rulePolicy := parseRulePolicy(rule)
		classifiedRules[rulePolicy] = append(classifiedRules[rulePolicy], rule+comment)
	}
}

func outputClassifiedRules() error {
	outputFile, err := os.Create(OUTPUT_FILENAME)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	rulePolicies := make([]string, 0, len(classifiedRules))
	for rulePolicy := range classifiedRules {
		rulePolicies = append(rulePolicies, rulePolicy)
	}
	sort.Strings(rulePolicies)

	writer := bufio.NewWriter(outputFile)
	for _, rulePolicy := range rulePolicies {
		if _, err = writer.WriteString(fmt.Sprintf("// %s\n", rulePolicy)); err != nil {
			return err
		}

		rules := classifiedRules[rulePolicy]
		sort.Strings(rules)
		for _, rule := range rules {
			if _, err = writer.WriteString(fmt.Sprintf("%s\n", rule)); err != nil {
				return err
			}
		}

		if _, err = writer.WriteString("\n"); err != nil {
			return err
		}
	}

	return nil
}
