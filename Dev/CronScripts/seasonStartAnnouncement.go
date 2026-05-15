{{/* --- CRON: 0 0 1 * * --- */}}
{{ sleep 5 }}
{{ $now := currentTime }}

{{/* 1. Setup Keys & Data */}}
{{ $keyGlobal := "TEST_banana_global" }}
{{ $keyFlag := "TEST_season_start_flag" }}
{{ $keyPin := "TEST_season_announcement" }}
{{ $chanID := 1438831937220378674 }}

{{/* Fetch the Season number */}}
{{ $global := sdict "season" 1 }}
{{ with (dbGet 0 $keyGlobal) }}{{ $global = dict .Value | sdict }}{{ end }}
{{ $seasonN := toInt ($global.Get "season") }}

{{/* 2. Unpin Management */}}
{{ with (dbGet 0 $keyPin) }}
    {{ $oldMsgID := toInt64 .Value }}
    {{/* Verify message exists before unpinning to avoid script errors */}}
    {{ if (getMessage (toInt64 $chanID) $oldMsgID) }}
        {{ unpinMessage (toInt64 $chanID) $oldMsgID }}
    {{ end }}
{{ end }}

{{/* 3. Date Calculation (End of Season) */}}
{{ $nextMonthYear := $now.Year }}
{{ $nextMonthVal := add (toInt $now.Month) 1 }}
{{ if eq $nextMonthVal 13 }}{{ $nextMonthYear = add $now.Year 1 }}{{ $nextMonthVal = 1 }}{{ end }}
{{ $nextMonthStart := (newDate $nextMonthYear $nextMonthVal 1 0 0 0) }}
{{ $endOfSeason := $nextMonthStart.Add (toDuration (mult 6 -1 | printf "%dh")) }}

{{/* 4. Send, Pin, and Save New ID */}}
{{ $msgContent := (printf "🍌 **SEASON %d HAS OFFICIALLY BEGUN!** 🍌\nThe floors are freshly waxed and the peels are scattered. It's time to slip away!\n\n🏁 **Season %d Ends:** <t:%d:F> (<t:%d:R>)\n\nGood luck to all slippers! May your falls be frequent and your prestige grow." $seasonN $seasonN $endOfSeason.Unix $endOfSeason.Unix) }}

{{/* Send message and capture ID */}}
{{ $newMsgID := sendMessageRetID (toInt64 $chanID) $msgContent }}

{{/* Pin the new message */}}
{{ pinMessage (toInt64 $chanID) $newMsgID }}

{{/* Store the ID for next season's unpinning */}}
{{ dbSet 0 $keyPin (str $newMsgID) }}

{{/* 5. Lift the gate */}}
{{ dbDel 0 $keyFlag }}