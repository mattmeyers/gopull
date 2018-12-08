package main

type GithubWebhook struct {
}

type BitbucketWebhook struct {
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
	Push struct {
		Changes []struct {
			New struct {
				Name string `json:"name"`
			} `json:"new"`
		} `json:"changes"`
	} `json:"push"`
}

type GitlabWebhook struct {
}
