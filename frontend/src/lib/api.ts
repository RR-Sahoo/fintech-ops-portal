export type TransactionStatus = "SUCCESS" | "PENDING" | "FAILED";
export type TransactionType = "DEBIT" | "CREDIT";

export interface Transaction {
  id: string;
  user_id: string;
  amount_paise: number;
  currency: string;
  status: TransactionStatus;
  reference_id: string;
  type: TransactionType;
  created_at: string;
  updated_at?: string;
  formatted_amount: string;
}

export interface TransactionSummary {
  total_volume_paise: number;
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

// Fallback Mock Transactions for seamless Vercel Demo / Offline Mode
let fallbackMockTransactions: Transaction[] = [
  {
    id: "a1029384-b56c-48de-9f12-000000000001",
    user_id: "usr_fin_9821",
    amount_paise: 125050,
    currency: "INR",
    status: "PENDING",
    reference_id: "UPI/428910293812/PAY_MERCHANT",
    type: "DEBIT",
    created_at: new Date(Date.now() - 2 * 3600 * 1000).toISOString(),
    formatted_amount: "₹1,250.50",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000002",
    user_id: "usr_fin_4412",
    amount_paise: 5000000,
    currency: "INR",
    status: "SUCCESS",
    reference_id: "NEFT/N90281203810/SALARY_CREDIT",
    type: "CREDIT",
    created_at: new Date(Date.now() - 4 * 3600 * 1000).toISOString(),
    formatted_amount: "₹50,000.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000003",
    user_id: "usr_fin_1092",
    amount_paise: 349900,
    currency: "INR",
    status: "PENDING",
    reference_id: "PG/TXN_881920384910",
    type: "DEBIT",
    created_at: new Date(Date.now() - 30 * 60 * 1000).toISOString(),
    formatted_amount: "₹3,499.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000004",
    user_id: "usr_fin_3389",
    amount_paise: 75000,
    currency: "INR",
    status: "FAILED",
    reference_id: "UPI/428919920192/BILL_UTILITY",
    type: "DEBIT",
    created_at: new Date(Date.now() - 60 * 60 * 1000).toISOString(),
    formatted_amount: "₹750.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000005",
    user_id: "usr_fin_9821",
    amount_paise: 1500000,
    currency: "INR",
    status: "PENDING",
    reference_id: "IMPS/2026092400192/VENDOR_PAY",
    type: "DEBIT",
    created_at: new Date(Date.now() - 15 * 60 * 1000).toISOString(),
    formatted_amount: "₹15,000.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000006",
    user_id: "usr_fin_7718",
    amount_paise: 220000,
    currency: "INR",
    status: "PENDING",
    reference_id: "UPI/428901928374/REFUND_ORDER",
    type: "CREDIT",
    created_at: new Date(Date.now() - 5 * 3600 * 1000).toISOString(),
    formatted_amount: "₹2,200.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000007",
    user_id: "usr_fin_5521",
    amount_paise: 12500000,
    currency: "INR",
    status: "PENDING",
    reference_id: "RTGS/R202609240092/EQUITY_DEP",
    type: "CREDIT",
    created_at: new Date(Date.now() - 10 * 60 * 1000).toISOString(),
    formatted_amount: "₹1,25,000.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000008",
    user_id: "usr_fin_4412",
    amount_paise: 49900,
    currency: "INR",
    status: "FAILED",
    reference_id: "CARD/AUTH_8910283019",
    type: "DEBIT",
    created_at: new Date(Date.now() - 6 * 3600 * 1000).toISOString(),
    formatted_amount: "₹499.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000009",
    user_id: "usr_fin_6634",
    amount_paise: 850000,
    currency: "INR",
    status: "SUCCESS",
    reference_id: "UPI/428900192834/LOAN_EMI",
    type: "DEBIT",
    created_at: new Date(Date.now() - 8 * 3600 * 1000).toISOString(),
    formatted_amount: "₹8,500.00",
  },
  {
    id: "a1029384-b56c-48de-9f12-000000000010",
    user_id: "usr_fin_2209",
    amount_paise: 187525,
    currency: "INR",
    status: "SUCCESS",
    reference_id: "UPI/428938491029/GROCERY_EXP",
    type: "DEBIT",
    created_at: new Date(Date.now() - 12 * 3600 * 1000).toISOString(),
    formatted_amount: "₹1,875.25",
  },
];

function getFallbackSummary(txs: Transaction[]): TransactionSummary {
  const totalVolumePaise = txs.reduce((acc, t) => acc + t.amount_paise, 0);
  const pendingCount = txs.filter((t) => t.status === "PENDING").length;
  const successCount = txs.filter((t) => t.status === "SUCCESS").length;
  const failedCount = txs.filter((t) => t.status === "FAILED").length;
  const successRate = txs.length > 0 ? (successCount / txs.length) * 100 : 0;

  // Format paise into standard Indian rupee string (integer arithmetic)
  const rupees = Math.floor(totalVolumePaise / 100);
  const remainder = totalVolumePaise % 100;
  const formattedRupees = rupees.toLocaleString("en-IN");
  const totalVolumeFormatted = `₹${formattedRupees}.${String(remainder).padStart(2, "0")}`;

  return {
    total_volume_paise: totalVolumePaise,
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
