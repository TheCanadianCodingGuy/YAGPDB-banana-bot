{{/* --- CRON: 0 0 * * * --- */}}
{{ range (dbTopEntries "TEST_banana_slips" 100 0) }}
    {{ if not (getMember .User.ID) }}
        {{ dbDel .User.ID "TEST_banana_slips" }}
        {{ dbDel .User.ID "TEST_banana_data" }}
    {{ end }}
{{ end }}