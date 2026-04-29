# 🍌 Banana Slips: Chaotic Mechanic Expansion Pack

This document contains the complete logic and descriptions for 7 chaotic "Glitch" mechanics for the YAGPDB Banana Slips script. 

**Note:** In this game, Slips are the desirable benefit.

---

### 1. The "Rank Swap" Rare Event
**Description:** A "Glitch" for players in the Top 20. Upon slipping, they swap their entire slip count with the person ranked exactly 1 to 5 spots below them.
**Target:** Top 20 Players | **Odds:** 3% | **Trigger:** On lucky slip only

---

### 2. The Poor Interest Rate (Pity System)
**Description:** Every turn that results in a "standard" slip from anyone (no bonus, rolled from anyone at any rank) increases the chance of the next slip from anyone at rank 11+ being a lucky roll bonus by X%, stacking until a bonus slip happens. This bonus resets resets upon hitting any bonus roll from any rank 11+ users.
**Target:** Rank 11+ | **Odds:** Stacking +2% | **Trigger:** On non-lucky slip only

---

### 3. The Halving Glitch
**Description:** A catastrophic error where the user's slip count is instantly halved.
**Target:** Anyone | **Odds:** 2% | **Trigger:** On lucky slip only

---

### 4. The "Oily Floor" Roll
**Description:** A chance that a slip "Greases the Floor." The very next person to use the command is guaranteed to slip, regardless of the 50/50 backflip odds.
**Target:** Anyone | **Odds:** 9% | **Trigger:** On lucky slip only

---

### 5. The "Market Crash" Glitch
**Description:** A ultra-rare event that breaks logic for the next 3 rolls from any users. Successful backflips grant +5 slips; slipping grants 0.
**Target:** Global | **Odds:** 2% Trigger | **Lasting:** 3 backflips | **Trigger:** On lucky slip only

---

### 6. The "Double Backflip" Glitch (Turbo)
**Description:** A chance that a slip resets the user's cooldown to 0 and doubles the slip value of their very next turn, including bonus slips, should they slip again.
**Target:** Anyone | **Odds:** 4% | **Trigger:** On lucky slip only

---

### 7. The Slip-plosion
**Description:** A rare chance on slip, that they cause a shockwave that reduces the total slip count of the person ranked directly above them by 1. If the roller is first, then nullifies their own slip.
**Target:** Player at [Current Rank - 1] | **Odds:** 5% | **Trigger:** On lucky slip only