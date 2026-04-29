{{/* Configuration */}}
{{ $dbKey := "TEST_banana_slips" }}
{{ $prestigeGlobal := "TEST_prestige_global" }}
{{ $limit := 10 }}
{{ $fetchAmount := 25 }}

{{ $top := dbTopEntries $dbKey $fetchAmount 0 }}

{{ if not $top }}
    [TEST] 🍌 **The floors are clean!** No one has slipped on a banana peel yet.
{{ else }}
    ### [TEST] 🍌 Banana Peel Hall of Shame (Top 10)
{{- /* 1. Fetch the Global Prestige Map */ -}}
{{- $prestigeMap := sdict -}}
{{- with (dbGet 0 $prestigeGlobal) }}{{ $prestigeMap = dict .Value | sdict }}{{ end -}}

{{- $displayCount := 0 -}}
{{- $prevValue := -1 -}}
{{- $rank := 0 -}}

{{- range $i, $entry := $top -}}
    {{- if lt $displayCount $limit -}}
        {{- $member := getMember $entry.User.ID -}}
        {{- if $member -}}
            {{- $val := toInt $entry.Value -}}
            {{- $displayCount = add $displayCount 1 -}}
            {{- if ne $val $prevValue -}}
                {{- $rank = $displayCount -}}
            {{- end -}}
            {{- $prevValue = $val -}}
            
            {{- /* 2. Sanitize Name */ -}}
            {{- $rawName := $entry.User.Username -}}
            {{- if $member.Nick }}{{ $rawName = $member.Nick -}}
            {{- else if $entry.User.Globalname }}{{ $rawName = $entry.User.Globalname }}{{ end -}}
            {{- $name := reReplace `([*_~>|\x60])` $rawName `\$1` -}}

            {{- /* 3. Lookup Prestige */ -}}
            {{- $userPrestige := index $prestigeMap (str $entry.User.ID) | or 0 -}}
            {{- if gt (toInt $userPrestige) 0 -}}
                {{- $name = printf "(🏆%d) %s" (toInt $userPrestige) $name -}}
            {{- end }}
**#{{ $rank }}:** {{ $name }} — `{{ $val }} slip{{ if ne $val 1 }}s{{ end }}`
        {{- end -}}
    {{- end -}}
{{- end }}

*⚠️ Floor slippery when banana.*

{{- /* 4. Calling User's Personal Stats */ -}}
{{- $myRank := "Unranked" -}}
{{- $mySlips := 0 -}}
{{- with (dbGet .User.ID $dbKey) }}{{ $mySlips = toInt .Value }}{{ end -}}

{{- if gt $mySlips 0 -}}
    {{- $rankCounter := 1 -}}
    {{- $allEntries := dbTopEntries $dbKey 100 0 -}}
    {{- $prevV := -1 -}}
    {{- range $idx, $ent := $allEntries -}}
        {{- $v := toInt .Value -}}
        {{- if ne $v $prevV }}{{ $rankCounter = add $idx 1 }}{{ end -}}
        {{- $prevV = $v -}}
        {{- if eq .User.ID $.User.ID }}{{ $myRank = str $rankCounter }}{{ end -}}
    {{- end -}}
{{- end -}}

{{- /* 5. Final Personalized Footer */ -}}
{{- $myPrestige := index $prestigeMap (str .User.ID) | or 0 -}}
{{- $myDisplayName := .User.Username -}}
{{- if .Member.Nick }}{{ $myDisplayName = .Member.Nick -}}
{{- else if .User.Globalname }}{{ $myDisplayName = .User.Globalname }}{{ end -}}
{{- $myDisplayName = reReplace `([*_~>|\x60])` $myDisplayName `\$1` -}}

{{- if gt (toInt $myPrestige) 0 -}}
    {{- $myDisplayName = printf "(🏆%d) %s" (toInt $myPrestige) $myDisplayName -}}
{{- end }}

**{{ $myDisplayName }}** is currently at rank **#{{ $myRank }}** with **{{ $mySlips }}** slip{{ if ne $mySlips 1 }}s{{ end }}.
{{ end }}