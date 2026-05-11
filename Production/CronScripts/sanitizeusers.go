{{/* Daily Cleanup: Slips & Flips (PROD) */}}

{{/* 1. Clean up based on Slips Leaderboard */}}
{{ range (dbTopEntries "banana_slips" 100 0) }}
    {{ if not (getMember .User.ID) }}
        {{ dbDel .User.ID "banana_slips" }}
        {{ dbDel .User.ID "banana_cooldown" }}
        {{ dbDel .User.ID "banana_current_backflip" }}
        {{ dbDel .User.ID "banana_record_backflip" }}
    {{ end }}
{{ end }}

{{/* 2. Clean up based on Record Backflip Leaderboard */}}
{{ range (dbTopEntries "banana_record_backflip" 100 0) }}
    {{ if not (getMember .User.ID) }}
        {{ dbDel .User.ID "banana_slips" }}
        {{ dbDel .User.ID "banana_cooldown" }}
        {{ dbDel .User.ID "banana_current_backflip" }}
        {{ dbDel .User.ID "banana_record_backflip" }}
    {{ end }}
{{ end }}