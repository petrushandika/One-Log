import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Lock, Mail, AlertCircle, Eye, EyeOff, ShieldCheck, Activity } from 'lucide-react';
import { authApi } from '../shared/lib/api';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);

    try {
      const { data } = await authApi.login({ email, password });
      const token = data.data?.token || data.data;
      if (token) {
        localStorage.setItem('token', token);
        navigate('/');
      } else {
        setError('Failed to retrieve token. Please try again.');
      }
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any) {
      setError(err.response?.data?.message || 'Login failed. Please check your credentials.');
    } finally {
      setIsLoading(false);
    }
  };

  // Creative Mock Streaming Layout triggers set boards
  const mockStreams = [
    { time: '14:55:01', label: 'Auth Service', msg: 'Admin Login initialized', level: 'SUCCESS', color: 'text-emerald-400' },
    { time: '14:55:02', label: 'Gateway Node 1', msg: 'Route check matching /api/logs', level: 'INFO', color: 'text-purple-400' },
    { time: '14:55:04', label: 'Datastore', msg: 'Connection pool warm (92%)', level: 'WARNING', color: 'text-amber-400' },
    { time: '14:55:05', label: 'Aggregator', msg: 'Dumping 1.4k nodes to memory', level: 'INFO', color: 'text-purple-400' },
  ];

  return (
    <div className="min-h-screen grid grid-cols-1 lg:grid-cols-2 bg-[#040405] text-white">
      {/* Left Column: Branding (Desktop Only) */}
      <div className="hidden lg:flex flex-col items-center justify-center p-12 bg-gradient-to-br from-[#09090b] to-[#040405] relative border-r border-white/[0.03] overflow-hidden">
        {/* Glow Spheres */}
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-purple-500/10 rounded-full blur-3xl -z-10 animate-pulse" />
        <div className="absolute bottom-1/4 right-1/4 w-80 h-80 bg-fuchsia-500/10 rounded-full blur-3xl -z-10" />

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.8 }}
          className="max-w-lg text-center flex flex-col items-center relative w-full"
        >
          <h1 className="text-4xl font-extrabold tracking-tight bg-gradient-to-r from-white via-zinc-200 to-zinc-500 bg-clip-text text-transparent">
            One Log Central
          </h1>
          <p className="text-zinc-400 mt-3 text-base">
            Consolidate high-throughput log streams, trace microservices and monitor status endpoints layouts.
          </p>

          {/* Creative Feature: Animated Live Stream Window triggers set boards */}
          <div className="w-full mt-10 rounded-2xl bg-zinc-950/60 border border-white/[0.04] backdrop-blur-md overflow-hidden shadow-2xl flex flex-col">
            <div className="p-3 bg-zinc-900/50 border-b border-white/[0.03] flex items-center justify-between">
              <div className="flex items-center gap-1.5">
                <div className="w-3 h-3 rounded-full bg-red-500/80" />
                <div className="w-3 h-3 rounded-full bg-amber-500/80" />
                <div className="w-3 h-3 rounded-full bg-green-500/80" />
              </div>
              <div className="flex items-center gap-1 text-emerald-400 text-xs font-semibold">
                <div className="w-1.5 h-1.5 bg-emerald-400 rounded-full animate-ping" />
                Live Network
              </div>
            </div>
            
            <motion.div 
              initial={{ opacity: 0, y: 30 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.3 }}
              className="p-4 font-mono text-left text-[11px] space-y-2 flex flex-col leading-relaxed min-h-[160px]"
            >
              {mockStreams.map((row, idx) => (
                <motion.div 
                  key={idx}
                  initial={{ opacity: 0, x: -10 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: 0.5 + idx * 0.15 }}
                  className="flex items-center gap-2 text-zinc-400 tracking-tight"
                >
                  <span className="text-zinc-600">{row.time}</span>
                  <span className="px-1.5 py-0.5 rounded bg-white/[0.03] text-zinc-300 font-semibold">{row.label}</span>
                  <span className={row.color}>{row.msg}</span>
                </motion.div>
              ))}
              <motion.div 
                animate={{ opacity: [1, 0, 1] }}
                transition={{ repeat: Infinity, duration: 0.7 }}
                className="w-1.5 h-3.5 bg-purple-400 mt-1"
              />
            </motion.div>
          </div>

          <div className="grid grid-cols-2 gap-3 mt-6 w-full max-w-sm">
            <div className="p-4 rounded-xl bg-white/[0.01] border border-white/[0.03] flex items-center gap-3">
              <ShieldCheck className="text-purple-400 shrink-0" size={22} />
              <div className="text-left">
                <p className="text-xs font-bold text-zinc-200">Autopilot</p>
                <p className="text-[10px] text-zinc-500">Anomaly alerts</p>
              </div>
            </div>
            <div className="p-4 rounded-xl bg-white/[0.01] border border-white/[0.03] flex items-center gap-3">
              <Activity className="text-emerald-400 shrink-0" size={22} />
              <div className="text-left">
                <p className="text-xs font-bold text-zinc-200">Realtime</p>
                <p className="text-[10px] text-zinc-500">Live Metric sync</p>
              </div>
            </div>
          </div>
        </motion.div>
      </div>

      {/* Right Column: Login Form */}
      <div className="flex flex-col items-center justify-center p-6 md:p-12 relative overflow-hidden">
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-80 h-80 bg-purple-500/5 rounded-full blur-3xl lg:hidden -z-10" />

        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="w-full max-w-sm"
        >
          <div className="flex flex-col mb-8 text-center lg:text-left">
            <h2 className="text-3xl font-bold tracking-tight text-white flex items-center gap-2 justify-center lg:justify-start">
              Welcome back
            </h2>
            <p className="text-sm text-zinc-400 mt-1.5">Sign in to manage your system stats</p>
          </div>

          <div className="p-7 md:p-8 rounded-3xl bg-zinc-900/40 border border-white/[0.04] backdrop-blur-xl shadow-2xl flex flex-col gap-5">
            {error && (
              <motion.div
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                className="p-3.5 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-400 text-sm flex items-center gap-2"
              >
                <AlertCircle size={18} className="shrink-0" />
                <span>{error}</span>
              </motion.div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <label className="text-xs font-medium text-zinc-400 ml-1">Email Address</label>
                <div className="relative">
                  <Mail className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-500" size={18} />
                  <input
                    type="email"
                    required
                    placeholder="admin@example.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="w-full pl-11 pr-4 py-2.5 rounded-xl bg-zinc-800/40 border border-white/[0.05] text-zinc-100 placeholder-zinc-600 focus:outline-none focus:border-purple-500/30 focus:bg-zinc-800/60 transition-all text-sm"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <label className="text-xs font-medium text-zinc-400 ml-1">Password</label>
                <div className="relative">
                  <Lock className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-500" size={18} />
                  <input
                    type={showPassword ? 'text' : 'password'}
                    required
                    placeholder="••••••••"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full pl-11 pr-12 py-2.5 rounded-xl bg-zinc-800/40 border border-white/[0.05] text-zinc-100 placeholder-zinc-600 focus:outline-none focus:border-purple-500/30 focus:bg-zinc-800/60 transition-all text-sm"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-4 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-300 transition-colors focus:outline-none"
                  >
                    {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                  </button>
                </div>
              </div>

              <button
                type="submit"
                disabled={isLoading}
                className="w-full mt-3 py-2.5 px-4 rounded-xl bg-gradient-to-r from-purple-500 to-fuchsia-500 hover:from-purple-600 hover:to-fuchsia-600 hover:shadow-purple-500/20 disabled:opacity-60 text-white font-semibold shadow-lg shadow-purple-500/10 transition-all duration-300 flex items-center justify-center gap-2 group text-sm disabled:cursor-not-allowed"
              >
                {isLoading ? (
                  <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                ) : (
                  <>
                    Sign In
                  </>
                )}
              </button>
            </form>
          </div>
        </motion.div>
      </div>
    </div>
  );
}
