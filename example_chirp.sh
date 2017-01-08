#!/bin/bash

export CHIRPETTER_CONSUMER_KEY=key
export CHIRPETTER_CONSUMER_SECRET=secret
export CHIRPETTER_ACCESS_TOKEN=token
export CHIRPETTER_ACCESS_SECRET=secret
export CHIRPETTER_LAST_ID=/path/to/last_tweet_id
export CHIRPETTER_EXTRACT_LINKS=false

/path/to/chirpetter/chirpetter | mail -s "tweets for `date`" mail@example.com
