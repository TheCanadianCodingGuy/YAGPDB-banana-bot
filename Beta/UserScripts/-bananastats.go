{{/* 1. Determine who we are looking up (Self vs Mention) */}}
{{$target := .User}}
{{$targetMember := .Member}} {{/* Default to the person running the command */}}

{{if .Message.Mentions}}
    {{$target = index .Message.Mentions 0}}
    {{$targetMember = getMember $target.ID}} {{/* Fetch the specific member data for the target */}}
{{end}}

{{/* 2. Identify the target's name (Nickname > Global Name > Username) */}}
{{$rawName := $target.Username}}
{{if $targetMember}}
    {{if $targetMember.Nick}}
        {{$rawName = $targetMember.Nick}}
    {{else if $target.Globalname}}
        {{$rawName = $target.Globalname}}
    {{end}}
{{end}}

{{/* 3. Sanitize the name for Markdown */}}
{{$userName := reReplace `([*_~>|\x60])` $rawName `\$1`}}

{{/* 4. Fetch the database data */}}
{{$keyData := "BETA_banana_data"}}
{{$userData := sdict}}
{{$hasData := false}}

{{with (dbGet $target.ID $keyData)}}
    {{$userData = dict .Value | sdict}}
    {{$hasData = true}}
{{end}}

{{/* 5. Handle the scenario where the user hasn't played yet */}}
{{if not $hasData}}
    ❌ **[BETA]{{$userName}}** has not interacted with any bananas yet! 
    {{return}}
{{end}}

{{/* 6. Extract stats */}}
{{$c := toInt ($userData.Get "c")}}
{{$r := toInt ($userData.Get "r")}}
{{$f := toInt ($userData.Get "f")}}
{{$ns := toInt ($userData.Get "ns")}}
{{$gs := toInt ($userData.Get "gs")}}
{{$cs := toInt ($userData.Get "cs")}}
{{$ms := toInt ($userData.Get "ms")}}
{{$ts := toInt ($userData.Get "ts")}}
{{$rss := toInt ($userData.Get "rss")}}
{{$ss := toInt ($userData.Get "ss")}}
{{$os := toInt ($userData.Get "os")}}
{{$mcs := toInt ($userData.Get "mcs")}}
{{$hs := toInt ($userData.Get "hs")}}

{{/* 7. Format and Send */}}
{{$desc := printf "**🍌 Slip History**\nNormal Slips: **%d**\nGold Slips: **%d**\nCosmic Slips: **%d**\nMythic Slips: **%d**\n\n**💥 Chaos Events Triggered**\nTurbos: **%d**\nRank Swaps: **%d**\nSlipsplosions: **%d**\nOily Floors: **%d**\nMarket Crashes: **%d**\nHalvings: **%d**\n\n**🤸 Flip Stats**\nTotal Flips: **%d**\nCurrent Streak: **%d**\nRecord Streak: **%d**" $ns $gs $cs $ms $ts $rss $ss $os $mcs $hs $f $c $r}}

{{$embed := cembed 
    "title" (printf "[BETA] 📊 Banana Stats: %s" $userName)
    "description" $desc
    "color" 16766720 
    "thumbnail" (sdict "url" ($target.AvatarURL "256"))
}}

{{sendMessage nil $embed}}