import { motion } from 'framer-motion';
import { AreaChart, Area, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';
import { ArrowUpRight, ArrowDownRight, Activity, Terminal, Shield, AlertCircle } from 'lucide-react';

export default function Overview() {
  const stats = [
    { name: 'Total Logs (24h)', value: '1,234', change: '+12%', positive: true, icon: Terminal, color: 'text-purple-400', bg: 'bg-purple-500/10', border: 'border-purple-500/20' },
    { name: 'System Errors', value: '42', change: '-4%', positive: true, icon: AlertCircle, color: 'text-rose-400', bg: 'bg-rose-500/10', border: 'border-rose-500/20' },
    { name: 'Active Sources', value: '5', change: '0%', positive: true, icon: Activity, color: 'text-emerald-400', bg: 'bg-emerald-500/10', border: 'border-emerald-500/20' },
    { name: 'Security Alerts', value: '3', change: '+1', positive: false, icon: Shield, color: 'text-amber-400', bg: 'bg-amber-500/10', border: 'border-amber-500/20' },
  ];

  const data = [
    { time: '00:00', logs: 400, errors: 24 },
    { time: '04:00', logs: 300, errors: 13 },
    { time: '08:00', logs: 200, errors: 10 },
    { time: '12:00', logs: 278, errors: 6 },
    { time: '16:00', logs: 189, errors: 4 },
    { time: '20:00', logs: 239, errors: 12 },
    { time: '24:00', logs: 349, errors: 18 },
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold tracking-tight text-white">Overview</h1>
        <p className="text-sm text-zinc-400">Everything happening on your systems</p>
      </div>

      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat, i) => (
          <motion.div
            key={stat.name}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.1, duration: 0.4 }}
            className={`p-6 rounded-2xl bg-white/[0.02] border border-white/[0.05] backdrop-blur-sm box-border flex flex-col justify-between`}
          >
            <div className="flex items-start justify-between">
              <div className={`p-2 rounded-xl ${stat.bg} ${stat.color} border ${stat.border}`}>
                <stat.icon size={20} />
              </div>
              <span className={`text-xs font-medium flex items-center gap-1 ${stat.positive ? 'text-emerald-400' : 'text-rose-400'}`}>
                {stat.change}
                {stat.positive ? <ArrowUpRight size={14} /> : <ArrowDownRight size={14} />}
              </span>
            </div>
            <div className="mt-4">
              <p className="text-sm text-zinc-400">{stat.name}</p>
              <p className="text-3xl font-bold tracking-tight text-zinc-100 mt-1">{stat.value}</p>
            </div>
          </motion.div>
        ))}
      </div>

      <motion.div
        initial={{ opacity: 0, y: 30 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.4, duration: 0.5 }}
        className="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.05] backdrop-blur-sm"
      >
        <h2 className="text-lg font-semibold text-zinc-100 mb-6">Log Ingestion Trend</h2>
        <div className="h-80 w-full">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={data}>
              <defs>
                <linearGradient id="colorLogs" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#a855f7" stopOpacity={0.3}/>
                  <stop offset="95%" stopColor="#a855f7" stopOpacity={0}/>
                </linearGradient>
                <linearGradient id="colorErrors" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#f43f5e" stopOpacity={0.3}/>
                  <stop offset="95%" stopColor="#f43f5e" stopOpacity={0}/>
                </linearGradient>
              </defs>
              <XAxis dataKey="time" stroke="#525252" fontSize={12} tickLine={false} axisLine={false} />
              <YAxis stroke="#525252" fontSize={12} tickLine={false} axisLine={false} />
              <Tooltip 
                contentStyle={{ background: '#09090b', border: '1px solid rgba(255,255,255,0.08)', borderRadius: '12px' }} 
                labelStyle={{ color: '#a1a1aa' }}
              />
              <Area type="monotone" dataKey="logs" name="Total Logs" stroke="#a855f7" strokeWidth={2} fillOpacity={1} fill="url(#colorLogs)" />
              <Area type="monotone" dataKey="errors" name="Errors" stroke="#f43f5e" strokeWidth={2} fillOpacity={1} fill="url(#colorErrors)" />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </motion.div>
    </div>
  );
}
