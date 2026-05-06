$roll := 1 		  // Crash ✅
$roll := 4      // Halving ✅
$roll := 6      // Rank Swap ✅
$roll := 10     // Turbo ✅
$roll := 20     // Slip-splosion ✅
$roll := 30     // Oily
$roll := 43     // Mythic
$roll := 46     // Cosmic
$roll := 100    // Golden


# Banana Slip Test Suite

This is a comprehensive manual test suite for your script, formatted as a Markdown checklist.

Because your script relies heavily on RNG (Random Number Generation), trying to trigger specific events like a **Market Crash** or **Rank Swap** naturally could take hours.

> 💡 **Testing Tip for Glitches:**
> To test the "Lucky" glitches quickly without waiting for the exact RNG, temporarily hardcode the roll values in your script for testing.

Example:

```go
{{ $isLucky := le (randInt 1 10001) (mult $luckyChance 100) }}
{{ if $isLucky }}
    {{ $roll := randInt 1 101 }}
```

Change it to:

```go
{{ $isLucky := true }}
{{ $roll := 3 }} {{/* Change this number to test specific events: 3=Crash, 4=Halving, 8=Swap, etc. */}}
```

---

## 1. Core Mechanics & Cooldowns

- [X] **TC 1.1: Basic Slip**  
  Action: Roll and get a normal slip (not lucky).  
  Expected: You gain +1 slip. Your streak resets to 0. Pity increases by 2. Cooldown is applied.

- [X] **TC 1.2: Basic Backflip (New Record)**  
  Action: Roll and get a backflip. Ensure it's higher than your current record.  
  Expected: You gain 0 slips. Text says "setting a new all-time personal record". Cooldown applied.

- [X] **TC 1.3: Basic Backflip (Tied Record)**  
  Action: Roll a backflip that matches your exact record.  
  Expected: Text says "pulling even with your all-time personal record." 

- [X] **TC 1.4: Basic Backflip (Chasing Record)**  
  Action: Roll a backflip while your streak is lower than your record.  
  Expected: Text says "chasing your all-time personal record." 

- [X] **TC 1.5: Cooldown Enforcement**  
  Action: Try to use the command immediately after rolling (without Turbo active).  
  Expected: Bot blocks the command and displays the "Try again <t:...:R>!" message. No slips or streaks are modified.

## 2. Global Economy Glitches

- [X] **TC 2.1: Trigger Market Crash**  
  Action: Hit a lucky roll between 1-3.  
  Expected: Global crash variable is set to 3. Text displays "MARKET CRASH!". Cooldown applied.

- [X] **TC 2.2: Market Crash - Slip Interaction**  
  Action: Roll a slip while Market Crash is active.  
  Expected: Text displays "MARKET CRASH ACTIVE! Your slip was worthless!". You gain 0 slips (even if you hit a multiplier). Crash counter decays.

- [X] **TC 2.3: Market Crash - Backflip Interaction**  
  Action: Roll a backflip while Market Crash is active.  
  Expected: Text displays "MARKET CRASH! Stole 5 Slips!". You gain +5 slips. Streak increases. Crash counter decays.

- [X] **TC 2.4: Market Crash - Expiration**  
  Action: Roll 3 times after a crash is triggered.  
  Expected: On the 4th roll, the economy returns to normal, and normal slips/backflips resume.

## 3. Turbo System

- [ ] **TC 3.1: Trigger Turbo Overdrive**  
  Action: Hit a lucky roll between 10-17.  
  Expected: Text displays "TURBO OVERDRIVE!". DB turbo is set to true. No cooldown is applied (you can roll immediately).

- [ ] **TC 3.2: Consume Turbo (Slip)**  
  Action: Roll a slip while Turbo is active.  
  Expected: Cooldown is bypassed. Earned slips are doubled. Turbo state resets to false. Cooldown is applied normally afterward.

- [ ] **TC 3.3: Consume Turbo (Backflip)**  
  Action: Roll a backflip while Turbo is active.  
  Expected: Cooldown is bypassed. Turbo state resets to false. Streak behaves normally.

## 4. Rank Swap Mechanics

- [ ] **TC 4.1: Top 10 Swap (Swapping Down)**  
  Precondition: User is Rank 1-10. There is at least one person below them on the leaderboard.  
  Action: Hit a Rank Swap roll (5-9).  
  Expected: User swaps with someone below them. Text says "You fumbled and swapped places... from below." DB updates both users' slips.

- [ ] **TC 4.2: Lower Bracket Swap (Swapping Up)**  
  Precondition: User is Rank 11+.
  Action: Hit a Rank Swap roll (5-9).
  Expected: User swaps with someone above them. Text says "You vaulted upwards and stole the position." DB updates both users' slips.

- [ ] **TC 4.3: Swap Boundary Protection (Bottom)**
  Precondition: User is in the top 10 but at the very bottom of the leaderboard (no one below them).
  Action: Hit a Rank Swap roll.
  Expected: Target index defaults to the last person. If the target is themselves, the swap fails and they receive a Golden Peel (x2 slips) instead.

- [ ] **TC 4.4: Solo Player Swap (No Target)**
  Precondition: The leaderboard is empty except for the testing user.
  Action: Hit a Rank Swap roll.
  Expected: Text says "tried to swap, but no one was there!". User gets a Golden Peel (x2 multiplier) instead.

## 5. Destructive & Multiplier Glitches

- [X] **TC 5.1: Halving Glitch**  
  Action: Hit a lucky roll of 4.  
  Expected: Text says "CATASTROPHIC ERROR!". User's current slip total is divided by 2 (rounded down).

- [ ] **TC 5.2: Slip-plosion (Rank 1)**  
  Precondition: User is Rank #1.  
  Action: Hit a Slip-plosion roll (18-27).  
  Expected: Text says "nullified your own slip." Earned slips are 0. Multiplier is 0.

- [ ] **TC 5.3: Slip-plosion (Rank 2+)**  
  Precondition: User is Rank 2 or lower.  
  Action: Hit a Slip-plosion roll (18-27).  
  Expected: Text says "destroyed 1 slip from <@Target>." The user directly above the roller loses 1 slip in the DB.

- [ ] **TC 5.4: Multiplier Verification**  
  Action: Trigger Mythic (x10), Cosmic (x5), and Golden (x2) peels.  
  Expected: The base slip (+1) is multiplied correctly and added to the user's total.

## 6. Pity System & Scaling Logic

- [ ] **TC 6.1: Pity Accumulation**  
  Action: Get normal (non-lucky) slips.  
  Expected: Check DB `TEST_banana_global` → pity should increase by 2 for every normal slip.

- [ ] **TC 6.2: Pity Consumption**  
  Precondition: User is Rank 11+ and Global Pity is > 0.  
  Action: Hit any lucky roll.  
  Expected: Global Pity resets to 0.

- [ ] **TC 6.3: Pity Ignored for Top 10**  
  Precondition: User is Rank 1-10 and Global Pity is > 0.  
  Action: Hit a lucky roll.  
  Expected: Global Pity remains untouched.

## 7. Dynamic Hype Messages

- [ ] **TC 7.1: Verify Rank 1 Hype**  
  Action: Have the highest slip count and roll a slip.  
  Expected: Message includes "🏆 The floor has officially given up all hope."

- [ ] **TC 7.2: Verify Rank 2-5 Hype**  
  Action: Be rank 2-5 and roll a slip.  
  Expected: Message includes "🍌 The floor is starting to fear."

- [ ] **TC 7.3: Verify Rank 6-10 Hype**  
  Action: Be rank 6-10 and roll a slip.  
  Expected: Message includes "⚠️ The floor is permanently dented."

- [ ] **TC 7.4: Verify Rank 11-20 Hype**  
  Action: Be rank 11-20 and roll a slip.  
  Expected: Message includes "👀 The floor has started wearing a helmet."

- [ ] **TC 7.5: Verify Rank 21+ Hype**  
  Action: Be rank 21+ and roll a slip.  
  Expected: Message includes "The floor doesn't know your name yet."

## 8. Prestige Display

- [X] **TC 8.1: Verify Prestige 0 display**  
  Action: Validate the username in flip, slip, topflip, topslip shows correctly, without spaces or prestiqge text wherever the username is displayed.  
  Expected: Username only "Username"

- [X] **TC 8.2: Verify Prestige 1+ display**  
  Action: Validate the username in flip, slip, topflip, topslip shows correctly with prestige text wherever the username is displayed.  
  Expected: Username includes "(🏆1+)Username"