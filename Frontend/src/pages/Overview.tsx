import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { AreaChart, Area, XAxis, YAxis, Tooltip, ResponsiveContainer, BarChart, Bar, Legend } from 'recharts';
import { ArrowUpRight, ArrowDownRight, Activity, Terminal, Shield, AlertCircle, RefreshCw, LayoutDashboard } from 'lucide-react';
import { statsApi, sourcesApi } from '../shared/lib/api';

interface StatsData {
  total: number;
  errors: number;
  active: number;
  security: number;
  CRITICAL?: number;
  ERROR?: number;
  WARN?: number;
  INFO?: number;
}

export default function Overview() {
  const [liveStats, setLiveStats] = useState<StatsData>({ total: 0, errors: 0, active: 0, security: 0 });
  const [isLoading, setIsLoading] = useState(true);
  const [lastRefresh, setLastRefresh] = useState(new Date());

  const fetchStats = async () => {
    setIsLoading(true);
    try {
      const [statsRes, sourcesRes] = await Promise.all([
        statsApi.getOverview(),
        sourcesApi.getAll(),
      ]);

      const d = statsRes.data?.data ?? {};
      const sources: { status: string }[] = sourcesRes.data?.data ?? [];
      const activeCount = sources.filter((s) => s.status === 'ONLINE').length;

      setLiveStats({
        total: d.total ?? 0,
        errors: (d.ERROR ?? 0) + (d.CRITICAL ?? 0),
        active: activeCount,
        security: d.SECURITY ?? d.CRITICAL ?? 0,
        CRITICAL: d.CRITICAL ?? 0,
        ERROR: d.ERROR ?? 0,
        WARN: d.WARN ?? 0,
        INFO: d.INFO ?? 0,
      });
      setLastRefresh(new Date());
    } catch (err) {
      console.error('Failed to fetch stats', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchStats();
  }, []);

  const statCards = [
    {
      name: 'Total Logs',
      value: liveStats.total.toLocaleString(),
      change: 'All time',
      positive: true,
      icon: Terminal,
      color: 'text-purple-400',
      bg: 'bg-purple-500/10',
      border: 'border-purple-500/20',
    },
    {
      name: 'Errors & Critical',
      value: liveStats.errors.toLocaleString(),
      change: liveStats.errors > 0 ? 'Needs attention' : 'All clear',
      positive: liveStats.errors === 0,
      icon: AlertCircle,
      color: 'text-rose-400',
      bg: 'bg-rose-500/10',
      border: 'border-rose-500/20',
    },
    {
      name: 'Online Sources',
      value: liveStats.active.toLocaleString(),
      change: 'Currently active',
      positive: true,
      icon: Activity,
      color: 'text-emerald-400',
      bg: 'bg-emerald-500/10',
      border: 'border-emerald-500/20',
    },
    {
      name: 'Security Alerts',
      value: liveStats.security.toLocaleString(),
      change: liveStats.security > 0 ? 'Review required' : 'No threats',
      positive: liveStats.security === 0,
      icon: Shield,
      color: 'text-amber-400',
      bg: 'bg-amber-500/10',
      border: 'border-amber-500/20',
    },
  ];

  const breakdownData = [
    { name: 'CRITICAL', count: liveStats.CRITICAL ?? 0, fill: '#f43f5e' },
    { name: 'ERROR', count: liveStats.ERROR ?? 0, fill: '#ef4444' },
    { name: 'WARN', count: liveStats.WARN ?? 0, fill: '#f59e0b' },
    { name: 'INFO', count: liveStats.INFO ?? 0, fill: '#10b981' },
  ];

  return (
    <div className="space-y-8">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-white flex items-center gap-2.5">
            <LayoutDashboard size={22} className="text-purple-400" />
            Overview
          </h1>
          <p className="text-sm text-zinc-400">Everything happening on your systems</p>
        </div>
        <button
          onClick={fetchStats}
          disabled={isLoading}
          className="flex items-center gap-2 px-3 py-2 text-xs rounded-xl bg-white/[0.03] border border-white/[0.06] text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.06] transition-all disabled:opacity-50"
        >
          <RefreshCw size={14} className={isLoading ? 'animate-spin' : ''} />
          {isLoading ? 'Loading...' : `Refreshed ${lastRefresh.toLocaleTimeString()}`}
        </button>
      </div>

      {/* Stat Cards */}
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {statCards.map((stat, i) => (
          <motion.div
            key={stat.name}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.1, duration: 0.4 }}
            className="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.05] backdrop-blur-sm flex flex-col justify-between"
          >
            <div className="flex items-start justify-between">
              <div className={`p-2 rounded-xl ${stat.bg} ${stat.color} border ${stat.border}`}>
                <stat.icon size={20} />
              </div>
              <span
                className={`text-xs font-medium flex items-center gap-1 ${
                  stat.positive ? 'text-emerald-400' : 'text-rose-400'
                }`}
              >
                {stat.change}
                {stat.positive ? <ArrowUpRight size={14} /> : <ArrowDownRight size={14} />}
              </span>
            </div>
            <div className="mt-4">
              <p className="text-sm text-zinc-400">{stat.name}</p>
              <p className="text-3xl font-bold tracking-tight text-zinc-100 mt-1">
                {isLoading ? (
                  <span className="inline-block w-16 h-8 rounded-lg bg-white/5 animate-pulse" />
                ) : (
                  stat.value
                )}
              </p>
            </div>
          </motion.div>
        ))}
      </div>

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Level Breakdown Bar Chart */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.45, duration: 0.5 }}
          className="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.05] backdrop-blur-sm"
        >
          <h2 className="text-base font-semibold text-zinc-100 mb-4">Logs by Level</h2>
          <div className="h-52">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={breakdownData} layout="vertical" margin={{ left: 8, right: 16 }}>
                <XAxis type="number" stroke="#525252" fontSize={11} tickLine={false} axisLine={false} />
                <YAxis type="category" dataKey="name" stroke="#525252" fontSize={11} tickLine={false} axisLine={false} width={60} />
                <Tooltip
                  cursor={{ fill: 'rgba(255,255,255,0.03)' }}
                  contentStyle={{ background: '#09090b', border: '1px solid rgba(255,255,255,0.08)', borderRadius: '10px' }}
                  labelStyle={{ color: '#a1a1aa' }}
                />
                <Bar dataKey="count" name="Count" radius={[0, 6, 6, 0]}>
                  {breakdownData.map((entry, idx) => (
                    <rect key={idx} fill={entry.fill} />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>
        </motion.div>

        {/* Trend Chart (static placeholder — real chart needs time-series endpoint) */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5, duration: 0.5 }}
          className="lg:col-span-2 p-6 rounded-2xl bg-white/[0.02] border border-white/[0.05] backdrop-blur-sm"
        >
          <h2 className="text-base font-semibold text-zinc-100 mb-4">Log Ingestion Trend</h2>
          <div className="h-52 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart
                data={[
                  { time: '00:00', logs: liveStats.total > 0 ? Math.round(liveStats.total * 0.12) : 0, errors: liveStats.errors > 0 ? Math.round(liveStats.errors * 0.1) : 0 },
                  { time: '04:00', logs: liveStats.total > 0 ? Math.round(liveStats.total * 0.08) : 0, errors: liveStats.errors > 0 ? Math.round(liveStats.errors * 0.08) : 0 },
                  { time: '08:00', logs: liveStats.total > 0 ? Math.round(liveStats.total * 0.15) : 0, errors: liveStats.errors > 0 ? Math.round(liveStats.errors * 0.15) : 0 },
                  { time: '12:00', logs: liveStats.total > 0 ? Math.round(liveStats.total * 0.2) : 0, errors: liveStats.errors > 0 ? Math.round(liveStats.errors * 0.2) : 0 },
                  { time: '16:00', logs: liveStats.total > 0 ? Math.round(liveStats.total * 0.22) : 0, errors: liveStats.errors > 0 ? Math.round(liveStats.errors * 0.22) : 0 },
                  { time: '20:00', logs: liveStats.total > 0 ? Math.round(liveStats.total * 0.18) : 0, errors: liveStats.errors > 0 ? Math.round(liveStats.errors * 0.18) : 0 },
                  { time: '24:00', logs: liveStats.total > 0 ? Math.round(liveStats.total * 0.05) : 0, errors: liveStats.errors > 0 ? Math.round(liveStats.errors * 0.07) : 0 },
                ]}
              >
                <defs>
                  <linearGradient id="colorLogs" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#a855f7" stopOpacity={0.3} />
                    <stop offset="95%" stopColor="#a855f7" stopOpacity={0} />
                  </linearGradient>
                  <linearGradient id="colorErrors" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#f43f5e" stopOpacity={0.3} />
                    <stop offset="95%" stopColor="#f43f5e" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <XAxis dataKey="time" stroke="#525252" fontSize={11} tickLine={false} axisLine={false} />
                <YAxis stroke="#525252" fontSize={11} tickLine={false} axisLine={false} />
                <Tooltip
                  contentStyle={{ background: '#09090b', border: '1px solid rgba(255,255,255,0.08)', borderRadius: '12px' }}
                  labelStyle={{ color: '#a1a1aa' }}
                />
                <Legend wrapperStyle={{ fontSize: '12px', color: '#71717a' }} />
                <Area type="monotone" dataKey="logs" name="Total Logs" stroke="#a855f7" strokeWidth={2} fillOpacity={1} fill="url(#colorLogs)" />
                <Area type="monotone" dataKey="errors" name="Errors" stroke="#f43f5e" strokeWidth={2} fillOpacity={1} fill="url(#colorErrors)" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </motion.div>
      </div>
    </div>
  );
}
