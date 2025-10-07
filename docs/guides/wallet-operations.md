# Wallet Operations

Complete guide to managing your Decred wallet through the Decred Pulse dashboard, including importing xpub keys, rescanning the blockchain, and monitoring sync progress.

## 📖 Overview

Wallet operations allow you to:
- Import extended public keys (xpub) for watch-only monitoring
- Rescan the blockchain to discover transactions
- Track sync progress in real-time
- Manage wallet connectivity and configuration

---

## 🔑 Import Extended Public Key (Xpub)

Import an **xpub key** to monitor wallet addresses without private key access. This creates a **watch-only wallet** that can view balances and transactions but cannot spend funds.

### What is an Xpub Key?

An **Extended Public Key (xpub)** is a master public key that can derive all public addresses for a wallet account. It's safe to share for monitoring purposes because:
- ✅ Can generate all receiving addresses
- ✅ Can view all transactions
- ✅ Cannot spend funds
- ✅ Cannot access private keys

**Use cases**:
- Monitor cold storage wallets
- Track balances without exposing private keys
- Audit wallet activity
- Portfolio monitoring

---

### How to Import Xpub

#### 1. Get Your Xpub Key

**From `dcrwallet` CLI:**
```bash
dcrctl --wallet getmasterpubkey default
```

**From Decrediton:**
1. Go to Accounts tab
2. Select account
3. Click "Export"
4. Copy the xpub key

**Example xpub**:
```
dpubZF6ScrXjYgjGdVL2FzAWMYpRbWbUk7VJT9JZjNGjqB9p5KMkJyKhGv8xv8riFP8...
```

**Security Notes**:
- ⚠️ Xpub reveals all addresses and transactions
- ⚠️ Can compromise privacy if shared carelessly
- ✅ Cannot spend funds (safe for viewing)
- ✅ Store securely but less critical than private keys

---

#### 2. Open Import Modal

In the Wallet Dashboard:
1. Click **"Import Xpub"** button (top right)
2. Modal dialog appears

---

#### 3. Enter Xpub Information

**Xpub Key** (required)
- Paste your extended public key
- Starts with `dpub` for Decred
- Long alphanumeric string

**Gap Limit** (required)
- Default: `200`
- Range: `20` - `1000`
- Recommended: `200` for most users

**What is Gap Limit?**
- Number of consecutive unused addresses to monitor
- Higher = discovers funds at higher address indices
- Lower = faster scanning, may miss funds

**When to adjust**:
- Low balance found: Increase to `500` or `1000`
- Missing transactions: Increase gap limit
- Slow scanning: Decrease to `100`
- Normal use: Keep at `200`

---

#### 4. Import Process

After clicking **"Import"**:

1. **Validation** (instant)
   - Xpub format checked
   - Gap limit validated

2. **Import** (~1-2 seconds)
   - Xpub registered with wallet
   - Account created

3. **Blockchain Rescan** (automatic)
   - Starts automatically
   - Progress bar appears
   - Duration: 5-30 minutes depending on:
     - Blockchain height
     - Gap limit
     - Transaction count
     - System performance

4. **Completion**
   - Progress bar disappears at 99%
   - Dashboard cards appear
   - Balances and transactions visible

---

### Import Modal Reference

```
┌─────────────────────────────────────────────┐
│  Import Extended Public Key                 │
├─────────────────────────────────────────────┤
│                                             │
│  Extended Public Key (xpub)                 │
│  ┌─────────────────────────────────────┐  │
│  │ dpubZF6ScrX...                      │  │
│  └─────────────────────────────────────┘  │
│                                             │
│  Gap Limit                                  │
│  ┌─────────────────────────────────────┐  │
│  │ 200                                 │  │
│  └─────────────────────────────────────┘  │
│                                             │
│  The gap limit determines how many          │
│  consecutive unused addresses to monitor.   │
│  Higher values find more transactions but   │
│  take longer to scan. Recommended: 200      │
│                                             │
│         [ Cancel ]    [ Import ]            │
└─────────────────────────────────────────────┘
```

---

## 🔄 Wallet Rescan

**Rescan** re-examines the blockchain to discover transactions and update balances. This is necessary when:
- Addresses were used while wallet was offline
- Gap limit was increased
- Transactions are missing
- Balance appears incorrect

### When to Rescan

✅ **After importing xpub** (automatic)
✅ **Increased gap limit** (manual)
✅ **Missing transactions** (manual)
✅ **Incorrect balance** (manual)
✅ **Wallet restored from seed** (automatic in dcrwallet)

⚠️ **Not needed for**:
- Regular operation
- New transactions (auto-detected)
- Wallet already synced

---

### How to Rescan

#### Automatic Rescan
- Triggered automatically when importing xpub
- No manual action needed
- Progress bar appears automatically

#### Manual Rescan

**Method 1: Dashboard Button**
1. Navigate to Wallet Dashboard
2. Click **"Rescan"** button (if available)
3. Confirm action
4. Progress bar appears

**Method 2: CLI** (Advanced)
```bash
dcrctl --wallet rescan
```

**Method 3: RPC Endpoint**
```bash
curl -X POST http://localhost:8080/api/wallet/rescan \
  -H "Content-Type: application/json"
```

---

### Rescan Process

#### Phase 1: Initialization
- Wallet prepares for rescan
- Dashboard cards hidden
- Progress bar appears
- Normal polling paused

#### Phase 2: Scanning
- Blockchain examined block by block
- Transactions discovered and indexed
- Balances calculated
- Progress updates every 2 seconds

**What you see**:
```
┌────────────────────────────────────────────┐
│  Scanning Blockchain                       │
│  ━━━━━━━━━━━━━━━━━━░░░░░░░░  68%          │
│                                            │
│  Block 1,016,234 / 1,016,401              │
│  Finding your transactions...              │
└────────────────────────────────────────────┘
```

**Duration**:
- Empty wallet: 5-10 minutes
- Active wallet: 10-30 minutes
- Large gap limit (1000): 30-60 minutes

#### Phase 3: Completion
- Progress reaches 99%
- Progress bar auto-hides
- Dashboard cards reappear
- Data refreshed immediately
- Normal polling resumes

---

## 📊 Sync Progress Monitoring

Real-time progress tracking during wallet rescan operations.

### Progress Bar Features

#### Visual Display
- **Percentage**: 0-99%
- **Progress Bar**: Animated fill
- **Block Count**: Current / Total
- **Status Message**: Operation description

#### Smart Polling
- **Frequency**: Every 2 seconds during active sync
- **Log Parsing**: Reads `dcrwallet.log` for progress
- **Stale Detection**: Stops if logs older than 2 minutes
- **Auto-Hide**: Disappears at 99% completion

#### Dashboard Behavior
- **During Sync**: Cards hidden, progress visible
- **After Sync**: Cards appear, progress hidden
- **On Navigation**: Persists if rescan active
- **Background**: Pauses other polling

---

### Sync Progress States

#### Active Sync
```
Status: Rescanning...
Progress: 42%
Action: Wait for completion
```

**What happens**:
- Progress bar visible
- Dashboard cards hidden
- Backend polls logs every 2s
- Frontend displays real-time updates

---

#### Stale Sync Detection
```
Status: Rescan inactive (logs stale)
Progress: N/A
Action: Rescan completed or stalled
```

**What happens**:
- Last log entry > 2 minutes old
- Considered inactive
- Polling stops automatically
- Dashboard resumes normal operation

**Triggers**:
- Rescan genuinely completed
- Wallet crashed or stopped
- Log rotation
- File permission issues

---

#### Sync Completion
```
Status: Complete
Progress: 99%
Action: Dashboard refreshes
```

**What happens**:
- Progress bar auto-hides
- Dashboard cards appear
- Data fetched immediately
- Normal 30s polling resumes

---

### Progress Tracking Technical Details

#### Backend: Log Parsing
**Location**: `/backend/handlers/wallet.go`

**Process**:
1. Opens `dcrwallet.log`
2. Reads last 500 lines
3. Searches for rescan messages:
   ```
   Rescanning blockchain for address...
   ```
4. Extracts progress data
5. Checks timestamp (< 2 minutes = active)
6. Returns JSON response

**Response Example**:
```json
{
  "isRescanning": true,
  "progress": 68.5,
  "currentBlock": 1016234,
  "totalBlocks": 1016401,
  "message": "Rescanning blockchain..."
}
```

#### Frontend: Progress Display
**Location**: `/frontend/src/pages/WalletDashboard.tsx`

**Process**:
1. Checks sync status on load
2. If active: Shows progress, hides cards
3. Polls every 2 seconds
4. Updates progress bar
5. On completion: Hides progress, shows cards
6. Fetches fresh data

---

## 🎯 Best Practices

### Import Xpub
✅ **Do**:
- Use gap limit of 400 for normal wallets
- Increase to 500-1000 if funds missing
- Wait for full rescan before using
- Keep xpub secure (privacy concern)

⚠️ **Don't**:
- Share xpub publicly (reveals all addresses)
- Set gap limit too low (may miss transactions)
- Navigate away during import
- Interrupt rescan process

---

### Wallet Rescan
✅ **Do**:
- Wait for blockchain sync completion first
- Use appropriate gap limit
- Monitor progress through dashboard
- Let rescan complete fully

⚠️ **Don't**:
- Rescan unnecessarily (wastes time)
- Stop rescan midway
- Rescan while blockchain syncing
- Use extremely high gap limits (>1000) unless needed

---

### Gap Limit Selection

| Scenario | Recommended Gap Limit | Reasoning |
|----------|----------------------|-----------|
| New wallet | 20-50 | Few addresses used |
| Normal use | 200 | Default, handles most cases |
| Active wallet | 200-500 | Many transactions |
| Missing funds | 500-1000 | High address indices |
| Legacy wallet | 1000+ | Very old or heavily used |

---

## 🐛 Troubleshooting

### Xpub Import Failed

**Problem**: Import button does nothing or shows error

**Solutions**:
1. **Invalid xpub**:
   - Verify xpub format (starts with `dpub`)
   - Check for copy/paste errors
   - Ensure complete string copied

2. **Wallet not connected**:
   - Check RPC connection
   - Verify wallet is running
   - Check backend logs

3. **Gap limit out of range**:
   - Use value between 20-1000
   - Default 400 is recommended

---

### Rescan Stuck or Frozen

**Problem**: Progress bar not moving

**Solutions**:
1. **Check logs**:
   ```bash
   docker compose logs -f dcrwallet
   ```

2. **Verify wallet running**:
   ```bash
   docker compose ps dcrwallet
   ```

3. **Stale detection active**:
   - Logs older than 2 minutes
   - Rescan may have completed
   - Refresh page to check

4. **Wallet crashed**:
   ```bash
   docker compose restart dcrwallet
   ```
   - Rescan will resume from checkpoint

---

### Progress Bar Won't Disappear

**Problem**: Stuck at high percentage (90%+)

**Solutions**:
1. **Wait**: May be finalizing (can take 1-2 minutes)

2. **Check completion**:
   - Look for transactions in history
   - Check if balance updated
   - May have completed despite display

3. **Refresh page**:
   - Force refresh (Ctrl+Shift+R)
   - Progress should clear

4. **Manual clear**:
   - Navigate away and back
   - Dashboard will check real status

---

### Cards Not Appearing After Rescan

**Problem**: Progress bar gone but no dashboard cards

**Solutions**:
1. **Force refresh**:
   ```
   Ctrl + Shift + R (Windows/Linux)
   Cmd + Shift + R (Mac)
   ```

2. **Check browser console** (F12):
   - Look for JavaScript errors
   - Check network requests
   - Verify API responses

3. **Check wallet status**:
   ```bash
   docker compose logs backend | grep -i wallet
   ```

4. **Restart backend**:
   ```bash
   docker compose restart backend
   ```

---

### Missing Transactions After Import

**Problem**: Expected transactions not showing

**Solutions**:
1. **Increase gap limit**:
   - Used addresses may be at high indices
   - Try 500 or 1000
   - Re-import xpub with higher limit

2. **Wrong account**:
   - Ensure correct account xpub
   - Check account number in source wallet
   - Import additional account xpubs if needed

3. **Blockchain not fully synced**:
   - Check node sync status
   - Wait for full blockchain sync
   - Then rescan wallet

4. **Wrong network**:
   - Verify mainnet vs testnet
   - Check xpub corresponds to correct network

---

### Low Balance Found

**Problem**: Balance lower than expected

**Solutions**:
1. **Increase gap limit**:
   - Current limit may be too low
   - High address indices not scanned
   - Try doubling the gap limit

2. **Multiple accounts**:
   - Import xpubs for all accounts
   - Check source wallet account list
   - Each account needs separate xpub

3. **Rescan incomplete**:
   - Wait for 100% completion
   - Check progress bar reached 99%
   - Allow background finalization

4. **Verify in source wallet**:
   - Compare with actual wallet balance
   - Check transaction history matches
   - Confirm correct xpub exported

---

## 📊 Gap Limit Explained

### What is a Gap Limit?

**Definition**: The number of consecutive unused addresses the wallet will monitor before assuming no more transactions exist.

**BIP0044 Standard**: Defines gap limit for HD wallets.

**Example**:
```
Address 0: Used ✅
Address 1: Used ✅
Address 2: Unused
Address 3: Unused
Address 4: Used ✅
Address 5: Unused
...
Address 201: Unused

Gap Limit = 200
```

**With gap limit 200**:
- Monitors up to address 201
- Finds address 4 (used)
- Finds any usage up to address 204

**With gap limit 3**:
- Stops at address 5
- Misses address 4 (gap of 3)
- Incomplete balance

---

### Choosing Gap Limit

**Factors to consider**:
- **Wallet age**: Older = higher limit
- **Usage pattern**: Random = higher limit
- **Address reuse**: Sequential = lower limit
- **Scan time**: Higher = slower

**Scan time impact**:
```
Gap Limit 20:   ~2 minutes
Gap Limit 200:  ~10 minutes
Gap Limit 500:  ~25 minutes
Gap Limit 1000: ~50 minutes
```

---

## 🔐 Security Considerations

### Xpub Safety

**What xpub reveals**:
- ✅ All public addresses
- ✅ All transactions
- ✅ Complete balance history

**What xpub cannot**:
- ❌ Spend funds
- ❌ Access private keys
- ❌ Sign transactions

**Privacy impact**:
- ⚠️ Links all addresses together
- ⚠️ Reveals transaction patterns
- ⚠️ Shows complete financial history

**Best practices**:
- Share xpub only with trusted parties
- Use different xpub for different purposes
- Consider privacy implications
- Store securely (less critical than private keys)

---

### Watch-Only Wallet Limitations

**Can do**:
- ✅ View balances
- ✅ Monitor transactions
- ✅ Generate receiving addresses
- ✅ Track transaction history

**Cannot do**:
- ❌ Send transactions
- ❌ Purchase tickets
- ❌ Sign messages
- ❌ Access private keys
- ❌ Vote with tickets
- ❌ Revoke tickets

**Use case**: Safe monitoring without spending risk.

---

## 🚀 Advanced Usage

### Multiple Account Monitoring

Import xpubs for multiple accounts:

1. **Export each account xpub**:
   ```bash
   dcrctl --wallet getmasterpubkey default
   dcrctl --wallet getmasterpubkey mixed
   dcrctl --wallet getmasterpubkey unmixed
   ```

2. **Import each individually**:
   - Use Import Xpub modal for each
   - Set appropriate gap limit per account
   - Wait for each rescan to complete

3. **View in dashboard**:
   - All accounts appear in Accounts card
   - Each shows individual balances
   - Cumulative total includes all

---

### Automated Monitoring

For automated balance checking (advanced):

**API endpoint**:
```bash
curl http://localhost:8080/api/wallet/dashboard
```

**Response includes**:
```json
{
  "accountInfo": {...},
  "accounts": [...],
  "stakingInfo": {...},
  "walletStatus": {...}
}
```

See [API Reference](../api/wallet-endpoints.md) for details.

---

## 📚 Related Documentation

- **[Wallet Dashboard](../features/wallet-dashboard.md)** - Dashboard overview
- **[Wallet Setup](../wallet-setup.md)** - Initial configuration
- **[Staking Guide](../features/staking-guide.md)** - Staking information
- **[API Reference](../api/wallet-endpoints.md)** - API documentation
- **[Troubleshooting](troubleshooting.md)** - Common issues

---

## ✅ Operations Checklist

### Before Importing Xpub
- [ ] Wallet RPC connected
- [ ] Blockchain fully synced
- [ ] Xpub key copied correctly
- [ ] Gap limit decided (400 recommended)
- [ ] Time allocated (10-30 minutes)

### During Import
- [ ] Modal shows progress
- [ ] Don't navigate away
- [ ] Monitor progress bar
- [ ] Watch for completion

### After Import
- [ ] Verify balance matches expected
- [ ] Check transaction history
- [ ] Review all accounts
- [ ] Test dashboard features

### If Issues
- [ ] Check troubleshooting section
- [ ] Review logs
- [ ] Try higher gap limit
- [ ] Contact support if needed

---

**Need Help?** Check the [FAQ](../reference/faq.md) or [Troubleshooting Guide](troubleshooting.md)

