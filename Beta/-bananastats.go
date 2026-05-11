{{$target := .User}}
{{if .Message.Mentions}}
    {{$target = index .Message.Mentions 0}}
{{end}}

{{$keyData := "TEST_banana_data"}}
{{$userData := sdict}}
{{$hasData := false}}

{{with (dbGet $target.ID $keyData)}}
    {{$userData = dict .Value | sdict}}
    {{$hasData = true}}
{{end}}

{{if not $hasData}}
    ❌ **{{$target.Username}}** has not interacted with any bananas yet! 
    {{return}}
{{end}}

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

{{$desc := printf "**🍌 Slip History**\nNormal Slips: **%d**\nGold Slips: **%d**\nCosmic Slips: **%d**\nMythic Slips: **%d**\n\n**💥 Chaos Events Triggered**\nTurbos: **%d**\nRank Swaps: **%d**\nSlipsplosions: **%d**\nOily Floors: **%d**\nMarket Crashes: **%d**\nHalvings: **%d**\n\n**🤸 Flip Stats**\nTotal Flips: **%d**\nCurrent Streak: **%d**\nRecord Streak: **%d**" $ns $gs $cs $ms $ts $rss $ss $os $mcs $hs $f $c $r}}

{{$embed := cembed 
    "title" (printf "📊 Banana Stats: %s" $target.Username)
    "description" $desc
    "color" 16766720 
    "thumbnail" (sdict "url" ($target.AvatarURL "256"))
}}

{{sendMessage nil $embed}}