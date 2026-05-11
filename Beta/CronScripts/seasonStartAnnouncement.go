{{/* --- CRON: 0 0 1 * * --- */}}
{{ sleep 5 }}
{{ $now := currentTime }}

{{/* 1. Setup Keys & Data */}}
{{ $keyGlobal := "TEST_banana_global" }}
{{ $keyFlag := "TEST_season_start_flag" }}

{{/* Fetch the Season number (already incremented by the Purge script) */}}
{{ $global := sdict "season" 1 }}
{{ with (dbGet 0 $keyGlobal) }}
    {{ $global = dict .Value | sdict }}
{{ end }}
{{ $seasonN := toInt ($global.Get "season") }}

{{/* 2. Calculate End of Season (6 hours before the next month starts) */}}
{{ $nextMonthYear := $now.Year }}
{{ $nextMonthVal := add (toInt $now.Month) 1 }}
{{ if eq $nextMonthVal 13 }}
    {{ $nextMonthYear = add $now.Year 1 }}
    {{ $nextMonthVal = 1 }}
{{ end }}

{{/* First of NEXT month at 00:00:00 */}}
{{ $nextMonthStart := (newDate $nextMonthYear $nextMonthVal 1 0 0 0) }}

{{/* Subtract 6 hours to get 18:00 on the last day of the current month */}}
{{ $endOfSeason := $nextMonthStart.Add (toDuration (mult 6 -1 | printf "%dh")) }}

{{/* 3. Execution: Send message and lift the gate */}}
{{ $chanID := 1438831937220378674 }}
{{ sendMessage (toInt64 $chanID) (printf "🍌 **SEASON %d HAS OFFICIALLY BEGUN!** 🍌\nThe floors are freshly waxed and the peels are scattered. It's time to slip away!\n\n🏁 **Season %d Ends:** <t:%d:F> (<t:%d:R>)\n\nGood luck to all slippers! May your falls be frequent and your prestige grow." $seasonN $seasonN $endOfSeason.Unix $endOfSeason.Unix) }}

{{ dbDel 0 $keyFlag }}