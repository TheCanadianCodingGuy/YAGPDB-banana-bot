{{$target := .User}}
{{$targetMember := .Member}} {{/* Default to the person running the command */}}

{{if .Message.Mentions}}
    {{$target = index .Message.Mentions 0}}
    {{$targetMember = getMember $target.ID}} 
{{end}}

{{$rawName := $target.Username}}
{{if $targetMember}}
    {{if $targetMember.Nick}}
        {{$rawName = $targetMember.Nick}}
    {{else if $target.Globalname}}
        {{$rawName = $target.Globalname}}
    {{end}}
{{end}}

{{$userName := reReplace `([*_~>|\x60])` $rawName `\$1`}}

{{$keyData := "banana_data"}}
{{$userData := sdict}}
{{$hasData := false}}

{{with (dbGet $target.ID $keyData)}}
    {{$userData = dict .Value | sdict}}
    {{$hasData = true}}
{{end}}

{{if not $hasData}}
    ❌ **{{$userName}}** has not interacted with any bananas yet! 
    {{return}}
{{end}}

{{$f := toInt ($userData.Get "g_f")}}
{{$ns := toInt ($userData.Get "g_ns")}}
{{$gs := toInt ($userData.Get "g_gs")}}
{{$cs := toInt ($userData.Get "g_cs")}}
{{$ms := toInt ($userData.Get "g_ms")}}
{{$ts := toInt ($userData.Get "g_ts")}}
{{$rss := toInt ($userData.Get "g_rss")}}
{{$ss := toInt ($userData.Get "g_ss")}}
{{$os := toInt ($userData.Get "g_os")}}
{{$mcs := toInt ($userData.Get "g_mcs")}}
{{$hs := toInt ($userData.Get "g_hs")}}

{{$desc := printf "**🍌 Slip History **\nNormal Slips: **%d**\nGold Slips: **%d**\nCosmic Slips: **%d**\nMythic Slips: **%d**\n\n**💥 Chaos Events Triggered**\nTurbos: **%d**\nRank Swaps: **%d**\nSlipsplosions: **%d**\nOily Floors: **%d**\nMarket Crashes: **%d**\nBig Rips: **%d**\n\n**🤸 Flip Stats**\nTotal Flips: **%d**" $ns $gs $cs $ms $ts $rss $ss $os $mcs $hs $f}}

{{$embed := cembed 
    "title" (printf "📊 Lifetime Stats: %s" $userName)
    "description" $desc
    "color" 16766720 
    "thumbnail" (sdict "url" ($target.AvatarURL "256"))
}}

{{sendMessage nil $embed}}