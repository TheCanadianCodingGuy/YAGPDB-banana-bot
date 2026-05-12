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

{{ $dbKey := "TEST_banana_slips" }}
{{ $pg := "TEST_prestige_global" }}
{{ $kg := "TEST_banana_global" }}
{{ $limit := 10 }}
{{ $fetchAmount := 25 }}

{{ $global := sdict "season" 1 }}{{ with (dbGet 0 $kg) }}{{ $global = dict .Value | sdict }}{{ end }}
{{ $prestigeMap := sdict }}{{ with (dbGet 0 $pg) }}{{ $prestigeMap = dict .Value | sdict }}{{ end }}

{{ $startOfMonth := (newDate currentTime.Year (toInt currentTime.Month) 1 0 0 0) }}

{{ $nextMonth := (add $startOfMonth.Month 1) }}
{{ $nextYear := $startOfMonth.Year }}
{{ if gt $nextMonth 12 }}{{ $nextMonth = 1 }}{{ $nextYear = add $nextYear 1 }}{{ end }}
{{ $endOfSeason := (newDate $nextYear (toInt $nextMonth) 1 0 0 0).Add (mult -6 3600 | toDuration) }}

{{ $top := dbTopEntries $dbKey $fetchAmount 0 }}

{{ if not $top }}
    [TEST] 🍌 **The floors are clean!** No one has slipped on a banana peel yet.
{{ else }}
    ### 🏆 Banana Season {{ $global.season }} Hall of Shame (Top 10)
    **Started on:** <t:{{ $startOfMonth.Unix }}:F>
    **Ends on:** <t:{{ $endOfSeason.Unix }}:F> (<t:{{ $endOfSeason.Unix }}:R>)
{{- "\n\u200b" -}}
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
            
            {{- $rawName := $entry.User.Username -}}
            {{- if $member.Nick }}{{ $rawName = $member.Nick -}}
            {{- else if $entry.User.Globalname }}{{ $rawName = $entry.User.Globalname }}{{ end -}}
            {{- $name := reReplace `([*_~>|\x60])` $rawName `\$1` -}}

            {{- $userPrestige := index $prestigeMap (str $entry.User.ID) | or 0 -}}
            {{- if gt (toInt $userPrestige) 0 -}}
                {{- $name = printf "(🏆%d) %s" (toInt $userPrestige) $name -}}
            {{- end }}
**#{{ $rank }}:** {{ $name }} — `{{ $val }} slip{{ if ne $val 1 }}s{{ end }}`
        {{- end -}}
    {{- end -}}
{{- end }}

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
*⚠️ Floor slippery when banana.*