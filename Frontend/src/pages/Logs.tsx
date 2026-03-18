import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Filter, Search, ChevronRight, X, Sparkles, ChevronLeft } from 'lucide-react';
import { logsApi } from '../shared/lib/api';

interface Log {
  id: number;
  time: string;
  level: string;
  category: string;
  message: string;
  source: string;
  stack: string;
}

export default function Logs() {
  const [selectedLog, setSelectedLog] = useState<Log | null>(null);
  const [itemsPerPage, setItemsPerPage] = useState<number | 'all'>(10);
  const [currentPage, setCurrentPage] = useState(1);
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [filters, setFilters] = useState({ level: '', source: '' });
  const [searchQuery, setSearchQuery] = useState('');
  
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [aiResponse, setAiResponse] = useState<string | null>(null);

  const [logs, setLogs] = useState<Log[]>([]);
  useEffect(() => {
    const fetchLogs = async () => {
      try {
        const { data } = await logsApi.getLogs({
          level: filters.level || undefined,
          source: filters.source || undefined,
          search: searchQuery || undefined,
        });
        setLogs(data.data?.items || []);
      } catch (error) {
         console.error("Failed to fetch logs", error);
      }
    };
    fetchLogs();
  }, [filters, searchQuery]);

  const levelStyles: Record<string, string> = {
    CRITICAL: 'bg-rose-500/10 text-rose-400 border-rose-500/20',
    ERROR: 'bg-red-500/10 text-red-400 border-red-500/20',
    WARN: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
    INFO: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
  };

  const filteredLogs = logs.filter(log => {
    const matchesSearch = log.message.toLowerCase().includes(searchQuery.toLowerCase()) || 
                         log.source.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         log.category.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesLevel = !filters.level || log.level === filters.level;
    const matchesSource = !filters.source || log.source === filters.source;
    return matchesSearch && matchesLevel && matchesSource;
  });

  const totalLogs = filteredLogs.length; 
  const maxPage = itemsPerPage === 'all' ? 1 : Math.ceil(totalLogs / itemsPerPage);

  const handleAiAnalysis = () => {
    setIsAnalyzing(true);
    setAiResponse(null);
    setTimeout(() => {
      setAiResponse(`### AI Diagnostics
**Root Cause Hypothesis:** The database connection pool was saturated by sudden peak traffic of 30,000 req/sec.

**Recommended Actions:**
1. Increase \`max_connections\` configuration from 100 to 250 explicitly.
2. Check leak buffers inside package \`main.connectDB()\`.
3. Enable redis-cache layer setups on high usage read tables headers.`);
      setIsAnalyzing(false);
    }, 1500);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-white">Log Explorer</h1>
          <p className="text-sm text-zinc-400">Search and filter system logs</p>
        </div>
      </div>

      <div className="flex flex-col gap-4">
        <div className="flex flex-col gap-4 md:flex-row md:items-center">
          <div className="flex-1 relative">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-500" size={18} />
            <input
              type="text"
              placeholder="Search messages, sources, IP..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-11 pr-4 py-2.5 rounded-xl bg-white/[0.03] border border-white/[0.05] text-zinc-200 placeholder-zinc-500 focus:outline-none focus:border-purple-500/30 transition-all text-sm"
            />
          </div>
          <div className="flex items-center gap-2">
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
            </button>
          </div>
        </div>

        <AnimatePresence>
          {isFilterOpen && (
            <motion.div
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              className="p-4 rounded-xl bg-white/2 border border-white/5 backdrop-blur-sm grid grid-cols-1 md:grid-cols-2 gap-4"
            >
              <div>
                <label className="block text-xs font-medium text-zinc-400 mb-1">Level</label>
                <select 
                  value={filters.level}
                  onChange={(e) => setFilters({ ...filters, level: e.target.value })}
                  className="w-full px-3 py-2 rounded-lg bg-white/3 border border-white/5 text-zinc-200 focus:outline-none text-sm"
                >
                  <option value="">All Levels</option>
                  <option value="CRITICAL">Critical</option>
                  <option value="ERROR">Error</option>
                  <option value="WARN">Warning</option>
                  <option value="INFO">Info</option>
                </select>
              </div>
              <div>
                <label className="block text-xs font-medium text-zinc-400 mb-1">Source</label>
                <select 
                  value={filters.source}
                  onChange={(e) => setFilters({ ...filters, source: e.target.value })}
                  className="w-full px-3 py-2 rounded-lg bg-white/3 border border-white/5 text-zinc-200 focus:outline-none text-sm"
                >
                  <option value="">All Sources</option>
                  <option value="Auth Service">Auth Service</option>
                  <option value="Gateway">Gateway</option>
                  <option value="DB Analytics">DB Analytics</option>
                </select>
              </div>
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
                <th className="px-6 py-4">Source</th>
                <th className="px-6 py-4">Level</th>
                <th className="px-6 py-4">Message</th>
                <th className="px-6 py-4"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/[0.03] text-sm text-zinc-300">
              {filteredLogs.map((log) => (
                <tr
                  key={log.id}
                  onClick={() => setSelectedLog(log)}
                  className="hover:bg-white/[0.01] cursor-pointer transition-colors group"
                >
                  <td className="px-6 py-4 text-xs font-mono text-zinc-500">{log.time}</td>
                  <td className="px-6 py-4 font-medium text-zinc-200">{log.source}</td>
                  <td className="px-6 py-4">
                    <span className={`px-2 py-1 rounded-md text-xs font-semibold border ${levelStyles[log.level]}`}>
                      {log.level}
                    </span>
                  </td>
                  <td className="px-6 py-4 truncate max-w-md text-zinc-400 group-hover:text-zinc-200 transition-colors">
                    {log.message}
                  </td>
                  <td className="px-4 py-4 text-right">
                    <ChevronRight size={16} className="text-zinc-600 opacity-0 group-hover:opacity-100 transition-all" />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Pagination Footer */}
        <div className="p-4 border-t border-white/[0.05] flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2 text-sm text-zinc-400">
            <span>Show</span>
            <select
              value={itemsPerPage}
              onChange={(e) => {
                const val = e.target.value === 'all' ? 'all' : Number(e.target.value);
                setItemsPerPage(val);
                setCurrentPage(1);
              }}
              className="px-2 py-1 rounded bg-white/[0.04] border border-white/[0.08] text-zinc-200 focus:outline-none"
            >
              <option value="10">10</option>
              <option value="50">50</option>
              <option value="100">100</option>
              <option value="all">All</option>
            </select>
            <span>entries</span>
          </div>

          <div className="flex items-center gap-4">
            <span className="text-sm text-zinc-400">
              Page <span className="text-zinc-100">{currentPage}</span> of <span className="text-zinc-100">{maxPage}</span>
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
              className="relative w-full max-w-2xl p-6 rounded-2xl bg-[#0c0c0e] border border-white/[0.08] shadow-2xl space-y-4"
            >
              <div className="flex items-center justify-between border-b border-white/[0.05] pb-4">
                <div className="flex items-center gap-3">
                  <span className={`px-2.5 py-1 rounded-md text-xs font-bold border ${levelStyles[selectedLog.level]}`}>
                    {selectedLog.level}
                  </span>
                  <span className="text-sm text-zinc-500">{selectedLog.time}</span>
                </div>
                <button onClick={() => { setSelectedLog(null); setAiResponse(null); }} className="p-1 rounded-lg hover:bg-white/10 text-zinc-400"><X size={20}/></button>
              </div>

              <div className="space-y-4">
                <div>
                  <p className="text-xs text-zinc-500">Message</p>
                  <p className="text-lg font-semibold text-zinc-100 mt-0.5">{selectedLog.message}</p>
                </div>
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <p className="text-xs text-zinc-500">Source</p>
                    <p className="font-medium text-zinc-200 mt-0.5">{selectedLog.source}</p>
                  </div>
                  <div>
                    <p className="text-xs text-zinc-500">Category</p>
                    <p className="font-medium text-zinc-200 mt-0.5">{selectedLog.category}</p>
                  </div>
                </div>

                {selectedLog.stack && (
                  <div>
                    <p className="text-xs text-zinc-500 mb-1">Stack Trace</p>
                    <pre className="p-4 rounded-xl bg-black/40 border border-white/[0.05] text-xs font-mono text-rose-300 overflow-x-auto">
                      {selectedLog.stack}
                    </pre>
                  </div>
                )}

                {aiResponse && (
                  <div className="p-4 rounded-xl bg-purple-500/5 border border-purple-500/10 text-zinc-300 text-sm space-y-2">
                    <p className="font-semibold text-purple-400 flex items-center gap-1.5"><Sparkles size={14}/> AI Diagnostics Result</p>
                    <pre className="text-xs font-sans whitespace-pre-wrap text-zinc-400">{aiResponse}</pre>
                  </div>
                )}

                <div className="pt-4 border-t border-white/[0.05]">
                  <button 
                    onClick={handleAiAnalysis}
                    disabled={isAnalyzing}
                    className="flex items-center justify-center w-full gap-2 px-4 py-3 text-sm font-semibold text-white bg-gradient-to-r from-purple-500 to-indigo-500 rounded-xl hover:opacity-90 transition-opacity disabled:opacity-50"
                  >
                    <Sparkles size={16} className={isAnalyzing ? 'animate-spin' : ''} />
                    {isAnalyzing ? 'Analyzing Log...' : 'Run AI Analysis'}
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
