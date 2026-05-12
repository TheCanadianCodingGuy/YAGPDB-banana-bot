{{/* 1. Configuration */}}
{{ $allowedCommands := cslice "-banana" "-bananastats" "-commandlist" "-topslips" "-topprestiges" "-what" }}

{{/* FIXED: IDs must be strings to avoid 'overflows int' error */}}
{{ $modRoleID := "1196392951333523476" }}
{{ $yagID := "204255221017214977" }}

{{/* 2. Check Permissions */}}
{{ $isMod := hasRoleID (toInt64 $modRoleID) }}
{{ $isYAG := eq .User.ID (toInt64 $yagID) }}

{{/* 3. Logic Execution */}}
{{ if or $isMod $isYAG }}
    {{/* Authorized user/bot: Do nothing */}}
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