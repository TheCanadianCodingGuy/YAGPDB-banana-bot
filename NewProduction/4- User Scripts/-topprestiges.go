{{/* 1. Setup Data */}}
{{ $pg := "prestige_global" }}
{{ $limit := 10 }}
{{ $prestigeMap := sdict }}{{ with (dbGet 0 $pg) }}{{ $prestigeMap = dict .Value | sdict }}{{ end }}

{{/* 2. Sort Logic */}}
{{ $list := cslice }}
{{ range $id, $score := $prestigeMap }}
    {{ $list = $list.Append (sdict "id" (toInt64 $id) "score" (toInt $score)) }}
{{ end }}
{{ range $i, $e := $list }}{{ range $j, $e2 := $list }}
    {{ if gt (index $list $i).score (index $list $j).score }}
        {{ $temp := index $list $i }}{{ $list.Set $i (index $list $j) }}{{ $list.Set $j $temp }}
    {{ end }}
{{ end }}{{ end }}

### 🏆 Global Prestige Leaderboard (Top 10)
{{ if not $list }}*The prestige registry is currently empty.*
{{ else }}
{{- $prevScore := -1 -}}{{ $rank := 0 }}
{{ range $idx, $entry := $list }}
    {{- if lt $idx $limit -}}
        {{- $score := $entry.score -}}
        {{- if ne $score $prevScore }}{{ $rank = add $idx 1 }}{{ end -}}
        {{- $prevScore = $score -}}
        {{- $user := userArg $entry.id -}}
        {{- $rawName := str $entry.id -}}
        {{- if $user -}}
            {{- $rawName = $user.Username -}}
            {{- if $user.Globalname }}{{ $rawName = $user.Globalname }}{{ end -}}
            {{- $mem := getMember $entry.id }}{{ if $mem }}{{ if $mem.Nick }}{{ $rawName = $mem.Nick }}{{ end }}{{ end -}}
        {{- end -}}
        {{- $name := reReplace `([*_~>|\x60])` $rawName `\$1` -}}
**#{{ $rank }}:** (🏆{{ $score }}) {{ $name }} — `{{ $score }} prestige`
    {{- end }}
{{ end }}

{{- $myID := str .User.ID -}}
{{- $myScore := 0 -}}
{{- $myRank := "Unranked" -}}

{{/* THE FIX: Use 'with' to safely check the map */}}
{{- with (index $prestigeMap $myID) -}}
    {{- $myScore = toInt . -}}
{{- end -}}

{{- if gt $myScore 0 -}}
    {{- $count := 1 -}}{{ range $list }}{{ if gt .score (toInt $myScore) }}{{ $count = add $count 1 }}{{ end }}{{ end -}}
    {{- $myRank = printf "#%d" $count -}}
{{- end -}}

{{- $myDisplayName := .User.Username -}}
{{- if .Member.Nick }}{{ $myDisplayName = .Member.Nick }}{{ else if .User.Globalname }}{{ $myDisplayName = .User.Globalname }}{{ end -}}
{{ $myDisplayName = reReplace `([*_~>|\x60])` $myDisplayName `\$1` }}
{{ if gt $myScore 0 }}(🏆{{ $myScore }}) {{ end }}**{{ $myDisplayName }}** is currently ranked **{{ $myRank }}** with **{{ $myScore }}** total prestige.
{{ end }}
*💪 True prestige is earned in the slip.*