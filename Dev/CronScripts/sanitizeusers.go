{{/* --- CRON: */30 * * * * --- */}}
{{ $batchSize := 8 }}
{{ $ptrKey := "cleanup_pointer" }}
{{ $lbKey := "banana_slips" }}

{{ $offset := 0 }}
{{ with (dbGet 0 $ptrKey) }}
    {{ $offset = (toInt .Value) }}
{{ end }}

{{ $entries := dbTopEntries $lbKey $batchSize $offset }}

{{ if not $entries }}
    {{ dbSet 0 $ptrKey 0 }}
    {{ return }}
{{ end }}

{{/* 4. Process the batch */}}
{{ range $entries }}
    {{ if not (getMember .User.ID) }}
        {{ dbDel .User.ID "banana_slips" }}
        {{ dbDel .User.ID "banana_data" }}
    {{ end }}
{{ end }}

{{ $newOffset := add $offset $batchSize }}
{{ dbSet 0 $ptrKey $newOffset }}