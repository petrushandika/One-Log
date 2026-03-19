import { useState } from 'react';
import { FileDown, Download, Clock, CheckCircle, AlertCircle, RefreshCw } from 'lucide-react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { activityApi, sourcesApi } from '../shared/lib/api';
import SelectField from '../shared/components/SelectField';

interface ComplianceExport {
  id: number;
  source_id: string;
  format: string;
  date_from: string;
  date_to: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  file_url: string;
  created_by: string;
  created_at: string;
}

interface Source {
  id: string;
  name: string;
}

export default function ComplianceExport() {
  const [sourceId, setSourceId] = useState('');
  const [format, setFormat] = useState('PDF');
  const [dateFrom, setDateFrom] = useState('');
  const [dateTo, setDateTo] = useState('');
  const queryClient = useQueryClient();

  // Fetch compliance exports
  const { data: exports, isLoading: exportsLoading } = useQuery({
    queryKey: ['compliance-exports'],
    queryFn: () => activityApi.getComplianceExports({ page: 1, limit: 20 }),
  });

  // Fetch sources for dropdown
  const { data: sourcesData, isLoading: sourcesLoading } = useQuery({
    queryKey: ['sources'],
    queryFn: () => sourcesApi.getAll(),
  });

  const exportList: ComplianceExport[] = exports?.data?.items || [];
  const sources: Source[] = sourcesData?.data?.data || [];

  const requestExport = useMutation({
    mutationFn: (data: { source_id: string; format: string; date_from: string; date_to: string }) =>
      activityApi.requestComplianceExport(data),
    onSuccess: () => {
      // Refresh exports list after successful request
      queryClient.invalidateQueries({ queryKey: ['compliance-exports'] });
    },
  });

  const isLoading = exportsLoading || sourcesLoading;

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return <CheckCircle className="w-4 h-4 text-green-400" />;
      case 'failed':
        return <AlertCircle className="w-4 h-4 text-red-400" />;
      case 'processing':
        return <Clock className="w-4 h-4 text-yellow-400" />;
      default:
        return <Clock className="w-4 h-4 text-zinc-400" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'text-green-400';
      case 'failed':
        return 'text-red-400';
      case 'processing':
        return 'text-yellow-400';
      default:
        return 'text-zinc-400';
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        <div className="p-2 rounded-xl bg-purple-500/10 text-purple-400">
          <FileDown size={24} />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-white">Compliance Export</h1>
          <p className="text-sm text-zinc-400">Export audit trails for compliance and regulatory requirements</p>
        </div>
      </div>

      {/* Export Request Form */}
      <div className="bg-white/2 border border-white/5 rounded-xl p-6">
        <h2 className="text-lg font-semibold text-white mb-4">Request New Export</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm text-zinc-400 mb-2">Source</label>
            <SelectField
              value={sourceId}
              onChange={(e) => setSourceId(e.target.value)}
              disabled={sourcesLoading}
            >
              <option value="">{sourcesLoading ? 'Loading sources...' : 'Select Source'}</option>
              <option value="all">All Sources</option>
              {sources.map((source) => (
                <option key={source.id} value={source.id}>
                  {source.name}
                </option>
              ))}
            </SelectField>
          </div>
          <div>
            <label className="block text-sm text-zinc-400 mb-2">Format</label>
            <SelectField
              value={format}
              onChange={(e) => setFormat(e.target.value)}
            >
              <option value="PDF">PDF</option>
              <option value="CSV">CSV</option>
            </SelectField>
          </div>
          <div>
            <label className="block text-sm text-zinc-400 mb-2">Date From</label>
            <input
              type="date"
              value={dateFrom}
              onChange={(e) => setDateFrom(e.target.value)}
              className="w-full bg-white/3 border border-white/8 text-zinc-200 rounded-xl px-3 py-2 text-sm focus:outline-none focus:border-purple-500/40"
            />
          </div>
          <div>
            <label className="block text-sm text-zinc-400 mb-2">Date To</label>
            <input
              type="date"
              value={dateTo}
              onChange={(e) => setDateTo(e.target.value)}
              className="w-full bg-white/3 border border-white/8 text-zinc-200 rounded-xl px-3 py-2 text-sm focus:outline-none focus:border-purple-500/40"
            />
          </div>
        </div>
        <button
          onClick={() => {
            if (dateFrom && dateTo) {
              requestExport.mutate({
                source_id: sourceId || 'all',
                format,
                date_from: dateFrom,
                date_to: dateTo,
              });
            }
          }}
          disabled={!dateFrom || !dateTo || requestExport.isPending}
          className="mt-4 px-6 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {requestExport.isPending ? 'Requesting...' : 'Request Export'}
        </button>
      </div>

      {/* Export History */}
      <div className="bg-white/2 border border-white/5 rounded-xl overflow-hidden">
        <div className="px-6 py-4 border-b border-white/5">
          <h2 className="text-lg font-semibold text-white">Export History</h2>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-white/5">
              <tr>
                <th className="text-left px-4 py-3 text-xs font-semibold uppercase text-zinc-500">ID</th>
                <th className="text-left px-4 py-3 text-xs font-semibold uppercase text-zinc-500">Source</th>
                <th className="text-left px-4 py-3 text-xs font-semibold uppercase text-zinc-500">Format</th>
                <th className="text-left px-4 py-3 text-xs font-semibold uppercase text-zinc-500">Date Range</th>
                <th className="text-left px-4 py-3 text-xs font-semibold uppercase text-zinc-500">Status</th>
                <th className="text-left px-4 py-3 text-xs font-semibold uppercase text-zinc-500">Created</th>
                <th className="text-left px-4 py-3 text-xs font-semibold uppercase text-zinc-500">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {isLoading ? (
                <tr>
                  <td colSpan={7} className="px-4 py-8 text-center">
                    <RefreshCw size={24} className="animate-spin text-zinc-400 mx-auto" />
                  </td>
                </tr>
              ) : exportList.length === 0 ? (
                <tr>
                  <td colSpan={7} className="px-4 py-8 text-center text-zinc-400">
                    No exports found
                  </td>
                </tr>
              ) : (
                exportList.map((item) => (
                  <tr key={item.id} className="hover:bg-white/5 transition-colors">
                    <td className="px-4 py-3 text-sm text-zinc-400">#{item.id}</td>
                    <td className="px-4 py-3 text-sm text-white">{item.source_id}</td>
                    <td className="px-4 py-3 text-sm text-zinc-300">{item.format}</td>
                    <td className="px-4 py-3 text-sm text-zinc-400">
                      {new Date(item.date_from).toLocaleDateString()} - {new Date(item.date_to).toLocaleDateString()}
                    </td>
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-2">
                        {getStatusIcon(item.status)}
                        <span className={`text-sm capitalize ${getStatusColor(item.status)}`}>
                          {item.status}
                        </span>
                      </div>
                    </td>
                    <td className="px-4 py-3 text-sm text-zinc-400">
                      {new Date(item.created_at).toLocaleString()}
                    </td>
                    <td className="px-4 py-3">
                      {item.status === 'completed' && item.file_url && (
                        <button
                          onClick={() => window.open(item.file_url, '_blank')}
                          className="flex items-center gap-1 text-sm text-purple-400 hover:text-purple-300"
                        >
                          <Download className="w-4 h-4" />
                          Download
                        </button>
                      )}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
