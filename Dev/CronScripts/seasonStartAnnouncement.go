{{/* --- CRON: 0 0 1 * * --- */}}
{{ sleep 5 }}
{{ $now := currentTime }}

{{ $keyGlobal := "banana_global" }}
{{ $keyFlag := "season_start_flag" }}
{{ $keyPin := "season_announcement" }}
{{ $chanID := 1497775909036494848 }}

{{ $global := sdict "season" 1 }}
{{ with (dbGet 0 $keyGlobal) }}{{ $global = dict .Value | sdict }}{{ end }}
{{ $seasonN := toInt ($global.Get "season") }}

{{ with (dbGet 0 $keyPin) }}
    {{ $oldMsgID := toInt64 .Value }}
    {{ if (getMessage (toInt64 $chanID) $oldMsgID) }}
        {{ unpinMessage (toInt64 $chanID) $oldMsgID }}
    {{ end }}
{{ end }}

{{ $nextMonthYear := $now.Year }}
{{ $nextMonthVal := add (toInt $now.Month) 1 }}
{{ if eq $nextMonthVal 13 }}{{ $nextMonthYear = add $now.Year 1 }}{{ $nextMonthVal = 1 }}{{ end }}
{{ $nextMonthStart := (newDate $nextMonthYear $nextMonthVal 1 0 0 0) }}
{{ $endOfSeason := $nextMonthStart.Add (toDuration (mult 6 -1 | printf "%dh")) }}

{{ $msgContent := (printf "🍌 **SEASON %d HAS OFFICIALLY BEGUN!** 🍌\nThe floors are freshly waxed and the peels are scattered. It's time to slip away!\n\n🏁 **Season %d Ends:** <t:%d:F> (<t:%d:R>)\n\nGood luck to all slippers! May your falls be frequent and your prestige grow." $seasonN $seasonN $endOfSeason.Unix $endOfSeason.Unix) }}

{{ $newMsgID := sendMessageRetID (toInt64 $chanID) $msgContent }}

{{ pinMessage (toInt64 $chanID) $newMsgID }}

{{ dbSet 0 $keyPin (str $newMsgID) }}

{{ dbDel 0 $keyFlag }}