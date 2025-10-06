// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 25000, // 25 seconds to accommodate wallet rescans
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface NodeStatus {
  status: string;
  syncProgress: number;
  version: string;
  syncPhase: string;
  syncMessage: string;
}

export interface BlockchainInfo {
  blockHeight: number;
  blockHash: string;
  difficulty: number;
  chainSize: number;
  blockTime: string;
}

export interface NetworkInfo {
  peerCount: number;
  hashrate: string;
  networkHashPS: number;
}

export interface Peer {
  id: number;
  address: string;
  protocol: string;
  latency: string;
  connTime: string;
  traffic: string;
  version: string;
  isSyncNode: boolean;
}

export interface SupplyInfo {
  circulatingSupply: string;
  stakedSupply: string;
  stakedPercent: number;
  exchangeRate: string;
  treasurySize: string;
  mixedPercent: string;
}

export interface StakingInfo {
  ticketPrice: number;
  poolSize: number;
  lockedDCR: number;
  participationRate: number;
  allMempoolTix: number;
  immature: number;
  live: number;
  voted: number;
  missed: number;
  revoked: number;
}

export interface MempoolInfo {
  size: number;
  bytes: number;
  txCount: number;
  totalFee: number;
  averageFeeRate: number;
}

export interface DashboardData {
  nodeStatus: NodeStatus;
  blockchainInfo: BlockchainInfo;
  networkInfo: NetworkInfo;
  peers: Peer[];
  supplyInfo: SupplyInfo;
  stakingInfo: StakingInfo;
  mempoolInfo: MempoolInfo;
  lastUpdate: string;
}

export interface RPCConnectionRequest {
  host: string;
  port: string;
  username: string;
  password: string;
}

export interface RPCConnectionResponse {
  success: boolean;
  message: string;
}

// API functions
export const getDashboardData = async (): Promise<DashboardData> => {
  const response = await api.get<DashboardData>('/dashboard');
  return response.data;
};

export const getNodeStatus = async (): Promise<NodeStatus> => {
  const response = await api.get<NodeStatus>('/node/status');
  return response.data;
};

export const getBlockchainInfo = async (): Promise<BlockchainInfo> => {
  const response = await api.get<BlockchainInfo>('/blockchain/info');
  return response.data;
};

export const getPeers = async (): Promise<Peer[]> => {
  const response = await api.get<Peer[]>('/network/peers');
  return response.data;
};

export const connectRPC = async (config: RPCConnectionRequest): Promise<RPCConnectionResponse> => {
  const response = await api.post<RPCConnectionResponse>('/connect', config);
  return response.data;
};

export const checkHealth = async (): Promise<any> => {
  const response = await api.get('/health');
  return response.data;
};

// Wallet Types
export interface WalletStatus {
  status: string;
  syncProgress: number;
  syncHeight: number;
  bestBlockHash: string;
  version: string;
  unlocked: boolean;
  rescanInProgress: boolean;
  syncMessage: string;
}

export interface AccountInfo {
  accountName: string;
  totalBalance: number;
  spendableBalance: number;
  immatureBalance: number;
  unconfirmedBalance: number;
  accountNumber: number;
}

export interface WalletTransaction {
  txid: string;
  amount: number;
  fee: number;
  confirmations: number;
  time: string;
  type: string;
  comment: string;
}

export interface WalletAddress {
  address: string;
  account: string;
  used: boolean;
  path: string;
}

export interface WalletDashboardData {
  walletStatus: WalletStatus;
  accountInfo: AccountInfo;
  accounts: AccountInfo[];
  lastUpdate: string;
}

export interface ImportXpubRequest {
  xpub: string;
  accountName: string;
  rescan: boolean;
}

export interface ImportXpubResponse {
  success: boolean;
  message: string;
  accountNum?: number;
}

// Wallet API Functions
export const getWalletStatus = async (): Promise<WalletStatus> => {
  const response = await api.get<WalletStatus>('/wallet/status');
  return response.data;
};

export const getWalletDashboard = async (): Promise<WalletDashboardData> => {
  const response = await api.get<WalletDashboardData>('/wallet/dashboard');
  return response.data;
};

export const importXpub = async (xpub: string, accountName: string, rescan: boolean): Promise<ImportXpubResponse> => {
  const response = await api.post<ImportXpubResponse>('/wallet/importxpub', {
    xpub,
    accountName,
    rescan
  });
  return response.data;
};

export const triggerRescan = async (): Promise<any> => {
  const response = await api.post('/wallet/rescan', { beginHeight: 0 });
  return response.data;
};

export interface SyncProgressData {
  isRescanning: boolean;
  scanHeight: number;
  chainHeight: number;
  progress: number;
  message: string;
}

export const getSyncProgress = async (): Promise<SyncProgressData> => {
  const response = await api.get<SyncProgressData>('/wallet/sync-progress');
  return response.data;
};

export default api;

