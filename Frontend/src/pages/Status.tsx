import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { Radio, CheckCircle, AlertTriangle, XCircle, Clock, RefreshCw, ExternalLink, Signal } from 'lucide-react';
import { statusApi } from '../shared/lib/api';

interface Source {
  id: string;
  name: string;
  status: string;
  health_url: string;
  updated_at: string;
}

const STATUS_CONFIG: Record<string, { label: string; color: string; bg: string; border: string; icon: typeof CheckCircle }> = {
  ONLINE:      { label: 'Online',      color: 'text-emerald-400', bg: 'bg-emerald-500/10', border: 'border-emerald-500/20', icon: CheckCircle },
  DEGRADED:    { label: 'Degraded',    color: 'text-amber-400',   bg: 'bg-amber-500/10',   border: 'border-amber-500/20',   icon: AlertTriangle },
  OFFLINE:     { label: 'Offline',     color: 'text-red-400',     bg: 'bg-red-500/10',     border: 'border-red-500/20',     icon: XCircle },
  MAINTENANCE: { label: 'Maintenance', color: 'text-zinc-400',    bg: 'bg-zinc-500/10',    border: 'border-zinc-500/20',    icon: Clock },
};

function formatRelative(dateStr: string): string {
  if (!dateStr) return '—';
  const diff = Math.floor((Date.now() - new Date(dateStr).getTime()) / 1000);
  if (diff < 60) return `${diff}s ago`;
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
  return new Date(dateStr).toLocaleDateString();
}

export default function Status() {
  const [sources, setSources] = useState<Source[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [lastRefresh, setLastRefresh] = useState<Date | null>(null);

  const fetchStatus = async () => {
    setIsLoading(true);
    try {
      const { data } = await statusApi.getPublic();
      setSources(data.data?.sources ?? []);
      setLastRefresh(new Date());
    } catch (err) {
      console.error('Failed to fetch status', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchStatus();
    const interval = setInterval(fetchStatus, 60_000);
    return () => clearInterval(interval);
  }, []);

  const onlineCount = sources.filter((s) => s.status === 'ONLINE').length;
  const degradedCount = sources.filter((s) => s.status === 'DEGRADED').length;
  const offlineCount = sources.filter((s) => s.status === 'OFFLINE').length;
  const allOperational = offlineCount === 0 && degradedCount === 0;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white tracking-tight flex items-center gap-2.5">
            <Signal size={22} className="text-purple-400" />
            Status Page
          </h1>
          <p className="text-sm text-zinc-400 mt-0.5">Real-time health status of all registered sources</p>
        </div>
        <div className="flex items-center gap-3">
          {lastRefresh && (
            <span className="text-xs text-zinc-600">Auto-refreshes every 60s · Last: {lastRefresh.toLocaleTimeString()}</span>
          )}
          <button
            onClick={fetchStatus}
            disabled={isLoading}
            className="flex items-center gap-1.5 px-3 py-2 text-xs font-medium rounded-lg bg-white/5 border border-white/10 text-zinc-300 hover:bg-white/10 transition-all disabled:opacity-50"
          >
            <RefreshCw size={13} className={isLoading ? 'animate-spin' : ''} />
            Refresh
          </button>
        </div>
      </div>

      {/* Overall Banner */}
      {sources.length > 0 && (
        <motion.div
          initial={{ opacity: 0, y: -8 }}
          animate={{ opacity: 1, y: 0 }}
          className={`flex items-center gap-3 px-5 py-4 rounded-2xl border ${
            allOperational
              ? 'bg-emerald-500/10 border-emerald-500/20'
              : offlineCount > 0
              ? 'bg-red-500/10 border-red-500/20'
              : 'bg-amber-500/10 border-amber-500/20'
          }`}
        >
          {allOperational ? (
            <CheckCircle size={20} className="text-emerald-400 shrink-0" />
          ) : offlineCount > 0 ? (
            <XCircle size={20} className="text-red-400 shrink-0" />
          ) : (
            <AlertTriangle size={20} className="text-amber-400 shrink-0" />
          )}
          <div>
            <p className={`font-semibold text-sm ${allOperational ? 'text-emerald-300' : offlineCount > 0 ? 'text-red-300' : 'text-amber-300'}`}>
              {allOperational
                ? 'All systems operational'
                : offlineCount > 0
                ? `${offlineCount} source${offlineCount > 1 ? 's' : ''} offline`
                : `${degradedCount} source${degradedCount > 1 ? 's' : ''} degraded`}
            </p>
            <p className="text-xs text-zinc-400 mt-0.5">
              {onlineCount} online · {degradedCount} degraded · {offlineCount} offline · {sources.length} total
            </p>
          </div>
        </motion.div>
      )}

      {/* Source Grid */}
      {isLoading && sources.length === 0 ? (
        <div className="flex items-center justify-center h-48 text-zinc-500 text-sm gap-2">
          <RefreshCw size={16} className="animate-spin" />
          Loading status...
        </div>
      ) : sources.length === 0 ? (
        <div className="flex flex-col items-center justify-center h-48 text-center">
          <div className="p-3 rounded-2xl bg-purple-500/10 text-purple-400 mb-3">
            <Radio size={28} />
          </div>
          <p className="text-zinc-300 font-semibold">No sources registered</p>
          <p className="text-zinc-500 text-sm mt-1">Add a source from the Sources page to track its health.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {sources.map((source, i) => {
            const cfg = STATUS_CONFIG[source.status] ?? STATUS_CONFIG.MAINTENANCE;
            const Icon = cfg.icon;
            return (
              <motion.div
                key={source.id}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.05 }}
                className="p-5 rounded-2xl bg-white/[0.02] border border-white/[0.05] hover:bg-white/[0.03] transition-colors"
              >
                <div className="flex items-start justify-between gap-3 mb-3">
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-semibold text-zinc-100 truncate">{source.name}</p>
                    {source.health_url && (
                      <a
                        href={source.health_url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center gap-1 text-xs text-zinc-500 hover:text-purple-400 transition-colors mt-0.5 truncate"
                      >
                        <ExternalLink size={10} />
                        {source.health_url}
                      </a>
                    )}
                  </div>
                  <span className={`flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold shrink-0 border ${cfg.bg} ${cfg.color} ${cfg.border}`}>
                    <Icon size={11} />
                    {cfg.label}
                  </span>
                </div>
                <div className="flex items-center gap-1.5 text-xs text-zinc-600">
                  <Clock size={11} />
                  <span>Checked {formatRelative(source.updated_at)}</span>
                </div>
              </motion.div>
            );
          })}
        </div>
      )}
    </div>
  );
}
