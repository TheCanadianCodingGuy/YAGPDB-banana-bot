{{$now := currentTime}}
{{$stM := (newDate $now.Year (toInt $now.Month) 1 0 0 0)}}
{{$ss := $stM.AddDate 0 1 0}}
{{$ssu := $ss.Unix}}
{{$lsu := sub $ssu 21600}}
{{if ge (toInt $now.Unix) (toInt $lsu)}}
    🚫 **The season has ended!**
	Tallying is currently underway and preparations for the new season are in progress.
    **Next season starts at:** <t:{{$ssu}}:F> (<t:{{$ssu}}:R>)
    {{return}}
{{end}}
{{$ks := "banana_slips"}}
{{$keyData := "banana_data"}}
{{$keyGlobal := "banana_global"}}
{{$prestigeGlobal := "prestige_global"}}
{{$cooldownDuration := 18000}}
{{$milc := 20}}
{{$malc := 70}}
{{$oMarketCrash := 7}}
{{$oHalving := 3}}
{{$oRankSwap := 4}}
{{$oTurbo := 12}}
{{$oPlosion := 14}}
{{$oOily := 14}}
{{$oMythic := 6}}
{{$oCosmic := 11}}
{{$rawName := .User.Username}}
{{if .Member.Nick}}{{$rawName = .Member.Nick}}
{{else if .User.Globalname}}{{$rawName = .User.Globalname}}{{end}}
{{$userName := reReplace `([*_~>|\x60])` $rawName `\$1`}}
{{$oldSlips := 0}}{{with (dbGet .User.ID $ks)}}{{$oldSlips = toInt .Value}}{{end}}
{{$userData := sdict "c" 0 "r" 0 "turbo" false "cd" 0 "mcs" 0 "hs" 0 "rss" 0 "gs" 0 "ts" 0 "ss" 0 "os" 0 "ms" 0 "cs" 0 "ns" 0 "f" 0}}{{with (dbGet .User.ID $keyData)}}{{$userData = dict .Value | sdict}}{{end}}
{{$global := sdict "pity" 0 "oily" false "crash" 0 "season" 0}}{{with (dbGet 0 $keyGlobal)}}{{$global = dict .Value | sdict}}{{end}}
{{$gc := false}}
{{$prestigeMap := sdict}}{{with (dbGet 0 $prestigeGlobal)}}{{$prestigeMap = dict .Value | sdict}}{{end}}
{{$myPrestige := index $prestigeMap (str .User.ID) | or 0}}
{{if gt (toInt $myPrestige) 0}}
    {{$userName = printf "(🏆%d) %s" (toInt $myPrestige) $userName}}
{{end}}
{{$te := dbTopEntries $ks 100 0}}
{{if and (gt (toInt $userData.cd) (toInt $now.Unix)) (not $userData.turbo)}}
    {{$userName}}, you look around but see no banana peel to slip on! How sad!
    **Try again <t:{{toInt $userData.cd}}:R>!**
{{else}}
    {{$currentRank := 101}}{{$myIndex := -1}}{{$prevVal := -1}}{{$rankTracker := 0}}
    {{range $i, $entry := $te}}
        {{$val := toInt .Value}}{{if ne $val $prevVal}}{{$rankTracker = add $i 1}}{{end}}{{$prevVal = $val}}
        {{if eq .User.ID $.User.ID}}{{$currentRank = $rankTracker}}{{$myIndex = $i}}{{end}}
    {{end}}
    {{$isCrashActive := gt $global.crash 0}}
    {{if $isCrashActive}}{{$global.Set "crash" (sub $global.crash 1)}}{{$gc = true}}{{end}}
    {{$isSlip := eq (randInt 2) 1}}{{if $global.oily}}{{$isSlip = true}}{{$global.Set "oily" false}}{{$gc = true}}{{end}}
    {{$addedSlips := 0}}{{$targetID := 0}}{{$targetSlips := 0}}{{$header := "**Oh no!** 🍌"}}{{$body := ""}}{{$isSwap := false}}{{$isPlosion := false}}
    {{$slipType := "ns"}}
    {{if $isSlip}}
        {{$luckyChance := $malc}}
        {{if le $currentRank 20}}
            {{$range := sub $malc $milc}}
            {{$luckyChance = add $milc (div (mult (sub $currentRank 1) $range) 19)}}
        {{end}}
        {{if ge $currentRank 3}}{{$luckyChance = add $luckyChance $global.pity}}{{end}}
        {{$wasTurbo := $userData.turbo}}
        {{$multiplier := 1}}{{if $wasTurbo}}{{$multiplier = 2}}{{$userData.Set "turbo" false}}{{end}}
        {{$isLucky := le (randInt 1 10001) (mult $luckyChance 100)}}
        {{if and $isLucky (not $isCrashActive)}}
            {{if ge $currentRank 3}}{{$global.Set "pity" 0}}{{$gc = true}}{{end}}
            {{$tCrash := $oMarketCrash}}
            {{$tHalving := add $tCrash $oHalving}}
            {{$tSwap := add $tHalving $oRankSwap}}
            {{$tTurbo := add $tSwap $oTurbo}}
            {{$tPlosion := add $tTurbo $oPlosion}}
            {{$tOily := add $tPlosion $oOily}}
            {{$tMythic := add $tOily $oMythic}}
            {{$tCosmic := add $tMythic $oCosmic}}
            {{$roll := randInt 1 101}}
            {{if $wasTurbo}}
                {{$roll = randInt (add $tOily 1) 101}}
            {{end}}
            {{if le $roll $tCrash}}
                {{$slipType = "mcs"}}
                {{$newCrashVal := add $global.crash 3}}
                {{$global.Set "crash" $newCrashVal}}{{$gc = true}}
                {{$body = printf "📉 **MARKET CRASH!** %s caused the economy to collapse! Slips grant 0 and flips grant 5 for the next %d rolls! 🚨" $userName (toInt $newCrashVal)}}
            {{else if le $roll $tHalving}}
                {{$slipType = "hs"}}
                {{$addedSlips := sub (div $oldSlips 2) $oldSlips}}
                {{$multiplier = 0}}
                {{$body = printf "📉💥 **BIG RIP!** %s's entire slip count was just **HALVED**!" $userName}}
            {{else if and (le $roll $tSwap) (ge $myIndex 0)}}
                {{$isSwap = true}}{{$multiplier = 0}}
                {{$offset := randInt 1 3}}
                {{$targetIdx := 0}}
                {{if le $currentRank 10}}
                    {{$targetIdx = add $myIndex $offset}}
                    {{if ge $targetIdx (len $te)}}{{$targetIdx = sub (len $te) 1}}{{end}}
                {{else}}
                    {{$targetIdx = sub $myIndex $offset}}
                    {{if lt $targetIdx 0}}{{$targetIdx = 0}}{{end}}
                {{end}}
                {{if ne $targetIdx $myIndex}}
                    {{$slipType = "rss"}}
                    {{$target := index $te $targetIdx}}
                    {{$targetID = $target.User.ID}}{{$targetSlips = toInt $target.Value}}
                    {{$addedSlips = sub $targetSlips $oldSlips}}
                    {{if lt $targetIdx $myIndex}}
                        {{$body = printf "🚀 **RANK SWAP!** %s vaulted upwards and swapped their **%d** slips with <@%d>'s **%d** slips! 🏎️💨" $userName $oldSlips $targetID $targetSlips}}
                    {{else}}
                        {{$body = printf "📉 **RANK SWAP!** %s fumbled and swapped their **%d** slips with <@%d>'s **%d** slips from below! 🤼‍♂️" $userName $oldSlips $targetID $targetSlips}}
                    {{end}}
                {{else}}
                    {{$slipType = "gs"}}
                    {{$isSwap = false}}{{$multiplier = 2}}
                    {{$body = printf "😮 %s tried to swap, but no one was there! A **Golden Peel** was found instead." $userName}}
                {{end}}
            {{else if le $roll $tTurbo}}
                {{$slipType = "ts"}}
                {{$userData.Set "turbo" true}}
                {{$body = printf "🔥 **TURBO OVERDRIVE!** %s hit top speeds! Cooldown is reset and next slip is worth **DOUBLE**! 🏎️💨" $userName}}
            {{else if le $roll $tPlosion}}
                {{$slipType = "ss"}}
                {{$isPlosion = true}}
                {{if eq $currentRank 1}}
                    {{$addedSlips = -3}}{{$multiplier = 0}}
                    {{$body = printf "💣 **SLIPSPLOSION!** %s hit the floor so hard they destroyed 3 of their own slips! 🌋" $userName}}
                {{else if gt $myIndex 0}}
                    {{$target := index $te (sub $myIndex 1)}}
                    {{$targetID = $target.User.ID}}{{$targetSlips = sub (toInt $target.Value) 3}}
                    {{$body = printf "💣 **SLIPSPLOSION!** %s sent a shockwave that destroyed 3 slips from <@%d>! 🌋" $userName $targetID}}
                {{end}}
            {{else if le $roll $tOily}}
                {{$slipType = "os"}}
                {{$global.Set "oily" true}}{{$gc = true}}
                {{$body = printf "🛢️💀 **OILY FLOOR!** %s spilled grease! The next roller is **DOOMED** to slip." $userName}}
            {{else if le $roll $tMythic}}
                {{$slipType = "ms"}}
                {{$header = "**HOLY CRAP!** 😱😱😱"}}
                {{$multiplier = mult $multiplier 10}}
                {{$body = printf "🎆🎆🎆 ***%s JUST SLIPPED ON A DANG MYTHIC PEEL WORTH %d SLIPS!*** 🎆🎆🎆" $userName $multiplier}}
            {{else if le $roll $tCosmic}}
                {{$slipType = "cs"}}
                {{$multiplier = mult $multiplier 5}}
                {{$body = printf "🌟🌌🌟 %s just slipped on a rare ***cosmic peel*** worth ***%d slips!*** 🌟🌌🌟" $userName $multiplier}}
            {{else}}
                {{$slipType = "gs"}}
                {{$multiplier = mult $multiplier 2}}
                {{$body = printf "😮 %s just slipped on a peel made of **pure gold** worth **%d slips!**" $userName $multiplier}}
            {{end}}
            {{if $wasTurbo}}
                {{$body = printf "%s\n🏎️💨 **TURBO BONUS!** %s's slip was worth **DOUBLE**!" $body $userName}}
            {{end}}
        {{else}}
            {{$global.Set "pity" (add $global.pity 10)}}{{$gc = true}}
            {{if $wasTurbo}}
                {{$body = printf "🏎️💨 **TURBO BONUS!** %s slipped on a regular peel, but was worth **DOUBLE**!" $userName}}
            {{else}}
                {{$body = printf "   %s just slipped on a peel!" $userName}}
            {{end}}
        {{end}}
        {{if $isCrashActive}}{{$addedSlips = 0}}{{$body = printf "\n🚨 **MARKET CRASH ACTIVE!** %s's slip was worthless! +0 Slips. (%d market crash rolls remaining)" $userName (toInt $global.crash)}}
        {{else if not $isSwap}}{{$addedSlips = add $addedSlips $multiplier}}{{end}}
        {{$newSlips := add $oldSlips $addedSlips}}{{$userData.Set "c" 0}}
        {{$userData.Set $slipType (add (toInt ($userData.Get $slipType)) 1)}}
        {{dbSet .User.ID $ks (str $newSlips)}}
        {{if and $isSwap $targetID}}{{dbSet $targetID $ks (str $oldSlips)}}
        {{else if and $isPlosion $targetID}}{{dbSet $targetID $ks (str $targetSlips)}}{{end}}
        {{$finalRank := 1}}{{range $te}}{{if gt (toInt .Value) $newSlips}}{{$finalRank = add $finalRank 1}}{{end}}{{end}}
        {{$header}}
        {{$body}}
        You just **{{if lt (toInt $addedSlips) 0}}lost {{mult (toInt $addedSlips) -1}}{{else}}gained {{toInt $addedSlips}}{{end}} slip{{if or (gt (toInt $addedSlips) 1) (lt (toInt $addedSlips) -1)}}s{{end}}** for a new total of **{{$newSlips}}** slips! Watch your step!
        Your clumsiness places you at rank **#{{$finalRank}}**!
    {{else}}
        {{$earnedSlips := 0}}{{if $isCrashActive}}{{$earnedSlips = 5}}{{end}}
        {{$newStreak := add $userData.c 1}}{{$userData.Set "c" $newStreak}}
        {{$userData.Set "f" (add (toInt ($userData.Get "flip")) 1)}}
        {{$streakStatus := "chasing your"}}{{if gt $newStreak $userData.r}}{{$streakStatus = "setting a new"}}{{$userData.Set "r" $newStreak}}
        {{else if eq $newStreak $userData.r}}{{$streakStatus = "pulling even with your"}}{{end}}
        {{$turboWaste := ""}}
        {{if $userData.turbo}}
            {{$userData.Set "turbo" false}}
            {{$turboWaste = printf "\n🏎️💨 **TURBO WASTED!** %s dodged the peel, meaning the Turbo boost was for nothing!" $userName}}
        {{end}}
        **CLEAN!** 🤸
        {{$userName}} dodged the peel with a backflip!
        You’ve now evaded gravity **{{ $newStreak }}** time{{ if ne $newStreak 1 }}s in a row{{ end }},
        {{ $streakStatus }} all-time personal record of **{{ $userData.r }}**.
        {{$turboWaste}}
        {{- if $isCrashActive -}}
            {{- $newTotal := add $oldSlips $earnedSlips -}}
            {{- dbSet .User.ID $ks (str $newTotal) -}}
            {{"\n"}}🚨 **MARKET CRASH ACTIVE!** {{$userName}} stole **5 Slips** from the market! Adding them to their new total of **{{$newTotal}}** slips! ({{$global.crash}} market crash rolls remaining)
        {{- end -}}
    {{end}}
    {{if not $userData.turbo}}{{$userData.Set "cd" (add currentTime.Unix $cooldownDuration)}}{{end}}
    {{dbSet .User.ID $keyData $userData}}
    {{if $gc}}{{dbSet 0 $keyGlobal $global}}{{end}}
{{end}}