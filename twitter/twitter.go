package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/kyleconroy/grain/archive"
	"github.com/kyleconroy/grain/gen/twitter"
	toml "github.com/pelletier/go-toml"
)

type Archiver struct {
	username          string
	consumerKey       string
	consumerSecret    string
	accessToken       string
	accessTokenSecret string
	csvPath           string
	basePath          string
	httpClient        *http.Client
}

func NewArchiver(c *toml.Tree) (*Archiver, error) {
	// TODO: Use a well-known location for the cache
	transport := httpcache.NewTransport(diskcache.New("httpcache"))

	a := Archiver{
		basePath:   filepath.Join("archive", "twitter"),
		httpClient: transport.Client(),
	}

	var ok bool
	pairs := []struct {
		prop     *string
		key      string
		required bool
	}{
		{&a.username, "username", true},
		{&a.consumerKey, "consumer-key", true},
		{&a.consumerSecret, "consumer-secret", true},
		{&a.accessToken, "access-token", true},
		{&a.accessTokenSecret, "access-token-secret", true},
		{&a.csvPath, "tweet-csv", false},
	}

	for _, pair := range pairs {
		if pair.required && !c.Has(pair.key) {
			return nil, fmt.Errorf("Missing `%s` key in [twitter] section", pair.key)
		}
		*pair.prop, ok = c.Get(pair.key).(string)
		if !ok {
			return nil, fmt.Errorf("Expected `%s` to be a string", pair.key)
		}
	}

	return &a, nil
}

func marshal(pb proto.Message, path string) error {
	m := jsonpb.Marshaler{Indent: "  ", OrigName: true}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return m.Marshal(f, pb)
}

func ReadN(c *Reader, n int, lookup map[int64]struct{}) ([]CSVTweet, error) {
	ids := []CSVTweet{}
	for len(ids) < n {
		tweet, err := c.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return ids, err
		}
		if _, ok := lookup[tweet.TweetID]; !ok {
			ids = append(ids, tweet)
		}
	}
	return ids, nil
}

func (a *Archiver) Sync(ctx context.Context) error {
	api := NewTwitterApi(a.accessToken, a.accessTokenSecret, a.consumerKey, a.consumerSecret)
	api.HttpClient = a.httpClient
	api.SetDelay(10 * time.Second)

	if err := os.MkdirAll(a.basePath, 0644); err != nil {
		return err
	}

	listPath := filepath.Join(a.basePath, "lists.json")
	tweetPath := filepath.Join(a.basePath, "tweets.json")
	friendPath := filepath.Join(a.basePath, "friends.json")
	followerPath := filepath.Join(a.basePath, "followers.json")
	favoritePath := filepath.Join(a.basePath, "favorites.json")
	dmPath := filepath.Join(a.basePath, "dms.json")

	g, ctx := errgroup.WithContext(ctx)

	if _, err := os.Stat(tweetPath); os.IsNotExist(err) {
		g.Go(func() error {
			tweets, err := a.tweets(ctx, api)
			if err != nil {
				return err
			}

			if err := marshal(&twitterpb.Archive{Timeline: tweets}, tweetPath); err != nil {
				return err
			}

			if a.csvPath != "" {
				h, err := os.Open(a.csvPath)
				if err != nil {
					return nil
				}

				// Build out a lookup table of existing tweets
				lookup := map[int64]struct{}{}
				for _, t := range tweets {
					lookup[t.Id] = struct{}{}
				}

				// Load the Tweet archive
				reader := NewCSVReader(h)

				for {
					// For each tweet, see if we have it
					subset, err := ReadN(reader, 100, lookup)
					if err != nil {
						return err
					}

					if len(subset) == 0 {
						break
					}

					ids := []string{}
					for _, t := range subset {
						ids = append(ids, strconv.Itoa(int(t.TweetID)))
					}

					params := url.Values{}
					params.Set("include_entities", "true")
					params.Set("id", strings.Join(ids, ","))
					statuses, err := api.LookupStatuses(params)
					if err != nil {
						return err
					}

					fmt.Printf("Processed %d tweets…\n", len(statuses))

					for _, s := range statuses {
						tweets = append(tweets, s)
					}
				}

				if err := marshal(&twitterpb.Archive{Timeline: tweets}, tweetPath); err != nil {
					return err
				}
			}

			for _, tweet := range tweets {
				if err := a.archiveTweet(ctx, tweet); err != nil {
					return err
				}

				if err := a.archiveTweet(ctx, tweet.RetweetedStatus); err != nil {
					return err
				}

				if err := a.archiveTweet(ctx, tweet.QuotedStatus); err != nil {
					return err
				}
			}
			return nil
		})
	}

	if _, err := os.Stat(favoritePath); os.IsNotExist(err) {
		g.Go(func() error {
			favs, err := a.favorites(ctx, api)
			if err != nil {
				return err
			}

			b, err := json.MarshalIndent(twitterpb.Archive{Favorites: favs}, "", "  ")
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(favoritePath, b, 0644); err != nil {
				return err
			}

			for _, fav := range favs {
				if err := a.archiveTweet(ctx, fav); err != nil {
					return err
				}

				if err := a.archiveTweet(ctx, fav.RetweetedStatus); err != nil {
					return err
				}

				if err := a.archiveTweet(ctx, fav.QuotedStatus); err != nil {
					return err
				}
			}
			return nil
		})
	}

	if _, err := os.Stat(listPath); os.IsNotExist(err) {
		g.Go(func() error {
			lists, err := a.lists(ctx, api)
			if err != nil {
				return err
			}

			b, err := json.MarshalIndent(twitterpb.Archive{Lists: lists}, "", "  ")
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(listPath, b, 0644); err != nil {
				return err
			}
			for _, l := range lists {
				for _, u := range l.Members {
					if err := a.archiveUser(ctx, u); err != nil {
						return err
					}
				}
			}
			return nil
		})
	}

	if _, err := os.Stat(friendPath); os.IsNotExist(err) {
		g.Go(func() error {
			friends, err := a.friends(ctx, api)
			if err != nil {
				return err
			}

			b, err := json.MarshalIndent(twitterpb.Archive{Friends: friends}, "", "  ")
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(friendPath, b, 0644); err != nil {
				return err
			}

			for _, u := range friends {
				if err := a.archiveUser(ctx, u); err != nil {
					return err
				}
			}
			return nil
		})
	}

	if _, err := os.Stat(followerPath); os.IsNotExist(err) {
		g.Go(func() error {
			followers, err := a.followers(ctx, api)
			if err != nil {
				return err
			}

			b, err := json.MarshalIndent(twitterpb.Archive{Followers: followers}, "", "  ")
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(followerPath, b, 0644); err != nil {
				return err
			}

			for _, u := range followers {
				if err := a.archiveUser(ctx, u); err != nil {
					return err
				}
			}
			return nil
		})
	}

	if _, err := os.Stat(dmPath); os.IsNotExist(err) {
		g.Go(func() error {
			dms, err := a.dms(ctx, api)
			if err != nil {
				return err
			}

			b, err := json.MarshalIndent(twitterpb.Archive{DirectMessages: dms}, "", "  ")
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(dmPath, b, 0644); err != nil {
				return err
			}
			return nil
		})
	}

	return g.Wait()
}

func (a *Archiver) tweets(ctx context.Context, api *TwitterApi) ([]*twitterpb.Tweet, error) {
	params := url.Values{}
	params.Set("screen_name", a.username)
	params.Set("count", "1000")
	tweets := []*twitterpb.Tweet{}

	for {
		timeline, err := api.GetUserTimeline(params)
		if err != nil {
			return tweets, err
		}

		fmt.Printf("Processed %d tweets…\n", len(timeline))

		// tweets -> file -> blob
		for _, tweet := range timeline {
			tweets = append(tweets, tweet)
		}

		if len(timeline) <= 1 {
			return tweets, nil
		}

		params.Set("max_id", strconv.Itoa(int(timeline[len(timeline)-1].Id-1)))
	}
	return tweets, nil
}

func (a *Archiver) favorites(ctx context.Context, api *TwitterApi) ([]*twitterpb.Tweet, error) {
	params := url.Values{}
	params.Set("screen_name", a.username)
	params.Set("count", "1000")
	favorites := []*twitterpb.Tweet{}

	for {
		timeline, err := api.GetFavorites(params)
		if err != nil {
			return favorites, err
		}

		fmt.Printf("Processed %d favorites…\n", len(timeline))

		// favorites -> file -> blob
		for _, tweet := range timeline {
			favorites = append(favorites, tweet)
		}

		if len(timeline) <= 1 {
			return favorites, nil
		}

		params.Set("max_id", strconv.Itoa(int(timeline[len(timeline)-1].Id-1)))
	}
	return favorites, nil
}

func (a *Archiver) lists(ctx context.Context, api *TwitterApi) ([]*twitterpb.List, error) {
	params := url.Values{}
	params.Set("screen_name", a.username)
	params.Set("count", "1000")
	params.Set("cursor", "-1")
	lists := []*twitterpb.List{}
	for {
		resp, err := api.GetListsOwnedBy(params)
		if err != nil {
			return lists, err
		}

		fmt.Printf("Processed %d lists…\n", len(resp.Lists))

		for _, list := range resp.Lists {
			l := list

			m, err := a.members(ctx, api, l.Id)
			if err != nil {
				return lists, err
			}

			l.Members = m

			lists = append(lists, &l)
		}

		if resp.NextCursor == "0" {
			return lists, nil
		}

		params.Set("cursor", resp.NextCursor)
	}
	return lists, nil
}

func (a *Archiver) friends(ctx context.Context, api *TwitterApi) ([]*twitterpb.User, error) {
	params := url.Values{}
	params.Set("screen_name", a.username)
	params.Set("count", "1000")
	params.Set("cursor", "-1")
	members := []*twitterpb.User{}
	for {
		resp, err := api.GetFriends(params)
		if err != nil {
			return members, err
		}
		for _, member := range resp.Users {
			m := member
			members = append(members, &m)
		}
		if resp.NextCursor == "0" {
			return members, nil
		}
		params.Set("cursor", resp.NextCursor)
	}
	return members, nil
}

func (a *Archiver) followers(ctx context.Context, api *TwitterApi) ([]*twitterpb.User, error) {
	params := url.Values{}
	params.Set("screen_name", a.username)
	params.Set("count", "1000")
	params.Set("cursor", "-1")
	members := []*twitterpb.User{}
	for {
		resp, err := api.GetFollowers(params)
		if err != nil {
			return members, err
		}
		for _, member := range resp.Users {
			m := member
			members = append(members, &m)
		}
		if resp.NextCursor == "0" {
			return members, nil
		}
		params.Set("cursor", resp.NextCursor)
	}
	return members, nil
}

func (a *Archiver) members(ctx context.Context, api *TwitterApi, list int64) ([]*twitterpb.User, error) {
	params := url.Values{}
	params.Set("list_id", strconv.Itoa(int(list)))
	params.Set("count", "1000")
	params.Set("cursor", "-1")
	members := []*twitterpb.User{}
	for {
		resp, err := api.GetMembers(params)
		if err != nil {
			return members, err
		}

		for _, member := range resp.Users {
			m := member
			members = append(members, &m)
		}

		if resp.NextCursor == "0" {
			return members, nil
		}

		params.Set("cursor", resp.NextCursor)
	}
	return members, nil
}

func (a *Archiver) dms(ctx context.Context, api *TwitterApi) ([]*twitterpb.DirectMessageEvent, error) {
	params := url.Values{}
	params.Set("count", "1000")
	events := []*twitterpb.DirectMessageEvent{}
	for {
		resp, err := api.GetDirectMessageEvents(params)
		if err != nil {
			return events, err
		}

		for _, event := range resp.Events {
			e := event
			events = append(events, e)
		}

		if resp.NextCursor == "" {
			return events, nil
		}

		params.Set("cursor", resp.NextCursor)
	}
	return events, nil
}

func (a *Archiver) archiveTweet(ctx context.Context, tweet *twitterpb.Tweet) error {
	if tweet == nil {
		return nil
	}

	if err := a.archiveUser(ctx, tweet.User); err != nil {
		return err
	}

	urls := []string{}
	// For each media item, attempt to save it to the
	if tweet.ExtendedEntities != nil {
		for _, media := range tweet.ExtendedEntities.Media {
			urls = append(urls, media.MediaUrlHttps)
		}
	}

	for _, uri := range urls {
		if uri == "" {
			// TODO: Understand why this happens
			continue
		}

		if err := archive.ArchiveURL(ctx, "twitter", "media", uri); err != nil {
			return err
		}
	}
	return nil
}

func (a *Archiver) archiveUser(ctx context.Context, user *twitterpb.User) error {
	if user == nil {
		return nil
	}
	urls := []string{
		user.ProfileBackgroundImageUrlHttps,
		user.ProfileImageUrlHttps,
	}

	for _, uri := range urls {
		if uri == "" {
			// TODO: Understand why this happens
			continue
		}

		if err := archive.ArchiveURL(ctx, "twitter", "media", uri); err != nil {
			return err
		}
	}
	return nil
}
