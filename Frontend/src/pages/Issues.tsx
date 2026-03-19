import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { BugPlay, ChevronLeft, ChevronRight, Filter, X, CheckCircle2, EyeOff, AlertTriangle, Clock, RefreshCw, BarChart3, AlertCircle } from 'lucide-react';
import { issuesApi } from '../shared/lib/api';
import SelectField from '../shared/components/SelectField';
import { categoryLabel } from '../shared/lib/utils';


interface Issue {
  fingerprint: string;
  source_id: string;
  status: string;
  category: string;
  level: string;
  message_sample: string;
  occurrence_count: number;
  first_seen_at: string;
  last_seen_at: string;
}

interface IssueLog {
  id: number;
  message: string;
  level: string;
  created_at: string;
  ip_address: string;
}

interface ToastState {
  message: string;
  type: 'success' | 'error';
}

const LEVEL_STYLES: Record<string, string> = {
  CRITICAL: 'bg-rose-500/10 text-rose-400 border-rose-500/20',
  ERROR: 'bg-red-500/10 text-red-400 border-red-500/20',
  WARN: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
  INFO: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
  DEBUG: 'bg-zinc-500/10 text-zinc-400 border-zinc-500/20',
};

const STATUS_STYLES: Record<string, string> = {
  OPEN: 'bg-rose-500/10 text-rose-400 border-rose-500/20',
  RESOLVED: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
  IGNORED: 'bg-zinc-500/10 text-zinc-500 border-zinc-500/20',
};

function timeAgo(dateStr: string): string {
  const diff = Date.now() - new Date(dateStr).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 60) return `${mins}m ago`;
  const hrs = Math.floor(mins / 60);
  if (hrs < 24) return `${hrs}h ago`;
  return `${Math.floor(hrs / 24)}d ago`;
}

export default function Issues() {
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState<'issues' | 'analytics'>('issues');
  const [currentPage, setCurrentPage] = useState(1);
  const [limit, setLimit] = useState(20);
  const [statusFilter, setStatusFilter] = useState('');
  const [selectedIssue, setSelectedIssue] = useState<Issue | null>(null);
  const [toast, setToast] = useState<ToastState | null>(null);

  const showToast = (message: string, type: 'success' | 'error' = 'success') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 3000);
  };

  // Fetch issues with pagination
  const {
    data: issuesData,
    isLoading,
    error,
    refetch
  } = useQuery<{ items: Issue[]; total: number }>({
    queryKey: ['issues', currentPage, limit, statusFilter],
    queryFn: async () => {
      const { data } = await issuesApi.list({
        status: statusFilter || undefined,
        page: currentPage,
        limit,
      });
      return {
        items: (data.data?.items ?? []) as Issue[],
        total: data.data?.meta?.total ?? 0,
      };
    },
  });

  const issues = issuesData?.items ?? [];
  const totalItems = issuesData?.total ?? 0;

  // Fetch all issues for analytics (limited to 100)
  const { data: allIssuesData } = useQuery<Issue[]>({
    queryKey: ['issues-analytics'],
    queryFn: async () => {
      const { data } = await issuesApi.list({ limit: 100 });
      return (data.data?.items ?? []) as Issue[];
    },
  });

  const allIssues = allIssuesData ?? [];

  // Fetch issue logs when a specific issue is selected
  const {
    data: issueLogs = [],
    isLoading: isLoadingLogs
  } = useQuery<IssueLog[]>({
    queryKey: ['issue-logs', selectedIssue?.fingerprint],
    queryFn: async () => {
      if (!selectedIssue) return [];
      const { data } = await issuesApi.logs(selectedIssue.fingerprint, { limit: 10 });
      return (data.data?.items ?? []) as IssueLog[];
    },
    enabled: !!selectedIssue,
  });

  // Update status mutation
  const updateStatusMutation = useMutation({
    mutationFn: async ({ fingerprint, status }: { fingerprint: string; status: string }) => {
      await issuesApi.updateStatus(fingerprint, status);
    },
    onSuccess: (_, variables) => {
      showToast(`Issue status updated to ${variables.status}`);
      queryClient.invalidateQueries({ queryKey: ['issues'] });
      queryClient.invalidateQueries({ queryKey: ['issues-analytics'] });
      if (selectedIssue?.fingerprint === variables.fingerprint) {
        setSelectedIssue({ ...selectedIssue, status: variables.status });
      }
    },
    onError: () => {
      showToast('Failed to update issue status', 'error');
    },
  });

  const handleUpdateStatus = async (fingerprint: string, newStatus: string) => {
    updateStatusMutation.mutate({ fingerprint, status: newStatus });
  };

  const maxPage = Math.max(1, Math.ceil(totalItems / limit));

  // Analytics computations
  const top10 = [...allIssues].sort((a: Issue, b: Issue) => b.occurrence_count - a.occurrence_count).slice(0, 10);
  const bySource = allIssues.reduce<Record<string, number>>((acc: Record<string, number>, i: Issue) => {
    acc[i.source_id] = (acc[i.source_id] ?? 0) + i.occurrence_count;
    return acc;
  }, {});
  const maxSourceCount = Math.max(1, ...Object.values(bySource).map((c: number) => c));
  const levelCounts = { CRITICAL: 0, ERROR: 0, WARN: 0 };
  allIssues.filter((i: Issue) => i.status === 'OPEN').forEach((i: Issue) => {
    if (i.level in levelCounts) levelCounts[i.level as keyof typeof levelCounts]++;
  });

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-xl bg-rose-500/10 text-rose-400">
            <BugPlay size={24} />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-white">Issues</h1>
            <p className="text-sm text-zinc-400">{totalItems.toLocaleString()} grouped error patterns</p>
          </div>
        </div>
        <button
          onClick={() => refetch()}
          disabled={isLoading}
          className="flex items-center gap-2 px-3 py-2 text-xs rounded-xl bg-white/3 border border-white/5 text-zinc-400 hover:text-zinc-200 transition-all disabled:opacity-50"
        >
          <RefreshCw size={14} className={isLoading ? 'animate-spin' : ''} />
          Refresh
        </button>
      </div>

      {/* Error Message */}
      {error && (
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          className="p-4 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 flex items-center gap-3"
        >
          <AlertCircle size={20} />
          <div className="flex-1">
            <p className="font-medium">Failed to load issues</p>
          </div>
          <button
            onClick={() => refetch()}
            className="px-3 py-1.5 text-sm bg-red-500/20 hover:bg-red-500/30 rounded-lg transition-colors flex items-center gap-1"
          >
            <RefreshCw size={14} />
            Retry
          </button>
        </motion.div>
      )}

      {/* Tab Switcher */}
      <div className="flex border-b border-white/5">
        {([
          { id: 'issues', label: 'Issues', icon: BugPlay },
          { id: 'analytics', label: 'Analytics', icon: BarChart3 },
        ] as const).map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`flex items-center gap-2 px-5 py-3 text-sm font-medium transition-all border-b-2 -mb-px ${
              activeTab === tab.id
                ? 'border-purple-500 text-purple-400'
                : 'border-transparent text-zinc-500 hover:text-zinc-300'
            }`}
          >
            <tab.icon size={15} />
            {tab.label}
          </button>
        ))}
      </div>

      {activeTab === 'analytics' && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="space-y-6">
          {/* Level Breakdown */}
          <div className="grid grid-cols-3 gap-4">
            {[
              { level: 'CRITICAL', count: levelCounts.CRITICAL, color: 'text-rose-400', bg: 'bg-rose-500/10', border: 'border-rose-500/20' },
              { level: 'ERROR',    count: levelCounts.ERROR,    color: 'text-red-400',  bg: 'bg-red-500/10',  border: 'border-red-500/20' },
              { level: 'WARN',     count: levelCounts.WARN,     color: 'text-amber-400',bg: 'bg-amber-500/10',border: 'border-amber-500/20' },
            ].map((l) => (
              <div key={l.level} className={`p-4 rounded-xl border ${l.bg} ${l.border}`}>
                <p className={`text-2xl font-bold ${l.color}`}>{l.count}</p>
                <p className="text-xs text-zinc-400 mt-0.5">Open {l.level}</p>
              </div>
            ))}
          </div>

          {/* Top 10 Error Messages */}
          <div className="bg-white/2 border border-white/5 rounded-xl p-5">
            <h3 className="text-sm font-semibold text-zinc-200 mb-4">Top 10 Most Frequent Errors</h3>
            {top10.length === 0 ? (
              <p className="text-zinc-500 text-sm">No issues found.</p>
            ) : (
              <div className="space-y-3">
                {top10.map((issue, i) => (
                  <div key={issue.fingerprint} className="flex items-start gap-3">
                    <span className="text-xs font-bold text-zinc-600 w-5 shrink-0 pt-0.5">{i + 1}</span>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm text-zinc-200 truncate">{issue.message_sample}</p>
                      <div className="flex items-center gap-2 mt-1">
                        <span className={`text-[10px] font-semibold px-1.5 py-0.5 rounded border ${LEVEL_STYLES[issue.level] ?? 'text-zinc-400 border-zinc-500/20'}`}>{issue.level}</span>
                        <span className="text-xs text-zinc-500">{categoryLabel(issue.category)}</span>
                      </div>
                    </div>
                    <span className="text-sm font-semibold text-rose-400 shrink-0">{issue.occurrence_count.toLocaleString()}×</span>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Errors by Source */}
          <div className="bg-white/2 border border-white/5 rounded-xl p-5">
            <h3 className="text-sm font-semibold text-zinc-200 mb-4">Total Occurrences by Source</h3>
            {Object.keys(bySource).length === 0 ? (
              <p className="text-zinc-500 text-sm">No data.</p>
            ) : (
              <div className="space-y-3">
                {(Object.entries(bySource) as [string, number][])
                  .sort(([, a], [, b]) => b - a)
                  .map(([sourceId, count]) => (
                    <div key={sourceId} className="space-y-1.5">
                      <div className="flex justify-between text-xs">
                        <span className="text-zinc-400 font-mono truncate max-w-[60%]">{sourceId}</span>
                        <span className="text-zinc-300 font-semibold">{count.toLocaleString()}</span>
                      </div>
                      <div className="h-2 rounded-full bg-white/5 overflow-hidden">
                        <motion.div
                          initial={{ width: 0 }}
                          animate={{ width: `${(count / maxSourceCount) * 100}%` }}
                          transition={{ duration: 0.6, ease: 'easeOut' }}
                          className="h-full rounded-full bg-rose-500/60"
                        />
                      </div>
                    </div>
                  ))}
              </div>
            )}
          </div>
        </motion.div>
      )}

      {activeTab === 'issues' && (
        <>
      {/* Status Filter */}
      <div className="flex items-center gap-2 flex-wrap">
        <Filter size={14} className="text-zinc-500" />
        {(['', 'OPEN', 'RESOLVED', 'IGNORED'] as const).map((s) => (
          <button
            key={s || 'all'}
            onClick={() => setStatusFilter(s)}
            className={`px-3 py-1.5 rounded-lg text-xs font-semibold border transition-all ${
              statusFilter === s
                ? 'bg-purple-500/15 text-purple-400 border-purple-500/30'
                : 'bg-white/2 text-zinc-400 border-white/5 hover:bg-white/5'
            }`}
          >
            {s || 'All'}
          </button>
        ))}
      </div>

      <div className="bg-white/2 border border-white/5 rounded-xl overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead>
              <tr className="border-b border-white/5 text-xs font-semibold uppercase text-zinc-500">
                <th className="px-4 py-3">Issue</th>
                <th className="px-4 py-3">Level</th>
                <th className="px-4 py-3">Status</th>
                <th className="px-4 py-3">Occurrences</th>
                <th className="px-4 py-3">Last Seen</th>
                <th className="px-4 py-3">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5 text-sm">
              {isLoading ? (
                Array.from({ length: 5 }).map((_, i) => (
                  <tr key={i}>
                    {Array.from({ length: 6 }).map((__, j) => (
                      <td key={j} className="px-4 py-3">
                        <div className="h-4 rounded bg-white/3 animate-pulse" />
                      </td>
                    ))}
                  </tr>
                ))
              ) : issues.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-4 py-16 text-center text-zinc-500">
                    <CheckCircle2 size={32} className="mx-auto mb-3 text-emerald-500/40" />
                    No issues found. Everything looks clean!
                  </td>
                </tr>
              ) : (
                issues.map((issue) => (
                  <tr
                    key={issue.fingerprint}
                    onClick={() => setSelectedIssue(issue)}
                    className="hover:bg-white/5 cursor-pointer transition-colors group"
                  >
                    <td className="px-4 py-3 max-w-xs">
                      <p className="text-zinc-200 text-sm font-medium truncate group-hover:text-white transition-colors">
                        {issue.message_sample}
                      </p>
                      <p className="text-xs text-zinc-500 mt-0.5 font-mono truncate">
                        {issue.fingerprint.slice(0, 16)}... · {categoryLabel(issue.category)}
                      </p>
                    </td>
                    <td className="px-4 py-3">
                      <span className={`px-2 py-0.5 rounded-md text-xs font-semibold border ${LEVEL_STYLES[issue.level] ?? LEVEL_STYLES.DEBUG}`}>
                        {issue.level}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      <span className={`px-2 py-0.5 rounded-md text-xs font-semibold border ${STATUS_STYLES[issue.status] ?? STATUS_STYLES.OPEN}`}>
                        {issue.status}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      <span className="flex items-center gap-1 text-sm font-bold text-zinc-200">
                        <AlertTriangle size={13} className="text-amber-400" />
                        {issue.occurrence_count.toLocaleString()}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      <span className="flex items-center gap-1 text-xs text-zinc-500">
                        <Clock size={12} />
                        {timeAgo(issue.last_seen_at)}
                      </span>
                    </td>
                    <td className="px-4 py-3" onClick={(e) => e.stopPropagation()}>
                      <div className="flex items-center gap-1">
                        {issue.status !== 'RESOLVED' && (
                          <button
                            onClick={() => handleUpdateStatus(issue.fingerprint, 'RESOLVED')}
                            disabled={updateStatusMutation.isPending}
                            className="p-1.5 rounded-lg hover:bg-emerald-500/10 text-zinc-500 hover:text-emerald-400 transition-colors disabled:opacity-50"
                            title="Mark Resolved"
                          >
                            <CheckCircle2 size={15} />
                          </button>
                        )}
                        {issue.status !== 'IGNORED' && (
                          <button
                            onClick={() => handleUpdateStatus(issue.fingerprint, 'IGNORED')}
                            disabled={updateStatusMutation.isPending}
                            className="p-1.5 rounded-lg hover:bg-zinc-500/10 text-zinc-500 hover:text-zinc-300 transition-colors disabled:opacity-50"
                            title="Ignore"
                          >
                            <EyeOff size={15} />
                          </button>
                        )}
                        {issue.status !== 'OPEN' && (
                          <button
                            onClick={() => handleUpdateStatus(issue.fingerprint, 'OPEN')}
                            disabled={updateStatusMutation.isPending}
                            className="p-1.5 rounded-lg hover:bg-rose-500/10 text-zinc-500 hover:text-rose-400 transition-colors disabled:opacity-50"
                            title="Reopen"
                          >
                            <RefreshCw size={15} />
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        <div className="p-4 border-t border-white/5 flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2 text-sm text-zinc-400">
            <span>Show</span>
            <SelectField
              value={limit}
              onChange={(e) => { setLimit(Number(e.target.value)); setCurrentPage(1); }}
              wrapperClassName="w-24"
            >
              <option value="10">10</option>
              <option value="20">20</option>
              <option value="50">50</option>
              <option value="100">100</option>
              <option value="99999">All</option>
            </SelectField>
            <span>of {totalItems.toLocaleString()} issues</span>
          </div>
          <div className="flex items-center gap-4">
            <span className="text-sm text-zinc-400">
              Page <span className="text-zinc-100">{currentPage}</span> of{' '}
              <span className="text-zinc-100">{maxPage}</span>
            </span>
            <div className="flex items-center gap-1">
              <button disabled={currentPage === 1} onClick={() => setCurrentPage((p) => Math.max(1, p - 1))} className="p-2 rounded-lg border border-white/5 hover:bg-white/5 disabled:opacity-40 text-zinc-300 disabled:cursor-not-allowed">
                <ChevronLeft size={16} />
              </button>
              <button disabled={currentPage === maxPage} onClick={() => setCurrentPage((p) => Math.min(maxPage, p + 1))} className="p-2 rounded-lg border border-white/5 hover:bg-white/5 disabled:opacity-40 text-zinc-300 disabled:cursor-not-allowed">
                <ChevronRight size={16} />
              </button>
            </div>
          </div>
        </div>
      </div>
        </>
      )}

      {/* Issue Detail Modal */}
      <AnimatePresence>
        {selectedIssue && (
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setSelectedIssue(null)}
              className="absolute inset-0 bg-black/60 backdrop-blur-sm"
            />
            <motion.div
              initial={{ opacity: 0, scale: 0.95, y: 10 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95, y: 10 }}
              className="relative w-full max-w-2xl max-h-[90vh] overflow-y-auto p-6 rounded-xl bg-[#0c0c0e] border border-white/5 shadow-2xl space-y-5"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="space-y-1 flex-1 min-w-0">
                  <div className="flex items-center gap-2 flex-wrap">
                    <span className={`px-2 py-0.5 rounded-md text-xs font-bold border ${LEVEL_STYLES[selectedIssue.level] ?? LEVEL_STYLES.DEBUG}`}>
                      {selectedIssue.level}
                    </span>
                    <span className={`px-2 py-0.5 rounded-md text-xs font-bold border ${STATUS_STYLES[selectedIssue.status] ?? STATUS_STYLES.OPEN}`}>
                      {selectedIssue.status}
                    </span>
                    <span className="text-xs text-purple-400 bg-purple-500/10 px-2 py-0.5 rounded-md">
                      {categoryLabel(selectedIssue.category)}
                    </span>
                  </div>
                  <p className="text-zinc-100 font-semibold leading-relaxed">{selectedIssue.message_sample}</p>
                </div>
                <button onClick={() => setSelectedIssue(null)} className="p-1.5 rounded-lg hover:bg-white/10 text-zinc-400 shrink-0">
                  <X size={20} />
                </button>
              </div>

              <div className="grid grid-cols-3 gap-3">
                <div className="bg-white/2 border border-white/5 rounded-xl p-3">
                  <p className="text-xs text-zinc-500">Occurrences</p>
                  <p className="text-xl font-bold text-zinc-100 mt-0.5">{selectedIssue.occurrence_count.toLocaleString()}</p>
                </div>
                <div className="bg-white/2 border border-white/5 rounded-xl p-3">
                  <p className="text-xs text-zinc-500">First Seen</p>
                  <p className="text-xs font-mono text-zinc-300 mt-0.5">{timeAgo(selectedIssue.first_seen_at)}</p>
                </div>
                <div className="bg-white/2 border border-white/5 rounded-xl p-3">
                  <p className="text-xs text-zinc-500">Last Seen</p>
                  <p className="text-xs font-mono text-zinc-300 mt-0.5">{timeAgo(selectedIssue.last_seen_at)}</p>
                </div>
              </div>

              <div>
                <p className="text-xs text-zinc-500 font-mono mb-1">Fingerprint</p>
                <p className="font-mono text-xs text-zinc-400 bg-black/30 px-3 py-2 rounded-lg break-all">{selectedIssue.fingerprint}</p>
              </div>

              {/* Action buttons */}
              <div className="flex gap-2 flex-wrap">
                {selectedIssue.status !== 'RESOLVED' && (
                  <button onClick={() => handleUpdateStatus(selectedIssue.fingerprint, 'RESOLVED')} disabled={updateStatusMutation.isPending} className="flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 hover:bg-emerald-500/20 transition-colors disabled:opacity-50">
                    <CheckCircle2 size={13} /> Mark Resolved
                  </button>
                )}
                {selectedIssue.status !== 'IGNORED' && (
                  <button onClick={() => handleUpdateStatus(selectedIssue.fingerprint, 'IGNORED')} disabled={updateStatusMutation.isPending} className="flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-semibold bg-zinc-500/10 text-zinc-400 border border-zinc-500/20 hover:bg-zinc-500/20 transition-colors disabled:opacity-50">
                    <EyeOff size={13} /> Ignore
                  </button>
                )}
                {selectedIssue.status !== 'OPEN' && (
                  <button onClick={() => handleUpdateStatus(selectedIssue.fingerprint, 'OPEN')} disabled={updateStatusMutation.isPending} className="flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-semibold bg-rose-500/10 text-rose-400 border border-rose-500/20 hover:bg-rose-500/20 transition-colors disabled:opacity-50">
                    <RefreshCw size={13} /> Reopen
                  </button>
                )}
              </div>

              {/* Related Logs */}
              <div>
                <h4 className="text-sm font-semibold text-zinc-300 mb-3">Recent Occurrences</h4>
                {isLoadingLogs ? (
                  <div className="space-y-2">
                    {Array.from({ length: 3 }).map((_, i) => (
                      <div key={i} className="h-10 rounded-lg bg-white/2 animate-pulse" />
                    ))}
                  </div>
                ) : issueLogs.length === 0 ? (
                  <p className="text-xs text-zinc-500">No individual log entries found for this issue.</p>
                ) : (
                  <div className="space-y-2">
                    {issueLogs.map((log) => (
                      <div key={log.id} className="p-3 rounded-lg bg-white/2 border border-white/5 flex items-center gap-3">
                        <span className={`shrink-0 px-2 py-0.5 rounded text-xs font-semibold border ${LEVEL_STYLES[log.level] ?? LEVEL_STYLES.DEBUG}`}>
                          {log.level}
                        </span>
                        <span className="flex-1 text-xs text-zinc-400 truncate">{log.message}</span>
                        <span className="shrink-0 text-xs font-mono text-zinc-600">{timeAgo(log.created_at)}</span>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>

      {/* Toast */}
      <AnimatePresence>
        {toast && (
          <motion.div
            initial={{ opacity: 0, y: 40, scale: 0.95 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: 40, scale: 0.95 }}
            className={`fixed bottom-24 right-6 px-4 py-3 rounded-xl border shadow-2xl flex items-center gap-2 z-60 text-sm font-medium ${
              toast.type === 'success'
                ? 'bg-[#111113] text-emerald-400 border-emerald-500/20'
                : 'bg-[#111113] text-rose-400 border-rose-500/20'
            }`}
          >
            {toast.type === 'success' ? <CheckCircle2 size={15} /> : <AlertCircle size={15} />}
            {toast.message}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
