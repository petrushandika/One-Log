import { useState } from 'react';
import { motion } from 'framer-motion';
import { 
  Shield, 
  Users, 
  RefreshCw,
  PieChart,
  Calendar,
  Grid3X3
} from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { activityApi } from '../shared/lib/api';
import SelectField from '../shared/components/SelectField';
import {
  PieChart as RePieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  Tooltip,
  Legend,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
} from 'recharts';

interface AuthMethodData {
  [key: string]: number;
}

interface TimelineData {
  date: string;
  login_success: number;
  login_failed: number;
}

interface HeatmapData {
  day_of_week: number;
  hour_of_day: number;
  failed_count: number;
}

interface Session {
  id: number;
  user_id: string;
  auth_method: string;
  ip_address: string;
  browser: string;
  device: string;
  created_at: string;
  last_activity: string;
}

const DAYS = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
const COLORS = ['#8b5cf6', '#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#6b7280'];

export default function ActivityAnalytics() {
  const [activeTab, setActiveTab] = useState<'methods' | 'timeline' | 'heatmap' | 'sessions'>('methods');
  const [days, setDays] = useState(30);
  const [sessionPage, setSessionPage] = useState(1);
  const [sessionLimit, setSessionLimit] = useState(20);

  // Fetch auth methods breakdown
  const { data: methodsData, isLoading: methodsLoading } = useQuery({
    queryKey: ['auth-methods', days],
    queryFn: async () => {
      const { data } = await activityApi.getAuthMethodBreakdown({ days });
      return data.data as AuthMethodData;
    },
  });

  // Fetch login timeline
  const { data: timelineData, isLoading: timelineLoading } = useQuery({
    queryKey: ['login-timeline', days],
    queryFn: async () => {
      const { data } = await activityApi.getLoginTimeline({ days });
      return data.data as TimelineData[];
    },
  });

  // Fetch failed login heatmap
  const { data: heatmapData, isLoading: heatmapLoading } = useQuery({
    queryKey: ['failed-login-heatmap', days],
    queryFn: async () => {
      const { data } = await activityApi.getFailedLoginHeatmap({ days });
      return data.data as HeatmapData[];
    },
  });

  // Fetch recent sessions
  const { data: sessionsData, isLoading: sessionsLoading } = useQuery({
    queryKey: ['sessions', sessionPage, sessionLimit],
    queryFn: async () => {
      const { data } = await activityApi.getSessions({ 
        page: sessionPage, 
        limit: sessionLimit 
      });
      return data.data;
    },
  });

  // Transform methods data for pie chart
  const pieChartData = methodsData 
    ? Object.entries(methodsData).map(([name, value]) => ({ name, value }))
    : [];

  // Transform heatmap data
  const heatmapMatrix = Array(7).fill(null).map((_, day) => 
    Array(24).fill(null).map((_, hour) => {
      const point = heatmapData?.find(
        (d) => d.day_of_week === day && d.hour_of_day === hour
      );
      return {
        day,
        hour,
        value: point?.failed_count || 0,
      };
    })
  );

  const maxFailedCount = Math.max(
    1,
    ...(heatmapData?.map((d) => d.failed_count) || [0])
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        <div className="p-2 rounded-xl bg-purple-500/10 text-purple-400">
          <Shield size={24} />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-white">Activity Analytics</h1>
          <p className="text-sm text-zinc-400">Monitor authentication and user activity</p>
        </div>
      </div>

      {/* Controls */}
      <div className="flex items-center gap-4 flex-wrap">
        <div className="flex rounded-xl overflow-hidden border border-white/10 bg-white/2">
          {[
            { id: 'methods', label: 'Auth Methods', icon: PieChart },
            { id: 'timeline', label: 'Timeline', icon: Calendar },
            { id: 'heatmap', label: 'Heatmap', icon: Grid3X3 },
            { id: 'sessions', label: 'Sessions', icon: Users },
          ].map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as 'methods' | 'timeline' | 'heatmap' | 'sessions')}
              className={`flex items-center gap-2 px-4 py-2 text-sm font-medium transition-all ${
                activeTab === tab.id
                  ? 'bg-purple-500/20 text-purple-400 border-r border-purple-500/20'
                  : 'text-zinc-400 hover:text-zinc-200 border-r border-white/5'
              } last:border-r-0`}
            >
              <tab.icon size={16} />
              {tab.label}
            </button>
          ))}
        </div>

        <SelectField
          value={days}
          onChange={(e) => setDays(Number(e.target.value))}
          wrapperClassName="w-32"
        >
          <option value={7}>7 Days</option>
          <option value={30}>30 Days</option>
          <option value={90}>90 Days</option>
        </SelectField>
      </div>

      {/* Auth Methods Tab */}
      {activeTab === 'methods' && (
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className="grid grid-cols-1 lg:grid-cols-2 gap-6"
        >
          <div className="bg-white/[0.02] border border-white/5 rounded-xl p-6">
            <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
              <PieChart size={20} className="text-purple-400" />
              Authentication Methods Breakdown
            </h3>
            {methodsLoading ? (
              <div className="h-64 flex items-center justify-center">
                <RefreshCw className="animate-spin text-zinc-500" size={24} />
              </div>
            ) : pieChartData.length === 0 ? (
              <div className="h-64 flex flex-col items-center justify-center text-zinc-500">
                <PieChart size={48} className="mb-2 opacity-30" />
                <p>No authentication data available</p>
              </div>
            ) : (
              <div className="h-64">
                <ResponsiveContainer width="100%" height="100%">
                  <RePieChart>
                    <Pie
                      data={pieChartData}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      label={({ name, percent }) => `${name}: ${((percent || 0) * 100).toFixed(0)}%`}
                      outerRadius={80}
                      fill="#8884d8"
                      dataKey="value"
                    >
                      {pieChartData.map((_entry, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip 
                      contentStyle={{ 
                        backgroundColor: '#18181b', 
                        border: '1px solid rgba(255,255,255,0.1)',
                        borderRadius: '8px'
                      }}
                    />
                    <Legend />
                  </RePieChart>
                </ResponsiveContainer>
              </div>
            )}
          </div>

          <div className="bg-white/[0.02] border border-white/5 rounded-xl p-6">
            <h3 className="text-lg font-semibold text-white mb-4">Method Statistics</h3>
            {methodsLoading ? (
              <div className="space-y-3">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="h-12 bg-white/5 rounded-lg animate-pulse" />
                ))}
              </div>
            ) : pieChartData.length === 0 ? (
              <p className="text-zinc-500 text-center py-8">No data available</p>
            ) : (
              <div className="space-y-3">
                {pieChartData
                  .sort((a, b) => b.value - a.value)
                  .map((method, index) => (
                    <div
                      key={method.name}
                      className="flex items-center justify-between p-3 rounded-lg bg-white/5"
                    >
                      <div className="flex items-center gap-3">
                        <div
                          className="w-3 h-3 rounded-full"
                          style={{ backgroundColor: COLORS[index % COLORS.length] }}
                        />
                        <span className="text-white font-medium">{method.name}</span>
                      </div>
                      <span className="text-zinc-400">{method.value.toLocaleString()} logins</span>
                    </div>
                  ))}
              </div>
            )}
          </div>
        </motion.div>
      )}

      {/* Timeline Tab */}
      {activeTab === 'timeline' && (
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-white/[0.02] border border-white/5 rounded-xl p-6"
        >
          <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <Calendar size={20} className="text-purple-400" />
            Login Timeline (Last {days} Days)
          </h3>
          {timelineLoading ? (
            <div className="h-80 flex items-center justify-center">
              <RefreshCw className="animate-spin text-zinc-500" size={24} />
            </div>
          ) : timelineData?.length === 0 ? (
            <div className="h-80 flex flex-col items-center justify-center text-zinc-500">
              <Calendar size={48} className="mb-2 opacity-30" />
              <p>No login data available</p>
            </div>
          ) : (
            <div className="h-80">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={timelineData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.05)" />
                  <XAxis 
                    dataKey="date" 
                    stroke="#52525b"
                    fontSize={12}
                    tickFormatter={(value) => new Date(value).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
                  />
                  <YAxis stroke="#52525b" fontSize={12} />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: '#18181b',
                      border: '1px solid rgba(255,255,255,0.1)',
                      borderRadius: '8px',
                    }}
                  />
                  <Legend />
                  <Bar dataKey="login_success" name="Success" fill="#10b981" />
                  <Bar dataKey="login_failed" name="Failed" fill="#ef4444" />
                </BarChart>
              </ResponsiveContainer>
            </div>
          )}
        </motion.div>
      )}

      {/* Heatmap Tab */}
      {activeTab === 'heatmap' && (
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-white/[0.02] border border-white/5 rounded-xl p-6"
        >
          <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <Grid3X3 size={20} className="text-purple-400" />
            Failed Login Heatmap (Last {days} Days)
          </h3>
          {heatmapLoading ? (
            <div className="h-96 flex items-center justify-center">
              <RefreshCw className="animate-spin text-zinc-500" size={24} />
            </div>
          ) : heatmapData?.length === 0 ? (
            <div className="h-96 flex flex-col items-center justify-center text-zinc-500">
              <Grid3X3 size={48} className="mb-2 opacity-30" />
              <p>No failed login data available</p>
            </div>
          ) : (
            <div className="space-y-4">
              {/* Heatmap Grid */}
              <div className="overflow-x-auto">
                <div className="min-w-[600px]">
                  {/* Header - Hours */}
                  <div className="flex">
                    <div className="w-16" /> {/* Day label spacer */}
                    {Array.from({ length: 24 }, (_, i) => (
                      <div
                        key={i}
                        className="flex-1 text-center text-xs text-zinc-500 py-1"
                      >
                        {i}
                      </div>
                    ))}
                  </div>
                  
                  {/* Grid */}
                  {DAYS.map((day, dayIndex) => (
                    <div key={day} className="flex items-center">
                      <div className="w-16 text-sm text-zinc-400 py-1">{day}</div>
                      <div className="flex flex-1">
                        {Array.from({ length: 24 }, (_, hour) => {
                          const value = heatmapMatrix[dayIndex]?.[hour]?.value || 0;
                          const intensity = value / maxFailedCount;
                          return (
                            <div
                              key={hour}
                              className="flex-1 h-8 border border-white/5"
                              style={{
                                backgroundColor: value > 0 
                                  ? `rgba(239, 68, 68, ${0.1 + intensity * 0.9})`
                                  : 'rgba(255,255,255,0.02)',
                              }}
                              title={`${day} ${hour}:00 - ${value} failed attempts`}
                            />
                          );
                        })}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
              
              {/* Legend */}
              <div className="flex items-center gap-4 justify-center text-sm text-zinc-400">
                <span>Low</span>
                <div className="flex gap-1">
                  {[0.1, 0.3, 0.5, 0.7, 0.9].map((opacity, i) => (
                    <div
                      key={i}
                      className="w-6 h-4 rounded"
                      style={{ backgroundColor: `rgba(239, 68, 68, ${opacity})` }}
                    />
                  ))}
                </div>
                <span>High</span>
              </div>
            </div>
          )}
        </motion.div>
      )}

      {/* Sessions Tab */}
      {activeTab === 'sessions' && (
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className="space-y-4"
        >
          <div className="bg-white/[0.02] border border-white/5 rounded-xl overflow-hidden">
            <table className="w-full">
              <thead className="bg-white/[0.03] border-b border-white/5">
                <tr>
                  <th className="text-left text-xs font-semibold uppercase text-zinc-500 px-4 py-3">User</th>
                  <th className="text-left text-xs font-semibold uppercase text-zinc-500 px-4 py-3">Auth Method</th>
                  <th className="text-left text-xs font-semibold uppercase text-zinc-500 px-4 py-3">IP Address</th>
                  <th className="text-left text-xs font-semibold uppercase text-zinc-500 px-4 py-3">Browser</th>
                  <th className="text-left text-xs font-semibold uppercase text-zinc-500 px-4 py-3">Device</th>
                  <th className="text-left text-xs font-semibold uppercase text-zinc-500 px-4 py-3">Last Activity</th>
                </tr>
              </thead>
              <tbody>
                {sessionsLoading ? (
                  <tr>
                    <td colSpan={6} className="text-center py-8">
                      <RefreshCw className="animate-spin mx-auto text-zinc-500" size={24} />
                    </td>
                  </tr>
                ) : sessionsData?.items?.length === 0 ? (
                  <tr>
                    <td colSpan={6} className="text-center py-8 text-zinc-500">
                      <Users size={32} className="mx-auto mb-2 opacity-30" />
                      <p>No active sessions found</p>
                    </td>
                  </tr>
                ) : (
                  sessionsData?.items?.map((session: Session) => (
                    <tr key={session.id} className="border-b border-white/5 hover:bg-white/[0.02]">
                      <td className="px-4 py-3 text-sm text-white font-medium">
                        {session.user_id}
                      </td>
                      <td className="px-4 py-3">
                        <span className="px-2 py-1 rounded-full text-xs bg-purple-500/10 text-purple-400">
                          {session.auth_method}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-sm text-zinc-400 font-mono">
                        {session.ip_address}
                      </td>
                      <td className="px-4 py-3 text-sm text-zinc-400">
                        {session.browser || 'Unknown'}
                      </td>
                      <td className="px-4 py-3 text-sm text-zinc-400">
                        {session.device || 'Unknown'}
                      </td>
                      <td className="px-4 py-3 text-sm text-zinc-400">
                        {new Date(session.last_activity).toLocaleString()}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>

          {sessionsData?.meta && (
            <div className="flex items-center justify-between px-4 py-3 bg-white/[0.02] border border-white/5 rounded-xl">
              <div className="flex items-center gap-2 text-sm text-zinc-400">
                <span>Show</span>
                <select
                  value={sessionLimit}
                  onChange={(e) => {
                    setSessionLimit(Number(e.target.value));
                    setSessionPage(1);
                  }}
                  className="bg-white/3 border border-white/8 text-zinc-200 rounded-lg px-2 py-1 text-sm"
                >
                  <option value={10}>10</option>
                  <option value={20}>20</option>
                  <option value={50}>50</option>
                </select>
                <span>of {sessionsData.meta.total} items</span>
              </div>
              <div className="flex items-center gap-2">
                <button
                  onClick={() => setSessionPage((p) => Math.max(1, p - 1))}
                  disabled={sessionPage === 1}
                  className="px-3 py-1.5 text-sm rounded-lg bg-white/5 text-zinc-300 hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Previous
                </button>
                <span className="text-sm text-zinc-400">
                  Page {sessionPage} of {Math.ceil(sessionsData.meta.total / sessionLimit)}
                </span>
                <button
                  onClick={() => setSessionPage((p) => Math.min(Math.ceil(sessionsData.meta.total / sessionLimit), p + 1))}
                  disabled={sessionPage === Math.ceil(sessionsData.meta.total / sessionLimit)}
                  className="px-3 py-1.5 text-sm rounded-lg bg-white/5 text-zinc-300 hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Next
                </button>
              </div>
            </div>
          )}
        </motion.div>
      )}
    </div>
  );
}
