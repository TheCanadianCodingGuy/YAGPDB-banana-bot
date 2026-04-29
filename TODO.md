1) Test -banana
2) adapt -topslips for new dataset and prestige
4) adapt userpurge for new dataset + limit of 5 calls db
5) Script database purge.
7) Script top 3 reward with a +3, 2, 1, at end of season.
8) Script end of seasons
    - Calculate points
    - Add points to users (add to data json)
    - lock the banana, topslips, topflips command for 6 hours
    - have the cronjob run every minute to delete the user list (have it cron??? or check.)
    - start new season at the start of the month with message (cron every month)
9) for everything, show end of season time (6h before end of each month) in scripts that show leaderboards.
10) end of season points award must give the point attribution award.
12) make new channel only respond to commands and erase other messages.







------------------script erase messages (custom command regex trigger .*)
{{/* 1. Check if the message starts with the dash */}}
{{if hasPrefix .Message.Content "-"}}
    
    {{/* 2. Strip the dash and get only the first word (the command) */}}
    {{$args := split (slice .Message.Content 1) " "}}
    {{$cmd := lower (index $args 0)}}

    {{/* 3. Check if that word is in our allowed list */}}
    {{range $allowed}}
        {{if eq $cmd .}}
            {{$isValid = true}}
        {{end}}
    {{end}}
{{end}}

{{/* 4. If it's not a valid command, delete it */}}
{{if not $isValid}}
    {{deleteTrigger 0}}
{{end}}








------------------------  Just put that in banana itself to remove the db entry, blocks it 6h before end of month
{{ $now := currentTime }}
{{ $nextMonth := (printf "%d-%02d-01T00:00:00Z" (int $now.Year) (add (int $now.Month) 1 | slice 0 12) | community.parseTime) }}
{{ if (int $now.Month) == 12 }}{{ $nextMonth = (printf "%d-01-01T00:00:00Z" (add $now.Year 1) | community.parseTime) }}{{ end }}

{{ if $now.After ($nextMonth.Add ( community.parseDuration "-6h" )) }}
    {{ sendMessage nil "⚠️ **Season Over!** Command locked for tallying. Results in 6 hours!" }}
    {{ return }}
{{ end }}





-------------------- prestige points example (CRON: 30 18 28-31 * *)
{{/* Check if today is the last day of the month */}}
{{ $now := currentTime }}
{{ $isLastDay := eq $now.Day ($now.AddDate 0 0 (int (neg $now.Day))).AddDate 0 1 -1.Day }}
{{ if not $isLastDay }}{{ return }}{{ end }}

{{/* 1. Get all entries (Limit 100 for safety) */}}
{{ $entries := dbTopEntries "TEST_banana_data" 100 0 }}
{{ $leaderboard := cslice }}

{{/* 2. Sort them by 'r' (record) */}}
{{ range $entries }}
    {{ $leaderboard = $leaderboard.Append (sdict "UID" .UserID "Score" (toInt .Value.r)) }}
{{ end }}

{{/* 3. Reward Top 3 */}}
{{ range $i, $user := $leaderboard }}
    {{ if eq $i 0 }}{{/* Gold: 50pts */}}
        {{ dbIncr $user.UID "Prestige_Points" 50 }}
    {{ else if eq $i 1 }}{{/* Silver: 25pts */}}
        {{ dbIncr $user.UID "Prestige_Points" 25 }}
    {{ else if eq $i 2 }}{{/* Bronze: 10pts */}}
        {{ dbIncr $user.UID "Prestige_Points" 10 }}
    {{ end }}
{{ end }}

{{/* TODO: system message for points, and when the new season stares (hammertime) */}}




-------------------------- purge users example (CRON: * 19-23 28-31 * *)
{{ $now := currentTime }}
{{ $isLastDay := eq $now.Day ($now.AddDate 0 0 (int (neg $now.Day))).AddDate 0 1 -1.Day }}

{{/* Only run on the last day during the 19:00-23:59 window */}}
{{ if and $isLastDay (ge $now.Hour 19) }}
    {{/* CONFIG */}}
    {{ $batchSize := 2 }}
    {{ $progressKey := "purge_index" }}

    {{/* 1. Check if we should even be purging */}}
    {{/* (Reuse Phase 1 logic here to check if within 5-hour window) */}}

    {{/* 2. Get Data */}}
    {{ $entries := dbTopEntries "TEST_banana_data" 100 0 }}
    {{ if not $entries }}
        {{ dbDel 0 "TEST_banana_global" }}
        {{ dbDel 0 $progressKey }}
        {{ return }}
    {{ end }}

    {{/* 3. Delete Batch */}}
    {{ range $i, $entry := $entries }}
        {{ if lt $i $batchSize }}
            {{ dbDel $entry.UserID "TEST_banana_slips" }}
            {{ dbDel $entry.UserID "TEST_banana_data" }}
            {{ dbDel $entry.UserID "TEST_banana_cooldown" }}
            {{/* TODO: DO NOT DELETE THE DATA, RESET r, c TO ZERO, KEEP p (prestige) */}}
        {{ end }}
    {{ end }}
{{ else }}
    {{ return }}
{{ end }}





----------------------- example add prestige
{{/* Configuration */}}
{{ $keyGlobal := "TEST_prestige_global" }}
{{ $expiration := 31536000 }}

{{/* 1. Define the Batch Data (UserID: Increment) */}}
{{ $updates := sdict 
    "180128558776057856" 3
    "600533807132835851" 2
    "149357474695217152" 1
}}

{{/* 2. Fetch the existing Global Map */}}
{{ $prestigeMap := sdict }}
{{ with (dbGet 0 $keyGlobal) }}
    {{ $prestigeMap = dict .Value | sdict }}
{{ end }}

{{/* 3. Apply the updates */}}
{{ range $id, $inc := $updates }}
    {{ $current := index $prestigeMap $id | or 0 }}
    {{ $prestigeMap.Set $id (add $current $inc) }}
{{ end }}

{{/* 4. Save the updated map */}}
{{ dbSet 0 $keyGlobal $prestigeMap }}

✅ **Batch update complete!**
{{ range $id, $inc := $updates }}
- `<@{{ $id }}>`: +{{ $inc }} prestige (Total: {{ index $prestigeMap $id }})
{{ end }}