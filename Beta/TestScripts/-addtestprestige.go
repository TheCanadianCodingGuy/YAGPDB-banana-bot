{{/* Configuration */}}
{{ $kg := "TEST_prestige_global" }}
{{ $ex := 31536000 }}

{{/* 1. Define the Batch Data (UserID: Increment) */}}
{{ $updates := sdict 
    "180128558776057856" 3
    "600533807132835851" 2
    "149357474695217152" 1
}}

{{/* 2. Fetch the existing Global Map */}}
{{ $prestigeMap := sdict }}
{{ with (dbGet 0 $kg) }}
    {{ $prestigeMap = dict .Value | sdict }}
{{ end }}

{{/* 3. Apply the updates */}}
{{ range $id, $inc := $updates }}
    {{ $current := index $prestigeMap $id | or 0 }}
    {{ $prestigeMap.Set $id (add $current $inc) }}
{{ end }}

{{/* 4. Save the updated map */}}
{{ dbSet 0 $kg $prestigeMap }}

✅ **Batch update complete!**
{{ range $id, $inc := $updates }}
- `<@{{ $id }}>`: +{{ $inc }} prestige (Total: {{ index $prestigeMap $id }})
{{ end }}