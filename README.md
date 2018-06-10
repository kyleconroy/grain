# Grain

The Entire History of You

## Usage

Grain looks for API access keys and secrets in `config.toml`. 

## Twitter

Grain will download a full archive of your Twitter account, unlike the official
Twitter Archive. Your Grain archive includes the following records:

- Direct messages
- Favorites
- Followers
- Friends
- Lists
- Tweets

The archive also includes all media associated with the above records.

For privacy and performance reasons, you'll need to obtain your own API
credentials.

1. Log onto [Twitter](https://twitter.com) with the account you'd like to archive
2. Request your archive from your [Settings page](https://twitter.com/settings/account)
3. Create a new [application](https://apps.twitter.com)
4. Go to the "Keys and Access Tokens" section for your application
5. Generate access tokens via "Create my access token"
6. Fill in `config.toml` with the access and secret tokens

```
[twitter]
username = ""
tweet-csv = "path/to/tweets.csv"
consumer-key = ""
consumer-secret = ""
access-token = ""
access-token-secret = ""
```

Running `grain` will download records to `archive/twitter` in the current
working directory.

### Limitations

The Twitter API has a set of limitations which makes archiving certain records
difficult.

- The API only returns the most recent 3,200 tweets from your timeline. You'll
  need to download your official Twitter Archive for Grain to successfully
archive all of your tweets.
- Only the last 30 days of direct message activity is available via the API.
- Rate limits are very aggressive which means certain records take forever to
  archive.

## Facebook

Coming soon…
