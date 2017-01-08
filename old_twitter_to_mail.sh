#!/bin/bash
#############

# change these
TWITTER_ACCOUNT="youraccounthere"
EMAIL_ADDRESS="you@example.com"

# probably don't need to change these
TWEETS_FILE=~/.t2m_tweets
LATEST_TWEETS_ID=~/.t2m_latest_id

# install https://github.com/sferik/t
T=/usr/local/bin/t

# set the right acct
$T set account $TWITTER_ACCOUNT

# check if we have a last tweeted id, otherwise choose one 200 ago
if [ ! -e $LATEST_TWEETS_ID ]
then
        $T timeline -l -n=200 | tail -n 1 | cut -f1 -d' ' > $LATEST_TWEETS_ID
fi

# fetch all tweets since the latest twitter id | remove first line > store in tweet file
$T timeline -r -s=`cat $LATEST_TWEETS_ID` -d -n 1000 | tail -n +2 > $TWEETS_FILE

# if the tweet file is non-empty, mail it out
if [ -s $TWEETS_FILE ]
then
	cat $TWEETS_FILE | mail -s "tweets for `date`" $EMAIL_ADDRESS

	# store the latest tweet id
	$T timeline -l -n=1 | tail -n 1 | cut -f1 -d' ' > $LATEST_TWEETS_ID
fi
