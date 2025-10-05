// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

import { Activity, Database, DollarSign, TrendingUp } from 'lucide-react';
import { MempoolInfo } from '../services/api';

interface MempoolActivityProps {
  data?: MempoolInfo;
}

export const MempoolActivity = ({ data }: MempoolActivityProps) => {
  const formatBytes = (bytes: number) => {
    if (bytes === 0) return 'N/A';
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
  };

  const formatDCR = (amount: number) => {
    if (amount === 0) return 'N/A';
    return amount.toFixed(4);
  };

  const formatFeeRate = (rate: number) => {
    if (rate === 0) return 'N/A';
    return rate.toFixed(6);
  };

  const mempoolItems = [
    {
      label: 'Pending Transactions',
      value: data ? `${data.txCount} txs` : 'Loading...',
      icon: Activity,
      color: 'text-primary',
    },
    {
      label: 'Mempool Size',
      value: data ? formatBytes(data.bytes) : 'Loading...',
      icon: Database,
      color: 'text-blue-500',
    },
    {
      label: 'Total Fees',
      value: data ? `${formatDCR(data.totalFee)} DCR` : 'Loading...',
      icon: DollarSign,
      color: 'text-success',
    },
    {
      label: 'Avg. Fee Rate',
      value: data ? `${formatFeeRate(data.averageFeeRate)} DCR/KB` : 'Loading...',
      icon: TrendingUp,
      color: 'text-purple-500',
    },
  ];

  return (
    <div className="p-6 rounded-xl bg-gradient-card backdrop-blur-sm border border-border/50 animate-fade-in">
      <h3 className="text-lg font-semibold mb-4">Mempool Activity</h3>
      <div className="space-y-4">
        {mempoolItems.map((item, index) => {
          const Icon = item.icon;
          return (
            <div key={index} className="flex items-center justify-between p-3 rounded-lg bg-muted/30 hover:bg-muted/50 transition-all duration-200">
              <div className="flex items-center gap-3">
                <div className={`p-2 rounded-lg bg-primary/10 ${item.color}`}>
                  <Icon className="w-4 h-4" />
                </div>
                <span className="text-sm text-muted-foreground">{item.label}</span>
              </div>
              <span className="font-medium">{item.value}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
};

