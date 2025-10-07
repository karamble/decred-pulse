// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

import { Link, useLocation } from 'react-router-dom';
import { Wallet, Compass } from 'lucide-react';

interface HeaderProps {
  nodeVersion?: string;
}

export const Header = ({ nodeVersion }: HeaderProps) => {
  const location = useLocation();
  const isWalletPage = location.pathname === '/wallet';
  const isExplorerPage = location.pathname.startsWith('/explorer');
  const isNodePage = location.pathname === '/';

  const getPageTitle = () => {
    if (isWalletPage) return 'Wallet';
    if (isExplorerPage) return 'Explorer';
    return 'Node';
  };

  const getPageDescription = () => {
    if (isWalletPage) return 'Monitor your watch-only wallet and transactions';
    if (isExplorerPage) return 'Search and explore the Decred blockchain';
    return 'Monitor your dcrd node performance and network status';
  };

  return (
    <div className="flex items-center justify-between mb-8 animate-fade-in">
      <div className="flex items-center gap-4">
        <div className="h-16 w-16 rounded-xl animate-pulse-glow flex items-center justify-center text-2xl font-bold bg-gradient-primary shadow-glow-primary">
          DCR
        </div>
        <div>
          <h1 className="text-4xl font-bold bg-gradient-primary bg-clip-text text-transparent mb-2">
            Decred {getPageTitle()}
          </h1>
          <p className="text-muted-foreground">
            {getPageDescription()}
          </p>
        </div>
      </div>
      <div className="flex items-center gap-4">
        {/* Navigation Buttons */}
        <Link
          to="/"
          className={`px-4 py-3 rounded-lg border transition-all duration-300 flex items-center gap-2 ${
            isNodePage
              ? 'bg-primary/20 border-primary/40 shadow-glow-primary'
              : 'bg-primary/10 border-primary/20 hover:bg-primary/20 hover:shadow-glow-primary'
          }`}
        >
          <div className="h-10 w-10 rounded-lg flex items-center justify-center font-bold text-sm bg-gradient-primary">
            DCR
          </div>
          <span className="text-primary font-semibold">Node</span>
        </Link>

        <Link
          to="/wallet"
          className={`px-4 py-3 rounded-lg border transition-all duration-300 flex items-center gap-2 ${
            isWalletPage
              ? 'bg-primary/20 border-primary/40 shadow-glow-primary'
              : 'bg-primary/10 border-primary/20 hover:bg-primary/20 hover:shadow-glow-primary'
          }`}
        >
          <div className="h-10 w-10 rounded-lg flex items-center justify-center bg-gradient-primary">
            <Wallet className="h-5 w-5 text-white" />
          </div>
          <span className="text-primary font-semibold">Wallet</span>
        </Link>

        <Link
          to="/explorer"
          className={`px-4 py-3 rounded-lg border transition-all duration-300 flex items-center gap-2 ${
            isExplorerPage
              ? 'bg-primary/20 border-primary/40 shadow-glow-primary'
              : 'bg-primary/10 border-primary/20 hover:bg-primary/20 hover:shadow-glow-primary'
          }`}
        >
          <div className="h-10 w-10 rounded-lg flex items-center justify-center bg-gradient-primary">
            <Compass className="h-5 w-5 text-white" />
          </div>
          <span className="text-primary font-semibold">Explorer</span>
        </Link>

        <a
          href="https://nodes.jholdstock.uk/"
          target="_blank"
          rel="noopener noreferrer"
          className="px-4 py-3 rounded-lg bg-primary/10 border border-primary/20 hover:bg-primary/20 transition-all duration-300 hover:shadow-glow-primary flex items-center gap-2"
        >
          <div className="h-10 w-10 rounded-lg flex items-center justify-center font-bold text-sm bg-gradient-primary">
            DCR
          </div>
          <span className="text-primary font-semibold">Node Mapper</span>
        </a>

        {nodeVersion && (
          <div className="px-4 py-2 rounded-lg bg-primary/10 border border-primary/20">
            <p className="text-sm text-muted-foreground">Version</p>
            <p className="text-lg font-semibold text-primary">{nodeVersion}</p>
          </div>
        )}
      </div>
    </div>
  );
};

