package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	githubv1 "github.com/roboslone/github-oauth-exchange/proto/github/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RawResolveResponse struct {
	Login             string `json:"login"`
	ID                uint32 `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	UserViewType      string `json:"user_view_type"`
	SiteAdmin         bool   `json:"site_admin"`
	Name              string `json:"name"`
	Blog              string `json:"blog"`
	Email             string `json:"email"`
	NotificationEmail string `json:"notification_email"`
	PublicRepos       uint32 `json:"public_repos"`
	PublicGists       uint32 `json:"public_gists"`
	Followers         uint32 `json:"followers"`
	Following         uint32 `json:"following"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

func parseTime(s string) (*timestamppb.Timestamp, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}
	return timestamppb.New(t), nil
}

func (r *RawResolveResponse) ToProto() (*githubv1.Account, error) {
	p := &githubv1.Account{
		Login:      r.Login,
		Id:         r.ID,
		NodeId:     r.NodeID,
		GravatarId: r.GravatarID,
		Urls: &githubv1.Account_URLs{
			Avatar:         r.AvatarURL,
			Api:            r.URL,
			Html:           r.HTMLURL,
			Followers:      r.FollowersURL,
			Following:      r.FollowingURL,
			Gists:          r.GistsURL,
			Starred:        r.StarredURL,
			Subscriptions:  r.SubscriptionsURL,
			Organizations:  r.OrganizationsURL,
			Repos:          r.ReposURL,
			Events:         r.EventsURL,
			ReceivedEvents: r.ReceivedEventsURL,
		},
		Type:              r.Type,
		UserViewType:      r.UserViewType,
		SiteAdmin:         r.SiteAdmin,
		Name:              r.Name,
		Blog:              r.Blog,
		Email:             r.Email,
		NotificationEmail: r.NotificationEmail,
		PublicRepos:       r.PublicRepos,
		PublicGists:       r.PublicGists,
		Followers:         r.Followers,
		Following:         r.Following,
	}

	var err error
	p.CreatedAt, err = parseTime(r.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing `created_at`: %q: %w", r.CreatedAt, err)
	}
	p.UpdatedAt, err = parseTime(r.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing `updated_at`: %q: %w", r.UpdatedAt, err)
	}

	return p, nil
}

func Resolve(ctx context.Context, accessToken string) (*githubv1.Account, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("constructing user info request: %w", err)
	}

	request.Header = http.Header{
		"accept":        []string{"application/json"},
		"content-type":  []string{"application/json"},
		"authorization": []string{fmt.Sprintf("Bearer %s", accessToken)},
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("resolving token: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var info RawResolveResponse
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}
	return info.ToProto()
}
