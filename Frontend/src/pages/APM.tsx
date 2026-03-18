import { useState, useEffect, useCallback } from 'react';
import { motion } from 'framer-motion';
import { Activity, Clock, AlertTriangle, TrendingUp, RefreshCw, Gauge } from 'lucide-react';
import SelectField from '../shared/components/SelectField';
import { apmApi, sourcesApi } from '../shared/lib/api';

interface EndpointStat {
  endpoint: string;
  count: number;
  p50: number;
  p95: number;
  p99: number;
}

interface Source {
  id: string;
  name: string;
}

const PERIODS = [
  { label: '1 Hour', value: '1h' },
  { label: '24 Hours', value: '24h' },
  { label: '7 Days', value: '7d' },
];

function formatMs(ms: number | null | undefined): string {
  if (ms == null || isNaN(ms)) return '—';
  if (ms < 1) return '<1 ms';
  return `${Math.round(ms)} ms`;
}

function latencyBadge(ms: number | null | undefined): string {
  if (ms == null || isNaN(ms)) return 'text-zinc-500';
  if (ms >= 2000) return 'text-red-400 font-semibold';
  if (ms >= 1000) return 'text-amber-400 font-semibold';
  if (ms >= 500) return 'text-yellow-400';
  return 'text-emerald-400';
}

export default function APM() {
  const [period, setPeriod] = useState('24h');
  const [sourceId, setSourceId] = useState('');
  const [stats, setStats] = useState<EndpointStat[]>([]);
  const [sources, setSources] = useState<Source[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [lastRefresh, setLastRefresh] = useState<Date | null>(null);

  const fetchStats = useCallback(async () => {
    setIsLoading(true);
    try {
      const { data } = await apmApi.endpointStats({
        period,
        source_id: sourceId || undefined,
      });
      const rows = data.data;
      setStats(Array.isArray(rows) ? rows : []);
      setLastRefresh(new Date());
    } catch (err) {
      console.error('Failed to fetch APM stats', err);
      setStats([]);
    } finally {
      setIsLoading(false);
    }
  }, [period, sourceId]);

  useEffect(() => {
    fetchStats();
  }, [fetchStats]);

  useEffect(() => {
    sourcesApi.getAll().then(({ data }) => setSources(data.data ?? [])).catch(console.error);
  }, []);

  const safeStats = Array.isArray(stats) ? stats : [];
  const avgP50 = safeStats.length ? safeStats.reduce((s, r) => s + (r.p50 ?? 0), 0) / safeStats.length : null;
  const avgP95 = safeStats.length ? safeStats.reduce((s, r) => s + (r.p95 ?? 0), 0) / safeStats.length : null;
  const avgP99 = safeStats.length ? safeStats.reduce((s, r) => s + (r.p99 ?? 0), 0) / safeStats.length : null;
  const totalRequests = safeStats.reduce((s, r) => s + (r.count ?? 0), 0);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white tracking-tight flex items-center gap-2.5">
            <Gauge size={22} className="text-purple-400" />
            APM
          </h1>
          <p className="text-sm text-zinc-400 mt-0.5">Application Performance Monitoring — endpoint latency percentiles</p>
        </div>
        <div className="flex items-center gap-3">
          {lastRefresh && (
            <span className="text-xs text-zinc-600">
              Updated {lastRefresh.toLocaleTimeString()}
            </span>
          )}
          <button
            onClick={fetchStats}
            disabled={isLoading}
            className="flex items-center gap-1.5 px-3 py-2 text-xs font-medium rounded-lg bg-white/5 border border-white/10 text-zinc-300 hover:bg-white/10 transition-all disabled:opacity-50"
          >
            <RefreshCw size={13} className={isLoading ? 'animate-spin' : ''} />
            Refresh
          </button>
        </div>
      </div>

      {/* Controls */}
      <div className="flex flex-col sm:flex-row gap-3">
        <div className="flex rounded-xl overflow-hidden border border-white/10 bg-white/2">
          {PERIODS.map((p) => (
            <button
              key={p.value}
              onClick={() => setPeriod(p.value)}
              className={`px-4 py-2 text-sm font-medium transition-all ${
                period === p.value
                  ? 'bg-purple-500/20 text-purple-300 border-r border-purple-500/20'
                  : 'text-zinc-400 hover:text-zinc-200 border-r border-white/5'
              } last:border-r-0`}
            >
              {p.label}
            </button>
          ))}
        </div>
        <SelectField
          value={sourceId}
          onChange={(e) => setSourceId(e.target.value)}
          wrapperClassName="min-w-[160px]"
        >
          <option value="">All Sources</option>
          {sources.map((s) => (
            <option key={s.id} value={s.id}>{s.name}</option>
          ))}
        </SelectField>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        {[
          { label: 'Total Requests', value: totalRequests.toLocaleString(), icon: TrendingUp, color: 'text-purple-400' },
          { label: 'Avg P50 Latency', value: formatMs(avgP50), icon: Clock, color: 'text-emerald-400' },
          { label: 'Avg P95 Latency', value: formatMs(avgP95), icon: Clock, color: avgP95 != null && avgP95 >= 1000 ? 'text-amber-400' : 'text-blue-400' },
          { label: 'Avg P99 Latency', value: formatMs(avgP99), icon: AlertTriangle, color: avgP99 != null && avgP99 >= 1000 ? 'text-red-400' : 'text-zinc-300' },
        ].map((card) => (
          <motion.div
            key={card.label}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            className="p-4 rounded-2xl bg-white/2 border border-white/5"
          >
            <div className="flex items-center gap-2 mb-2">
              <card.icon size={15} className={card.color} />
              <span className="text-xs text-zinc-500 font-medium">{card.label}</span>
            </div>
            <p className={`text-2xl font-bold tracking-tight ${card.color}`}>{card.value}</p>
          </motion.div>
        ))}
      </div>

      {/* Table */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="rounded-2xl bg-white/2 border border-white/5 overflow-hidden"
      >
        {isLoading ? (
          <div className="flex items-center justify-center h-48 text-zinc-500 text-sm gap-2">
            <RefreshCw size={16} className="animate-spin" />
            Loading performance data...
          </div>
        ) : safeStats.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-48 text-center px-6">
            <div className="p-3 rounded-2xl bg-purple-500/10 text-purple-400 mb-3">
              <Activity size={28} />
            </div>
            <p className="text-zinc-300 font-semibold">No performance data found</p>
            <p className="text-zinc-500 text-sm mt-1 max-w-sm">
              Send logs with <code className="bg-zinc-800 text-purple-300 px-1 rounded text-xs">category: PERFORMANCE</code> and a{' '}
              <code className="bg-zinc-800 text-purple-300 px-1 rounded text-xs">duration_ms</code> field in the context object.
            </p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead>
                <tr className="border-b border-white/5 text-xs font-semibold uppercase tracking-wider text-zinc-400">
                  <th className="px-6 py-4">Endpoint</th>
                  <th className="px-6 py-4 text-right">Requests</th>
                  <th className="px-6 py-4 text-right">P50</th>
                  <th className="px-6 py-4 text-right">P95</th>
                  <th className="px-6 py-4 text-right">P99</th>
                </tr>
              </thead>
              <tbody>
                {safeStats.map((row, i) => (
                  <motion.tr
                    key={row.endpoint}
                    initial={{ opacity: 0, x: -10 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: i * 0.03 }}
                    className="border-b border-white/3 hover:bg-white/2 transition-colors"
                  >
                    <td className="px-6 py-3.5">
                      <code className="text-sm text-zinc-200 font-mono">{row.endpoint}</code>
                    </td>
                    <td className="px-6 py-3.5 text-right text-sm text-zinc-400">
                      {(row.count ?? 0).toLocaleString()}
                    </td>
                    <td className={`px-6 py-3.5 text-right text-sm ${latencyBadge(row.p50)}`}>
                      {formatMs(row.p50)}
                    </td>
                    <td className={`px-6 py-3.5 text-right text-sm ${latencyBadge(row.p95)}`}>
                      {formatMs(row.p95)}
                    </td>
                    <td className={`px-6 py-3.5 text-right text-sm ${latencyBadge(row.p99)}`}>
                      {formatMs(row.p99)}
                    </td>
                  </motion.tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </motion.div>

      {/* Legend */}
      {safeStats.length > 0 && (
        <div className="flex flex-wrap gap-4 text-xs text-zinc-500">
          <span><span className="text-emerald-400">●</span> &lt;500ms — Fast</span>
          <span><span className="text-yellow-400">●</span> 500–999ms — Acceptable</span>
          <span><span className="text-amber-400">●</span> 1000–1999ms — Slow</span>
          <span><span className="text-red-400">●</span> ≥2000ms — Critical</span>
        </div>
      )}
    </div>
  );
}
