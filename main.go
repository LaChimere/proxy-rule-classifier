package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"testing"

	rule2 "github.com/LaChimere/proxy-rule-classifier/rule"
)

var (
	rulesDirPath string
	ruleFiles    []os.DirEntry

	existedRules    = make(map[string]bool)
	classifiedRules = make(map[string][]string)

	// Special rules should be written at the end of the file.
	geoRule, finalRule *rule2.Rule
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
		if err := readRuleStrings(filePath); err != nil {
			log.Fatalln(err.Error())
		}
	}

	if err := classifyRules(); err != nil {
		log.Fatalln(err.Error())
	}

	if err := outputClassifiedRules(); err != nil {
		log.Fatalln(err.Error())
	}
}

func readRuleStrings(filename string) error {
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

func classifyRules() error {
	for ruleStr := range existedRules {
		rule, err := rule2.NewRuleFromString(ruleStr)
		if err != nil {
			return err
		}

		switch rule.Type {
		case rule2.GEOIP:
			geoRule = rule
		case rule2.FINAL:
			finalRule = rule
		default:
			classifiedRules[rule.Policy] = append(classifiedRules[rule.Policy], rule.String())
		}
	}

	return nil
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

	// TODO: write special rules at the end of the file.
	return writer.Flush()
}
