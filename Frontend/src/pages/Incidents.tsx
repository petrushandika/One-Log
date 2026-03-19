import { useState } from 'react';
import { AlertTriangle, CheckCircle, Clock, Calendar } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { incidentsApi } from '../shared/lib/api';
import SelectField from '../shared/components/SelectField';
import { formatDistanceToNow } from 'date-fns';

interface Incident {
  id: number;
  source_id: string;
  status: 'OPEN' | 'RESOLVED';
  started_at: string;
  resolved_at?: string;
  duration_sec: number;
  message: string;
}

export default function Incidents() {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);
  const [statusFilter, setStatusFilter] = useState('');

  // Fetch incidents
  const { data: incidentsData, isLoading } = useQuery({
    queryKey: ['incidents', page, limit, statusFilter],
    queryFn: async () => {
      const { data } = await incidentsApi.list({
        page,
        limit,
        status: statusFilter || undefined,
      });
      return data.data;
    },
  });

  const incidents = incidentsData?.items || [];
  const total = incidentsData?.meta?.total || 0;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        <div className="p-2 rounded-xl bg-red-500/10 text-red-400">
          <AlertTriangle size={24} />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-white">Incidents</h1>
          <p className="text-sm text-zinc-400">Track and monitor system downtime incidents</p>
        </div>
      </div>

      {/* Stats Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-white/[0.02] border border-white/5 rounded-xl p-4">
          <div className="flex items-center gap-2 text-red-400 mb-2">
            <AlertTriangle size={18} />
            <span className="text-sm font-medium">Open Incidents</span>
          </div>
          <p className="text-2xl font-bold text-white">
            {incidents.filter((i: Incident) => i.status === 'OPEN').length}
          </p>
        </div>
        
        <div className="bg-white/[0.02] border border-white/5 rounded-xl p-4">
          <div className="flex items-center gap-2 text-emerald-400 mb-2">
            <CheckCircle size={18} />
            <span className="text-sm font-medium">Total Incidents</span>
          </div>
          <p className="text-2xl font-bold text-white">{total}</p>
        </div>
        
        <div className="bg-white/[0.02] border border-white/5 rounded-xl p-4">
          <div className="flex items-center gap-2 text-orange-400 mb-2">
            <Clock size={18} />
            <span className="text-sm font-medium">Page</span>
          </div>
          <p className="text-2xl font-bold text-white">{page}</p>
        </div>
        
        <div className="bg-white/[0.02] border border-white/5 rounded-xl p-4">
          <div className="flex items-center gap-2 text-blue-400 mb-2">
            <Calendar size={18} />
            <span className="text-sm font-medium">Limit</span>
          </div>
          <p className="text-2xl font-bold text-white">{limit}</p>
        </div>
      </div>

      {/* Filters */}
      <div className="flex items-center gap-4">
        <div>
          <label className="block text-xs text-zinc-400 mb-1">Status</label>
          <SelectField
            value={statusFilter}
            onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
              setStatusFilter(e.target.value);
              setPage(1);
            }}
          >
            <option value="">All Status</option>
            <option value="OPEN">Open</option>
            <option value="RESOLVED">Resolved</option>
          </SelectField>
        </div>
      </div>

      {/* Incidents Table */}
      <div className="bg-white/[0.02] border border-white/5 rounded-xl overflow-hidden">
        <table className="w-full">
          <thead className="bg-white/[0.03] border-b border-white/5">
            <tr>
              <th className="text-left text-xs font-semibold uppercase text-zinc-500 px-4 py-3">Status</th>
              <th className="text-left text-xs font-semibold uppercase text-zinc-500 px-4 py-3">Message</th>
              <th className="text-left text-xs font-semibold uppercase text-zinc-500 px-4 py-3">Started</th>
              <th className="text-left text-xs font-semibold uppercase text-zinc-500 px-4 py-3">Resolved</th>
            </tr>
          </thead>
          <tbody>
            {isLoading ? (
              <tr>
                <td colSpan={4} className="text-center py-8 text-zinc-500">Loading incidents...</td>
              </tr>
            ) : incidents.length === 0 ? (
              <tr>
                <td colSpan={4} className="text-center py-8 text-zinc-500">
                  <div className="flex flex-col items-center gap-2">
                    <CheckCircle size={32} className="text-emerald-500" />
                    <p>No incidents found</p>
                  </div>
                </td>
              </tr>
            ) : (
              incidents.map((incident: Incident) => (
                <tr key={incident.id} className="border-b border-white/5 hover:bg-white/[0.02]">
                  <td className="px-4 py-3">
                    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium ${
                      incident.status === 'OPEN' 
                        ? 'bg-red-500/10 text-red-400' 
                        : 'bg-emerald-500/10 text-emerald-400'
                    }`}>
                      {incident.status === 'OPEN' ? (
                        <AlertTriangle size={12} />
                      ) : (
                        <CheckCircle size={12} />
                      )}
                      {incident.status}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-sm text-white max-w-md truncate">
                    {incident.message}
                  </td>
                  <td className="px-4 py-3 text-sm text-zinc-400">
                    {formatDistanceToNow(new Date(incident.started_at), { addSuffix: true })}
                  </td>
                  <td className="px-4 py-3 text-sm text-zinc-400">
                    {incident.resolved_at 
                      ? formatDistanceToNow(new Date(incident.resolved_at), { addSuffix: true })
                      : '-'
                    }
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      {total > 0 && (
        <div className="flex items-center justify-between px-4 py-3 bg-white/[0.02] border border-white/5 rounded-xl">
          <div className="flex items-center gap-2 text-sm text-zinc-400">
            <span>Show</span>
            <select
              value={limit}
              onChange={(e) => {
                setLimit(Number(e.target.value));
                setPage(1);
              }}
              className="bg-white/3 border border-white/8 text-zinc-200 rounded-lg px-2 py-1 text-sm"
            >
              <option value={10}>10</option>
              <option value={20}>20</option>
              <option value={50}>50</option>
            </select>
            <span>of {total} items</span>
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page === 1}
              className="px-3 py-1.5 text-sm rounded-lg bg-white/5 text-zinc-300 hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Previous
            </button>
            <span className="text-sm text-zinc-400">
              Page {page} of {Math.ceil(total / limit)}
            </span>
            <button
              onClick={() => setPage((p) => Math.min(Math.ceil(total / limit), p + 1))}
              disabled={page === Math.ceil(total / limit)}
              className="px-3 py-1.5 text-sm rounded-lg bg-white/5 text-zinc-300 hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
