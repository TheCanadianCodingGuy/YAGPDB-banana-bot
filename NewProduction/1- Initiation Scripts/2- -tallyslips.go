{{ $now := currentTime }}

{{ $nextMonth := (newDate $now.Year (add (toInt $now.Month) 1) 1 0 0 0) }}
{{ if eq (toInt $now.Month) 12 }}{{ $nextMonth = (newDate (add $now.Year 1) 1 1 0 0 0) }}{{ end }}
{{ $lastDayOfMonth := $nextMonth.Add (toDuration (mult 24 -1 | printf "%dh")) }}

{{ $keySlips := "banana_slips" }}
{{ $keyPrestige := "prestige_global" }}
{{ $keyGlobal := "banana_global" }}
{{ $keyPin := "season_announcement" }}
{{ $chanID := 1497775909036494848 }}

{{ with (dbGet 0 $keyPin) }}
    {{ $oldMsgID := toInt64 .Value }}
    {{ if (getMessage (toInt64 $chanID) $oldMsgID) }}
        {{ unpinMessage (toInt64 $chanID) $oldMsgID }}
    {{ end }}
{{ end }}

{{ $entries := dbTopEntries $keySlips 100 0 }}
{{ $prestigeMap := sdict }}
{{ with (dbGet 0 $keyPrestige) }}{{ $prestigeMap = dict .Value | sdict }}{{ end }}

{{ $global := sdict "season" 1 }}{{ with (dbGet 0 $keyGlobal) }}{{ $global = dict .Value | sdict }}{{ end }}
{{ $seasonN := toInt ($global.Get "season") }}

{{ $currentRank := 0 }}{{ $rewardedCount := 0 }}{{ $prevScore := -1 }}{{ $hasTies := false }}{{ $winnersList := cslice }}
{{ $pointsMap := dict 1 5 2 3 3 1 }}
{{ $medalMap := dict 1 "🥇" 2 "🥈" 3 "🥉" }}
{{ $textMap := dict 1 "1st" 2 "2nd" 3 "3rd" }}

{{ range $entries }}
    {{ $uid := .UserID }}{{ $score := toInt .Value }}
    
    {{ if ne $score $prevScore }}
        {{ if ge $rewardedCount 3 }}{{ break }}{{ end }}
        {{ $currentRank = add $currentRank 1 }}
    {{ else }}
        {{ $hasTies = true }}
    {{ end }}

    {{ $points := index $pointsMap $currentRank | or 0 }}
    {{ if gt $points 0 }}
        {{ $oldTotal := toInt (index $prestigeMap (str $uid) | or 0) }}
        {{ $newTotal := add $oldTotal $points }}
        {{ $prestigeMap.Set (str $uid) $newTotal }}
        
        {{ $medal := index $medalMap $currentRank }}
        {{ $rankText := index $textMap $currentRank }}
        
        {{ $line := printf "%s %s Place: <@%d> — %d Slips | Prestige: %d (+%d) | New Total: %d" $medal $rankText $uid $score $oldTotal $points $newTotal }}
        {{ $winnersList = $winnersList.Append $line }}
        
        {{ $rewardedCount = add $rewardedCount 1 }}
        {{ $prevScore = $score }}
    {{ end }}
{{ end }}

{{ dbSet 0 $keyPrestige $prestigeMap }}

{{ $header := "🏆 **THE BANANA POINTS HAVE BEEN TALLIED! TOP 3 SLIPPERS AWARDED** 🏆\n" }}
{{ $tieNote := "" }}{{ if $hasTies }}{{ $tieNote = "*Note: Ties were detected. Slippers with the same score share the same rank level and prizes.*\n" }}{{ end }}
{{ $footer := printf "\nCongratulations to our winners of **season %d**.\n\n🛠️ Preparations for **season %d** are underway! The new season will begin at <t:%d:F> (<t:%d:R>)." $seasonN (add $seasonN 1) $nextMonth.Unix $nextMonth.Unix }}

{{ $announcement := printf "%s\n%s\n%s\n%s" $header $tieNote (joinStr "\n" $winnersList) $footer }}

{{ $newMsgID := sendMessageRetID (toInt64 $chanID) $announcement }}
{{ pinMessage (toInt64 $chanID) $newMsgID }}
{{ dbSet 0 $keyPin (str $newMsgID) }}