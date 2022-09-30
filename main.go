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

	"github.com/LaChimere/proxy-rule-classifier/global"
	"github.com/LaChimere/proxy-rule-classifier/rule"
)

var (
	rulesDirPath string
	outputPath   string

	existedRules    = make(map[string]bool)
	classifiedRules = make(map[string][]string)

	// Special rules should be written at the end of the file.
	geoRule, finalRule *rule.Rule

	inputCount, outputCount int
)

func init() {
	flag.StringVar(&rulesDirPath, "i", "", "input rule directory")
	flag.StringVar(&outputPath, "o", "output", "output file path")
	flag.Parse()

	if rulesDirPath == "" {
		log.Fatalln("empty input rule directory")
	}
}

func main() {
	ruleFiles, err := os.ReadDir(rulesDirPath)
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range ruleFiles {
		filePath := fmt.Sprintf("%s/%s", rulesDirPath, file.Name())
		log.Printf("reading rule file: %s", filePath)

		if err := readRuleStrings(filePath); err != nil {
			log.Fatalln(err)
		}
	}

	log.Printf("%d rule(s) read", inputCount)

	if err := classifyRules(); err != nil {
		log.Fatalln(err)
	}

	if err := outputClassifiedRules(); err != nil {
		log.Fatalln(err)
	}

	log.Printf("%d rule(s) written", outputCount)
	log.Printf("the classified rules have been written into file: %s", outputPath)
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
		ru := string(line)
		line, _, err = reader.ReadLine()
		if _, ok := existedRules[ru]; ok {
			continue
		}

		// Disregard empty and comment lines.
		if ru == "" || ru == "\n" || strings.HasPrefix(ru, "#") {
			continue
		}

		inputCount++
		existedRules[ru] = true
	}

	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func classifyRules() error {
	for ruleStr := range existedRules {
		ru, err := rule.NewRuleFromString(ruleStr)
		if err != nil {
			return err
		}

		switch ru.Type {
		case rule.GEOIP:
			geoRule = ru
		case rule.FINAL:
			finalRule = ru
		default:
			classifiedRules[ru.Policy] = append(classifiedRules[ru.Policy], ru.String())
		}
	}

	return nil
}

func outputClassifiedRules() error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	for _, policy := range global.Policies {
		rules, ok := classifiedRules[policy]
		if !ok {
			continue
		}

		if _, err = writer.WriteString(fmt.Sprintf("# %s\n", policy)); err != nil {
			return err
		}

		sort.Strings(rules)
		for _, ru := range rules {
			if _, err = writer.WriteString(fmt.Sprintf("%s\n", ru)); err != nil {
				return err
			}
			outputCount++
		}
	}

	if geoRule != nil || finalRule != nil {
		if _, err = writer.WriteString("# Rules should be at last\n"); err != nil {
			return err
		}
	}
	if geoRule != nil {
		if _, err = writer.WriteString(fmt.Sprintf("%s\n", geoRule)); err != nil {
			return err
		}
		outputCount++
	}
	if finalRule != nil {
		if _, err = writer.WriteString(fmt.Sprintf("%s\n", finalRule)); err != nil {
			return err
		}
		outputCount++
	}

	return writer.Flush()
}
