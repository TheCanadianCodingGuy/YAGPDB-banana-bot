{{/* --- CRON: */15 18-23 15 5 * --- */}}
{{ sleep 1 }}
{{ $now := currentTime }}

{{/* 1. Date Check */}}
{{ $nextMonth := (newDate $now.Year (add (toInt $now.Month) 1) 1 0 0 0) }}
{{ if eq (toInt $now.Month) 12 }}{{ $nextMonth = (newDate (add $now.Year 1) 1 1 0 0 0) }}{{ end }}
{{ $lastDayOfMonth := $nextMonth.Add (toDuration (mult 24 -1 | printf "%dh")) }}
{{ if ne $now.Day $lastDayOfMonth.Day }}{{ return }}{{ end }}

{{/* 2. Time Gate: Don't start until after the Reward script (18:30) has finished */}}
{{ $totalMinutes := add (mult $now.Hour 60) $now.Minute }}
{{ if lt $totalMinutes 1125 }}{{ return }}{{ end }}

{{/* 3. Circuit Breaker: If season is already flagged as prepared, stop */}}
{{if (dbGet 0 "TEST_season_start_flag")}}{{return}}{{end}}

{{/* 4. Identify 1 user to purge */}}
{{$entries := dbTopEntries "TEST_banana_slips" 1 0}}

{{if $entries}}
    {{/* 5. Delete user data (2 calls) */}}
    {{$target := index $entries 0}}
    {{dbDel $target.UserID "TEST_banana_slips"}}
    {{dbDel $target.UserID "TEST_banana_data"}}
{{else}}
    {{/* 6. Cleanup finished! Reset Global & Increment Season */}}
    {{$global := sdict "crash" 0 "oily" false "pity" 0 "season" 1}}
    {{with (dbGet 0 "TEST_banana_global")}}
        {{$global = dict .Value | sdict}}
    {{end}}

    {{$newSeason := add (toInt ($global.Get "season")) 1}}
    {{$global.Set "season" $newSeason}}
    {{$global.Set "crash" 0}}
    {{$global.Set "oily" false}}
    {{$global.Set "pity" 0}}

    {{/* Save Global and set the Start Flag (2 calls) */}}
    {{dbSet 0 "TEST_banana_global" $global}}
    {{dbSet 0 "TEST_season_start_flag" true}}
{{end}}