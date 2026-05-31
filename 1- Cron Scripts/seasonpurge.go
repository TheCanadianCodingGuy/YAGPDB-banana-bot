{{/* --- CRON: *\/15 18-23 15 5 * --- */}}
{{ sleep 1 }}
{{ $now := currentTime }}

{{ $nextMonth := (newDate $now.Year (add (toInt $now.Month) 1) 1 0 0 0) }}
{{ if eq (toInt $now.Month) 12 }}{{ $nextMonth = (newDate (add $now.Year 1) 1 1 0 0 0) }}{{ end }}
{{ $lastDayOfMonth := $nextMonth.Add (toDuration (mult 24 -1 | printf "%dh")) }}
{{ if ne $now.Day $lastDayOfMonth.Day }}{{ return }}{{ end }}

{{ $totalMinutes := add (mult $now.Hour 60) $now.Minute }}
{{ if lt $totalMinutes 1125 }}{{ return }}{{ end }}

{{if (dbGet 0 "season_start_flag")}}{{return}}{{end}}

{{/* Fetch or initialize the iteration offset position */}}
{{$offset := 0}}
{{with (dbGet 0 "banana_purge_offset")}}
    {{$offset = toInt .Value}}
{{end}}

{{/* Fetch the single user entry at the current offset position */}}
{{$entries := dbTopEntries "banana_slips" 1 $offset}}

{{if $entries}}
    {{$target := index $entries 0}}
    
    {{/* --- AUTOMATED SEASONAL DATA OVERRIDE --- */}}
    {{$userData := sdict "c" 0 "r" 0 "turbo" false "cd" 0 "s_mcs" 0 "s_hs" 0 "s_rss" 0 "s_gs" 0 "s_ts" 0 "s_ss" 0 "s_os" 0 "s_ms" 0 "s_cs" 0 "s_ns" 0 "s_f" 0 "g_mcs" 0 "g_hs" 0 "g_rss" 0 "g_gs" 0 "g_ts" 0 "g_ss" 0 "g_os" 0 "g_ms" 0 "g_cs" 0 "g_ns" 0 "g_f" 0}}
    {{with (dbGet $target.UserID "banana_data")}}
        {{$userData = dict .Value | sdict}}
    {{end}}

    {{range $key, $value := $userData}}
        {{if hasPrefix $key "s_"}}
            {{$userData.Set $key 0}}
        {{else if eq $key "turbo"}}
            {{$userData.Set $key false}}
        {{else if eq $key "c" "r" "cd"}}
            {{$userData.Set $key 0}}
        {{end}}
    {{end}}
    {{dbSet $target.UserID "banana_data" $userData}}

    {{/* Increments and saves the position pointer for the next execution flight */}}
    {{dbSet 0 "banana_purge_offset" (add $offset 1)}}

{{else}}
    {{/* No more entries found: Clear the offset pointer and increment global seasonal data */}}
    {{dbDel 0 "banana_purge_offset"}}

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