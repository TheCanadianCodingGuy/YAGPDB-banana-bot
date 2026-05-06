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
        {{/* Use 'or' to handle nil/zero values */}}
        {{ $displayPoints := or $points 0 }}
        - **User ID:** `{{ $id }}` | **Prestige:** `{{ $displayPoints }}`
    {{ end }}
{{ end }}

*Total entries in map: {{ len $prestigeMap }}*