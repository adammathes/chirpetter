# chirpetter

a deliberate twitter client

## what?

Stop checking twitter all the time.

Instead just have whatever has happened emailed to you once a day.

This is mostly a proof of concept, the code here is trivial, I just wanted to play around with Go a bit.

This will email you all the tweets posted since the last run. It also helpfully pulls out any non-twitter/social media links to the top and fetches their title.

## how?

I should probably distribute binaries for this but that is left as an exercise for the reader.

Alternatively, the shell script old version of this is included in this repo (old_twitter_to_mail.sh)

   1. Install go https://golang.org
   2. git clone https://github.com/adammathes/chirpetter.git
   3. go get relevant libraries, go build this repo
   4. cp example_chirp.sh chirp.sh
   5. Edit chirp.sh to include your twitter keys/auth/etc. Consider following the instructions and using the authorization tool https://github.com/sferik/t
   6. run chirp.sh

Make sure everything looks reasonable and add to cron something like

   0 8 * * * ~/chirp.sh



## TODO (maybe)

   * Add mail integration to decrease independence on system mail/mailx
   * HTML output
