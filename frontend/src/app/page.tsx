"use client";

import React, { useState, useMemo } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  fetchTransactions,
  reconcileTransaction,
  type Transaction,
} from "@/lib/api";
import {
  DollarSign,
  AlertTriangle,
  CheckCircle2,
  Clock,
  ArrowUpRight,
  ArrowDownLeft,
  Search,
  RefreshCw,
  Copy,
  Check,
  Activity,
  Layers,
  Sparkles,
} from "lucide-react";

export default function DashboardPage() {
  const queryClient = useQueryClient();
  const [selectedStatus, setSelectedStatus] = useState<string>("ALL");
  const [searchQuery, setSearchQuery] = useState<string>("");
  const [copiedId, setCopiedId] = useState<string | null>(null);
  const [notification, setNotification] = useState<{
    type: "success" | "error";
    message: string;
  } | null>(null);

  // Fetch transactions query
  const {
    data,
    isLoading,
    isError,
    error,
    isFetching,
    refetch,
  } = useQuery({
    queryKey: ["transactions", selectedStatus],
    queryFn: () => fetchTransactions(selectedStatus),
  });

  // Reconcile mutation
  const reconcileMutation = useMutation({
    mutationFn: (id: string) => reconcileTransaction(id),
    onSuccess: (res) => {
      setNotification({
        type: "success",
        message: `Transaction ${res.transaction.id.slice(0, 8)}... successfully reconciled to SUCCESS!`,
      });
      // Invalidate queries to auto-refresh data
      queryClient.invalidateQueries({ queryKey: ["transactions"] });
      setTimeout(() => setNotification(null), 5000);
    },
    onError: (err: Error) => {
      setNotification({
        type: "error",
        message: err.message || "Failed to reconcile transaction",
      });
      setTimeout(() => setNotification(null), 5000);
    },
  });

  // Copy helper
  const handleCopy = (id: string) => {
    navigator.clipboard.writeText(id);
    setCopiedId(id);
    setTimeout(() => setCopiedId(null), 2000);
  };

  // Client-side search filtering
  const filteredTransactions = useMemo(() => {
    if (!data?.transactions) return [];
    if (!searchQuery.trim()) return data.transactions;

    const q = searchQuery.toLowerCase().trim();
    return data.transactions.filter(
      (tx) =>
        tx.id.toLowerCase().includes(q) ||
        tx.user_id.toLowerCase().includes(q) ||
        tx.reference_id.toLowerCase().includes(q)
    );
  }, [data?.transactions, searchQuery]);

  const summary = data?.summary;

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 font-sans">
      {/* Background Decorative Gradients */}
      <div className="fixed inset-0 pointer-events-none overflow-hidden opacity-30">
        <div className="absolute -top-40 left-1/4 w-96 h-96 bg-indigo-600/30 rounded-full blur-3xl" />
        <div className="absolute top-1/3 right-10 w-96 h-96 bg-emerald-600/20 rounded-full blur-3xl" />
      </div>

      <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-8">
        {/* Header */}
        <header className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 border-b border-slate-800/80 pb-6">
          <div className="space-y-1">
            <div className="flex items-center gap-2.5">
              <div className="p-2 bg-gradient-to-tr from-indigo-600 to-violet-500 rounded-xl shadow-lg shadow-indigo-500/20">
                <Layers className="w-6 h-6 text-white" />
              </div>
              <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-white">
                Fintech Ops & Reconciliation Portal
              </h1>
            </div>
            <p className="text-sm text-slate-400 pl-11">
              Real-time payment settlement monitoring, liquidity metrics, and instant transaction reconciliation.
            </p>
          </div>

          <div className="flex items-center gap-3 pl-11 md:pl-0">
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-slate-900 border border-slate-800 text-xs text-slate-300">
              <span className="relative flex h-2 w-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
              </span>
              <span>Ledger Active</span>
            </div>

            <button
              onClick={() => refetch()}
              disabled={isFetching}
              className="inline-flex items-center gap-2 px-3.5 py-2 text-xs font-medium rounded-lg bg-slate-800/80 hover:bg-slate-700 text-slate-200 border border-slate-700/60 transition-colors shadow-sm disabled:opacity-50 cursor-pointer"
              title="Refresh Transactions"
            >
              <RefreshCw className={`w-3.5 h-3.5 ${isFetching ? "animate-spin" : ""}`} />
              <span>Refresh</span>
            </button>
          </div>
        </header>

        {/* Toast / Notification Banner */}
        {notification && (
          <div
            className={`p-4 rounded-xl border flex items-center justify-between transition-all shadow-lg backdrop-blur-md ${
              notification.type === "success"
                ? "bg-emerald-950/60 border-emerald-500/40 text-emerald-200"
                : "bg-rose-950/60 border-rose-500/40 text-rose-200"
            }`}
          >
            <div className="flex items-center gap-3">
              {notification.type === "success" ? (
                <CheckCircle2 className="w-5 h-5 text-emerald-400 shrink-0" />
              ) : (
                <AlertTriangle className="w-5 h-5 text-rose-400 shrink-0" />
              )}
              <p className="text-sm font-medium">{notification.message}</p>
            </div>
            <button
              onClick={() => setNotification(null)}
              className="text-xs opacity-75 hover:opacity-100 underline ml-4 cursor-pointer"
            >
              Dismiss
            </button>
          </div>
        )}

        {/* Top Metrics Cards */}
        <section className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
          {/* Card 1: Total Volume */}
          <div className="relative overflow-hidden rounded-2xl bg-gradient-to-b from-slate-900/90 to-slate-900/50 border border-slate-800 p-5 shadow-sm hover:border-slate-700 transition-all">
            <div className="flex items-center justify-between">
              <p className="text-xs font-semibold uppercase tracking-wider text-slate-400">
                Total Gross Volume
              </p>
              <div className="p-2 bg-indigo-500/10 rounded-lg text-indigo-400 border border-indigo-500/20">
                <DollarSign className="w-4 h-4" />
              </div>
            </div>
            <div className="mt-4">
              <div className="text-3xl font-extrabold tracking-tight text-white">
                {summary?.total_volume_formatted || "$0.00"}
              </div>
              <p className="text-xs text-slate-400 mt-1">
                {(summary?.total_volume_cents ?? 0).toLocaleString("en-US")} cents in total throughput
              </p>
            </div>
          </div>

          {/* Card 2: Pending Settlements */}
          <div className="relative overflow-hidden rounded-2xl bg-gradient-to-b from-slate-900/90 to-slate-900/50 border border-slate-800 p-5 shadow-sm hover:border-slate-700 transition-all">
            <div className="flex items-center justify-between">
              <p className="text-xs font-semibold uppercase tracking-wider text-slate-400">
                Pending Settlements
              </p>
              <div className="p-2 bg-amber-500/10 rounded-lg text-amber-400 border border-amber-500/20">
                <Clock className="w-4 h-4" />
              </div>
            </div>
            <div className="mt-4 flex items-baseline justify-between">
              <div className="text-3xl font-extrabold tracking-tight text-white">
                {summary?.pending_count ?? 0}
              </div>
              {(summary?.pending_count ?? 0) > 0 ? (
                <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-amber-500/10 border border-amber-500/30 text-amber-400">
                  <AlertTriangle className="w-3 h-3" />
                  Action Required
                </span>
              ) : (
                <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-500/10 border border-emerald-500/30 text-emerald-400">
                  <CheckCircle2 className="w-3 h-3" />
                  Settled
                </span>
              )}
            </div>
            <p className="text-xs text-slate-400 mt-1">
              Awaiting ledger clearance or manual reconciliation
            </p>
          </div>

          {/* Card 3: Settlement Success Rate */}
          <div className="relative overflow-hidden rounded-2xl bg-gradient-to-b from-slate-900/90 to-slate-900/50 border border-slate-800 p-5 shadow-sm hover:border-slate-700 transition-all">
            <div className="flex items-center justify-between">
              <p className="text-xs font-semibold uppercase tracking-wider text-slate-400">
                Settlement Success Rate
              </p>
              <div className="p-2 bg-emerald-500/10 rounded-lg text-emerald-400 border border-emerald-500/20">
                <Activity className="w-4 h-4" />
              </div>
            </div>
            <div className="mt-4 flex items-baseline justify-between">
              <div className="text-3xl font-extrabold tracking-tight text-white">
                {summary?.success_rate_percentage !== undefined
                  ? `${summary.success_rate_percentage.toFixed(1)}%`
                  : "0.0%"}
              </div>
              <span className="text-xs text-slate-400">
                {summary?.success_count ?? 0} of {summary?.total_count ?? 0} cleared
              </span>
            </div>
            {/* Progress bar */}
            <div className="mt-3 w-full bg-slate-800 rounded-full h-1.5 overflow-hidden">
              <div
                className="bg-emerald-500 h-full rounded-full transition-all duration-500"
                style={{
                  width: `${Math.min(100, summary?.success_rate_percentage ?? 0)}%`,
                }}
              />
            </div>
          </div>
        </section>

        {/* Controls: Filter Tabs + Search */}
        <section className="bg-slate-900/70 border border-slate-800 rounded-2xl p-4 flex flex-col md:flex-row md:items-center justify-between gap-4 backdrop-blur-md">
          {/* Status Tabs */}
          <div className="flex flex-wrap items-center gap-1.5 bg-slate-950/80 p-1 rounded-xl border border-slate-800/80">
            {(["ALL", "SUCCESS", "PENDING", "FAILED"] as const).map((status) => {
              const isActive = selectedStatus === status;
              return (
                <button
                  key={status}
                  onClick={() => setSelectedStatus(status)}
                  className={`px-3.5 py-1.5 rounded-lg text-xs font-medium transition-all cursor-pointer ${
                    isActive
                      ? "bg-indigo-600 text-white shadow-md shadow-indigo-600/30"
                      : "text-slate-400 hover:text-slate-200 hover:bg-slate-800/60"
                  }`}
                >
                  {status === "ALL" ? "All Transactions" : status}
                </button>
              );
            })}
          </div>

          {/* Search Input */}
          <div className="relative w-full md:w-80">
            <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
            <input
              type="text"
              placeholder="Search ID, User, or Ref..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-9 pr-4 py-2 bg-slate-950/80 border border-slate-800 rounded-xl text-xs text-slate-200 placeholder:text-slate-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all"
            />
            {searchQuery && (
              <button
                onClick={() => setSearchQuery("")}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-slate-400 hover:text-slate-200 cursor-pointer"
              >
                Clear
              </button>
            )}
          </div>
        </section>

        {/* Transactions Table Section */}
        <section className="bg-slate-900/70 border border-slate-800 rounded-2xl overflow-hidden shadow-sm backdrop-blur-md">
          {isLoading ? (
            <div className="p-12 text-center space-y-3">
              <RefreshCw className="w-6 h-6 animate-spin mx-auto text-indigo-400" />
              <p className="text-sm text-slate-400 font-medium">Fetching transaction ledger records...</p>
            </div>
          ) : isError ? (
            <div className="p-12 text-center space-y-4">
              <div className="p-3 bg-rose-500/10 rounded-full w-12 h-12 flex items-center justify-center mx-auto text-rose-400 border border-rose-500/20">
                <AlertTriangle className="w-6 h-6" />
              </div>
              <div className="space-y-1">
                <p className="text-sm font-semibold text-rose-300">Failed to load transactions</p>
                <p className="text-xs text-slate-400">
                  {error instanceof Error ? error.message : "Ensure the Go backend is running on port 8080."}
                </p>
              </div>
              <button
                onClick={() => refetch()}
                className="px-4 py-2 text-xs font-semibold rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 cursor-pointer"
              >
                Retry Connection
              </button>
            </div>
          ) : filteredTransactions.length === 0 ? (
            <div className="p-12 text-center space-y-3">
              <Sparkles className="w-8 h-8 mx-auto text-slate-600" />
              <p className="text-sm font-semibold text-slate-300">No transactions match your criteria</p>
              <p className="text-xs text-slate-500">
                Try switching status filters or clearing the search query.
              </p>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-xs text-slate-300">
                <thead className="bg-slate-950/60 border-b border-slate-800/80 text-[11px] font-semibold uppercase tracking-wider text-slate-400">
                  <tr>
                    <th className="px-5 py-3.5">Transaction ID</th>
                    <th className="px-5 py-3.5">User ID</th>
                    <th className="px-5 py-3.5">Amount</th>
                    <th className="px-5 py-3.5">Type</th>
                    <th className="px-5 py-3.5">Status</th>
                    <th className="px-5 py-3.5">Timestamp</th>
                    <th className="px-5 py-3.5 text-right">Action</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-800/60">
                  {filteredTransactions.map((tx: Transaction) => (
                    <tr
                      key={tx.id}
                      className="hover:bg-slate-800/30 transition-colors group"
                    >
                      {/* Transaction ID & Reference */}
                      <td className="px-5 py-4">
                        <div className="space-y-0.5">
                          <div className="flex items-center gap-1.5 font-mono text-slate-200 font-medium">
                            <span>{tx.id.slice(0, 18)}...</span>
                            <button
                              onClick={() => handleCopy(tx.id)}
                              className="opacity-0 group-hover:opacity-100 transition-opacity p-1 text-slate-400 hover:text-slate-200 rounded cursor-pointer"
                              title="Copy full UUID"
                            >
                              {copiedId === tx.id ? (
                                <Check className="w-3 h-3 text-emerald-400" />
                              ) : (
                                <Copy className="w-3 h-3" />
                              )}
                            </button>
                          </div>
                          <div className="text-[11px] text-slate-500 font-mono truncate max-w-[200px]" title={tx.reference_id}>
                            {tx.reference_id}
                          </div>
                        </div>
                      </td>

                      {/* User ID */}
                      <td className="px-5 py-4 font-mono text-slate-300">
                        {tx.user_id}
                      </td>

                      {/* Amount */}
                      <td className="px-5 py-4">
                        <div className="font-semibold text-sm text-slate-100">
                          {tx.formatted_amount}
                        </div>
                        <div className="text-[10px] text-slate-500">
                          {(tx.amount_cents ?? 0).toLocaleString("en-US")} cents
                        </div>
                      </td>

                      {/* Type Badge */}
                      <td className="px-5 py-4">
                        {tx.type === "CREDIT" ? (
                          <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-md text-[11px] font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                            <ArrowDownLeft className="w-3 h-3" />
                            CREDIT
                          </span>
                        ) : (
                          <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-md text-[11px] font-semibold bg-sky-500/10 text-sky-400 border border-sky-500/20">
                            <ArrowUpRight className="w-3 h-3" />
                            DEBIT
                          </span>
                        )}
                      </td>

                      {/* Status Pill */}
                      <td className="px-5 py-4">
                        {tx.status === "SUCCESS" && (
                          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-emerald-500/15 text-emerald-400 border border-emerald-500/30">
                            <span className="w-1.5 h-1.5 rounded-full bg-emerald-400" />
                            SUCCESS
                          </span>
                        )}
                        {tx.status === "PENDING" && (
                          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-amber-500/15 text-amber-400 border border-amber-500/30">
                            <span className="w-1.5 h-1.5 rounded-full bg-amber-400 animate-pulse" />
                            PENDING
                          </span>
                        )}
                        {tx.status === "FAILED" && (
                          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-rose-500/15 text-rose-400 border border-rose-500/30">
                            <span className="w-1.5 h-1.5 rounded-full bg-rose-400" />
                            FAILED
                          </span>
                        )}
                      </td>

                      {/* Timestamp */}
                      <td className="px-5 py-4 text-slate-400 whitespace-nowrap">
                        <div>
                          {new Date(tx.created_at).toLocaleDateString("en-US", {
                            day: "2-digit",
                            month: "short",
                            year: "numeric",
                          })}
                        </div>
                        <div className="text-[10px] text-slate-500">
                          {new Date(tx.created_at).toLocaleTimeString("en-US", {
                            hour: "2-digit",
                            minute: "2-digit",
                            second: "2-digit",
                          })}
                        </div>
                      </td>

                      {/* Action */}
                      <td className="px-5 py-4 text-right">
                        {tx.status === "PENDING" ? (
                          <button
                            onClick={() => reconcileMutation.mutate(tx.id)}
                            disabled={reconcileMutation.isPending && reconcileMutation.variables === tx.id}
                            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 shadow-md shadow-amber-500/20 hover:shadow-amber-500/30 transition-all disabled:opacity-50 cursor-pointer"
                          >
                            {reconcileMutation.isPending && reconcileMutation.variables === tx.id ? (
                              <>
                                <RefreshCw className="w-3 h-3 animate-spin" />
                                <span>Reconciling...</span>
                              </>
                            ) : (
                              <>
                                <CheckCircle2 className="w-3.5 h-3.5" />
                                <span>Reconcile Now</span>
                              </>
                            )}
                          </button>
                        ) : (
                          <span className="text-[11px] text-slate-500 font-mono italic">
                            Completed
                          </span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      </div>
    </div>
  );
}
