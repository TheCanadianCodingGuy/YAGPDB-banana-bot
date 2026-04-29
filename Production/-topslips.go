{{/* Configuration */}}
{{ $dbKey := "banana_slips" }}
{{ $limit := 10 }}
{{ $fetchAmount := 25 }}

{{ $top := dbTopEntries $dbKey $fetchAmount 0 }}

{{ if not $top }}
    🍌 **The floors are clean!** No one has slipped on a banana peel yet.
{{ else }}
    ### 🍌 Banana Peel Hall of Shame (Top 10)
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
            {{- if $member.Nick -}}{{- $rawName = $member.Nick -}}
            {{- else if $entry.User.Globalname -}}{{- $rawName = $entry.User.Globalname -}}{{- end -}}
            {{- $name := reReplace `([*_~>|\x60])` $rawName `\$1` }}
			**#{{ $rank }}:** {{ $name }} — `{{ $val }} slip{{ if ne $val 1 }}s{{ end }}`
        {{- end -}}
    {{- end -}}
{{- end }}

*⚠️ Floor slippery when banana.*
{{ end }}