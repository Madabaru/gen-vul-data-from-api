package main

// CommitsByKeyword Struct to save all commits that contain the given keyword from github API
type CommitsByKeyword struct {
	TotalCount int `json:"total_count"`
	Commits    []struct {
	Sha string `json:"sha"`
	} `json:"items"`
}

// CommitDetails Struct to save the commit details retrieved from github API
type CommitDetails struct {
	Sha   string `json:"sha"`
	Files []struct {
	Filename string `json:"filename"`
	Patch    string `json:"patch"`
	} `json:"files"`
}

// ChangedFile Struct to save the changed file details from github API
type ChangedFile struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

// Output Struct to save the retrieved vulnerability data
type Output struct {
	Project string
	Commit  string
	Label   string
	Code    string
}