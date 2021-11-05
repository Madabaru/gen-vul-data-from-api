# gen-vulnerability-data-from-api
Generate vulnerability data from Github API

## Usage
Run ```go build && ./gen-vulnerability-data-from-api <Github Username> <AccessToken>```. It is strongly recommended to pass your Github username and access token to increase the number of permitted requests.

You can change the repositories and keywords as desired in ```data/```. By default, the program takes into account all commits that are not more than 1 month old from the current time.
