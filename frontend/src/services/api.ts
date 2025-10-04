// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface NodeStatus {
  status: string;
  syncProgress: number;
  version: string;
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
}

export interface SupplyInfo {
  circulatingSupply: string;
  stakedSupply: string;
  stakedPercent: number;
  exchangeRate: string;
  treasurySize: string;
  mixedPercent: string;
}

export interface DashboardData {
  nodeStatus: NodeStatus;
  blockchainInfo: BlockchainInfo;
  networkInfo: NetworkInfo;
  peers: Peer[];
  supplyInfo: SupplyInfo;
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

export default api;

