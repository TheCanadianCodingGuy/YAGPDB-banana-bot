{{ $cooldownHr := 5 }}
{{ $marketCrashRolls := 3 }}
{{ $oddsMarketCrash := 7 }}
{{ $oddsBigRip := 1}}
{{ $oddsRankSwap := 4 }}
{{ $oddsTurbo := 12 }}
{{ $oddsPlosion := 14 }}
{{ $oddsOily := 14 }}
{{ $oddsMythic := 6 }}
{{ $oddsCosmic := 11 }}
{{ $oddsGolden := 28 }} {{/* Remainder of 100 - others */}}

{{ $embed := sdict
    "title" "🍌 Welcome To The Official Guide on How to Banana"
    "color" 16773120
    "description" "The goal is simple: **Slip often, slip hard.** Collect slips to climb the leaderboard before the season resets."
    "fields" (cslice
        (sdict "name" "🎲 The Initial Toss-Up (50/50)" "value" "Every time you play, it's a **50/50 chance**:\n• **Slip:** You fall and gain slips (points).\n• **Flip:** You dodge the peel and gain a streak, but no points." "inline" false)
              
        (sdict "name" "🕒 Cooldown & Rules" "value" (printf "• **Cooldown:** Every action sets a **%d-hour** cooldown.\n• **The Season:** Ends monthly. The final 6 hours are reserved to crunch numbers and wax the floor for the next season (No slipping allowed!)." $cooldownHr) "inline" false)
        
        (sdict "name" "✨ The Lucky Roll" "value" "If you **Slip**, the bot rolls a second time for a **Lucky Slip**. Your odds of being lucky increase if you are ranked lower (Pity System). **All special events below ONLY trigger on a Lucky Slip.**" "inline" false)
        
        (sdict "name" "🎢 Lucky Peel Multipliers" "value" (printf "Hit a Lucky Slip to trigger these bonuses:\n• **Golden Peel (%d%%):** Gained 2 slips!\n• **Cosmic Peel (%d%%):** Gained 5 slips!\n• **Mythic Peel (%d%%):** Gained 10 slips!" $oddsGolden $oddsCosmic $oddsMythic) "inline" false)

        (sdict "name" "⚠️ Chaos Events (Triggered on Lucky Slips)" "value" (printf "• **Oily Floor (%d%%):** Dooms the next player to a guaranteed slip.\n• **Slipsplosion (%d%%):** Destroy 3 slips from the person ranked above you!\n• **Turbo Overdrive (%d%%):** Resets cooldown and makes your next slip x2.\n• **Rank Swap (%d%%):** Swap your total slips with someone nearby.\n• **Market Crash (%d%%):** Slips grant 0 for %d rolls, but flipping steals 5 slips from the market!\n• **The Big Rip (%d%%):** Your total slips are **HALVED** immediately." $oddsOily $oddsPlosion $oddsTurbo $oddsRankSwap $oddsMarketCrash $marketCrashRolls $oddsBigRip) "inline" false)

        (sdict "name" "🏆 Prestige" "value" "Winners of past seasons carry the **(🏆)** badge. Check the Global Prestige Leaderboard!" "inline" false)
    )
    "footer" (sdict "text" "See pinned messages for more information! For any issues, contact Ikaras.")
}}

{{ sendMessage nil (complexMessage "embed" $embed) }}