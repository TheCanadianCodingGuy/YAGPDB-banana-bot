{{ $embed := sdict
    "title" "🍌 Banana Bot: Official Command Registry"
    "color" 16773120
    "description" (printf "These are the currently available commands for interacting with the Banana Bot.")
    "fields" (cslice
        (sdict "name" "🕹️ Core Gameplay" "value" "• `-banana`: The main event. A 50/50 toss-up to either **Slip** (gain slips) or **Flip** (build a streak). Triggering a **Lucky Slip** can activate chaos events or massive multipliers." "inline" false)
        
        (sdict "name" "📊 Player Data" "value" "• `-bananastats [@user]`: View detailed analytics including total slips, longest flip streaks, and a breakdown of every special peel type encountered.\n• `-what`: A comprehensive deep-dive into the game's mechanics, event odds, and the prestige system." "inline" false)

        (sdict "name" "🏆 Competitive Boards" "value" "• `-topslips`: View the current season's Top 10 slippers. This board resets monthly.\n• `-topprestiges`: The hall of eternal glory. Displays the all-time leaders in Prestige Points earned from past season victories." "inline" false)
        
        (sdict "name" "ℹ️ Utility" "value" "• `-commandlist`: Displays this directory of available interactions." "inline" false)
    )
    "footer" (sdict "text" "Developed by Ikaras")
}}

{{ sendMessage nil (complexMessage "embed" $embed) }}