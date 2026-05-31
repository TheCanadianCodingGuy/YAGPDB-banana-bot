{{ $allowedCommands := cslice "-banana" "-lifetimestats" "-seasonstats" "-commandlist" "-topslips" "-topflips" "-topprestiges" "-what" }}

{{ $modRoleID := "1196392951333523476" }}
{{ $yagID := "204255221017214977" }}

{{ $isMod := hasRoleID (toInt64 $modRoleID) }}
{{ $isYAG := eq .User.ID (toInt64 $yagID) }}

{{ if or $isMod $isYAG }}
    {{/* Authorized */}}
{{ else }}
    {{ $isCommand := false }}
    {{ $input := lower .Message.Content }}
    
    {{ range $allowedCommands }}
        {{ if (hasPrefix $input .) }}
            {{ $isCommand = true }}
        {{ end }}
    {{ end }}

    {{ if not $isCommand }}
        {{ deleteMessage nil .Message.ID 0 }}
    {{ end }}
{{ end }}