// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

import { Blocks, Clock, HardDrive, Network } from 'lucide-react';
import { BlockchainInfo as BlockchainInfoType } from '../services/api';

interface BlockchainInfoProps {
  data?: BlockchainInfoType;
}

export const BlockchainInfo = ({ data }: BlockchainInfoProps) => {
  const infoItems = [
    { 
      label: 'Latest Block', 
      value: data?.blockHeight?.toLocaleString() || 'Loading...', 
      icon: Blocks 
    },
    { 
      label: 'Block Time', 
      value: data?.blockTime || 'Loading...', 
      icon: Clock 
    },
    { 
      label: 'Chain Size', 
      value: data?.chainSize ? `${(data.chainSize / 1e9).toFixed(1)} GB` : 'Loading...', 
      icon: HardDrive 
    },
    { 
      label: 'Difficulty', 
      value: data?.difficulty ? `${(data.difficulty / 1e12).toFixed(1)}T` : 'Loading...', 
      icon: Network 
    },
  ];

  return (
    <div className="p-6 rounded-xl bg-gradient-card backdrop-blur-sm border border-border/50 animate-fade-in">
      <h3 className="text-lg font-semibold mb-4">Blockchain Information</h3>
      <div className="space-y-4">
        {infoItems.map((item, index) => {
          const Icon = item.icon;
          return (
            <div key={item.label}>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="p-2 rounded-lg bg-primary/10 border border-primary/20">
                    <Icon className="h-4 w-4 text-primary" />
                  </div>
                  <span className="text-sm text-muted-foreground">{item.label}</span>
                </div>
                <span className="text-sm font-semibold">{item.value}</span>
              </div>
              {index < infoItems.length - 1 && (
                <div className="mt-4 h-px bg-border/50" />
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
};

