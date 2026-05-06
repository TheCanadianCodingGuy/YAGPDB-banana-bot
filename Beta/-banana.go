{{$now := currentTime}}
{{$startOfThisMonth := (newDate $now.Year (toInt $now.Month) 1 0 0 0)}}
{{$seasonStart := $startOfThisMonth.AddDate 0 1 0}}
{{$seasonStartUnix := $seasonStart.Unix}}
{{$lockoutStartUnix := sub $seasonStartUnix 21600}}
{{if ge (toInt $now.Unix) (toInt $lockoutStartUnix)}}
    🚫 **The season has ended!** 
    Tallying is currently underway and preparations for the new season are in progress.
    **Next season starts at:** <t:{{$seasonStartUnix}}:F> (<t:{{$seasonStartUnix}}:R>)
    {{return}}
{{end}}
{{$keySlips := "TEST_banana_slips"}}
{{$keyData := "TEST_banana_data"}}
{{$keyGlobal := "TEST_banana_global"}}
{{$prestigeGlobal := "TEST_prestige_global"}}
{{$cooldownDuration := 36000}}
{{$expiration := 31536000}}
{{$minLuckyChance := 5}}
{{$maxLuckyChance := 45}}
{{$oddsMarketCrash := 3}}
{{$oddsHalving := 1}}
{{$oddsRankSwap := 5}}
{{$oddsTurbo := 8}}
{{$oddsPlosion := 10}}
{{$oddsOily := 12}}
{{$oddsMythic := 3}}
{{$oddsCosmic := 9}}
{{$rawName := .User.Username}}
{{if .Member.Nick}}{{$rawName = .Member.Nick}}
{{else if .User.Globalname}}{{$rawName = .User.Globalname}}{{end}}
{{$userName := reReplace `([*_~>|\x60])` $rawName `\$1`}}
{{$oldSlips := 0}}{{with (dbGet .User.ID $keySlips)}}{{$oldSlips = toInt .Value}}{{end}}
{{$userData := sdict "c" 0 "r" 0 "turbo" false "cd" 0}}{{with (dbGet .User.ID $keyData)}}{{$userData = dict .Value | sdict}}{{end}}
{{$global := sdict "pity" 0 "oily" false "crash" 0}}{{with (dbGet 0 $keyGlobal)}}{{$global = dict .Value | sdict}}{{end}}
{{$globalChanged := false}}
{{$prestigeMap := sdict}}{{with (dbGet 0 $prestigeGlobal)}}{{$prestigeMap = dict .Value | sdict}}{{end}}
{{$myPrestige := index $prestigeMap (str .User.ID) | or 0}}
{{if gt (toInt $myPrestige) 0}}
    {{$userName = printf "(🏆%d) %s" (toInt $myPrestige) $userName}}
{{end}}
{{$topEntries := dbTopEntries $keySlips 100 0}}
{{$userData.Set "cd" (sub $now.Unix 5)}}
{{if and (gt (toInt $userData.cd) (toInt $now.Unix)) (not $userData.turbo)}}
    {{$userName}}, you look around but see no banana peel to slip on!
    **Try again <t:{{toInt $userData.cd}}:R>!**
{{else}}
    {{$currentRank := 101}}{{$myIndex := -1}}{{$prevVal := -1}}{{$rankTracker := 0}}
    {{range $i, $entry := $topEntries}}
        {{$val := toInt .Value}}{{if ne $val $prevVal}}{{$rankTracker = add $i 1}}{{end}}{{$prevVal = $val}}
        {{if eq .User.ID $.User.ID}}{{$currentRank = $rankTracker}}{{$myIndex = $i}}{{end}}
    {{end}}
    {{$isCrashActive := gt $global.crash 0}}
    {{if $isCrashActive}}{{$global.Set "crash" (sub $global.crash 1)}}{{$globalChanged = true}}{{end}}
    {{$isSlip := eq (randInt 2) 1}}{{if $global.oily}}{{$isSlip = true}}{{$global.Set "oily" false}}{{end}}
    {{$addedSlips := 0}}{{$targetID := 0}}{{$targetSlips := 0}}{{$header := "**Oh no!** 🍌"}}{{$body := ""}}{{$isSwap := false}}{{$isPlosion := false}}
    {{$isSlip := true}}
    {{if $isSlip}}
        {{$luckyChance := $maxLuckyChance}}
        {{if le $currentRank 20}}
            {{$range := sub $maxLuckyChance $minLuckyChance}}
            {{$luckyChance = add $minLuckyChance (div (mult (sub $currentRank 1) $range) 19)}}
        {{end}}
        {{if ge $currentRank 11}}{{$luckyChance = add $luckyChance $global.pity}}{{end}}
        {{$wasTurbo := $userData.turbo}}
        {{$multiplier := 1}}{{if $wasTurbo}}{{$multiplier = 2}}{{$userData.Set "turbo" false}}{{end}}
        {{$isLucky := le (randInt 1 10001) (mult $luckyChance 100)}}
        {{$isLucky := true}}
        {{if and $isLucky (not $isCrashActive)}}
            {{if ge $currentRank 11}}{{$global.Set "pity" 0}}{{$globalChanged = true}}{{end}}
            {{$tCrash := $oddsMarketCrash}}
            {{$tHalving := add $tCrash $oddsHalving}}
            {{$tSwap := add $tHalving $oddsRankSwap}}
            {{$tTurbo := add $tSwap $oddsTurbo}}
            {{$tPlosion := add $tTurbo $oddsPlosion}}
            {{$tOily := add $tPlosion $oddsOily}}
            {{$tMythic := add $tOily $oddsMythic}}
            {{$tCosmic := add $tMythic $oddsCosmic}}
            {{$roll := randInt 1 101}}
            {{$roll := 10}}
            {{if $wasTurbo}}
                {{$roll = randInt (add $tOily 1) 101}}
            {{end}}
            {{if le $roll $tCrash}}
                {{$newCrashVal := add $global.crash 3}}
                {{$global.Set "crash" $newCrashVal}}{{$globalChanged = true}}
                {{$body = printf "\n📉 **MARKET CRASH!** %s caused the economy to collapse! Slips grant 0 and flips grant 5 for the next %d rolls! 🚨" $userName (toInt $newCrashVal)}}
            {{else if le $roll $tHalving}}
                {{$preHalve := $oldSlips}}
                {{$oldSlips = div $oldSlips 2}}
                {{$addedSlips = sub $oldSlips $preHalve}}
                {{$multiplier = 0}}
                {{$body = printf "\n📉💥 **CATASTROPHIC ERROR!** %s's entire slip count was just **HALVED**!" $userName}}
            {{else if and (le $roll $tSwap) (ge $myIndex 0)}}
                {{$isSwap = true}}{{$multiplier = 0}}
                {{$offset := randInt 1 6}}
                {{$targetIdx := 0}}
                {{if le $currentRank 10}}
                    {{$targetIdx = add $myIndex $offset}}
                    {{if ge $targetIdx (len $topEntries)}}{{$targetIdx = sub (len $topEntries) 1}}{{end}}
                {{else}}
                    {{$targetIdx = sub $myIndex $offset}}
                    {{if lt $targetIdx 0}}{{$targetIdx = 0}}{{end}}
                {{end}}
                {{if ne $targetIdx $myIndex}}
                    {{$target := index $topEntries $targetIdx}}
                    {{$targetID = $target.User.ID}}{{$targetSlips = toInt $target.Value}}
                    {{$addedSlips = sub $targetSlips $oldSlips}}
                    {{if lt $targetIdx $myIndex}}
                        {{$body = printf "\n🚀 **RANK SWAP!** %s vaulted upwards and swapped their **%d** slips with <@%d>'s **%d** slips! 🏎️💨" $userName $oldSlips $targetID $targetSlips}}
                    {{else}}
                        {{$body = printf "\n📉 **RANK SWAP!** %s fumbled and swapped their **%d** slips with <@%d>'s **%d** slips from below! 🤼‍♂️" $userName $oldSlips $targetID $targetSlips}}
                    {{end}}
                {{else}}
                    {{$isSwap = false}}{{$multiplier = 2}}
                    {{$body = printf "\n😮 %s tried to swap, but no one was there! A **Golden Peel** was found instead." $userName}}
                {{end}}
            {{else if le $roll $tTurbo}}
                {{$userData.Set "turbo" true}}
                {{$body = printf "\n🔥 **TURBO OVERDRIVE!** %s hit top speeds! Cooldown is reset and next slip is worth **DOUBLE**! 🏎️💨" $userName}}
            {{else if le $roll $tPlosion}}
                {{$isPlosion = true}}
                {{if eq $currentRank 1}}
                    {{$addedSlips = -2}}{{$multiplier = 0}}
                    {{$body = printf "\n💣 **SLIP-PLOSION!** %s hit the floor so hard they destroyed 2 of their own slips! 🌋" $userName}}
                {{else if gt $myIndex 0}}
                    {{$target := index $topEntries (sub $myIndex 1)}}
                    {{$targetID = $target.User.ID}}{{$targetSlips = sub (toInt $target.Value) 2}}
                    {{$body = printf "\n💣 **SLIP-PLOSION!** %s sent a shockwave that destroyed 2 slips from <@%d>! 🌋" $userName $targetID}}
                {{end}}
            {{else if le $roll $tOily}}
                {{$global.Set "oily" true}}{{$globalChanged = true}}
                {{$body = printf "\n🛢️💀 **OILY FLOOR!** %s spilled grease! The next roller is **DOOMED** to slip." $userName}}
            {{else if le $roll $tMythic}}
                {{$header = "**HOLY CRAP!** 😱😱😱"}}
                {{$multiplier = mult $multiplier 10}}
                {{$body = printf "\n🎆🎆🎆 ***%s JUST SLIPPED ON A DANG MYTHIC PEEL WORTH %d SLIPS!*** 🎆🎆🎆" $userName $multiplier}}
            {{else if le $roll $tCosmic}}
                {{$multiplier = mult $multiplier 5}}
                {{$body = printf "\n🌟🌌🌟 %s just slipped on a rare ***cosmic peel*** worth ***%d slips!*** 🌟🌌🌟" $userName $multiplier}}
            {{else}}
                {{$multiplier = mult $multiplier 2}}
                {{$body = printf "\n😮 %s just slipped on a peel made of **pure gold** worth **%d slips!**" $userName $multiplier}}
            {{end}}
            {{if $wasTurbo}}
                {{$body = printf "%s\n🏎️💨 **TURBO BONUS!** %s's slip was worth **DOUBLE**!" $body $userName}}
            {{end}}
        {{else}}
            {{$global.Set "pity" (add $global.pity 2)}}{{$globalChanged = true}}
            {{if $wasTurbo}}
                {{$body = printf "\n🏎️💨 **TURBO BONUS!** %s slipped on a regular peel, but was worth **DOUBLE**!" $userName}}
            {{else}}
                {{$body = printf "\n%s just slipped on a peel!" $userName}}
            {{end}}
        {{end}}
        {{if $isCrashActive}}{{$addedSlips = 0}}{{$body = printf "\n🚨 **MARKET CRASH ACTIVE!** %s's slip was worthless! +0 Slips. (%d market crash rolls remaining)" $userName (toInt $global.crash)}}
        {{else if not $isSwap}}{{$addedSlips = add $addedSlips $multiplier}}{{end}}
        {{$newSlips := add $oldSlips $addedSlips}}{{$userData.Set "c" 0}}
        {{dbSetExpire .User.ID $keySlips (str $newSlips) $expiration}}
        {{if and $isSwap $targetID}}{{dbSetExpire $targetID $keySlips (str $oldSlips) $expiration}}
        {{else if and $isPlosion $targetID}}{{dbSetExpire $targetID $keySlips (str $targetSlips) $expiration}}{{end}}
        {{$finalRank := 1}}{{range $topEntries}}{{if gt (toInt .Value) $newSlips}}{{$finalRank = add $finalRank 1}}{{end}}{{end}}
        [TEST] {{$header}}
        {{$body}}
        {{$userName}} just **{{if lt (toInt $addedSlips) 0}}lost {{mult (toInt $addedSlips) -1}}{{else}}gained {{toInt $addedSlips}}{{end}} slip{{if or (gt (toInt $addedSlips) 1) (lt (toInt $addedSlips) -1)}}s{{end}}** for a new total of **{{$newSlips}}** slips! Watch your step!
        Rank **#{{$finalRank}}**
    {{else}}
        {{$earnedSlips := 0}}{{if $isCrashActive}}{{$earnedSlips = 5}}{{end}}
        {{$newStreak := add $userData.c 1}}{{$userData.Set "c" $newStreak}}
        {{$streakStatus := "chasing their"}}{{if gt $newStreak $userData.r}}{{$streakStatus = "setting a new"}}{{$userData.Set "r" $newStreak}}
        {{else if eq $newStreak $userData.r}}{{$streakStatus = "pulling even with their"}}{{end}}
        {{$turboWaste := ""}}
        {{if $userData.turbo}}
            {{$userData.Set "turbo" false}}
            {{$turboWaste = printf "\n🏎️💨 **TURBO WASTED!** %s dodged the peel, meaning the Turbo boost was for nothing!" $userName}}
        {{end}}
        [TEST] **CLEAN!** 🤸
        {{$userName}} dodged the peel! Streak: **{{$newStreak}}**, {{$streakStatus}} record of **{{$userData.r}}**. {{$turboWaste}}
        {{- if $isCrashActive -}}
            {{- $newTotal := add $oldSlips $earnedSlips -}}
            {{- dbSetExpire .User.ID $keySlips (str $newTotal) $expiration -}}
            {{"\n"}}🚨 **MARKET CRASH ACTIVE!** {{$userName}} stole **5 Slips** from the market! Adding them to their new total of **{{$newTotal}}** slips! ({{$global.crash}} market crash rolls remaining)
        {{- end -}}
    {{end}}
    {{if not $userData.turbo}}{{$userData.Set "cd" (add currentTime.Unix $cooldownDuration)}}{{end}}
    {{dbSetExpire .User.ID $keyData $userData $expiration}}
    {{if $globalChanged}}{{dbSet 0 $keyGlobal $global}}{{end}}
{{end}}