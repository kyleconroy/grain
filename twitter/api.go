package twitter

import (
	"compress/zlib"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ChimeraCoder/tokenbucket"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/kyleconroy/grain/gen/twitter"
)

const (
	_GET            = iota
	_POST           = iota
	ClientTimeout   = 20
	BaseUrlV1       = "https://api.twitter.com/1"
	BaseUrl         = "https://api.twitter.com/1.1"
	UploadBaseUrl   = "https://upload.twitter.com/1.1"
	DEFAULT_DELAY   = 0 * time.Second
	EFAULT_CAPACITY = 5
)

type query struct {
	url         string
	form        url.Values
	data        interface{}
	method      int
	response_ch chan response
}

type response struct {
	data interface{}
	err  error
}

type TwitterApi struct {
	oauthClient oauth.Client
	Credentials *oauth.Credentials
	HttpClient  *http.Client
	queryQueue  chan query
	bucket      *tokenbucket.Bucket

	// used for testing
	// defaults to BaseUrl
	baseUrl string
}

//NewTwitterApi takes an user-specific access token and secret and returns a TwitterApi struct for that user.
//The TwitterApi struct can be used for accessing any of the endpoints available.
func NewTwitterApi(access_token, access_token_secret, consumer_key, consumer_secret string) *TwitterApi {
	//TODO figure out how much to buffer this channel
	//A non-buffered channel will cause blocking when multiple queries are made at the same time
	queue := make(chan query)
	c := &TwitterApi{
		oauthClient: oauth.Client{
			TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
			ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
			TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
			Credentials: oauth.Credentials{
				Token:  consumer_key,
				Secret: consumer_secret,
			},
		},
		queryQueue: queue,
		bucket:     tokenbucket.NewBucket(5*time.Second, 3),
		Credentials: &oauth.Credentials{
			Token:  access_token,
			Secret: access_token_secret,
		},
		HttpClient: http.DefaultClient,
		baseUrl:    BaseUrl,
	}
	//Configure a timeout to HTTP client (DefaultClient has no default timeout, which may deadlock Mutex-wrapped uses of the lib.)
	c.HttpClient.Timeout = time.Duration(ClientTimeout * time.Second)
	go c.throttledQuery()
	return c
}

func defaultValues(v url.Values) url.Values {
	if v == nil {
		v = url.Values{}
	}
	v.Set("tweet_mode", "extended")
	return v
}

// SetDelay will set the delay between throttled queries
// To turn of throttling, set it to 0 seconds
func (c *TwitterApi) SetDelay(t time.Duration) {
	c.bucket.SetRate(t)
}

// apiGet issues a GET request to the Twitter API and decodes the response JSON to data.
func (c TwitterApi) apiGet(urlStr string, form url.Values, data interface{}) error {
	form = defaultValues(form)
	resp, err := c.oauthClient.Get(c.HttpClient, c.Credentials, urlStr, form)
	if err != nil {
		return err
	}
	fmt.Println(urlStr, resp.StatusCode)
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// apiPost issues a POST request to the Twitter API and decodes the response JSON to data.
func (c TwitterApi) apiPost(urlStr string, form url.Values, data interface{}) error {
	resp, err := c.oauthClient.Post(c.HttpClient, c.Credentials, urlStr, form)
	if err != nil {
		return err
	}
	fmt.Println(urlStr, resp.StatusCode)
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// decodeResponse decodes the JSON response from the Twitter API.
func decodeResponse(resp *http.Response, data interface{}) error {
	// Prevent memory leak in the case where the Response.Body is not used.
	// As per the net/http package, Response.Body still needs to be closed.
	defer resp.Body.Close()

	// Twitter returns deflate data despite the client only requesting gzip
	// data.  net/http automatically handles the latter but not the former:
	// https://github.com/golang/go/issues/18779
	if resp.Header.Get("Content-Encoding") == "deflate" {
		var err error
		resp.Body, err = zlib.NewReader(resp.Body)
		if err != nil {
			return err
		}
	}

	// according to dev.twitter.com, chunked upload append returns HTTP 2XX
	// so we need a special case when decoding the response
	if strings.HasSuffix(resp.Request.URL.String(), "upload.json") ||
		strings.Contains(resp.Request.URL.String(), "webhooks") {
		if resp.StatusCode == 204 {
			// empty response, don't decode
			return nil
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return newApiError(resp)
		}
	} else if resp.StatusCode != 200 {
		return newApiError(resp)
	}
	return json.NewDecoder(resp.Body).Decode(data)

}

//query executes a query to the specified url, sending the values specified by form, and decodes the response JSON to data
//method can be either _GET or _POST
func (c TwitterApi) execQuery(urlStr string, form url.Values, data interface{}, method int) error {
	switch method {
	case _GET:
		return c.apiGet(urlStr, form, data)
	case _POST:
		return c.apiPost(urlStr, form, data)
	default:
		return fmt.Errorf("HTTP method not yet supported")
	}
}

func (c *TwitterApi) throttledQuery() {
	for q := range c.queryQueue {
		url := q.url
		form := q.form
		data := q.data //This is where the actual response will be written
		method := q.method
		response_ch := q.response_ch

		if c.bucket != nil {
			<-c.bucket.SpendToken(1)
		}

		err := c.execQuery(url, form, data, method)

		// Check if Twitter returned a rate-limiting error
		if err != nil {
			if apiErr, ok := err.(*ApiError); ok {
				if isRateLimitError, nextWindow := apiErr.RateLimitCheck(); isRateLimitError {
					// c.Log.Info(apiErr.Error())

					// If this is a rate-limiting error, re-add the job to the queue
					// TODO it really should preserve order
					go func(q query) {
						c.queryQueue <- q
					}(q)

					delay := nextWindow.Sub(time.Now())
					<-time.After(delay)

					// Drain the bucket (start over fresh)
					if c.bucket != nil {
						c.bucket.Drain()
					}

					continue
				}
			}
		}

		response_ch <- response{data, err}
	}
}

func (a TwitterApi) LookupStatuses(params url.Values) ([]*twitterpb.Tweet, error) {
	resp := []*twitterpb.Tweet{}
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/statuses/lookup.json", params, &resp, _POST, response_ch}
	return resp, (<-response_ch).err
}

func (a TwitterApi) GetUserTimeline(params url.Values) ([]*twitterpb.Tweet, error) {
	resp := []*twitterpb.Tweet{}
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/statuses/user_timeline.json", params, &resp, _GET, response_ch}
	return resp, (<-response_ch).err
}

func (a TwitterApi) GetFavorites(params url.Values) ([]*twitterpb.Tweet, error) {
	resp := []*twitterpb.Tweet{}
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/favorites/list.json", params, &resp, _GET, response_ch}
	return resp, (<-response_ch).err
}

type ListsOwnershipsResp struct {
	NextCursor     string           `json:"next_cursor_str"`
	PreviousCursor string           `json:"previous_cursor_str"`
	Lists          []twitterpb.List `json:"lists"`
}

func (a TwitterApi) GetListsOwnedBy(params url.Values) (ListsOwnershipsResp, error) {
	resp := ListsOwnershipsResp{}
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/ownerships.json", params, &resp, _GET, response_ch}
	return resp, (<-response_ch).err
}

type UsersResp struct {
	NextCursor     string           `json:"next_cursor_str"`
	PreviousCursor string           `json:"previous_cursor_str"`
	Users          []twitterpb.User `json:"users"`
}

func (a TwitterApi) GetMembers(params url.Values) (UsersResp, error) {
	resp := UsersResp{}
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/members.json", params, &resp, _GET, response_ch}
	return resp, (<-response_ch).err
}

func (a TwitterApi) GetFriends(params url.Values) (UsersResp, error) {
	resp := UsersResp{}
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/friends/list.json", params, &resp, _GET, response_ch}
	return resp, (<-response_ch).err
}

func (a TwitterApi) GetFollowers(params url.Values) (UsersResp, error) {
	resp := UsersResp{}
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/followers/list.json", params, &resp, _GET, response_ch}
	return resp, (<-response_ch).err
}

// "{\"type\":\"message_create\",\"id\":\"1005146433668841476\",\"created_timestamp\":\"1528480559161\",\"message_create\":{\"target\":{\"recipient_id\":\"86892924\"},\"sender_id\":\"160248151\",\"source_app_id\":\"268278\",\"message_data\":{\"text\":\"Sorry for the random Twitter message, I need to test out the DM API\",\"entities\":{\"hashtags\":[],\"symbols\":[],\"user_mentions\":[],\"urls\":[]}}}}"
type DirectMessageEventsResp struct {
	NextCursor string                          `json:"next_cursor"`
	Events     []*twitterpb.DirectMessageEvent `json:"events"`
}

func (a TwitterApi) GetDirectMessageEvents(params url.Values) (DirectMessageEventsResp, error) {
	resp := DirectMessageEventsResp{}
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/direct_messages/events/list.json", params, &resp, _GET, response_ch}
	return resp, (<-response_ch).err
}
