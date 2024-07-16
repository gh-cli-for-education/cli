package owner

import (
	"net/http"
	"time"

	"github.com/cli/cli/v2/api"
)

type OrganizationList struct {
	Organizations []Organization
	TotalCount    int
	User          string
}

type Organization struct {
	Login string
}

func listAllOrgs(httpClient *http.Client, hostname string) (*OrganizationList, error) {
	type response struct {
		User struct {
			Login         string
			Organizations struct {
				TotalCount int
				Nodes      []Organization
				PageInfo   struct {
					HasNextPage bool
					EndCursor   string
				}
			} `graphql:"organizations(first: $limit, after: $endCursor)"`
		}
	}

	query := `query OrganizationList($user: String!, $limit: Int!, $endCursor: String) {
		user(login: $user) {
			login
			organizations(first: $limit, after: $endCursor) {
				totalCount
				nodes {
					login
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	}`

	cachedHTTPClient := api.NewCachedHTTPClient(httpClient, time.Hour*24*7)
	client := api.NewClientFromHTTP(cachedHTTPClient)

	user, err := api.CurrentLoginName(client, hostname)
	if err != nil {
		return nil, err
	}

	listResult := OrganizationList{}
	listResult.User = user
	listResult.Organizations = append(listResult.Organizations, Organization{Login: user})
	pageLimit := 5
	variables := map[string]interface{}{
		"user": user,
	}

	for {
		variables["limit"] = pageLimit
		var data response
		err := client.GraphQL(hostname, query, variables, &data)
		if err != nil {
			return nil, err
		}

		if listResult.TotalCount == 0 {
			listResult.TotalCount = data.User.Organizations.TotalCount + 1
		}

		listResult.Organizations = append(listResult.Organizations, data.User.Organizations.Nodes...)

		if data.User.Organizations.PageInfo.HasNextPage {
			variables["endCursor"] = data.User.Organizations.PageInfo.EndCursor
		} else {
			break
		}
	}

	return &listResult, nil
}
