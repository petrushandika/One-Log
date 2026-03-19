import { useState } from 'react';
import { motion } from 'framer-motion';
import { Users, Activity, FileText, Download, Filter } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { activityApi } from '../shared/lib/api';
import SelectField from '../shared/components/SelectField';

interface ActivityFeedItem {
  id: number;
  user_id: string;
  source_id: string;
  action: string;
  resource_type: string;
  resource_id: string;
  context: Record<string, unknown>;
  ip_address: string;
  created_at: string;
}

export default function ActivityFeed() {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(20);
  const [action, setAction] = useState('');

  const { data, isLoading, error } = useQuery({
    queryKey: ['activity-feed', page, limit, action],
    queryFn: () => activityApi.getFeed({ page, limit, action }),
  });

  const feedItems: ActivityFeedItem[] = data?.data?.items || [];
  const total = data?.data?.meta?.total || 0;
  const totalPages = Math.ceil(total / limit);

  const getActionIcon = (actionType: string) => {
    switch (actionType) {
      case 'create':
        return <FileText className="w-4 h-4 text-green-400" />;
      case 'update':
        return <Activity className="w-4 h-4 text-blue-400" />;
      case 'delete':
        return <Activity className="w-4 h-4 text-red-400" />;
      case 'export':
        return <Download className="w-4 h-4 text-purple-400" />;
      default:
        return <Activity className="w-4 h-4 text-zinc-400" />;
    }
  };

  const getActionColor = (actionType: string) => {
    switch (actionType) {
      case 'create':
        return 'text-green-400';
      case 'update':
        return 'text-blue-400';
      case 'delete':
        return 'text-red-400';
      case 'export':
        return 'text-purple-400';
      default:
        return 'text-zinc-400';
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="p-8 space-y-6"
    >
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white flex items-center gap-3">
            <Users className="w-7 h-7 text-purple-400" />
            Activity Feed
          </h1>
          <p className="text-zinc-400 mt-1">
            Track all user activities across your applications
          </p>
        </div>
        <div className="text-sm text-zinc-400">
          Total Activities: {total.toLocaleString()}
        </div>
      </div>

      {/* Filters */}
      <div className="flex gap-4 items-center bg-[#0c0c0c] border border-white/10 rounded-lg p-4">
        <div className="flex items-center gap-2 text-zinc-400">
          <Filter className="w-4 h-4" />
          <span className="text-sm">Filters:</span>
        </div>
        <SelectField
          value={action}
          onChange={(e) => setAction(e.target.value)}
          className="w-40"
        >
          <option value="">All Actions</option>
          <option value="create">Create</option>
          <option value="update">Update</option>
          <option value="delete">Delete</option>
          <option value="export">Export</option>
          <option value="view">View</option>
        </SelectField>
      </div>

      {/* Activity Feed Table */}
      <div className="bg-[#0c0c0c] border border-white/10 rounded-lg overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-white/5">
              <tr>
                <th className="text-left py-3 px-4 text-sm font-medium text-zinc-400">Action</th>
                <th className="text-left py-3 px-4 text-sm font-medium text-zinc-400">User</th>
                <th className="text-left py-3 px-4 text-sm font-medium text-zinc-400">Resource</th>
                <th className="text-left py-3 px-4 text-sm font-medium text-zinc-400">Source</th>
                <th className="text-left py-3 px-4 text-sm font-medium text-zinc-400">IP Address</th>
                <th className="text-left py-3 px-4 text-sm font-medium text-zinc-400">Time</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {isLoading ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-zinc-400">
                    Loading activities...
                  </td>
                </tr>
              ) : error ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-red-400">
                    Failed to load activities
                  </td>
                </tr>
              ) : feedItems.length === 0 ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-zinc-400">
                    No activities found
                  </td>
                </tr>
              ) : (
                feedItems.map((item) => (
                  <tr key={item.id} className="hover:bg-white/5 transition-colors">
                    <td className="py-3 px-4">
                      <div className="flex items-center gap-2">
                        {getActionIcon(item.action)}
                        <span className={`capitalize text-sm ${getActionColor(item.action)}`}>
                          {item.action}
                        </span>
                      </div>
                    </td>
                    <td className="py-3 px-4 text-sm text-white font-mono">
                      {item.user_id}
                    </td>
                    <td className="py-3 px-4">
                      <div className="text-sm text-white">{item.resource_type}</div>
                      <div className="text-xs text-zinc-500 font-mono">{item.resource_id}</div>
                    </td>
                    <td className="py-3 px-4 text-sm text-zinc-300">
                      {item.source_id}
                    </td>
                    <td className="py-3 px-4 text-sm text-zinc-400 font-mono">
                      {item.ip_address}
                    </td>
                    <td className="py-3 px-4 text-sm text-zinc-400">
                      {new Date(item.created_at).toLocaleString()}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        <div className="flex items-center justify-between px-4 py-3 border-t border-white/10">
          <div className="flex items-center gap-2 text-sm text-zinc-400">
            <span>Rows per page:</span>
            <SelectField
              value={limit.toString()}
              onChange={(e) => {
                setLimit(Number(e.target.value));
                setPage(1);
              }}
              className="w-20"
            >
              <option value="10">10</option>
              <option value="20">20</option>
              <option value="50">50</option>
              <option value="100">100</option>
            </SelectField>
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={() => setPage(p => Math.max(1, p - 1))}
              disabled={page === 1}
              className="px-3 py-1 text-sm bg-white/5 text-zinc-300 rounded hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Previous
            </button>
            <span className="text-sm text-zinc-400">
              Page {page} of {totalPages}
            </span>
            <button
              onClick={() => setPage(p => Math.min(totalPages, p + 1))}
              disabled={page >= totalPages}
              className="px-3 py-1 text-sm bg-white/5 text-zinc-300 rounded hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </motion.div>
  );
}
