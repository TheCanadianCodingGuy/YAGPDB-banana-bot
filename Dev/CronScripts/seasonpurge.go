{{/* --- CRON: */15 18-23 15 5 * --- */}}
{{ sleep 1 }}
{{ $now := currentTime }}

{{ $nextMonth := (newDate $now.Year (add (toInt $now.Month) 1) 1 0 0 0) }}
{{ if eq (toInt $now.Month) 12 }}{{ $nextMonth = (newDate (add $now.Year 1) 1 1 0 0 0) }}{{ end }}
{{ $lastDayOfMonth := $nextMonth.Add (toDuration (mult 24 -1 | printf "%dh")) }}
{{ if ne $now.Day $lastDayOfMonth.Day }}{{ return }}{{ end }}

{{ $totalMinutes := add (mult $now.Hour 60) $now.Minute }}
{{ if lt $totalMinutes 1125 }}{{ return }}{{ end }}

{{if (dbGet 0 "season_start_flag")}}{{return}}{{end}}

{{$entries := dbTopEntries "banana_slips" 1 0}}

{{if $entries}}
    {{$target := index $entries 0}}
    {{dbDel $target.UserID "banana_slips"}}
    {{dbDel $target.UserID "banana_data"}}
{{else}}
    {{$global := sdict "crash" 0 "oily" false "pity" 0 "season" 1}}
    {{with (dbGet 0 "banana_global")}}
        {{$global = dict .Value | sdict}}
    {{end}}

    {{$newSeason := add (toInt ($global.Get "season")) 1}}
    {{$global.Set "season" $newSeason}}
    {{$global.Set "crash" 0}}
    {{$global.Set "oily" false}}
    {{$global.Set "pity" 0}}

    {{dbSet 0 "banana_global" $global}}
    {{dbSet 0 "season_start_flag" true}}
{{end}}