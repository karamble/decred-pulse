// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

import { AccountInfo } from '../services/api';
import { Wallet } from 'lucide-react';

interface AccountsListProps {
  accounts: AccountInfo[];
}

export const AccountsList = ({ accounts }: AccountsListProps) => {
  return (
    <div className="p-6 rounded-xl bg-gradient-card backdrop-blur-sm border border-border/50 animate-fade-in">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-2">
          <Wallet className="w-5 h-5 text-primary" />
          <h3 className="text-lg font-semibold">Accounts</h3>
        </div>
        <span className="text-sm text-muted-foreground">{accounts.length} account(s)</span>
      </div>

      {accounts.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground">No accounts found</p>
          <p className="text-sm text-muted-foreground mt-2">
            Import an xpub to create a new account
          </p>
        </div>
      ) : (
        <div className="space-y-3">
          {accounts.map((account) => (
            <div
              key={account.accountName}
              className="p-4 rounded-lg bg-background/50 border border-border/30 hover:border-primary/30 transition-colors"
            >
              <div className="flex items-center justify-between mb-3">
                <div>
                  <h4 className="font-medium text-foreground">{account.accountName}</h4>
                </div>
                <div className="text-right">
                  <p className="text-lg font-semibold text-primary">
                    {account.totalBalance.toFixed(8)} DCR
                  </p>
                  <p className="text-xs text-muted-foreground">Total Balance</p>
                </div>
              </div>

              <div className="grid grid-cols-3 gap-2 text-xs">
                <div className="p-2 rounded bg-background/30">
                  <p className="text-muted-foreground mb-1">Spendable</p>
                  <p className="font-medium text-foreground">{account.spendableBalance.toFixed(8)}</p>
                </div>
                <div className="p-2 rounded bg-background/30">
                  <p className="text-muted-foreground mb-1">Immature</p>
                  <p className="font-medium text-foreground">{account.immatureBalance.toFixed(8)}</p>
                </div>
                <div className="p-2 rounded bg-background/30">
                  <p className="text-muted-foreground mb-1">Unconfirmed</p>
                  <p className="font-medium text-foreground">{account.unconfirmedBalance.toFixed(8)}</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

