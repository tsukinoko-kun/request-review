package linear

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tsukinoko-kun/request-review/internal/config"
)

type (
	GraphQLResponse struct {
		Data struct {
			Issues struct {
				Nodes []Issue `json:"nodes"`
			} `json:"issues"`
		} `json:"data"`
	}
	Issue struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Project     struct {
			Name string `json:"name"`
		} `json:"project"`
	}
)

var (
	ErrNoIssueFound        = fmt.Errorf("no issue found")
	ErrMultipleIssuesFound = fmt.Errorf("multiple issues found")
	ErrApiKeyNotSet        = fmt.Errorf("API key not set")
)

func FindIssueByBranchName(cfg config.Config, branchName string) (Issue, error) {
	if cfg.LinearPersonalApiKey == "" {
		return Issue{}, ErrApiKeyNotSet
	}

	query := fmt.Sprintf(`
query {
  	issues(filter: {
  		branchName: { eq: "%s" }
    }) {
   		nodes {
   			id
      		title
      		description
        	project {
                name
            }
        }
    }
}
`, branchName)

	jsonValue, _ := json.Marshal(map[string]string{
		"query": query,
	})

	req, err := http.NewRequest("POST", "https://api.linear.app/graphql", bytes.NewBuffer(jsonValue))
	if err != nil {
		return Issue{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.LinearPersonalApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Issue{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= http.StatusBadRequest {
		return Issue{}, fmt.Errorf("failed to query issue: %s", string(body))
	}

	var graphQLResponse GraphQLResponse
	err = json.Unmarshal(body, &graphQLResponse)
	if err != nil {
		return Issue{}, fmt.Errorf("failed to parse response: %w", err)
	}

	switch len(graphQLResponse.Data.Issues.Nodes) {
	case 0:
		return Issue{}, ErrNoIssueFound
	case 1:
		return graphQLResponse.Data.Issues.Nodes[0], nil
	default:
		return Issue{}, ErrMultipleIssuesFound
	}
}
