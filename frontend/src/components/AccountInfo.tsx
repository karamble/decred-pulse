// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

import { Wallet, TrendingUp, Clock, AlertCircle } from 'lucide-react';

interface AccountInfoProps {
  accountName: string;
  totalBalance: number;
  spendableBalance: number;
  immatureBalance: number;
  unconfirmedBalance: number;
}

export const AccountInfo = ({
  accountName,
  totalBalance,
  spendableBalance,
  immatureBalance,
  unconfirmedBalance
}: AccountInfoProps) => {
  const formatDCR = (amount: number) => {
    return amount.toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 8
    });
  };

  return (
    <div className="p-6 rounded-xl bg-gradient-card backdrop-blur-sm border border-border/50 animate-fade-in">
      <div className="flex items-center gap-3 mb-6">
        <div className="p-2 rounded-lg bg-primary/10 border border-primary/20">
          <Wallet className="h-5 w-5 text-primary" />
        </div>
        <div>
          <h3 className="text-lg font-semibold">Account Balance</h3>
          <p className="text-sm text-muted-foreground">{accountName}</p>
        </div>
      </div>

      {/* Total Balance */}
      <div className="mb-6 p-4 rounded-lg bg-primary/5 border border-primary/10">
        <p className="text-sm text-muted-foreground mb-2">Total Balance</p>
        <p className="text-4xl font-bold text-primary">
          {formatDCR(totalBalance)} <span className="text-2xl">DCR</span>
        </p>
      </div>

      {/* Balance Breakdown */}
      <div className="space-y-4">
        <div className="flex items-center justify-between p-3 rounded-lg bg-success/5 border border-success/10">
          <div className="flex items-center gap-2">
            <TrendingUp className="h-4 w-4 text-success" />
            <span className="text-sm font-medium">Spendable</span>
          </div>
          <span className="font-semibold text-success">{formatDCR(spendableBalance)} DCR</span>
        </div>

        {immatureBalance > 0 && (
          <div className="flex items-center justify-between p-3 rounded-lg bg-warning/5 border border-warning/10">
            <div className="flex items-center gap-2">
              <Clock className="h-4 w-4 text-warning" />
              <span className="text-sm font-medium">Immature</span>
            </div>
            <span className="font-semibold text-warning">{formatDCR(immatureBalance)} DCR</span>
          </div>
        )}

        {unconfirmedBalance > 0 && (
          <div className="flex items-center justify-between p-3 rounded-lg bg-muted/5 border border-muted/10">
            <div className="flex items-center gap-2">
              <AlertCircle className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm font-medium">Unconfirmed</span>
            </div>
            <span className="font-semibold text-muted-foreground">{formatDCR(unconfirmedBalance)} DCR</span>
          </div>
        )}
      </div>

      <div className="mt-4 p-3 rounded-lg bg-muted/5 border border-muted/10">
        <p className="text-xs text-muted-foreground">
          <strong>Watch-Only Mode:</strong> This wallet can monitor balances and transactions but cannot spend funds.
        </p>
      </div>
    </div>
  );
};

