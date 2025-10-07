// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft, AlertCircle, Wallet, TrendingUp, Clock, ArrowUpDown } from 'lucide-react';
import { CopyButton } from '../components/explorer/CopyButton';

export const AddressView = () => {
  const { address } = useParams<{ address: string }>();
  const navigate = useNavigate();

  return (
    <div className="space-y-6">
      {/* Navigation */}
      <button
        onClick={() => navigate('/explorer')}
        className="flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to Explorer
      </button>

      {/* Address Header */}
      <div className="p-6 rounded-lg bg-gradient-to-br from-background to-muted/20 border border-border">
        <div className="flex items-start justify-between mb-4">
          <div>
            <h1 className="text-3xl font-bold mb-2">Address</h1>
            <div className="flex items-center gap-2">
              <code className="text-sm sm:text-base font-mono bg-muted/50 px-3 py-1.5 rounded break-all">
                {address}
              </code>
              <CopyButton text={address || ''} />
            </div>
          </div>
        </div>
      </div>

      {/* Not Implemented Notice */}
      <div className="p-8 rounded-lg bg-gradient-to-br from-yellow-500/10 to-orange-500/5 border border-yellow-500/20">
        <div className="flex items-start gap-4">
          <AlertCircle className="h-8 w-8 text-warning flex-shrink-0 mt-1" />
          <div className="flex-1">
            <h2 className="text-xl font-semibold mb-3 text-warning">
              Address Lookup Not Yet Implemented
            </h2>
            <p className="text-muted-foreground mb-6">
              Address lookup requires <code className="px-1.5 py-0.5 bg-muted/50 rounded text-sm">addrindex</code> to 
              be enabled on your <code className="px-1.5 py-0.5 bg-muted/50 rounded text-sm">dcrd</code> node. 
              This feature is currently under development.
            </p>

            <div className="space-y-4">
              <h3 className="font-semibold text-foreground">Features coming soon:</h3>
              
              <div className="grid gap-3">
                <div className="flex items-start gap-3 p-3 rounded-lg bg-background/50">
                  <Wallet className="h-5 w-5 text-primary flex-shrink-0 mt-0.5" />
                  <div>
                    <p className="font-medium">Balance Information</p>
                    <p className="text-sm text-muted-foreground">View total balance, spendable, and locked amounts</p>
                  </div>
                </div>

                <div className="flex items-start gap-3 p-3 rounded-lg bg-background/50">
                  <ArrowUpDown className="h-5 w-5 text-success flex-shrink-0 mt-0.5" />
                  <div>
                    <p className="font-medium">Transaction History</p>
                    <p className="text-sm text-muted-foreground">Complete list of all transactions involving this address</p>
                  </div>
                </div>

                <div className="flex items-start gap-3 p-3 rounded-lg bg-background/50">
                  <TrendingUp className="h-5 w-5 text-warning flex-shrink-0 mt-0.5" />
                  <div>
                    <p className="font-medium">Staking Activity</p>
                    <p className="text-sm text-muted-foreground">Ticket purchases, votes, and revocations</p>
                  </div>
                </div>

                <div className="flex items-start gap-3 p-3 rounded-lg bg-background/50">
                  <Clock className="h-5 w-5 text-blue-500 flex-shrink-0 mt-0.5" />
                  <div>
                    <p className="font-medium">UTXO Set</p>
                    <p className="text-sm text-muted-foreground">Unspent transaction outputs available for spending</p>
                  </div>
                </div>
              </div>

              <div className="mt-6 p-4 rounded-lg bg-muted/30 border border-border/50">
                <p className="text-sm text-muted-foreground">
                  <strong className="text-foreground">Note:</strong> To enable address indexing, restart your dcrd node with the 
                  <code className="mx-1 px-1.5 py-0.5 bg-background rounded text-xs">--addrindex</code> flag. 
                  The first startup will require a full blockchain reindex (may take several hours).
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

