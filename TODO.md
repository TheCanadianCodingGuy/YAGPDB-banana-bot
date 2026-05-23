FOR RELEASE
1) Make -banana in general chat return a message to head over the other channel
2) remove all other commands
4) replace prod (with real channels and restrictions) with Dev-RealDates in YAGPDB
5) Clean the repo for prod scripts.
a) keep a copy of minified prod script for <10000 chars
6) make 15 cron jobs for purge season with 60 second sleep in between.


HOW TO RELEASE
1) run and pin -commandslist in #slippery-slope
2) run and pin -what in #slippery-slope
3) run 1- -globalsetup.go in #slippery-slope
4) run 2- -tallyslips.go in #slippery-slope
5) run 3- -oldseasonpurge.go enough to purge all db while in #bot-test
6) run 4- -incrementseason.go
7) activate cron for seasonstartannouncement (only that one)

ONCE SEASON START
1) activate all crons