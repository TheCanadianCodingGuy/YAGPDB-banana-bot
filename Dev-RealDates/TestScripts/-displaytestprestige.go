{{ $kg := "TEST_prestige_global" }}
{{ $prestigeMap := sdict }}
{{ with (dbGet 0 $kg) }}
    {{ $prestigeMap = dict .Value | sdict }}
{{ end }}

### 🏆 Global Prestige Registry
{{ if not $prestigeMap }}
    Registry is currently empty.
{{ else }}
    {{ range $id, $points := $prestigeMap }}
        {{ $displayPoints := or $points 0 }}
        
        {{/* Fetch user info. Note: This can be heavy on large lists */}}
        {{ $user := userArg $id }}
        {{ $name := (printf "Unknown (%s)" $id) }}
        
        {{ if $user }}
            {{ $name = $user.Username }}
            {{ if $user.Globalname }}{{ $name = $user.Globalname }}{{ end }}
        {{ end }}

        - **{{ $name }}**: `{{ $displayPoints }}` Prestige
    {{ end }}
{{ end }}

*Total entries in map: {{ len $prestigeMap }}*