{{ $now := currentTime }}
{{ $startOfThisMonth := (newDate $now.Year (toInt $now.Month) 1 0 0 0) }}
{{ $seasonStart := $startOfThisMonth.AddDate 0 1 0 }}
{{ $seasonStartUnix := $seasonStart.Unix }}
{{ $lockoutStartUnix := sub $seasonStartUnix 21600 }}
{{ if ge (toInt $now.Unix) (toInt $lockoutStartUnix) }}
    🚫 **The season has ended!** 
    Tallying is currently underway and preparations for the new season are in progress.
    **Next season starts at:** <t:{{ $seasonStartUnix }}:F> (<t:{{ $seasonStartUnix }}:R>)
    {{ return }}
{{ end }}

{{ $keySlips := "TEST_banana_slips" }}
{{ $keyCD := "TEST_banana_cooldown" }}
{{ $keyData := "TEST_banana_data" }}
{{ $keyGlobal := "TEST_banana_global" }}
{{ $prestigeGlobal := "TEST_prestige_global" }}
{{ $cooldownDuration := 36000 }}
{{ $expiration := 31536000 }}
{{ $minLuckyChance := 5 }}
{{ $maxLuckyChance := 45 }}
{{ $oddsMarketCrash := 3 }}
{{ $oddsHalving := 1 }}
{{ $oddsRankSwap := 5 }}
{{ $oddsTurbo := 8 }}
{{ $oddsPlosion := 10 }}
{{ $oddsOily := 12 }}
{{ $oddsMythic := 3 }}
{{ $oddsCosmic := 9 }}

{{ $rawName := .User.Username }}
{{ if .Member.Nick }}{{ $rawName = .Member.Nick }}
{{ else if .User.Globalname }}{{ $rawName = .User.Globalname }}{{ end }}
{{ $userName := reReplace `([*_~>|\x60])` $rawName `\$1` }}

{{ $cooldownData := dbGet .User.ID $keyCD }}
{{ $oldSlips := 0 }}{{ with (dbGet .User.ID $keySlips) }}{{ $oldSlips = toInt .Value }}{{ end }}
{{ $userData := sdict "c" 0 "r" 0 "turbo" false }}{{ with (dbGet .User.ID $keyData) }}{{ $userData = dict .Value | sdict }}{{ end }}
{{ $global := sdict "pity" 0 "oily" false "crash" 0 }}{{ with (dbGet 0 $keyGlobal) }}{{ $global = dict .Value | sdict }}{{ end }}

{{ $prestigeMap := sdict }}
{{ with (dbGet 0 $prestigeGlobal) }}{{ $prestigeMap = dict .Value | sdict }}{{ end }}
{{ $myPrestige := index $prestigeMap (str .User.ID) | or 0 }}

{{ if gt (toInt $myPrestige) 0 }}
    {{ $userName = printf "(🏆%d) %s" (toInt $myPrestige) $userName }}
{{ end }}

{{ $topEntries := dbTopEntries $keySlips 100 0 }}
{{ $cooldownData := false }}
{{ if and $cooldownData (not $userData.turbo) }}
    {{ $userName }}, you look around but see no banana peel to slip on! 
    **Try again <t:{{ toInt $cooldownData.Value }}:R>!**
{{ else }}
    {{ $currentRank := 101 }}{{ $myIndex := -1 }}{{ $prevVal := -1 }}{{ $rankTracker := 0 }}
    {{ range $i, $entry := $topEntries }}
        {{ $val := toInt .Value }}{{ if ne $val $prevVal }}{{ $rankTracker = add $i 1 }}{{ end }}{{ $prevVal = $val }}
        {{ if eq .User.ID $.User.ID }}{{ $currentRank = $rankTracker }}{{ $myIndex = $i }}{{ end }}
    {{ end }}

    {{ $isCrashActive := gt $global.crash 0 }}
    {{ if $isCrashActive }}{{ $global.Set "crash" (sub $global.crash 1) }}{{ end }}

    {{ $isSlip := eq (randInt 2) 1 }}{{ if $global.oily }}{{ $isSlip = true }}{{ $global.Set "oily" false }}{{ end }}

    {{ $addedSlips := 0 }}{{ $targetID := 0 }}{{ $targetSlips := 0 }}{{ $header := "**Oh no!** 🍌" }}{{ $body := "" }}{{ $isSwap := false }}{{ $isPlosion := false }}

    {{ if $isSlip }}
        {{ $luckyChance := $maxLuckyChance }}
        {{ if le $currentRank 20 }}
            {{ $range := sub $maxLuckyChance $minLuckyChance }}
            {{ $luckyChance = add $minLuckyChance (div (mult (sub $currentRank 1) $range) 19) }}
        {{ end }}
        {{ if ge $currentRank 11 }}{{ $luckyChance = add $luckyChance $global.pity }}{{ end }}

        {{ $multiplier := 1 }}{{ if $userData.turbo }}{{ $multiplier = 2 }}{{ $userData.Set "turbo" false }}{{ end }}
        {{ $isLucky := le (randInt 1 10001) (mult $luckyChance 100) }}

        {{ if $isLucky }}
            {{ if ge $currentRank 11 }}{{ $global.Set "pity" 0 }}{{ end }}
            {{ $roll := randInt 1 101 }}
            
            {{ $tCrash := $oddsMarketCrash }}
            {{ $tHalving := add $tCrash $oddsHalving }}
            {{ $tSwap := add $tHalving $oddsRankSwap }}
            {{ $tTurbo := add $tSwap $oddsTurbo }}
            {{ $tPlosion := add $tTurbo $oddsPlosion }}
            {{ $tOily := add $tPlosion $oddsOily }}
            {{ $tMythic := add $tOily $oddsMythic }}
            {{ $tCosmic := add $tMythic $oddsCosmic }}

            {{ if le $roll $tCrash }}
                {{ $global.Set "crash" 3 }}
                {{ $body = "\n📉 **MARKET CRASH!** The slip economy is broken! Backflips grant +5, Slips grant 0 for the next 3 rolls! 🚨" }}
            {{ else if le $roll $tHalving }}
                {{ $oldSlips = div $oldSlips 2 }}
                {{ $body = "\n📉💥 **CATASTROPHIC ERROR!** Your entire slip count was just **HALVED**!" }}
            {{ else if and (le $roll $tSwap) (ge $myIndex 0) }}
                {{ $isSwap = true }}
                {{ $offset := randInt 1 6 }}
                {{ $targetIdx := 0 }}
                {{ if le $currentRank 10 }}
                    {{ $targetIdx = add $myIndex $offset }}
                    {{ if ge $targetIdx (len $topEntries) }}{{ $targetIdx = sub (len $topEntries) 1 }}{{ end }}
                {{ else }}
                    {{ $targetIdx = sub $myIndex $offset }}
                    {{ if lt $targetIdx 0 }}{{ $targetIdx = 0 }}{{ end }}
                {{ end }}

                {{ if ne $targetIdx $myIndex }}
                    {{ $target := index $topEntries $targetIdx }}
                    {{ $targetID = $target.User.ID }}{{ $targetSlips = toInt $target.Value }}
                    {{ $addedSlips = sub $targetSlips $oldSlips }}
                    {{ if lt $targetIdx $myIndex }}
                        {{ $body = print "\n🚀 **RANK SWAP!** You vaulted upwards and stole the position of <@" $targetID ">! 🏎️💨" }}
                    {{ else }}
                        {{ $body = print "\n📉 **RANK SWAP!** You fumbled and swapped places with <@" $targetID "> from below! 🤼‍♂️" }}
                    {{ end }}
                {{ else }}
                    {{ $isSwap = false }}{{ $multiplier = mult $multiplier 2 }}
                    {{ $body = print "\n😮 " $userName " tried to swap, but no one was there! You got a **Golden Peel** instead." }}
                {{ end }}
            {{ else if le $roll $tTurbo }}
                {{ $userData.Set "turbo" true }}
                {{ $body = "\n🔥 **TURBO OVERDRIVE!** Cooldown instantly refreshed! Your next slip will be **DOUBLED**! 🏎️💨" }}
            {{ else if le $roll $tPlosion }}
                {{ $isPlosion = true }}
                {{ if eq $currentRank 1 }}
                    {{ $addedSlips = 0 }}{{ $multiplier = 0 }}
                    {{ $body = "\n💣 **SLIP-PLOSION!** You hit the floor so hard it nullified your own slip! 🌋" }}
                {{ else if gt $myIndex 0 }}
                    {{ $target := index $topEntries (sub $myIndex 1) }}
                    {{ $targetID = $target.User.ID }}{{ $targetSlips = sub (toInt $target.Value) 1 }}
                    {{ $body = print "\n💣 **SLIP-PLOSION!** The impact destroyed 1 slip from <@" $targetID ">! 🌋" }}
                {{ end }}
            {{ else if le $roll $tOily }}
                {{ $global.Set "oily" true }}
                {{ $body = "\n🛢️💀 **OILY FLOOR!** You spilled grease! The next roller is **DOOMED** to slip." }}
            {{ else if le $roll $tMythic }}
                {{ $header = "**HOLY CRAP!** 😱😱😱" }}
                {{ $multiplier = mult $multiplier 10 }}
                {{ $body = print "\n🎆🎆🎆 ***" $userName " JUST SLIPPED ON A DANG MYTHIC PEEL WORTH " $multiplier " SLIPS!*** 🎆🎆🎆" }}
            {{ else if le $roll $tCosmic }}
                {{ $multiplier = mult $multiplier 5 }}
                {{ $body = print "\n🌟🌌🌟 " $userName " just slipped on a rare ***cosmic peel*** worth ***" $multiplier " slips!*** 🌟🌌🌟" }}
            {{ else }}
                {{ $multiplier = mult $multiplier 2 }}
                {{ $body = print "\n😮 " $userName " just slipped on a peel made of **pure gold** worth **" $multiplier " slips!**" }}
            {{ end }}
        {{ else }}
            {{ $global.Set "pity" (add $global.pity 2) }}{{ $body = print "\n" $userName " just slipped on a peel!" }}
        {{ end }}

        {{ if $isCrashActive }}{{ $addedSlips = 0 }}{{ $body = print "\n🚨 **MARKET CRASH ACTIVE!** Your slip was worthless! +0 Slips." }}
        {{ else if not $isSwap }}{{ $addedSlips = add $addedSlips $multiplier }}{{ end }}

        {{ $newSlips := add $oldSlips $addedSlips }}{{ $userData.Set "c" 0 }}
        {{ dbSetExpire .User.ID $keySlips (str $newSlips) $expiration }}

        {{ if and $isSwap $targetID }}{{ dbSetExpire $targetID $keySlips (str $oldSlips) $expiration }}
        {{ else if and $isPlosion $targetID }}{{ dbSetExpire $targetID $keySlips (str $targetSlips) $expiration }}{{ end }}

        {{ $finalRank := 1 }}{{ range $topEntries }}{{ if gt (toInt .Value) $newSlips }}{{ $finalRank = add $finalRank 1 }}{{ end }}{{ end }}

        [TEST] {{ $header }}
        {{ $body }}
        You just **gained {{ $addedSlips }}slip{{ if ne $addedSlips 1 }}s{{ end }}** for a new total of **{{ $newSlips }}** slips! Watch your step!
        Rank **#{{ $finalRank }}**

    {{ else }}
        {{ $earnedSlips := 0 }}{{ if $isCrashActive }}{{ $earnedSlips = 5 }}{{ end }}
        {{ $newStreak := add $userData.c 1 }}{{ $userData.Set "c" $newStreak }}
        {{ $streakStatus := "chasing your" }}{{ if gt $newStreak $userData.r }}{{ $streakStatus = "setting a new" }}{{ $userData.Set "r" $newStreak }}
        {{ else if eq $newStreak $userData.r }}{{ $streakStatus = "pulling even with your" }}{{ end }}
        
        [TEST] **CLEAN!** 🤸
        {{ $userName }} dodged the peel! Streak: **{{ $newStreak }}**, {{ $streakStatus }} record of **{{ $userData.r }}**.
        {{ if $isCrashActive }}🚨 **MARKET CRASH!** Stole **5 Slips**!{{ dbSetExpire .User.ID $keySlips (str (add $oldSlips $earnedSlips)) $expiration }}{{ end }}
    {{ end }}

    {{ if not $userData.turbo }}{{ dbSetExpire .User.ID $keyCD (str (add currentTime.Unix $cooldownDuration)) $cooldownDuration }}{{ end }}
    {{ dbSetExpire .User.ID $keyData $userData $expiration }}{{ dbSet 0 $keyGlobal $global }}
{{ end }}