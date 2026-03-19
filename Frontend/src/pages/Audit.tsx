import { useState } from 'react';
import { ShieldAlert, Search, ChevronLeft, ChevronRight } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { activityApi } from '../shared/lib/api';
import SelectField from '../shared/components/SelectField';
import { categoryLabel } from '../shared/lib/utils';

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

export default function Audit() {
  const [currentPage, setCurrentPage] = useState(1);
  const [limit, setLimit] = useState(20);
  const [searchQuery, setSearchQuery] = useState('');

  // Fetch activities using React Query
  const { data, isLoading, error } = useQuery({
    queryKey: ['audit-logs', currentPage, limit],
    queryFn: async () => {
      const { data } = await activityApi.list({ page: currentPage, limit });
      return data.data;
    },
  });

  const activities: ActivityLog[] = data?.items || [];
  const totalItems = data?.meta?.total || 0;
  const maxPage = Math.max(1, Math.ceil(totalItems / limit));

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
      {/* Header */}
      <div className="flex items-center gap-3">
        <div className="p-2 rounded-xl bg-purple-500/10 text-purple-400">
          <ShieldAlert size={24} />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-white">Audit Trail</h1>
          <p className="text-sm text-zinc-400">Immutable activity log — {totalItems.toLocaleString()} events tracked</p>
        </div>
      </div>

      <div className="flex flex-col gap-4 md:flex-row md:items-center">
        <div className="flex-1 relative">
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-500" size={18} />
          <input
            type="text"
            placeholder="Search message, category, IP..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-11 pr-4 py-2.5 rounded-xl bg-white/3 border border-white/5 text-zinc-200 placeholder-zinc-500 focus:outline-none focus:border-purple-500/30 transition-all text-sm"
          />
        </div>
      </div>

      <div className="bg-white/2 border border-white/5 rounded-xl overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead className="bg-white/5">
              <tr>
                <th className="px-4 py-3 text-xs font-semibold uppercase text-zinc-500">Time</th>
                <th className="px-4 py-3 text-xs font-semibold uppercase text-zinc-500">Category</th>
                <th className="px-4 py-3 text-xs font-semibold uppercase text-zinc-500">Level</th>
                <th className="px-4 py-3 text-xs font-semibold uppercase text-zinc-500">Message</th>
                <th className="px-4 py-3 text-xs font-semibold uppercase text-zinc-500">Source</th>
                <th className="px-4 py-3 text-xs font-semibold uppercase text-zinc-500">IP Address</th>
                <th className="px-4 py-3 text-xs font-semibold uppercase text-zinc-500">Context</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {isLoading ? (
                <tr>
                  <td colSpan={7} className="px-4 py-8 text-center text-zinc-400">
                    <div className="flex items-center justify-center gap-2">
                      <div className="w-5 h-5 border-2 border-purple-500/30 border-t-purple-500 rounded-full animate-spin" />
                      Loading audit logs...
                    </div>
                  </td>
                </tr>
              ) : error ? (
                <tr>
                  <td colSpan={7} className="px-4 py-8 text-center text-red-400">
                    Failed to load audit logs. Please try again.
                  </td>
                </tr>
              ) : displayed.length === 0 ? (
                <tr>
                  <td colSpan={7} className="px-4 py-8 text-center text-zinc-500">
                    {searchQuery ? 'No matching audit logs found' : 'No audit logs available'}
                  </td>
                </tr>
              ) : (
                displayed.map((activity) => (
                  <tr key={activity.id} className="hover:bg-white/5 transition-colors">
                    <td className="px-4 py-3 text-sm text-zinc-400 whitespace-nowrap">
                      {new Date(activity.created_at).toLocaleString()}
                    </td>
                    <td className="px-4 py-3 text-sm text-zinc-300">
                      {categoryLabel(activity.category)}
                    </td>
                    <td className="px-4 py-3">
                      <span
                        className={`inline-flex items-center px-2 py-1 rounded-md text-xs font-medium border ${
                          LEVEL_STYLES[activity.level] || LEVEL_STYLES.INFO
                        }`}
                      >
                        {activity.level}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-sm text-zinc-200 max-w-md truncate" title={activity.message}>
                      {activity.message}
                    </td>
                    <td className="px-4 py-3 text-sm text-zinc-400 font-mono">
                      {activity.source_id}
                    </td>
                    <td className="px-4 py-3 text-sm text-zinc-400 font-mono">
                      {activity.ip_address || '—'}
                    </td>
                    <td className="px-4 py-3 text-sm text-zinc-500 max-w-xs truncate">
                      {getContextDisplay(activity.context)}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {!searchQuery && totalItems > 0 && (
          <div className="flex items-center justify-between px-4 py-3 border-t border-white/5">
            <div className="flex items-center gap-2 text-sm text-zinc-400">
              <span>Rows per page:</span>
              <SelectField
                value={limit.toString()}
                onChange={(e) => {
                  setLimit(Number(e.target.value));
                  setCurrentPage(1);
                }}
                className="w-20"
              >
                <option value="10">10</option>
                <option value="20">20</option>
                <option value="50">50</option>
                <option value="100">100</option>
              </SelectField>
              <span className="text-zinc-500">
                Showing {((currentPage - 1) * limit) + 1} - {Math.min(currentPage * limit, totalItems)} of {totalItems}
              </span>
            </div>
            <div className="flex items-center gap-2">
              <button
                onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                disabled={currentPage === 1}
                className="p-2 rounded-lg bg-white/5 text-zinc-400 hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                <ChevronLeft size={18} />
              </button>
              <span className="text-sm text-zinc-400">
                Page {currentPage} of {maxPage}
              </span>
              <button
                onClick={() => setCurrentPage((p) => Math.min(maxPage, p + 1))}
                disabled={currentPage >= maxPage}
                className="p-2 rounded-lg bg-white/5 text-zinc-400 hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                <ChevronRight size={18} />
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
