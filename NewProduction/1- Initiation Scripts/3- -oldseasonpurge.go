{{$entries := dbTopEntries "banana_slips" 1 0}}

{{if $entries}}
    {{$target := index $entries 0}}
    {{dbDel $target.UserID "banana_slips"}}
    {{dbDel $target.UserID "banana_record_backflip"}}
    {{dbDel $target.UserID "banana_current_backflip"}}
    {{dbDel $target.UserID "banana_cooldown"}}
{{end}}