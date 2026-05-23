{{/* Configuration */}}
{{ $kg := "banana_global" }}

{{/* 1. Fetch current or setup defaults */}}
{{ $global := sdict "pity" 0 "oily" false "crash" 0 "season" 1 }}
{{ with (dbGet 0 $kg) }}
    {{ $global = dict .Value | sdict }}
{{ end }}

{{/* 2. Edit Values Here */}}
{{ $global.Set "season" 1 }} {{/* Change this number to increment the season */}}
{{ $global.Set "pity" 0 }}
{{ $global.Set "oily" false }}
{{ $global.Set "crash" 0 }}

{{/* 3. Save to Database */}}
{{ dbSet 0 $kg $global }}