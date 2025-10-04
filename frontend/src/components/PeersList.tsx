// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

import { Wifi, ArrowUpDown } from 'lucide-react';
import { Peer } from '../services/api';

interface PeersListProps {
  peers?: Peer[];
}

export const PeersList = ({ peers = [] }: PeersListProps) => {
  return (
    <div className="p-6 rounded-xl bg-gradient-card backdrop-blur-sm border border-border/50 animate-fade-in">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold">Connected Peers</h3>
        <span className="px-3 py-1 rounded-md text-sm font-medium bg-primary/10 text-primary border border-primary/20">
          {peers.length} Active
        </span>
      </div>
      
      <div className="max-h-[300px] overflow-y-auto">
        <div className="space-y-3">
          {peers.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              No peers connected
            </div>
          ) : (
            peers.map((peer) => (
              <div
                key={peer.id}
                className="flex items-center justify-between p-3 rounded-lg bg-muted/30 border border-border/30 hover:border-primary/30 transition-all"
              >
                <div className="flex items-center gap-3">
                  <div className="p-2 rounded-lg bg-success/10 border border-success/20">
                    <Wifi className="h-4 w-4 text-success" />
                  </div>
                  <div>
                    <p className="text-sm font-medium">{peer.address}</p>
                    <p className="text-xs text-muted-foreground">{peer.protocol}</p>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <ArrowUpDown className="h-3 w-3 text-muted-foreground" />
                  <span className="text-xs text-muted-foreground">{peer.latency}</span>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
};

