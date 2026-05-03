{{ $uid := 180128558776057856 }}
{{ $slips := 14 }}
{{ $data := sdict "c" 22 "r" 41 "cd" 0"turbo" false }}
{{ dbDel $uid "TEST_banana_slips" }}
{{ dbDel $uid "TEST_banana_data" }}
✅ **[DELETED]** User `{{ $uid }}` data removed.
{{ dbSetExpire $uid "TEST_banana_slips" $slips 31536000 }}
{{ dbSetExpire $uid "TEST_banana_data" $data 31536000 }}
{{ $checkSlips := dbGet $uid "TEST_banana_slips" }}
{{ $checkData := dbGet $uid "TEST_banana_data" }}
{{ $sd := dict $checkData.Value | sdict }}
✅ **[INSERTED] Data for User:** `{{ $uid }}`
*   **Slips:** `{{ $checkSlips.Value }}`
*   **Current Streak (c):** `{{ $sd.c }}`
*   **Record Streak (r):** `{{ $sd.r }}`
*   **Turbo Active:** `{{ $sd.turbo }}`