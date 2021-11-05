package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func StringInList(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func RemoveDuplicateStrings(stringList []string) []string {
	keys := make(map[string]bool)
	list := make([]string, 0)
	for _, item := range stringList {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func getFromAPI(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("error: " + err.Error())
	} else if response.Status == "403 Forbidden" {
		fmt.Println("error: API rate limit exceeded. waiting for 60 min.")
		time.Sleep(60 * time.Minute)
		response, _ = http.Get(url)
		if response.Status == "403 Forbidden" {
			fmt.Println("error: API rate limit exceeded. try again later.")
			os.Exit(1)
		}
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	return body
}

func GetCommitsByKeywordFromAPI(url string) CommitsByKeyword {
	commitsByKeyword := CommitsByKeyword{}
	body := getFromAPI(url)
	err := json.Unmarshal(body, &commitsByKeyword)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	return commitsByKeyword
}

func GetCommitDetailsFromAPI(url string) CommitDetails {
	commitDetails := CommitDetails{}
	body := getFromAPI(url)

	err := json.Unmarshal(body, &commitDetails)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	return commitDetails
}

func GetChangedFileFromAPI(url string) ChangedFile {
	changedFile := ChangedFile{}
	body := getFromAPI(url)
	err := json.Unmarshal(body, &changedFile)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	return changedFile
}

func Base64Decode(str string) string {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	return string(data)
}

func RemoveEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func ReadTextFile(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer func() {
		if err = file.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}()
	content, err := ioutil.ReadAll(file)
	sanitizedContent := strings.ReplaceAll(string(content), "\r", "")
	extractedKeywords := strings.Split(sanitizedContent, "\n")
	return RemoveEmptyStrings(extractedKeywords)
}