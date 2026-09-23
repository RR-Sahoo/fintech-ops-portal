export type TransactionStatus = "SUCCESS" | "PENDING" | "FAILED";
export type TransactionType = "DEBIT" | "CREDIT";

export interface Transaction {
  id: string;
  user_id: string;
  amount_cents: number;
  currency: string;
  status: TransactionStatus;
  reference_id: string;
  type: TransactionType;
  created_at: string;
  updated_at?: string;
  formatted_amount: string;
}

export interface TransactionSummary {
  total_volume_cents: number;
  total_volume_formatted: string;
  total_count: number;
  pending_count: number;
  success_count: number;
  failed_count: number;
  success_rate_percentage: number;
}

export interface TransactionsResponse {
  transactions: Transaction[];
  summary: TransactionSummary;
}

export interface ReconcileResponse {
  message: string;
  transaction: Transaction;
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

// Fallback Mock Transactions in USD for seamless Vercel Demo / Offline Mode
let fallbackMockTransactions: Transaction[] = [
  {
    id: "a1029384-b56c-48de-9f12-000000000001",
    user_id: "usr_fin_9821",
    amount_cents: 12500, // $125.00
    currency: "USD",
    status: "PENDING",
    reference_id: "ACH/428910293812/PAY_MERCHANT",
    type: "DEBIT",
    created_at: new Date(Date.now() - 2 * 3600 * 1000).toISOString(),
    formatted_amount: "$125.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000002",
    user_id: "usr_fin_4412",
    amount_cents: 450000, // $4,500.00
    currency: "USD",
    status: "SUCCESS",
    reference_id: "WIRE/W90281203810/SALARY_CREDIT",
    type: "CREDIT",
    created_at: new Date(Date.now() - 4 * 3600 * 1000).toISOString(),
    formatted_amount: "$4,500.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000003",
    user_id: "usr_fin_1092",
    amount_cents: 125000, // $1,250.00
    currency: "USD",
    status: "PENDING",
    reference_id: "STRIPE/TXN_881920384910",
    type: "DEBIT",
    created_at: new Date(Date.now() - 30 * 60 * 1000).toISOString(),
    formatted_amount: "$1,250.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000004",
    user_id: "usr_fin_3389",
    amount_cents: 7500, // $75.00
    currency: "USD",
    status: "FAILED",
    reference_id: "CARD/428919920192/BILL_UTILITY",
    type: "DEBIT",
    created_at: new Date(Date.now() - 60 * 60 * 1000).toISOString(),
    formatted_amount: "$75.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000005",
    user_id: "usr_fin_9821",
    amount_cents: 50000, // $500.00
    currency: "USD",
    status: "PENDING",
    reference_id: "ACH/2026092400192/VENDOR_PAY",
    type: "DEBIT",
    created_at: new Date(Date.now() - 15 * 60 * 1000).toISOString(),
    formatted_amount: "$500.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000006",
    user_id: "usr_fin_7718",
    amount_cents: 22000, // $220.00
    currency: "USD",
    status: "PENDING",
    reference_id: "REFUND/428901928374/REFUND_ORDER",
    type: "CREDIT",
    created_at: new Date(Date.now() - 5 * 3600 * 1000).toISOString(),
    formatted_amount: "$220.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000007",
    user_id: "usr_fin_5521",
    amount_cents: 320000, // $3,200.00
    currency: "USD",
    status: "PENDING",
    reference_id: "FEDWIRE/F202609240092/EQUITY_DEP",
    type: "CREDIT",
    created_at: new Date(Date.now() - 10 * 60 * 1000).toISOString(),
    formatted_amount: "$3,200.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000008",
    user_id: "usr_fin_4412",
    amount_cents: 4999, // $49.99
    currency: "USD",
    status: "FAILED",
    reference_id: "CARD/AUTH_8910283019",
    type: "DEBIT",
    created_at: new Date(Date.now() - 6 * 3600 * 1000).toISOString(),
    formatted_amount: "$49.99",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000009",
    user_id: "usr_fin_6634",
    amount_cents: 85000, // $850.00
    currency: "USD",
    status: "SUCCESS",
    reference_id: "ACH/428900192834/LOAN_PAYMENT",
    type: "DEBIT",
    created_at: new Date(Date.now() - 8 * 3600 * 1000).toISOString(),
    formatted_amount: "$850.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000010",
    user_id: "usr_fin_2209",
    amount_cents: 18750, // $187.50
    currency: "USD",
    status: "SUCCESS",
    reference_id: "CARD/428938491029/SAAS_SUB",
    type: "DEBIT",
    created_at: new Date(Date.now() - 12 * 3600 * 1000).toISOString(),
    formatted_amount: "$187.50",
  },
];

function getFallbackSummary(txs: Transaction[]): TransactionSummary {
  const totalVolumeCents = txs.reduce((acc, t) => acc + (t.amount_cents || 0), 0);
  const pendingCount = txs.filter((t) => t.status === "PENDING").length;
  const successCount = txs.filter((t) => t.status === "SUCCESS").length;
  const failedCount = txs.filter((t) => t.status === "FAILED").length;
  const successRate = txs.length > 0 ? (successCount / txs.length) * 100 : 0;

  // Format cents into USD formatted string (pure integer arithmetic)
  const dollars = Math.floor(totalVolumeCents / 100);
  const remainder = totalVolumeCents % 100;
  const formattedDollars = dollars.toLocaleString("en-US");
  const totalVolumeFormatted = `$${formattedDollars}.${String(remainder).padStart(2, "0")}`;

  return {
    total_volume_cents: totalVolumeCents,
    total_volume_formatted: totalVolumeFormatted,
    total_count: txs.length,
    pending_count: pendingCount,
    success_count: successCount,
    failed_count: failedCount,
    success_rate_percentage: successRate,
  };
}

/**
 * Fetch list of transactions with summary metrics.
 * Automatically falls back to mock dataset if Go backend is unreachable.
 */
export async function fetchTransactions(status?: string): Promise<TransactionsResponse> {
  try {
    const url = new URL(`${API_BASE_URL}/api/v1/transactions`);
    if (status && status !== "ALL") {
      url.searchParams.set("status", status);
    }

    const res = await fetch(url.toString(), {
      method: "GET",
      headers: {
        "Accept": "application/json",
      },
      cache: "no-store",
    });

    if (!res.ok) {
      throw new Error(`HTTP ${res.status}`);
    }

    return await res.json();
  } catch {
    // Graceful fallback to mock data (e.g. for Vercel deployment before backend is hosted)
    let filtered = [...fallbackMockTransactions];
    if (status && status !== "ALL") {
      filtered = filtered.filter((t) => t.status === status);
    }

    return {
      transactions: filtered,
      summary: getFallbackSummary(fallbackMockTransactions),
    };
  }
}

/**
 * Reconcile a pending transaction, flipping its status to SUCCESS.
 * Automatically falls back to local in-memory mutation if Go backend is unreachable.
 */
export async function reconcileTransaction(transactionId: string): Promise<ReconcileResponse> {
  try {
    const res = await fetch(`${API_BASE_URL}/api/v1/transactions/reconcile`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Accept": "application/json",
      },
      body: JSON.stringify({ transaction_id: transactionId }),
    });

    if (!res.ok) {
      const err = await res.json().catch(() => null);
      throw new Error(err?.error || `HTTP ${res.status}`);
    }

    return await res.json();
  } catch {
    // Local fallback reconciliation
    const targetIdx = fallbackMockTransactions.findIndex((t) => t.id === transactionId);
    if (targetIdx === -1) {
      throw new Error("Transaction not found");
    }

    if (fallbackMockTransactions[targetIdx].status !== "PENDING") {
      throw new Error("Only PENDING transactions can be reconciled");
    }

    fallbackMockTransactions[targetIdx] = {
      ...fallbackMockTransactions[targetIdx],
      status: "SUCCESS",
      updated_at: new Date().toISOString(),
    };

    return {
      message: "Transaction successfully reconciled (Demo mode)",
      transaction: fallbackMockTransactions[targetIdx],
    };
  }
}
