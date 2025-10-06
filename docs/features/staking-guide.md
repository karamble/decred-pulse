# Staking Guide

Complete guide to Decred Proof-of-Stake (PoS) staking using the Decred Pulse dashboard. Learn how to monitor ticket pools, track your tickets, and understand staking rewards.

## 📖 What is Decred Staking?

**Decred staking** is the process of time-locking DCR to purchase **tickets** that participate in network governance and consensus. In return, you earn staking rewards when your tickets are called to vote.

### Key Concepts

**Ticket**: A special transaction that locks DCR for voting
**Ticket Price**: The cost to purchase a ticket (changes every 144 blocks)
**Ticket Pool**: All active tickets waiting to be called to vote
**Voting**: When your ticket is selected to validate a block
**VSP**: Voting Service Provider (optional, ensures ticket votes even if offline)

---

## 🎫 Ticket Lifecycle

### 1. Purchase (Mempool) ⏳
- Buy ticket at current ticket price
- Transaction enters mempool
- Awaiting confirmation
- **Duration**: ~5 minutes (1 block)

### 2. Immature 🐣
- Ticket confirmed on blockchain
- Not yet eligible to vote
- Gaining maturity
- **Duration**: ~21 hours (256 blocks)

### 3. Live/Unspent ✅
- Ticket enters the pool
- Eligible to be called for voting
- **Duration**: ~28 days average (142 days max)
- **Probability**: Random selection per block

### 4. Voting 🗳️
- Ticket is selected to vote
- Validates previous block
- Earns voting reward
- **Reward**: ~0.8% of ticket price

### 5. Voted ✅
- Vote confirmed
- Reward + ticket price returned
- Rewards are immature for 256 blocks
- **Total Return**: Ticket price + ~0.8% reward

### Alternative Outcomes

**Expired** ⌛
- Not selected within 142 days (~40,960 blocks)
- Must be revoked to recover funds
- No reward earned

**Revoked** ❌
- Manual or automatic revocation
- Original ticket price returned (minus small fee)
- No reward earned

---

## 📊 Understanding the Dashboard

### Ticket Pool Info Card

Monitor the global Decred ticket pool:

#### Pool Size 🎫
**Current**: Number of live tickets in the pool
**Target**: ~40,960 tickets
**Significance**: Higher pool = more decentralization

**What you see**:
```
Pool Size: 41,095 tickets
```

**What it means**:
- Network has 41,095 active tickets
- Your ticket competes with these for voting
- Healthy network participation

---

#### Current Difficulty 💎
**What it is**: The current ticket price in DCR
**Updates**: Every 144 blocks (~12 hours)
**Algorithm**: Adjusts to maintain ~40,960 pool size

**What you see**:
```
Current Difficulty: 293.08 DCR
```

**What it means**:
- Next ticket purchase costs 293.08 DCR
- This is the price you'll pay now
- Changes at next difficulty adjustment

---

#### Next Difficulty 🔄
**What it is**: Estimated next ticket price
**Timing**: Updates in next difficulty window
**Calculation**: Based on recent demand

**What you see**:
```
Next Difficulty: 293.08 DCR
```

**What it means**:
- Expected price for next window
- Same = stable demand
- Higher = increasing demand
- Lower = decreasing demand

---

#### Estimated Difficulty Range 📊
**What it is**: Price prediction range
**Based on**: Recent ticket purchases
**Algorithm**: `estimatestakediff` RPC

**What you see**:
```
Min: 291.54 DCR
Expected: 292.20 DCR
Max: 294.59 DCR
```

**What it means**:
- **Min**: Lowest possible next price
- **Expected**: Most likely next price  
- **Max**: Highest possible next price

**Usage**: Plan your ticket purchases

---

#### Mempool Tickets ⏳
**What it is**: Network-wide pending ticket purchases
**Status**: Awaiting confirmation
**Accuracy**: Calculated using stake difficulty

**What you see**:
```
All Mempool Tickets: 15
```

**What it means**:
- 15 tickets being purchased network-wide
- Will confirm in next block
- Will become immature tickets

**Technical Note**: Count calculated by:
```
Total Stake Submission Outputs / Current Stake Difficulty
```
This accurately counts tickets even in coinjoin transactions.

---

### My Tickets Card

Track your personal ticket statistics:

#### Mempool ⏳
**Your pending tickets**
- Just purchased
- Awaiting first confirmation
- Not yet immature

**Typical count**: 0-1 (unless batch purchasing)

**What to expect**:
- Moves to immature after ~5 minutes
- If stuck, check transaction fee

---

#### Immature 🐣
**Your maturing tickets**
- Confirmed but not yet live
- Requires 256 confirmations
- ~21 hours until live

**Typical count**: 0-5 (depending on purchase frequency)

**What to expect**:
- Becomes live after 256 blocks
- Cannot vote yet
- Funds are locked

---

#### Live/Unspent ✅
**Your active tickets**
- In the ticket pool
- Eligible to vote on every block
- Average lifetime: ~28 days

**Typical count**: Varies by staking strategy

**What to expect**:
- Random selection for voting
- Average 28-day wait
- Max 142-day expiry

**Probability per block**:
```
Chance = 5 votes needed / ~40,960 pool size
      ≈ 0.012% per block
      ≈ 28 days average
```

---

#### Voted 🗳️
**Your successful votes**
- Tickets that have voted
- Rewards earned
- Historical count

**Typical count**: Accumulates over time

**What to expect**:
- Rewards + ticket price returned
- Rewards immature for 256 blocks
- Then available as spendable

---

#### Revoked ❌
**Your revoked tickets**
- Manually revoked
- Auto-revoked after expiry
- Funds recovered (minus fee)

**Typical count**: Low (most tickets vote)

**What to expect**:
- ~5% of tickets expire (statistically)
- Auto-revoke with `enablevoting=1`
- Manual revoke if needed

---

#### Expired ⌛
**Your expired unspent tickets**
- Not selected within 142 days
- Awaiting revocation
- Funds still locked

**Typical count**: Should be 0 with auto-revoke

**Action needed**:
- Manually revoke if auto-revoke disabled
- Run: `dcrctl rebroadcastmissed`

---

#### Total Subsidy 💰
**Cumulative voting rewards**
- All-time earnings from voting
- Excludes original ticket prices
- Only the ~0.8% rewards

**Typical value**: Grows with each vote

**Calculation**:
```
Per Vote ≈ Ticket Price × 0.008
```

**Example**:
- Ticket Price: 293 DCR
- Reward: ~2.34 DCR per vote
- 10 votes: ~23.4 DCR total subsidy

---

## 💰 Staking Economics

### Costs

**Ticket Price**
- Current: Check "Current Difficulty" card
- Historical Range: 30-600 DCR
- Current Range: 200-350 DCR (typical)

**VSP Fee** (if using VSP)
- Typical: 0.5-5% of reward
- One-time per ticket
- Deducted from voting reward

**Transaction Fees**
- Purchase: ~0.01 DCR
- Revocation: ~0.01 DCR
- Voting: Paid by VSP or your wallet

### Returns

**Voting Reward**
- Approximate: ~0.8% of ticket price per vote
- Annual ROI: ~6-7% (varies)
- Compounds if restaked

**Example**:
```
Ticket Price: 300 DCR
Voting Reward: ~2.4 DCR
Time Locked: ~28 days average
Annual ROI: ~7.1%
```

**Calculation**:
```
ROI per vote = (Reward / Ticket Price) × 100
             = (2.4 / 300) × 100
             ≈ 0.8%

Annual ROI = 0.8% × (365 / 28)
           ≈ 10.4% theoretical

Actual ROI ≈ 6-7% (accounting for expiries, timing)
```

### Risk Factors

**Expiration Risk**
- ~5% chance ticket expires without voting
- No reward if expired
- Can revoke to recover funds

**Price Volatility**
- DCR price may fluctuate
- Locked for ~28 days average
- Consider your risk tolerance

**Opportunity Cost**
- Funds locked during staking
- Cannot use for other purposes
- Compare returns to alternatives

---

## 🎯 Staking Strategies

### Conservative Staking
**Goal**: Steady, predictable returns
**Approach**:
- Purchase tickets regularly (DCA)
- Use VSP for reliability
- Auto-revoke expired tickets
- Reinvest rewards

**Best for**: Long-term holders, passive stakers

---

### Active Staking
**Goal**: Maximize returns, active management
**Approach**:
- Time purchases at low difficulty
- Monitor pool size trends
- Solo stake (if always online)
- Optimize transaction fees

**Best for**: Experienced users, technical operators

---

### Accumulation Staking
**Goal**: Grow DCR holdings
**Approach**:
- Reinvest all rewards into new tickets
- Compound returns over time
- Long-term commitment
- Regular monitoring

**Best for**: HODLers, accumulation phase

---

## 📱 Monitoring Your Stakes

### Daily Checks
✅ Live ticket count (stable or growing?)
✅ Voted tickets (earning rewards?)
✅ Expired tickets (need revoking?)
✅ Immature tickets (recently purchased?)

### Weekly Checks
✅ Total subsidy (rewards accumulating?)
✅ Ticket difficulty trends (buy now or wait?)
✅ Pool size changes (network health?)
✅ Revoked count (within expected ~5%?)

### Monthly Reviews
✅ Total ROI calculation
✅ Reward reinvestment strategy
✅ Difficulty trend analysis
✅ VSP performance (if using)

---

## 🔧 Staking Operations

### Purchasing Tickets

**Via `dcrwallet` CLI:**
```bash
dcrctl --wallet purchaseticket default 300 1 $(dcrctl --wallet getnewaddress)
```

**Parameters**:
- `default`: Account to purchase from
- `300`: Maximum price willing to pay
- `1`: Number of tickets to purchase
- `getnewaddress`: Commitment address

**Via `dcrwallet` RPC:**
See [Wallet Operations Guide](../guides/wallet-operations.md) for API details.

---

### Automatic Ticket Buying

**Enable in `dcrwallet.conf`:**
```ini
enableticketbuyer=1
ticketbuyer.limit=10
ticketbuyer.maxpricerelative=1.25
```

**Settings**:
- `limit`: Max tickets to maintain
- `maxpricerelative`: Max price relative to average (1.25 = 125%)

**Caution**: Only for experienced users with sufficient balance.

---

### Using a VSP (Voting Service Provider)

**Why use a VSP?**
- Ensures votes even if you're offline
- Professional infrastructure
- Small fee for reliability

**Setup**:
1. Choose a VSP from [decred.org/vsp](https://decred.org/vsp/)
2. Register and get VSP pubkey
3. Configure in `dcrwallet`:
   ```ini
   enablevoting=1
   vsp.url=https://your-vsp.com
   vsp.pubkey=your-vsp-pubkey
   ```

**Dashboard Monitoring**:
- Voted tickets should appear normally
- VSP handles the voting process
- You earn rewards minus VSP fee

---

### Revoking Expired Tickets

**Automatic Revocation** (recommended):
```ini
# dcrwallet.conf
enablevoting=1
```

**Manual Revocation**:
```bash
dcrctl --wallet rebroadcastmissed
```

**Check for Missed**:
```bash
dcrctl --wallet getstakeinfo
```

Look for `unspentexpired` count.

---

## 🐛 Troubleshooting

### No Tickets Showing in Dashboard

**Problem**: My Tickets card shows zeros

**Solutions**:
1. **Watch-only wallet**: 
   - Cannot display tickets with xpub import
   - Need full RPC connection
   - See: [Wallet Setup](../wallet-setup.md)

2. **Recently purchased**:
   - Check mempool count
   - Wait for confirmation
   - May take 5+ minutes

3. **RPC connection issue**:
   - Verify wallet is connected
   - Check credentials
   - Restart backend/wallet

---

### Ticket Stuck in Mempool

**Problem**: Mempool ticket not confirming

**Solutions**:
1. **Low transaction fee**:
   - May take multiple blocks
   - Check mempool competition
   - Consider fee bump (advanced)

2. **Full mempool**:
   - Wait for next block
   - Typically confirms within 1-2 blocks

3. **Transaction error**:
   - Check wallet logs
   - Verify sufficient balance
   - Ensure correct ticket price

---

### Expired Tickets Not Auto-Revoking

**Problem**: Expired count growing

**Solutions**:
1. **Enable auto-revoke**:
   ```ini
   enablevoting=1
   ```

2. **Manual revoke**:
   ```bash
   dcrctl --wallet rebroadcastmissed
   ```

3. **Check wallet online**:
   - Wallet must be running
   - Must be synced
   - Must have connectivity

---

### Lower Than Expected ROI

**Problem**: Returns below advertised ~7%

**Possible causes**:
1. **Recent start**: Not enough data yet
2. **Expired tickets**: ~5% don't vote
3. **VSP fees**: Reduces net return
4. **Price timing**: Bought at high difficulty
5. **Short duration**: ROI averages over time

**Check**:
- Total voted vs expired ratio
- Average time to vote
- VSP fee structure
- Difficulty trend during purchases

---

## 📊 Advanced Metrics

### Effective ROI Calculation

```
Total Earned = Total Subsidy
Total Invested = (Avg Ticket Price) × (Total Tickets Bought)
Time Period = Days since first ticket

Annual ROI = (Total Earned / Total Invested) × (365 / Time Period) × 100
```

### Expected Vote Time

```
Average Time = (Pool Size / 5 votes per block) × 5 minutes per block

Current: (40,960 / 5) × 5 = 40,960 minutes
       = 28.4 days average
```

### Expiration Probability

```
Probability = 1 - (1 - 5/PoolSize)^MaxBlocks

Max Blocks = 40,960 (142 days)
Pool Size = 40,960
Probability ≈ 5%
```

---

## 🎓 Learning Resources

### Official Documentation
- [Decred Staking Guide](https://docs.decred.org/proof-of-stake/overview/)
- [Ticket Lifecycle](https://docs.decred.org/proof-of-stake/overview/#ticket-lifecycle)
- [VSP List](https://decred.org/vsp/)

### Dashboard Features
- [Wallet Dashboard](wallet-dashboard.md) - Balance and ticket overview
- [Wallet Operations](../guides/wallet-operations.md) - Manage your wallet
- [Transaction History](transaction-history.md) - Track transactions

### Community
- [Decred Discord](https://discord.gg/decred) - Ask questions
- [Decred Matrix](https://chat.decred.org) - Community chat
- [r/decred](https://reddit.com/r/decred) - Reddit community

---

## ✅ Staking Checklist

Before you start staking:

- [ ] Wallet fully synced
- [ ] Sufficient DCR balance (check current ticket price + fees)
- [ ] Understand ticket lifecycle
- [ ] Decided: Solo stake or VSP?
- [ ] Enabled auto-revoke (`enablevoting=1`)
- [ ] Reviewed current difficulty trends
- [ ] Calculated expected ROI
- [ ] Set up dashboard monitoring

During staking:

- [ ] Monitor dashboard daily
- [ ] Track live ticket count
- [ ] Watch for expired tickets
- [ ] Review voting rewards
- [ ] Adjust strategy as needed

---

## 🚀 Next Steps

Ready to start staking?

1. **[Setup Wallet](../wallet-setup.md)** - Configure your wallet
2. **[Purchase Tickets](../guides/wallet-operations.md)** - Buy your first ticket
3. **[Monitor Dashboard](wallet-dashboard.md)** - Track your stakes
4. **[Join Community](../reference/faq.md#community)** - Get support

---

**Happy Staking!** 🎫✨

Questions? Check the [FAQ](../reference/faq.md) or [Troubleshooting Guide](../guides/troubleshooting.md)

