cron tab that deletes users one by one every 30 seconds for as long as there are users to delete, 
then deletes the global data at the end. 
Acts like the testusers script.
only does so if there is a flag from tally.
once done, delete flag from tally so that this job does not do anything anymore.