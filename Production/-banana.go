{{/* --- CONFIGURATION --- */}}
{{ $keySlips := "banana_slips" }}
{{ $keyCD := "banana_cooldown" }}
{{ $keyCurrentStreak := "banana_current_backflip" }}
{{ $keyRecordStreak := "banana_record_backflip" }}
{{ $cooldownDuration := 36000 }} {{/* 10 Hours */}}
{{ $expiration := 31536000 }} {{/* 12 Months */}}

{{/* --- USER IDENTITY --- */}}
{{ $rawName := .User.Username }}
{{ if .Member.Nick }}{{ $rawName = .Member.Nick }}
{{ else if .User.Globalname }}{{ $rawName = .User.Globalname }}{{ end }}
{{ $userName := reReplace `([*_~>|\x60])` $rawName `\$1` }}

{{/* --- COOLDOWN CHECK --- */}}
{{ $cooldownData := dbGet .User.ID $keyCD }}

{{ if $cooldownData }}
    {{ $userName }}, you look around but see no banana peel to slip on! How sad!
    **Try again <t:{{ toInt $cooldownData.Value }}:R>!**
{{ else }}
    {{/* Pre-fetch current user data from the database */}}
    {{ $oldSlips := 0 }}{{ with (dbGet .User.ID $keySlips) }}{{ $oldSlips = toInt .Value }}{{ end }}
    {{ $oldRecord := 0 }}{{ with (dbGet .User.ID $keyRecordStreak) }}{{ $oldRecord = toInt .Value }}{{ end }}
    {{ $oldStreak := 0 }}{{ with (dbGet .User.ID $keyCurrentStreak) }}{{ $oldStreak = toInt .Value }}{{ end }}

    {{/* INITIAL RANK CHECK: Determine where they stand before the roll */}}
    {{ $currentRank := 101 }}{{ $prevVal := -1 }}{{ $rankTracker := 0 }}{{ $found := false }}
    {{ range $i, $entry := dbTopEntries $keySlips 100 0 }} {{/* Pull top 100 to find user's position */}}
        {{ if not $found }}
            {{ $val := toInt .Value }}
            {{ if ne $val $prevVal }}{{ $rankTracker = add $i 1 }}{{ end }} {{/* Handle tie logic */}}
            {{ $prevVal = $val }}
            {{ if eq .User.ID $.User.ID }}{{ $currentRank = $rankTracker }}{{ $found = true }}{{ end }}
        {{ end }}
    {{ end }}

    {{/* Set the 10-hour cooldown immediately */}}
    {{ $expiresAt := add currentTime.Unix $cooldownDuration }}
    {{ dbSetExpire .User.ID $keyCD (str $expiresAt) $cooldownDuration }}

    {{/* --- THE FLIP (50% Chance to Slip or Backflip) --- */}}
    {{ if eq (randInt 2) 1 }}
        {{/* --- SLIP BRANCH --- */}}
        {{/* Golden chance & Catch-up scaling */}}
        {{ $luckyChance := 45 }}
        {{ if le $currentRank 10 }}{{ $luckyChance = 5 }}
        {{ else if le $currentRank 20 }}{{ $luckyChance = 25 }}
        {{ else if le $currentRank 30 }}{{ $luckyChance = 35 }}
        {{ end }}

        {{ $multiplier := 1 }}
        {{ $isGolden := false }}
        {{ $isCosmic := false }}
        {{ $isMythic := false }}

        {{/* Evaluates Golden chance first */}}
        {{ if le (randInt 1 101) $luckyChance }}
            {{ $multiplier = 2 }}
            {{ $isGolden = true }}
            
            {{/* One roll to decide if it upgrades to Mythic, Cosmic, or stays Gold */}}
            {{ $upgradeRoll := randInt 1 101 }}
            
            {{ if le $upgradeRoll 5 }} 
                {{/* 1% chance for Mythic */}}
                {{ $multiplier = 10 }}
                {{ $isMythic = true }}
            {{ else if le $upgradeRoll 20 }} 
                {{/* Values 2-11 = 10% chance for Cosmic */}}
                {{ $multiplier = 5 }}
                {{ $isCosmic = true }}
            {{ end }}
        {{ end }}
        
        {{ $addedSlips := $multiplier }}
        {{ $newSlips := add $oldSlips $addedSlips }}
        
        {{/* Update all database keys and refresh expiration */}}
        {{ dbSetExpire .User.ID $keySlips (str $newSlips) $expiration }}
        {{ dbSetExpire .User.ID $keyCurrentStreak "0" $expiration }}
        {{ dbSetExpire .User.ID $keyRecordStreak (str $oldRecord) $expiration }}

        {{/* POST-SLIP RANK CHECK: See where they landed after the chaos */}}
        {{ $finalRank := "100+" }}{{ $prevVal = -1 }}{{ $rankTracker = 0 }}{{ $found = false }}
        {{ range $i, $entry := dbTopEntries $keySlips 100 0 }}
            {{ if not $found }}
                {{ $val := toInt .Value }}
                {{ if ne $val $prevVal }}{{ $rankTracker = add $i 1 }}{{ end }}
                {{ $prevVal = $val }}
                {{ if eq .User.ID $.User.ID }}{{ $finalRank = str $rankTracker }}{{ $found = true }}{{ end }}
            {{ end }}
        {{ end }}

        {{ $rankInt := toInt $finalRank }}
        {{ $hype := "The floor doesn't know your name yet... keep it that way." }}
        {{ if eq $rankInt 1 }}{{ $hype = "🏆 The floor has officially given up all hope when it sees you coming!" }}
        {{ else if le $rankInt 5 }}{{ $hype = "🍌 The floor is starting to fear your very existence." }}
        {{ else if le $rankInt 10 }}{{ $hype = "⚠️ The floor is permanently dented with your butt print." }}
        {{ else if le $rankInt 20 }}{{ $hype = "👀 The floor has started wearing a helmet in anticipation of your arrival." }}
        {{ end }}

        {{/* --- OUTPUT MESSAGES --- */}}
        {{ $header := "**Oh no!** 🍌"}}
        {{ $body := "" }}
        {{ if $isMythic }}
            {{ $header = "**HOLY CRAP!** 😱😱😱" }}
            {{ $body = print "🎆🎆🎆 ***" $userName " JUST SLIPPED ON A DANG MYTHIC PEEL WORTH 10 SLIPS!*** 🎆🎆🎆" }}
        {{- else if $isCosmic }}
            {{ $body = print "" $userName " just slipped on a rare ***cosmic peel*** worth ***5 slips!*** 🌟🌌🌟" }}
        {{- else if $isGolden }}
            {{ $body = print "" $userName " just slipped on a peel made of **pure gold** worth **2 slips!** 😮" }}
        {{- else }}
            {{ $body = print "" $userName " just slipped on a peel!" }}
        {{- end }}
        {{ $header }}
        {{ $body }}
        That is **{{ $newSlips }}** total **slip{{ if ne $newSlips 1 }}s{{ end }}** for you. Watch your step!

        {{ $hype }} (Your clumsiness places you at rank **#{{ $finalRank }}**!)
    {{ else }}
        {{/* --- BACKFLIP BRANCH --- */}}
        {{ $newStreak := add $oldStreak 1 }} 
        {{ $newRecord := $oldRecord }}
        {{ $streakStatus := "chasing your" }}

        {{ if gt $newStreak $oldRecord }}
            {{ $streakStatus = "setting a new" }}
            {{ $newRecord = $newStreak }}
        {{ else if eq $newStreak $oldRecord }}
            {{ $streakStatus = "pulling even with your" }}
        {{ end }}

        {{/* Update all keys and reset expiries (3 calls) */}}
        {{ dbSetExpire .User.ID $keyCurrentStreak (str $newStreak) $expiration }}
        {{ dbSetExpire .User.ID $keyRecordStreak (str $newRecord) $expiration }}
        {{ dbSetExpire .User.ID $keySlips (str $oldSlips) $expiration }}

        **CLEAN!** 🤸
        {{ $userName }} dodged the peel with a backflip! 
        You’ve now evaded gravity **{{ $newStreak }}** time{{ if ne $newStreak 1 }}s in a row{{ end }},
        {{ $streakStatus }} all-time personal record of **{{ $newRecord }}**.
    {{ end }}
{{ end }}