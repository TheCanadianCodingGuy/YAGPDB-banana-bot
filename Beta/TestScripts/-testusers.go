{{/* CONFIGURATION */}}
{{ $isDeleteMode := true }} {{/* Set to true to remove data, false to insert */}}

{{/* Dynamic Batch Size */}}
{{ $batchSize := 3 }}
{{ if $isDeleteMode }}{{ $batchSize = 2 }}{{ end }}

{{ $progressKey := "setup_progress_index" }}

{{/* DATASET: Matching your exact nested JSON structure */}}
{{ $data := cslice
    (sdict "UID" "180128558776057856" "TEST_banana_slips" "14" "TEST_banana_data" (sdict "c" 22 "r" 41 "turbo" false))
    (sdict "UID" "600533807132835851" "TEST_banana_slips" "48" "TEST_banana_data" (sdict "c" 5 "r" 37 "turbo" false))
    (sdict "UID" "149357474695217152" "TEST_banana_slips" "31" "TEST_banana_data" (sdict "c" 12 "r" 12 "turbo" false))
    (sdict "UID" "477305422026637312" "TEST_banana_slips" "7" "TEST_banana_data" (sdict "c" 19 "r" 44 "turbo" false))
    (sdict "UID" "1104816192088191048" "TEST_banana_slips" "22" "TEST_banana_data" (sdict "c" 8 "r" 29 "turbo" false))
    (sdict "UID" "1313088371484135496" "TEST_banana_slips" "41" "TEST_banana_data" (sdict "c" 33 "r" 48 "turbo" false))
    (sdict "UID" "894495967896817716" "TEST_banana_slips" "15" "TEST_banana_data" (sdict "c" 1 "r" 15 "turbo" false))
    (sdict "UID" "652960517006295040" "TEST_banana_slips" "3" "TEST_banana_data" (sdict "c" 27 "r" 30 "turbo" false))
    (sdict "UID" "452410789077450754" "TEST_banana_slips" "50" "TEST_banana_data" (sdict "c" 42 "r" 42 "turbo" false))
    (sdict "UID" "785205575695728651" "TEST_banana_slips" "26" "TEST_banana_data" (sdict "c" 10 "r" 36 "turbo" false))
    (sdict "UID" "1443025349905612842" "TEST_banana_slips" "18" "TEST_banana_data" (sdict "c" 25 "r" 49 "turbo" false))
    (sdict "UID" "706983040047513651" "TEST_banana_slips" "9" "TEST_banana_data" (sdict "c" 4 "r" 17 "turbo" false))
    (sdict "UID" "121862600391786498" "TEST_banana_slips" "37" "TEST_banana_data" (sdict "c" 38 "r" 45 "turbo" false))
    (sdict "UID" "138817710212513792" "TEST_banana_slips" "12" "TEST_banana_data" (sdict "c" 21 "r" 21 "turbo" false))
    (sdict "UID" "1470571194171265117" "TEST_banana_slips" "44" "TEST_banana_data" (sdict "c" 16 "r" 34 "turbo" false))
    (sdict "UID" "234760745974366220" "TEST_banana_slips" "29" "TEST_banana_data" (sdict "c" 2 "r" 11 "turbo" false))
    (sdict "UID" "241349892985978880" "TEST_banana_slips" "5" "TEST_banana_data" (sdict "c" 44 "r" 50 "turbo" false))
    (sdict "UID" "244774816987611136" "TEST_banana_slips" "21" "TEST_banana_data" (sdict "c" 30 "r" 31 "turbo" false))
    (sdict "UID" "223466384854614027" "TEST_banana_slips" "33" "TEST_banana_data" (sdict "c" 18 "r" 25 "turbo" false))
    (sdict "UID" "750804346593083462" "TEST_banana_slips" "10" "TEST_banana_data" (sdict "c" 6 "r" 40 "turbo" false))
    (sdict "UID" "190461465822625792" "TEST_banana_slips" "47" "TEST_banana_data" (sdict "c" 15 "r" 19 "turbo" false))
    (sdict "UID" "249936900658298880" "TEST_banana_slips" "2" "TEST_banana_data" (sdict "c" 39 "r" 47 "turbo" false))
    (sdict "UID" "347406407467008001" "TEST_banana_slips" "38" "TEST_banana_data" (sdict "c" 20 "r" 20 "turbo" false))
    (sdict "UID" "987430833797337159" "TEST_banana_slips" "16" "TEST_banana_data" (sdict "c" 3 "r" 28 "turbo" false))
    (sdict "UID" "358560331427217410" "TEST_banana_slips" "25" "TEST_banana_data" (sdict "c" 45 "r" 46 "turbo" false))
    (sdict "UID" "541779925468708879" "TEST_banana_slips" "1" "TEST_banana_data" (sdict "c" 11 "r" 14 "turbo" false))
    (sdict "UID" "547927726754103306" "TEST_banana_slips" "39" "TEST_banana_data" (sdict "c" 28 "r" 39 "turbo" false))
    (sdict "UID" "1016410717546610768" "TEST_banana_slips" "13" "TEST_banana_data" (sdict "c" 7 "r" 10 "turbo" false))
    (sdict "UID" "461695582822727713" "TEST_banana_slips" "20" "TEST_banana_data" (sdict "c" 50 "r" 50 "turbo" false))
    (sdict "UID" "420410318682849290" "TEST_banana_slips" "42" "TEST_banana_data" (sdict "c" 23 "r" 27 "turbo" false))
}}

{{/* PROGRESS LOGIC */}}
{{ $currentIndex := 0 }}
{{ $currentEntry := dbGet 0 $progressKey }}
{{ if $currentEntry }}
    {{ $currentIndex = toInt $currentEntry.Value }}
{{ end }}

{{/* FINAL CLEANUP STAGE (-1) */}}
{{ if eq $currentIndex -1 }}
    {{ if $isDeleteMode }}
        {{ dbDel 0 "TEST_banana_global" }}
        {{ dbDel 0 "TEST_prestige_global" }}
        {{ dbDel 0 $progressKey }}
        ✅ **[DELETED] Cleanup Complete!** Test Globals and progress key have been removed.
    {{ else }}
        {{ $g := sdict "pity" 0 "oily" false "crash" 0 "season" 3 }}
        {{ $m := sdict "180128558776057856" 3 "600533807132835851" 2 "149357474695217152" 1 }}
        {{ dbSet 0 "TEST_banana_global" $g }}
        {{ dbSet 0 "TEST_prestige_global" $m }}
        {{ dbDel 0 $progressKey }}
        ✅ **[ADDED] Success!** All 30 test users and Test Globals have been inserted.
    {{ end }}
{{ else }}
    {{ $endIndex := add $currentIndex $batchSize }}
    {{ if gt $endIndex (len $data) }} {{ $endIndex = len $data }} {{ end }}

    {{/* EXECUTION LOOP */}}
    {{ range $i, $user := $data }}
        {{ if and (ge $i $currentIndex) (lt $i $endIndex) }}
            {{ $uid := toInt64 $user.UID }}
            {{ if $isDeleteMode }}
                {{ dbDel $uid "TEST_banana_slips" }}
                {{ dbDel $uid "TEST_banana_data" }}
                {{ dbDel $uid "TEST_banana_cooldown" }}
            {{ else }}
                {{ dbSetExpire $uid "TEST_banana_slips" (toInt $user.TEST_banana_slips) 31536000 }}
                {{ dbSetExpire $uid "TEST_banana_data" $user.TEST_banana_data 31536000 }}
            {{ end }}
        {{ end }}
    {{ end }}

    {{/* UPDATE OR TRANSITION */}}
    {{ $modeLabel := "[ADDED]" }}{{ if $isDeleteMode }}{{ $modeLabel = "[DELETED]" }}{{ end }}

    {{ if ge $endIndex (len $data) }}
        {{ if $isDeleteMode }}
            {{ dbSet 0 $progressKey -1 }}
            🔄 **[DELETED] All test users processed.** Next run will wipe global data.
        {{ else }}
            {{ dbSet 0 $progressKey -1 }}
            🔄 **[ADDED] All test users processed.** Next run add global data.
        {{ end }}
    {{ else }}
        {{ dbSet 0 $progressKey $endIndex }}
        🔄 **{{ $modeLabel }} Processed test users {{ add $currentIndex 1 }} to {{ $endIndex }}.**
        Run again for the next batch ({{ add $endIndex 1 }}).
    {{ end }}
{{ end }}