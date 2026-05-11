5) Season start announcement cron script. 
7) Create prestige leaderboard
9) create command list script
10) A what command explaining the bot and events.
11) Create script to delete messages automatically that is not bot or allowed commands.

Transfer from beta to live TODO:
1) make -banana obselete in general, new message to go to the enw channel and use it there
2) ???


-------------------- prestige points example (CRON: 30 18 28-31 * *)





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
{{ $kg := "TEST_prestige_global" }}
{{ $ex := 31536000 }}

{{/* 1. Define the Batch Data (UserID: Increment) */}}
{{ $updates := sdict 
    "180128558776057856" 3
    "600533807132835851" 2
    "149357474695217152" 1
}}

{{/* 2. Fetch the existing Global Map */}}
{{ $prestigeMap := sdict }}
{{ with (dbGet 0 $kg) }}
    {{ $prestigeMap = dict .Value | sdict }}
{{ end }}

{{/* 3. Apply the updates */}}
{{ range $id, $inc := $updates }}
    {{ $current := index $prestigeMap $id | or 0 }}
    {{ $prestigeMap.Set $id (add $current $inc) }}
{{ end }}

{{/* 4. Save the updated map */}}
{{ dbSet 0 $kg $prestigeMap }}

✅ **Batch update complete!**
{{ range $id, $inc := $updates }}
- `<@{{ $id }}>`: +{{ $inc }} prestige (Total: {{ index $prestigeMap $id }})
{{ end }}