1) Test Cron jobs scripts (fake end of season from Dev-CronTest folder and make it end season)
a) Test multiple same purge scripts and give sleep in them to separate execution time since it is only each 15 minutes, not each minute.
b) make 15 scripts for purge with a 60 sec sleep each except first one.
c) redo crons for 12 hours instead of 6
d) redo scripts ifs for 12 hours, not 6 (end of season)
2) Migrate scripts from Dev-RealDates into Beta
3) For Beta:
a) Remove restrictions
b) Preface all the things with BETA (messages and variables)
c) Lower cooldown to 2 hours
d) Create Beta group in YAGBDP and scripts.
e) Alter a purge script for beta

ONCE BETA IS OVER
1) Alter Dev-RealDates for prod data and variables

FOR RELEASE
1) Make -banana in general chat return a message to head over the other channel
2) remove all other commands
3) remove pins
4) replace prod (with real channels and restrictions) with Dev-RealDates in YAGPDB
5) Clean the repo for prod scripts.
a) keep a copy of minified prod script for <10000 chars