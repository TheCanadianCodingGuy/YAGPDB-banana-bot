{{/* --- DATA PROCESSING --- */}}
{{ $dbKey := "banana_record_backflip" }}
{{ $keyCurrentStreak := "banana_current_backflip" }}
{{ $limit := 10 }}
{{ $fetchAmount := 100 }}
{{ $top := dbTopEntries $dbKey $fetchAmount 0 }}

{{ if not $top }}
    🤸 **Gravity remains unchallenged!** No one has landed a backflip yet.
{{ else }}
    {{- $callerRank := "100+" -}}
    {{- $rankTracker := 0 -}}
    {{- $prevValTracker := -1 -}}
    {{- $callerFoundInTop10 := false -}}
    {{- $displayCount := 0 -}}
    {{- $prevValue := -1 -}}
    {{- $rank := 0 -}}

    🤸 **Banana Backflip Pantheon (Top 10)**
    {{- range $i, $entry := $top -}}
        {{- $val := toInt .Value -}}
        {{- if ne $val $prevValTracker -}}{{- $rankTracker = add $i 1 -}}{{- end -}}
        {{- $prevValTracker = $val -}}
        {{- if eq .User.ID $.User.ID -}}{{- $callerRank = str $rankTracker -}}{{- end -}}

        {{- if lt $displayCount $limit -}}
            {{- $member := getMember .User.ID -}}
            {{- if $member -}}
                {{- $displayCount = add $displayCount 1 -}}
                {{- if ne $val $prevValue -}}{{- $rank = $displayCount -}}{{- end -}}
                {{- $prevValue = $val -}}
                {{- if eq .User.ID $.User.ID -}}{{- $callerFoundInTop10 = true -}}{{- end -}}
                
                {{- $rawName := .User.Username -}}
                {{- if $member.Nick -}}{{- $rawName = $member.Nick -}}
                {{- else if .User.Globalname -}}{{- $rawName = .User.Globalname -}}{{- end -}}
                {{- $name := reReplace `([*_~>|\x60])` $rawName `\$1` }}
            **#{{ $rank }}:** {{ $name }} — `{{ $val }} backflip{{ if ne $val 1 }}s{{ end }}`
            {{- end -}}
        {{- end -}}
    {{- end -}}

    {{/* --- FOOTER --- */}}
    {{- $currentData := dbGet .User.ID $keyCurrentStreak -}}
    {{- $recordData := dbGet .User.ID $dbKey -}}
    {{- $currentStreak := 0 -}}{{- if $currentData -}}{{- $currentStreak = toInt $currentData.Value -}}{{- end -}}
    {{- $personalRecord := 0 -}}{{- if $recordData -}}{{- $personalRecord = toInt $recordData.Value -}}{{- end -}}
    {{- $userName := .User.Username -}}{{- if .Member.Nick -}}{{- $userName = .Member.Nick -}}{{- else if .User.Globalname -}}{{- $userName = .User.Globalname -}}{{- end -}}
    {{- $footerRank := "" -}}{{- if not $callerFoundInTop10 -}}{{- $footerRank = printf "😎 Your current rank is **#%s**!" $callerRank -}}{{- end }}
        {{- $streakStatus := "chasing your" -}}
        {{- if gt $currentStreak $personalRecord -}}
            {{- $streakStatus = "setting a new" -}}
            {{- $personalRecord = $currentStreak -}}
        {{- else if eq $currentStreak $personalRecord -}}
            {{- $streakStatus = "pulling even with your" -}}
        {{ end }}
    {{ if gt $currentStreak 0 }}
You are evading gravity with an active streak of **{{ $currentStreak }}**,
		{{- if gt $personalRecord 0 }}
{{ $streakStatus }} all-time personal record of **{{ $personalRecord }}**! 🏆
		{{- else }}
and you have not made a single backflip. Bummer!
		{{- end -}}
    {{ else }}
You have no active streak, but your record stands at **{{ $personalRecord }}** backflips. Time to flip! 🤸
    {{- end }}
{{ $footerRank }}
{{ end }}