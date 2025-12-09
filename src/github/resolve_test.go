package github_test

import (
	_ "embed"
	"encoding/json"
	"testing"
	"time"

	"github.com/roboslone/github-oauth-exchange/src/github"
	"github.com/stretchr/testify/require"
)

//go:embed resolve.json
var rawResponse []byte

func TestResolveUnmarshalling(t *testing.T) {
	var s github.RawResolveResponse

	err := json.Unmarshal(rawResponse, &s)
	require.NoError(t, err)

	p, err := s.ToProto()
	require.NoError(t, err)

	require.EqualValues(t, "roboslone", p.Login)
	require.EqualValues(t, 1941275, p.Id)
	require.EqualValues(t, "MDQ6VXNlcjE5NDEyNzU=", p.NodeId)
	require.EqualValues(t, "https://avatars.githubusercontent.com/u/1941275?v=4", p.Urls.Avatar)
	require.EqualValues(t, "", p.GravatarId)
	require.EqualValues(t, "https://api.github.com/users/roboslone", p.Urls.Api)
	require.EqualValues(t, "https://github.com/roboslone", p.Urls.Html)
	require.EqualValues(t, "https://api.github.com/users/roboslone/followers", p.Urls.Followers)
	require.EqualValues(t, "https://api.github.com/users/roboslone/following{/other_user}", p.Urls.Following)
	require.EqualValues(t, "https://api.github.com/users/roboslone/gists{/gist_id}", p.Urls.Gists)
	require.EqualValues(t, "https://api.github.com/users/roboslone/starred{/owner}{/repo}", p.Urls.Starred)
	require.EqualValues(t, "https://api.github.com/users/roboslone/subscriptions", p.Urls.Subscriptions)
	require.EqualValues(t, "https://api.github.com/users/roboslone/orgs", p.Urls.Organizations)
	require.EqualValues(t, "https://api.github.com/users/roboslone/repos", p.Urls.Repos)
	require.EqualValues(t, "https://api.github.com/users/roboslone/events{/privacy}", p.Urls.Events)
	require.EqualValues(t, "https://api.github.com/users/roboslone/received_events", p.Urls.ReceivedEvents)
	require.EqualValues(t, "User", p.Type)
	require.EqualValues(t, "public", p.UserViewType)
	require.EqualValues(t, false, p.SiteAdmin)
	require.EqualValues(t, "roboslone", p.Name)
	require.EqualValues(t, "", p.Blog)
	require.EqualValues(t, "roboslone@gmail.com", p.Email)
	require.EqualValues(t, "roboslone@gmail.com", p.NotificationEmail)
	require.EqualValues(t, 10, p.PublicRepos)
	require.EqualValues(t, 1, p.PublicGists)
	require.EqualValues(t, 6, p.Followers)
	require.EqualValues(t, 9, p.Following)
	require.EqualValues(t, "2012-07-09T07:58:49Z", p.CreatedAt.AsTime().Format(time.RFC3339))
	require.EqualValues(t, "2025-11-28T08:07:55Z", p.UpdatedAt.AsTime().Format(time.RFC3339))
}
