package main

type GithubWebhook struct {
	Ref         string `json:"ref"`
	Respository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`
}

type BitbucketWebhook struct {
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`
	Push struct {
		Changes []struct {
			New struct {
				BranchName string `json:"name"`
			} `json:"new"`
		} `json:"changes"`
	} `json:"push"`
}

type GitlabWebhook struct {
	Project struct {
		Name     string `json:"name"`
		FullName string `json:"path_with_namespace"`
	} `json:"project"`
	Ref string `json:"ref"`
}
