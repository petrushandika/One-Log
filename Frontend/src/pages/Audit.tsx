import { useState, useEffect, useCallback } from 'react';
import { motion } from 'framer-motion';
import { ShieldAlert, Search, ChevronLeft, ChevronRight, AlertTriangle, Info, ShieldCheck } from 'lucide-react';
import { activityApi } from '../shared/lib/api';
import SelectField from '../shared/components/SelectField';

interface ActivityLog {
  id: number;
  source_id: string;
  category: string;
  level: string;
  message: string;
  ip_address: string;
  created_at: string;
  context: Record<string, unknown> | null;
}

const LEVEL_STYLES: Record<string, string> = {
  CRITICAL: 'bg-rose-500/10 text-rose-400 border-rose-500/20',
  ERROR: 'bg-red-500/10 text-red-400 border-red-500/20',
  WARN: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
  INFO: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
  DEBUG: 'bg-zinc-500/10 text-zinc-400 border-zinc-500/20',
};

const LEVEL_ICONS: Record<string, React.ReactNode> = {
  CRITICAL: <AlertTriangle size={14} />,
  ERROR: <AlertTriangle size={14} />,
  WARN: <AlertTriangle size={14} />,
  INFO: <Info size={14} />,
  DEBUG: <ShieldCheck size={14} />,
};

export default function Audit() {
  const [activities, setActivities] = useState<ActivityLog[]>([]);
  const [totalItems, setTotalItems] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [limit, setLimit] = useState(20);
  const [searchQuery, setSearchQuery] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const maxPage = Math.max(1, Math.ceil(totalItems / limit));

  const fetchActivity = useCallback(async () => {
    setIsLoading(true);
    try {
      const { data } = await activityApi.list({ page: currentPage, limit });
      setActivities(data.data?.items ?? []);
      setTotalItems(data.data?.meta?.total ?? 0);
    } catch (err) {
      console.error('Failed to fetch activity', err);
    } finally {
      setIsLoading(false);
    }
  }, [currentPage, limit]);

  useEffect(() => {
    fetchActivity();
  }, [fetchActivity]);

  const displayed = searchQuery
    ? activities.filter(
        (a) =>
          a.message.toLowerCase().includes(searchQuery.toLowerCase()) ||
          a.category.toLowerCase().includes(searchQuery.toLowerCase()) ||
          (a.ip_address ?? '').includes(searchQuery),
      )
    : activities;

  const getContextDisplay = (ctx: Record<string, unknown> | null): string => {
    if (!ctx) return '—';
    const keys = Object.keys(ctx);
    if (keys.length === 0) return '—';
    return keys
      .slice(0, 2)
      .map((k) => `${k}: ${ctx[k]}`)
      .join(', ');
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-white flex items-center gap-2.5">
          <ShieldAlert className="text-purple-400" size={22} />
          Audit Trail
        </h1>
        <p className="text-sm text-zinc-400">
          Immutable activity log — {totalItems.toLocaleString()} events tracked
        </p>
      </div>

      <div className="flex flex-col gap-4 md:flex-row md:items-center">
        <div className="flex-1 relative">
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-500" size={18} />
          <input
            type="text"
            placeholder="Search message, category, IP..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-11 pr-4 py-2.5 rounded-xl bg-white/[0.03] border border-white/5 text-zinc-200 placeholder-zinc-500 focus:outline-none focus:border-purple-500/30 transition-all text-sm"
          />
        </div>
      </div>

      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="rounded-2xl bg-white/[0.02] border border-white/5 backdrop-blur-sm overflow-hidden"
      >
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead>
              <tr className="border-b border-white/5 text-xs font-semibold uppercase tracking-wider text-zinc-400">
                <th className="px-6 py-4">Timestamp</th>
                <th className="px-6 py-4">Category</th>
                <th className="px-6 py-4">Level</th>
                <th className="px-6 py-4">Message</th>
                <th className="px-6 py-4">Context</th>
                <th className="px-6 py-4">IP Address</th>
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
              ) : displayed.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-6 py-12 text-center text-zinc-500">
                    No activity events found.
                  </td>
                </tr>
              ) : (
                displayed.map((item) => (
                  <tr key={item.id} className="hover:bg-white/[0.01] transition-colors">
                    <td className="px-6 py-4 text-xs font-mono text-zinc-500 whitespace-nowrap">
                      {new Date(item.created_at).toLocaleString()}
                    </td>
                    <td className="px-6 py-4">
                      <span className="text-xs font-semibold text-purple-400 bg-purple-500/10 px-2 py-0.5 rounded-md">
                        {item.category}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <span
                        className={`flex items-center gap-1 px-2 py-0.5 rounded-md text-xs font-semibold border w-fit ${
                          LEVEL_STYLES[item.level] ?? LEVEL_STYLES.DEBUG
                        }`}
                      >
                        {LEVEL_ICONS[item.level]}
                        {item.level}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-zinc-200 max-w-xs truncate">{item.message}</td>
                    <td className="px-6 py-4 text-xs text-zinc-500 max-w-xs truncate">
                      {getContextDisplay(item.context)}
                    </td>
                    <td className="px-6 py-4 text-xs font-mono text-zinc-500">
                      {item.ip_address || '—'}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination Footer */}
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
            <span>of {totalItems.toLocaleString()} events</span>
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
    </div>
  );
}
