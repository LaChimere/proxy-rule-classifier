# Proxy Rule Classifier

This is a simple program to classify proxy rules into their corresponding types, and then write the classified rules to the result file within a specified order.

## Specified Order

The order of output rules classified by the type of rule policies is specified in `global/policy.go`.

## Input Format

Examples of valid rule strings:

```
DOMAIN-SUFFIX,microsoft.com,Japan // Here's the comment.
DOMAIN-KEYWORD,google,Hong Kong
IP-CIDR,100.100.100.100/22,US,no-resolve
FINAL,DIRECT,dns-failed // FINAL is special as it has no value field.
```

The valid rule string may be split into three or four fields by the separator `,`.

## Use Guidance

### Build

```bash
git clone https://github.com/LaChimere/proxy-rule-classifier.git
cd proxy-rule-classifier
go build main.go
```

### Classify

Currently, the program only supports inputting a directory path where includes rule files. The output file include all unique rules from these files.

For instance, after building the project, we can specify `./input` as an input directory.

```
.
├───<other files>
├───input
|   ├───a.rule
│   └───b.rule
├───main
└───main.go
```

Run the command below will automatically write the result to a newly created file `./output`.

```bash
./main -i input
```

The output file path can be customized by using the argument `-o`.

```bash
./main -i input -o classfied_rules
```
