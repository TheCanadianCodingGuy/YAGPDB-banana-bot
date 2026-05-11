{{if hasPrefix .Message.Content "-"}}
    {{$args := split (slice .Message.Content 1) " "}}
    {{$cmd := lower (index $args 0)}}

    {{range $allowed}}
        {{if eq $cmd .}}
            {{$isValid = true}}
        {{end}}
    {{end}}
{{end}}

{{if not $isValid}}
    {{deleteTrigger 0}}
{{end}}