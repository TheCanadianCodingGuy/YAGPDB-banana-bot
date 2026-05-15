{{/* --- CRON: */30 * * * * --- */}}
{{ $batchSize := 8 }} {{/* Number of users to check per run */}}
{{ $ptrKey := "TEST_cleanup_pointer" }}
{{ $lbKey := "TEST_banana_slips" }}

{{/* 1. Retrieve the current offset pointer */}}
{{ $offset := 0 }}
{{ with (dbGet 0 $ptrKey) }}
    {{ $offset = (toInt .Value) }}
{{ end }}

{{/* 2. Fetch the next batch of entries using the offset */}}
{{ $entries := dbTopEntries $lbKey $batchSize $offset }}

{{/* 3. If we hit the end of the list, reset the pointer and exit */}}
{{ if not $entries }}
    {{ dbSet 0 $ptrKey 0 }}
    {{ return }}
{{ end }}

{{/* 4. Process the batch */}}
{{ range $entries }}
    {{ if not (getMember .User.ID) }}
        {{/* Target identified as having left the server: Wipe all related keys */}}
        {{ dbDel .User.ID "TEST_banana_slips" }}
        {{ dbDel .User.ID "TEST_banana_data" }}
        {{ printf "Successfully purged inactive user: %d" .User.ID }}
    {{ end }}
{{ end }}

{{/* 5. Update the pointer for the next execution */}}
{{ $newOffset := add $offset $batchSize }}
{{ dbSet 0 $ptrKey $newOffset }}