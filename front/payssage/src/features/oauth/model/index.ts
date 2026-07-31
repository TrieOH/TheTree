interface ConnectRequestBase {
  provider_redirect_url: string;
  final_redirect_url: string;
}

export interface CollectorConnectRequest extends ConnectRequestBase {
  flow: "collector";
  organization_id?: string;
  wallet_id?: never;
}

export interface SellerConnectRequest extends ConnectRequestBase {
  flow: "seller";
  wallet_id: string;
  organization_id?: never;
}

export type ConnectRequest = CollectorConnectRequest | SellerConnectRequest;

export interface RevokeRequest {
  flow: "collector" | "seller";
  id: string;
}

export interface Collector {
  id: string;
  owner_id: string;
  organization_id: string | null;
  provider: string;
  provider_user_id: string;
  created_at: string;
  revoked_at: string | null;
}

export interface Seller {
  id: string;
  wallet_id: string;
  provider: string;
  provider_user_id: string;
  created_at: string;
  revoked_at: string | null;
}
