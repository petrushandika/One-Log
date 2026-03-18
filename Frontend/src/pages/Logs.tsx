import { useState, useEffect, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Filter, Search, ChevronRight, X, Sparkles, ChevronLeft, Download, ScrollText } from 'lucide-react';
import SelectField from '../shared/components/SelectField';
import { logsApi, sourcesApi } from '../shared/lib/api';
import { categoryLabel } from '../shared/lib/utils';

interface Log {
  id: number;
  source_id: string;
  category: string;
  level: string;
  message: string;
  created_at: string;
  stack_trace: string;
  ai_insight: { analysis?: string } | null;
  fingerprint: string;
  ip_address: string;
}

interface Source {
  id: string;
  name: string;
}

const LEVEL_STYLES: Record<string, string> = {
  CRITICAL: 'bg-rose-500/10 text-rose-400 border-rose-500/20',
  ERROR: 'bg-red-500/10 text-red-400 border-red-500/20',
  WARN: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
  WARNING: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
  INFO: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
  DEBUG: 'bg-zinc-500/10 text-zinc-400 border-zinc-500/20',
};

export default function Logs() {
  const [selectedLog, setSelectedLog] = useState<Log | null>(null);
  const [limit, setLimit] = useState<number>(20);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalLogs, setTotalLogs] = useState(0);
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [filters, setFilters] = useState({ level: '', source_id: '', category: '', from: '', to: '' });
  const [searchQuery, setSearchQuery] = useState('');
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [aiResponse, setAiResponse] = useState<string | null>(null);
  const [logs, setLogs] = useState<Log[]>([]);
  const [sources, setSources] = useState<Source[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isExporting, setIsExporting] = useState(false);

  const maxPage = Math.max(1, Math.ceil(totalLogs / limit));

  const fetchLogs = useCallback(async () => {
    setIsLoading(true);
    try {
      const { data } = await logsApi.getLogs({
        source_id: filters.source_id || undefined,
        level: filters.level || undefined,
        category: filters.category || undefined,
        from: filters.from ? new Date(filters.from).toISOString() : undefined,
        to: filters.to ? new Date(filters.to).toISOString() : undefined,
        page: currentPage,
        limit,
      });
      setLogs(data.data?.items ?? []);
      setTotalLogs(data.data?.meta?.total ?? 0);
    } catch (err) {
      console.error('Failed to fetch logs', err);
    } finally {
      setIsLoading(false);
    }
  }, [filters, currentPage, limit]);

  useEffect(() => {
    fetchLogs();
  }, [fetchLogs]);

  useEffect(() => {
    sourcesApi.getAll().then(({ data }) => setSources(data.data ?? [])).catch(console.error);
  }, []);

  // Reset to page 1 on filter change
  useEffect(() => {
    setCurrentPage(1);
  }, [filters, limit]);

  const handleAiAnalysis = async () => {
    if (!selectedLog) return;
    setIsAnalyzing(true);
    setAiResponse(null);
    try {
      const { data } = await logsApi.analyze(selectedLog.id);
      const insight = data.data?.ai_insight;
      if (insight) {
        const parsed = typeof insight === 'string' ? JSON.parse(insight) : insight;
        setAiResponse(parsed?.analysis ?? JSON.stringify(parsed, null, 2));
      } else {
        setAiResponse('No analysis was generated. Please try again.');
      }
      // Refresh the log entry to show updated ai_insight badge
      setSelectedLog({ ...selectedLog, ai_insight: data.data?.ai_insight ? (typeof data.data.ai_insight === 'string' ? JSON.parse(data.data.ai_insight) : data.data.ai_insight) : null });
    } catch (err) {
      setAiResponse('Failed to run AI analysis. Please check your AI configuration.');
      console.error(err);
    } finally {
      setIsAnalyzing(false);
    }
  };

  const handleOpenLog = (log: Log) => {
    setSelectedLog(log);
    setAiResponse(null);
    // Show existing AI insight if already analyzed
    if (log.ai_insight) {
      const insight = typeof log.ai_insight === 'string' ? JSON.parse(log.ai_insight as unknown as string) : log.ai_insight;
      if (insight?.analysis) setAiResponse(insight.analysis);
    }
  };

  const handleExportCSV = async () => {
    setIsExporting(true);
    try {
      const response = await logsApi.export({
        source_id: filters.source_id || undefined,
        level: filters.level || undefined,
        category: filters.category || undefined,
        from: filters.from ? new Date(filters.from).toISOString() : undefined,
        to: filters.to ? new Date(filters.to).toISOString() : undefined,
      });
      const url = URL.createObjectURL(new Blob([response.data], { type: 'text/csv' }));
      const a = document.createElement('a');
      a.href = url;
      a.download = `logs-${new Date().toISOString().slice(0, 10)}.csv`;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      console.error('Export failed', err);
    } finally {
      setIsExporting(false);
    }
  };

  const displayedLogs = searchQuery
    ? logs.filter(
        (l) =>
          l.message.toLowerCase().includes(searchQuery.toLowerCase()) ||
          l.category.toLowerCase().includes(searchQuery.toLowerCase()) ||
          l.source_id.toLowerCase().includes(searchQuery.toLowerCase()),
      )
    : logs;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-white flex items-center gap-2.5">
            <ScrollText size={22} className="text-purple-400" />
            Log Explorer
          </h1>
          <p className="text-sm text-zinc-400">
            {totalLogs.toLocaleString()} total logs
          </p>
        </div>
        <button
          onClick={handleExportCSV}
          disabled={isExporting}
          className="flex items-center gap-2 px-4 py-2.5 rounded-xl bg-white/[0.03] border border-white/[0.06] text-zinc-300 hover:bg-white/[0.06] text-sm transition-all disabled:opacity-50"
        >
          <Download size={16} />
          {isExporting ? 'Exporting...' : 'Export CSV'}
        </button>
      </div>

      <div className="flex flex-col gap-4">
        <div className="flex flex-col gap-4 md:flex-row md:items-center">
          <div className="flex-1 relative">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-500" size={18} />
            <input
              type="text"
              placeholder="Search messages, source IDs, categories..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-11 pr-4 py-2.5 rounded-xl bg-white/[0.03] border border-white/[0.05] text-zinc-200 placeholder-zinc-500 focus:outline-none focus:border-purple-500/30 transition-all text-sm"
            />
          </div>
          <button
            onClick={() => setIsFilterOpen(!isFilterOpen)}
            className={`flex items-center gap-2 px-4 py-2.5 rounded-xl border text-sm transition-all ${
              isFilterOpen
                ? 'bg-purple-500/10 border-purple-500/30 text-purple-400'
                : 'bg-white/[0.03] border-white/[0.05] text-zinc-300 hover:bg-white/[0.06]'
            }`}
          >
            <Filter size={16} />
            Filters
            {(filters.level || filters.source_id || filters.category) && (
              <span className="w-2 h-2 rounded-full bg-purple-400" />
            )}
          </button>
        </div>

        <AnimatePresence>
          {isFilterOpen && (
            <motion.div
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              className="p-4 rounded-xl bg-white/[0.02] border border-white/5 backdrop-blur-sm grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-5 gap-4"
            >
              <div>
                <label className="block text-xs font-medium text-zinc-400 mb-1.5">Level</label>
                <SelectField value={filters.level} onChange={(e) => setFilters({ ...filters, level: e.target.value })}>
                  <option value="">All Levels</option>
                  <option value="CRITICAL">Critical</option>
                  <option value="ERROR">Error</option>
                  <option value="WARN">Warning</option>
                  <option value="INFO">Info</option>
                  <option value="DEBUG">Debug</option>
                </SelectField>
              </div>
              <div>
                <label className="block text-xs font-medium text-zinc-400 mb-1.5">Source</label>
                <SelectField value={filters.source_id} onChange={(e) => setFilters({ ...filters, source_id: e.target.value })}>
                  <option value="">All Sources</option>
                  {sources.map((s) => (
                    <option key={s.id} value={s.id}>{s.name}</option>
                  ))}
                </SelectField>
              </div>
              <div>
                <label className="block text-xs font-medium text-zinc-400 mb-1.5">Category</label>
                <SelectField value={filters.category} onChange={(e) => setFilters({ ...filters, category: e.target.value })}>
                  <option value="">All Categories</option>
                  <option value="SYSTEM_ERROR">System Error</option>
                  <option value="AUTH_EVENT">Auth Event</option>
                  <option value="USER_ACTIVITY">User Activity</option>
                  <option value="SECURITY">Security</option>
                  <option value="PERFORMANCE">Performance</option>
                  <option value="AUDIT_TRAIL">Audit Trail</option>
                </SelectField>
              </div>
              <div>
                <label className="block text-xs font-medium text-zinc-400 mb-1">From</label>
                <input
                  type="datetime-local"
                  value={filters.from}
                  onChange={(e) => setFilters({ ...filters, from: e.target.value })}
                  className="w-full px-3 py-2 rounded-lg bg-white/[0.03] border border-white/5 text-zinc-200 focus:outline-none text-sm [color-scheme:dark]"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-zinc-400 mb-1">To</label>
                <input
                  type="datetime-local"
                  value={filters.to}
                  onChange={(e) => setFilters({ ...filters, to: e.target.value })}
                  className="w-full px-3 py-2 rounded-lg bg-white/[0.03] border border-white/5 text-zinc-200 focus:outline-none text-sm [color-scheme:dark]"
                />
              </div>
              {(filters.from || filters.to) && (
                <div className="md:col-span-3 flex">
                  <button
                    onClick={() => setFilters({ ...filters, from: '', to: '' })}
                    className="text-xs text-zinc-500 hover:text-zinc-200 transition-colors"
                  >
                    Clear date range
                  </button>
                </div>
              )}
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="rounded-2xl bg-white/[0.02] border border-white/[0.05] backdrop-blur-sm overflow-hidden"
      >
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead>
              <tr className="border-b border-white/[0.05] text-xs font-semibold uppercase tracking-wider text-zinc-400">
                <th className="px-6 py-4">Timestamp</th>
                <th className="px-6 py-4">Category</th>
                <th className="px-6 py-4">Level</th>
                <th className="px-6 py-4">Message</th>
                <th className="px-6 py-4">AI</th>
                <th className="px-6 py-4" />
              </tr>
            </thead>
            <tbody className="divide-y divide-white/[0.03] text-sm text-zinc-300">
              {isLoading ? (
                Array.from({ length: 5 }).map((_, i) => (
                  <tr key={i}>
                    {Array.from({ length: 6 }).map((__, j) => (
                      <td key={j} className="px-6 py-4">
                        <div className="h-4 rounded bg-white/[0.03] animate-pulse" />
                      </td>
                    ))}
                  </tr>
                ))
              ) : displayedLogs.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-6 py-12 text-center text-zinc-500">
                    No logs found. Try adjusting your filters.
                  </td>
                </tr>
              ) : (
                displayedLogs.map((log) => (
                  <tr
                    key={log.id}
                    onClick={() => handleOpenLog(log)}
                    className="hover:bg-white/[0.01] cursor-pointer transition-colors group"
                  >
                    <td className="px-6 py-4 text-xs font-mono text-zinc-500 whitespace-nowrap">
                      {new Date(log.created_at).toLocaleString()}
                    </td>
                    <td className="px-6 py-4 text-xs text-zinc-400">{categoryLabel(log.category)}</td>
                    <td className="px-6 py-4">
                      <span className={`px-2 py-1 rounded-md text-xs font-semibold border ${LEVEL_STYLES[log.level] ?? LEVEL_STYLES.DEBUG}`}>
                        {log.level}
                      </span>
                    </td>
                    <td className="px-6 py-4 truncate max-w-xs text-zinc-400 group-hover:text-zinc-200 transition-colors">
                      {log.message}
                    </td>
                    <td className="px-6 py-4">
                      {log.ai_insight && (
                        <span className="flex items-center gap-1 text-purple-400 text-xs">
                          <Sparkles size={12} /> AI
                        </span>
                      )}
                    </td>
                    <td className="px-4 py-4 text-right">
                      <ChevronRight size={16} className="text-zinc-600 opacity-0 group-hover:opacity-100 transition-all" />
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination Footer */}
        <div className="p-4 border-t border-white/[0.05] flex flex-col md:flex-row items-center justify-between gap-4">
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
            <span>of {totalLogs.toLocaleString()} entries</span>
          </div>

          <div className="flex items-center gap-4">
            <span className="text-sm text-zinc-400">
              Page <span className="text-zinc-100">{currentPage}</span> of{' '}
              <span className="text-zinc-100">{maxPage}</span>
            </span>
            <div className="flex items-center gap-1">
              <button
                disabled={currentPage === 1}
                onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                className="p-2 rounded-lg border border-white/[0.04] hover:bg-white/[0.03] disabled:opacity-40 text-zinc-300 disabled:cursor-not-allowed"
              >
                <ChevronLeft size={16} />
              </button>
              <button
                disabled={currentPage === maxPage}
                onClick={() => setCurrentPage((p) => Math.min(maxPage, p + 1))}
                className="p-2 rounded-lg border border-white/[0.04] hover:bg-white/[0.03] disabled:opacity-40 text-zinc-300 disabled:cursor-not-allowed"
              >
                <ChevronRight size={16} />
              </button>
            </div>
          </div>
        </div>
      </motion.div>

      {/* Log Detail Modal */}
      <AnimatePresence>
        {selectedLog && (
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => { setSelectedLog(null); setAiResponse(null); }}
              className="absolute inset-0 bg-black/60 backdrop-blur-sm"
            />
            <motion.div
              initial={{ opacity: 0, scale: 0.95, y: 10 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95, y: 10 }}
              className="relative w-full max-w-2xl max-h-[90vh] overflow-y-auto p-6 rounded-2xl bg-[#0c0c0e] border border-white/[0.08] shadow-2xl space-y-4"
            >
              <div className="flex items-center justify-between border-b border-white/[0.05] pb-4">
                <div className="flex items-center gap-3">
                  <span className={`px-2.5 py-1 rounded-md text-xs font-bold border ${LEVEL_STYLES[selectedLog.level] ?? LEVEL_STYLES.DEBUG}`}>
                    {selectedLog.level}
                  </span>
                  <span className="text-xs text-zinc-400 bg-white/[0.03] px-2 py-1 rounded-lg">
                    {categoryLabel(selectedLog.category)}
                  </span>
                </div>
                <button
                  onClick={() => { setSelectedLog(null); setAiResponse(null); }}
                  className="p-1 rounded-lg hover:bg-white/10 text-zinc-400"
                >
                  <X size={20} />
                </button>
              </div>

              <div className="space-y-4">
                <div>
                  <p className="text-xs text-zinc-500 mb-1">Message</p>
                  <p className="text-base font-semibold text-zinc-100 leading-relaxed">{selectedLog.message}</p>
                </div>

                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <p className="text-xs text-zinc-500">Timestamp</p>
                    <p className="font-mono text-xs text-zinc-300 mt-0.5">
                      {new Date(selectedLog.created_at).toLocaleString()}
                    </p>
                  </div>
                  {selectedLog.ip_address && (
                    <div>
                      <p className="text-xs text-zinc-500">IP Address</p>
                      <p className="font-mono text-xs text-zinc-300 mt-0.5">{selectedLog.ip_address}</p>
                    </div>
                  )}
                  <div>
                    <p className="text-xs text-zinc-500">Fingerprint</p>
                    <p className="font-mono text-xs text-zinc-400 mt-0.5 truncate">{selectedLog.fingerprint}</p>
                  </div>
                  <div>
                    <p className="text-xs text-zinc-500">Log ID</p>
                    <p className="font-mono text-xs text-zinc-400 mt-0.5">#{selectedLog.id}</p>
                  </div>
                </div>

                {selectedLog.stack_trace && (
                  <div>
                    <p className="text-xs text-zinc-500 mb-1">Stack Trace</p>
                    <pre className="p-4 rounded-xl bg-black/40 border border-white/[0.05] text-xs font-mono text-rose-300 overflow-x-auto whitespace-pre-wrap leading-relaxed">
                      {selectedLog.stack_trace}
                    </pre>
                  </div>
                )}

                {/* AI Response with Markdown rendering */}
                {aiResponse && (
                  <motion.div
                    initial={{ opacity: 0, y: 8 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="p-4 rounded-xl bg-purple-500/5 border border-purple-500/15"
                  >
                    <p className="font-semibold text-purple-400 flex items-center gap-1.5 mb-3 text-sm">
                      <Sparkles size={14} /> AI Diagnostics
                    </p>
                    <div className="text-zinc-300 space-y-1 prose-sm">
                      <ReactMarkdown
                        remarkPlugins={[remarkGfm]}
                        components={{
                          h1: ({ children }) => <h1 className="text-purple-300 font-bold text-sm mt-3 mb-1">{children}</h1>,
                          h2: ({ children }) => <h2 className="text-purple-300 font-bold text-sm mt-3 mb-1">{children}</h2>,
                          h3: ({ children }) => <h3 className="text-purple-400 font-semibold text-xs mt-3 mb-1 uppercase tracking-wide">{children}</h3>,
                          strong: ({ children }) => <strong className="text-zinc-200 font-semibold">{children}</strong>,
                          em: ({ children }) => <em className="text-zinc-300 italic">{children}</em>,
                          p: ({ children }) => <p className="text-zinc-400 text-xs leading-relaxed mb-2">{children}</p>,
                          li: ({ children }) => <li className="text-zinc-400 text-xs ml-4 mb-0.5 list-disc">{children}</li>,
                          ol: ({ children }) => <ol className="list-decimal space-y-1 mb-2">{children}</ol>,
                          ul: ({ children }) => <ul className="list-disc space-y-1 mb-2">{children}</ul>,
                          code: ({ children }) => (
                            <code className="bg-black/50 px-1.5 py-0.5 rounded text-purple-300 font-mono text-xs">
                              {children}
                            </code>
                          ),
                          pre: ({ children }) => (
                            <pre className="bg-black/40 p-3 rounded-lg overflow-x-auto text-xs font-mono text-zinc-300 my-2">
                              {children}
                            </pre>
                          ),
                        }}
                      >
                        {aiResponse}
                      </ReactMarkdown>
                    </div>
                  </motion.div>
                )}

                <div className="pt-4 border-t border-white/[0.05]">
                  <button
                    onClick={handleAiAnalysis}
                    disabled={isAnalyzing}
                    className="flex items-center justify-center w-full gap-2 px-4 py-3 text-sm font-semibold text-white bg-gradient-to-r from-purple-500 to-indigo-500 rounded-xl hover:opacity-90 transition-opacity disabled:opacity-50"
                  >
                    <Sparkles size={16} className={isAnalyzing ? 'animate-spin' : ''} />
                    {isAnalyzing ? 'Analyzing with AI...' : aiResponse ? 'Re-analyze' : 'Run AI Analysis'}
                  </button>
                </div>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>
    </div>
  );
}
