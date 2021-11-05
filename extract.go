package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ExtractFuncNameFromTxt Extract all functions names from the log details
func ExtractFuncNameFromTxt(fileText string) []string {

	reFindAllFunction := regexp.MustCompile(`@@ [a-zA-Z0-9\s._]*\(`)
	functionNameList := reFindAllFunction.FindAllString(fileText, -1)

	sanitizedFunctionNameList := make([]string, 0)
	for _, functionName := range functionNameList {
		partlySanitizedFunctionName := strings.ReplaceAll(functionName, "@@ ", "")
		sanitizedFunctionName := strings.ReplaceAll(partlySanitizedFunctionName, "(", "")
		sanitizedFunctionNameList = append(sanitizedFunctionNameList, sanitizedFunctionName)
	}

	// Remove duplicate strings in list
	uniqueSanitizedFunctionList := RemoveDuplicateStrings(sanitizedFunctionNameList)
	return uniqueSanitizedFunctionList
}

// ExtractFunctionFromFile Extract the function source code from a given file
func ExtractFunctionFromFile(repoName string, commit string, fileText string, functionList []string, label int) []Output {
	var outputList []Output
	for _, functionName := range functionList {
		reReplace := regexp.MustCompile(`/\*[^*]*\*+(?:[^/*][^*]*\*+)*/`)
		sanitizedText := reReplace.ReplaceAllString(fileText, "")

		reFindIndex := regexp.MustCompile(functionName + `\([\w\d\s,\*\n\r]*\)[\r\n\s]*\{`)
		startEndIndex := reFindIndex.FindStringIndex(sanitizedText)

		if len(startEndIndex) == 2 {
			start := startEndIndex[0]
			balanced := 0
			end := 0
			init := false

			// Find the final closing curly bracket in the given text
			for pos, char := range sanitizedText[start:] {
				if string(char) == "{" {
					init = true
					balanced += 1
				} else if string(char) == "}" && init == false {
					fmt.Println("error when trying to find the closing curly bracket")
					break
				} else if string(char) == "}" && init == true {
					balanced -= 1
				}

				if balanced == 0 && init == true {
					end = start + pos + 1
					sourceCode := sanitizedText[start:end]
					if len(strings.Split(functionName, " ")) <= 1 {
						sourceCode = "void " + sourceCode
					}
					output := Output{Project: repoName, Commit: commit, Label: strconv.Itoa(label), Code: sourceCode}
					outputList = append(outputList, output)
					break
				}
			}
		}
	}
	return outputList
}
