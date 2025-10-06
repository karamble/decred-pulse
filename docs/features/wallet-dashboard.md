# Wallet Dashboard

The **Wallet Dashboard** provides comprehensive monitoring and management of your Decred wallet, including account balances, transaction history, staking information, and ticket management.

## 📊 Overview

The Wallet Dashboard displays real-time information about your wallet's financial status and staking activities. It's designed to work with both full wallets (via RPC) and watch-only wallets (via imported xpub keys).

**Access**: Click the **"Wallet"** button in the header navigation.

---

## 🎯 Dashboard Components

### 1. Account Balance Card

Displays wallet-wide balance information:

#### Cumulative Total
- **Total balance** across all accounts
- Displayed with **2 decimal places**
- Additional decimal places shown in smaller font

#### Total Spendable
- **Available balance** for transactions
- Excludes locked and immature funds
- Can be used for sending or ticket purchases

#### Total Locked by Tickets
- **Funds locked** in active tickets
- Automatically unlocked when tickets vote or expire
- Not spendable until tickets are spent

**Visual Features:**
- Large, prominent balance display
- Color-coded for quick scanning
- Real-time updates every 30 seconds

---

### 2. Transaction History Card

View your wallet's transaction activity:

#### Features
- **Initial Display**: Shows last 10 transactions
- **Progressive Loading**: Click "Load 50 More" to see additional transactions
- **Clickable Transactions**: Click any transaction to view details in block explorer

#### Transaction Information
Each transaction shows:
- **Type Icon**: Visual indicator (send, receive, ticket, vote, revocation, coinbase)
- **Category**: Send, receive, ticket, vote, revocation, immature, generate
- **Amount**: Transaction value in DCR
- **Confirmations**: Number of confirmations (or "Pending")
- **Time**: Relative time (e.g., "2 hours ago") or formatted date
- **Address**: Relevant address (truncated)
- **Transaction ID**: Truncated txid with copy functionality

#### Transaction Types

**Regular Transactions:**
- 🔽 **Receive** - Green arrow down
- 🔼 **Send** - Red arrow up

**Staking Transactions:**
- 🎫 **Ticket** - Yellow ticket icon (ticket purchase)
- ✅ **Vote** - Green checkmark (ticket voted)
- ❌ **Revocation** - Red X (ticket revoked)

**Mining/Generation:**
- 💰 **Coinbase** - Gold coins icon (mined block)
- ⏱️ **Immature** - Clock icon (immature rewards)

#### Status Badges
- **Pending**: Red badge for 0 confirmations
- **X confirmations**: Gray badge showing confirmation count

---

### 3. Accounts Card

Lists all wallet accounts with detailed balance breakdown:

#### Account Information
For each account, displays:
- **Account Name**: default, mixed, unmixed, imported, etc.
- **Total Balance**: Overall account balance

#### Granular Balance Details

Each account shows **six balance types**:

**Spendable** 💎
- Funds available for immediate use
- Can send or purchase tickets

**Locked by Tickets** 🎫
- Funds locked in active tickets
- Released when ticket votes/expires

**Immature Coinbase** ⛏️
- Mining rewards awaiting maturity
- Requires 256 confirmations

**Immature Stake Generation** 🎲
- Staking rewards awaiting maturity
- Requires 256 confirmations

**Unconfirmed** ⏳
- Transactions pending confirmation
- Not yet spendable

**Voting Authority** 🗳️
- Funds in tickets you have voting rights for
- May not own the ticket

#### Formatting
- All balances shown with **2 decimal places**
- Icons for quick visual identification
- Compact, space-efficient layout
- Only non-zero balances displayed

---

### 4. Ticket Pool Info Card

Global Decred network ticket pool information:

#### Current Network Statistics

**Pool Size** 🎫
- Total number of live tickets in the pool
- Target: ~40,960 tickets

**Current Difficulty** 💎
- **Current ticket price** in DCR
- Price you'll pay for next ticket purchase

**Next Difficulty** 🔄
- Estimated ticket price for next window
- Adjusts every 144 blocks

**Estimated Difficulty Range** 📊
- **Min**: Minimum estimated price
- **Expected**: Most likely price
- **Max**: Maximum estimated price

**Mempool Tickets** ⏳
- Pending ticket purchases across the network
- Awaiting confirmation

---

### 5. My Tickets Card

Your personal ticket statistics:

#### Your Ticket Counts

**Mempool** ⏳
- Your tickets waiting for confirmation
- Not yet active in pool

**Immature** 🐣
- Recently purchased tickets
- Require 256 confirmations (~21 hours)

**Live** ✅
- Active tickets in the pool
- Eligible to vote

**Voted** 🗳️
- Tickets that have voted
- Earned voting rewards

**Revoked** ❌
- Expired tickets that were revoked
- Funds returned (minus fee)

**Expired** ⌛
- Tickets that expired without voting
- Can be revoked to recover funds

**Total Subsidy** 💰
- Total voting rewards earned
- Accumulated from all votes

#### Watch-Only Wallet Disclaimer

When no tickets are shown:
```
🎫 No tickets found

Connect to an external wallet via RPC to see stats

Note: Tickets cannot be detected on watch-only wallets 
with imported xpub keys
```

**Why?** Watch-only wallets (xpub imports) cannot access ticket information because:
- Ticket data requires private key access
- Xpub only provides address monitoring
- Connect via RPC for full ticket tracking

---

## 🔄 Sync Progress Tracking

When performing wallet operations (rescan, xpub import), a **sync progress bar** appears:

### Features
- **Real-time Progress**: Shows percentage complete
- **Log Parsing**: Monitors `dcrwallet.log` for progress
- **Smart Polling**: Polls every 2 seconds during sync
- **Auto-Hide**: Disappears at 99% completion
- **Stale Detection**: Stops polling if logs become stale

### During Sync
- Dashboard cards are hidden
- Only sync progress bar visible
- Regular polling paused to prevent timeouts
- Background RPC calls suspended

### After Sync Completion
- Progress bar automatically hides
- Dashboard cards reappear
- Normal polling resumes (30-second intervals)
- Data refreshed immediately

---

## ⚙️ Configuration

### Auto-Refresh
Dashboard refreshes every **30 seconds** automatically.

To adjust refresh interval, modify:
```typescript
// frontend/src/pages/WalletDashboard.tsx
useEffect(() => {
  const interval = setInterval(fetchData, 30000); // Change 30000 to desired ms
  return () => clearInterval(interval);
}, []);
```

### Transaction Display Limits
By default:
- Initial display: **10 transactions**
- Load more: **50 transactions** per click
- Backend fetches: **100 transactions** total

To adjust:
```typescript
// frontend/src/components/TransactionHistory.tsx
const [visibleCount, setVisibleCount] = useState(10);  // Initial display
const loadMoreCount = 50;  // Load more increment

// Backend: modify count parameter
const data = await getWalletTransactions(100);  // Total to fetch
```

---

## 🔗 Related Documentation

- **[Wallet Setup](../wallet-setup.md)** - Initial wallet configuration
- **[Wallet Operations](../guides/wallet-operations.md)** - Import xpub, rescan, sync
- **[Staking Guide](staking-guide.md)** - Complete staking information
- **[Transaction History](transaction-history.md)** - Transaction management details
- **[API Reference](../api/wallet-endpoints.md)** - Wallet API endpoints

---

## 💡 Tips & Best Practices

### For Best Performance
1. ✅ Wait for full blockchain sync before rescanning
2. ✅ Use appropriate gap limit (200 recommended)
3. ✅ Don't navigate away during wallet rescan
4. ✅ Monitor sync progress through dashboard

### For Accurate Balance Display
1. ✅ Allow time for confirmations to mature
2. ✅ Check "immature" balances for pending rewards
3. ✅ Remember locked balance requires ticket voting/expiry
4. ✅ Refresh dashboard if balances seem outdated

### For Watch-Only Wallets
1. ✅ Import xpub for address monitoring
2. ✅ Use appropriate gap limit (200+)
3. ⚠️ Ticket info requires full RPC connection
4. ⚠️ Cannot send transactions (watch-only)

### For Transaction History
1. ✅ Click transactions to view blockchain details
2. ✅ Use "Load More" for older transactions
3. ✅ Check confirmations before considering final
4. ✅ Reference txid for support/tracking

---

## 🐛 Troubleshooting

### Balances Not Updating
**Problem**: Dashboard shows outdated balances

**Solutions:**
1. Check wallet is synced: Verify sync status in header
2. Wait for confirmations: Check transaction confirmations
3. Refresh manually: Click refresh button
4. Check RPC connection: Verify wallet connectivity

### No Transactions Showing
**Problem**: Transaction history is empty

**Solutions:**
1. **Watch-only wallet**: Import xpub with correct gap limit
2. **New wallet**: No transactions yet
3. **High address index**: Increase gap limit and rescan
4. **Sync incomplete**: Wait for wallet to finish syncing

### Ticket Counts Incorrect
**Problem**: Ticket statistics don't match expectations

**Solutions:**
1. **Watch-only wallet**: Connect via RPC for ticket data
2. **Sync in progress**: Wait for sync completion
3. **Expired tickets**: Check "expired" count
4. **Recently purchased**: Check "immature" count

### Sync Progress Stuck
**Problem**: Rescan progress bar frozen

**Solutions:**
1. Check dcrwallet logs: `docker compose logs -f dcrwallet`
2. Verify dcrwallet is running: `docker compose ps`
3. Check for stale logs: Backend auto-detects stale progress
4. Restart if necessary: `docker compose restart dcrwallet`

### Cards Not Appearing After Rescan
**Problem**: Dashboard cards don't show after rescan completes

**Solutions:**
1. Refresh page: Force browser refresh (Ctrl+Shift+R)
2. Check for errors: Open browser console (F12)
3. Verify completion: Check logs for 100% progress
4. Manual refresh: Click refresh button in dashboard

---

## 📊 Understanding Balance Types

### Cumulative Total
**What it is**: Sum of all balances across all accounts

**Includes:**
- Spendable
- Locked by tickets
- Immature rewards

**Excludes:**
- Unconfirmed transactions (until confirmed)

**Use case**: Total wallet value

---

### Total Spendable
**What it is**: Immediately available funds

**Can be used for:**
- Sending transactions
- Purchasing tickets
- Any wallet operation

**Excludes:**
- Locked funds
- Immature rewards
- Unconfirmed transactions

**Use case**: Available balance for spending

---

### Total Locked by Tickets
**What it is**: Funds committed to active tickets

**Locked in:**
- Immature tickets
- Live tickets in pool
- Recently voted tickets (awaiting maturity)

**Released when:**
- Ticket votes (rewards mature after 256 blocks)
- Ticket expires (must revoke first)
- Ticket is revoked

**Use case**: Staking commitment tracking

---

### Immature Coinbase Rewards
**What it is**: Block mining rewards awaiting maturity

**Requirements:**
- 256 confirmations (~21 hours)
- Mined by your wallet's addresses

**After maturity:**
- Moves to spendable balance
- Can be used immediately

**Use case**: Mining reward tracking

---

### Immature Stake Generation
**What it is**: Voting/staking rewards awaiting maturity

**Requirements:**
- 256 confirmations (~21 hours)
- Earned from ticket votes

**After maturity:**
- Moves to spendable balance
- Includes original ticket price + reward

**Use case**: Staking reward tracking

---

### Unconfirmed
**What it is**: Pending transactions

**Status:**
- Not yet in a block
- Or recently in mempool
- Awaiting confirmation

**Becomes spendable:**
- After first confirmation
- May require multiple confirmations for large amounts

**Use case**: Pending transaction tracking

---

### Voting Authority
**What it is**: Tickets you control voting rights for

**Scenarios:**
- Solo staking: Matches your tickets
- VSP staking: May include delegated tickets
- Split tickets: May own partial ticket

**Use case**: Governance participation tracking

---

## 🔐 Security Considerations

### RPC Credentials
- Stored in backend only
- Never exposed to frontend
- Use strong passwords
- Rotate regularly

### Watch-Only Wallets
- Safe for monitoring
- Cannot spend funds
- Limited information access
- Xpub can be shared (carefully)

### Wallet Access
- Secure RPC endpoint
- Use TLS for remote access
- Firewall RPC ports
- Monitor access logs

---

## 📈 Next Steps

After setting up your wallet dashboard:

1. **[Import Xpub](../guides/wallet-operations.md#import-xpub)** - Add watch-only addresses
2. **[Start Staking](staking-guide.md)** - Purchase tickets and earn rewards
3. **[Monitor Transactions](transaction-history.md)** - Track your wallet activity
4. **[Backup Wallet](../guides/backup-restore.md)** - Protect your funds

---

**Questions?** Check the [FAQ](../reference/faq.md) or [Troubleshooting Guide](../guides/troubleshooting.md)

