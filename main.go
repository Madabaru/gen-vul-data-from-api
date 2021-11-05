package main

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"
)

const SaveFile  = "data/dataset.json"
const KeywordFile = "data/keywords.txt"
const RepoFile = "data/repos.txt"

func main() {

	keywordList := ReadTextFile(KeywordFile)
	repoNameList := ReadTextFile(RepoFile)

	accessToken := ""
	username := ""
	authenticationString := ""

	if len(os.Args) == 3 {
		username = os.Args[1]
		accessToken = os.Args[2]
		authenticationString = username + ":" + accessToken + "@"
	}

	// Get query time (= 1 month ago)
	currentTime := time.Now()
	queryTime := currentTime.AddDate(0, -1, 0)
	queryDate := queryTime.Format("2006-01-02")

	// Final output list to store the generated vulnerability data
	outputList := make([]Output, 0)
	processedCommitList := make([]string, 0)

	for _, repoName := range repoNameList {
		for _, keywordName := range keywordList {

			url := "https://" + authenticationString + "api.github.com/search/commits?q=" + url.QueryEscape(keywordName) + "+repo:" + repoName + "+committer-date:>" + queryDate
			commitsByKeyword := GetCommitsByKeywordFromAPI(url)

			// Convert to list of strings
			cleanedCommitList := make([]string, 0)
			for _, commit := range commitsByKeyword.Commits {
				cleanedCommitList = append(cleanedCommitList, commit.Sha)
			}

			if len(cleanedCommitList) > 0 {

				// Keep track of all processed commits to avoid processing the same commit twice
				filteredCommitList := make([]string, 0)

				for _, commit := range cleanedCommitList {
					if !StringInList(commit, processedCommitList) {
						filteredCommitList = append(filteredCommitList, commit)
					}
				}

				for _, filteredCommit := range filteredCommitList {
					processedCommitList = append(processedCommitList, filteredCommit)

					url = "https://" + authenticationString + "api.github.com/repos/" + repoName + "/commits/" + filteredCommit
					commitDetails := GetCommitDetailsFromAPI(url)

					// Iterate over all changed files in the given commit
					for _, changedFile := range commitDetails.Files {

						if strings.HasSuffix(changedFile.Filename, ".c") {

							// Extract all functions names that are mentioned in the log details of a given commit
							extractedFunctionNameList := ExtractFuncNameFromTxt(changedFile.Patch)

							// Extract file right before the commit
							beforeUrl := "https://" + authenticationString + "api.github.com/repos/FFmpeg/FFmpeg/contents/" + changedFile.Filename + "?ref=" + filteredCommit + "~1"
							beforeFile := GetChangedFileFromAPI(beforeUrl)
							decodedBeforeFileContent := Base64Decode(beforeFile.Content)

							// Extract file after the commit
							afterUrl := "https://" + authenticationString + "api.github.com/repos/FFmpeg/FFmpeg/contents/" + changedFile.Filename + "?ref=" + filteredCommit
							afterFile := GetChangedFileFromAPI(afterUrl)
							decodedAfterFileContent := Base64Decode(afterFile.Content)

							// Extract the function's source code (if possible)
							extractedOutputBenign := ExtractFunctionFromFile(repoName, filteredCommit, decodedBeforeFileContent, extractedFunctionNameList, 0)
							if len(extractedOutputBenign) > 0 {
								outputList = append(outputList, extractedOutputBenign...)
							}
							// Extract the function's source code (if possible)
							extractedOutputMalicious := ExtractFunctionFromFile(repoName, filteredCommit, decodedAfterFileContent, extractedFunctionNameList, 1)
							if len(extractedOutputMalicious) > 0 {
								outputList = append(outputList, extractedOutputMalicious...)
							}
						}
					}
				}
			}
			file, _ := json.Marshal(outputList)
			_ = ioutil.WriteFile(SaveFile, file, 0644)
		}
	}
}
