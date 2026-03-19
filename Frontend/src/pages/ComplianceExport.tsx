import { useState } from 'react';
import { motion } from 'framer-motion';
import { FileDown, Download, Clock, CheckCircle, AlertCircle } from 'lucide-react';
import { useQuery, useMutation } from '@tanstack/react-query';
import { activityApi } from '../shared/lib/api';
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

export default function ComplianceExport() {
  const [sourceId, setSourceId] = useState('');
  const [format, setFormat] = useState('PDF');
  const [dateFrom, setDateFrom] = useState('');
  const [dateTo, setDateTo] = useState('');

  const { data: exports, isLoading } = useQuery({
    queryKey: ['compliance-exports'],
    queryFn: () => activityApi.getComplianceExports({ page: 1, limit: 20 }),
  });

  const exportList: ComplianceExport[] = exports?.data?.items || [];

  const requestExport = useMutation({
    mutationFn: (data: { source_id: string; format: string; date_from: string; date_to: string }) =>
      activityApi.requestComplianceExport(data),
  });

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
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="p-8 space-y-6"
    >
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-white flex items-center gap-3">
          <FileDown className="w-7 h-7 text-purple-400" />
          Compliance Export
        </h1>
        <p className="text-zinc-400 mt-1">
          Export audit trails for compliance and regulatory requirements
        </p>
      </div>

      {/* Export Request Form */}
      <div className="bg-[#0c0c0c] border border-white/10 rounded-lg p-6">
        <h2 className="text-lg font-semibold text-white mb-4">Request New Export</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm text-zinc-400 mb-2">Source</label>
            <SelectField
              value={sourceId}
              onChange={(e) => setSourceId(e.target.value)}
            >
              <option value="">Select Source</option>
              <option value="all">All Sources</option>
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
      <div className="bg-[#0c0c0c] border border-white/10 rounded-lg overflow-hidden">
        <div className="px-6 py-4 border-b border-white/10">
          <h2 className="text-lg font-semibold text-white">Export History</h2>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-white/5">
              <tr>
                <th className="text-left py-3 px-6 text-sm font-medium text-zinc-400">ID</th>
                <th className="text-left py-3 px-6 text-sm font-medium text-zinc-400">Source</th>
                <th className="text-left py-3 px-6 text-sm font-medium text-zinc-400">Format</th>
                <th className="text-left py-3 px-6 text-sm font-medium text-zinc-400">Date Range</th>
                <th className="text-left py-3 px-6 text-sm font-medium text-zinc-400">Status</th>
                <th className="text-left py-3 px-6 text-sm font-medium text-zinc-400">Created</th>
                <th className="text-left py-3 px-6 text-sm font-medium text-zinc-400">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {isLoading ? (
                <tr>
                  <td colSpan={7} className="py-8 text-center text-zinc-400">
                    Loading exports...
                  </td>
                </tr>
              ) : exportList.length === 0 ? (
                <tr>
                  <td colSpan={7} className="py-8 text-center text-zinc-400">
                    No exports found
                  </td>
                </tr>
              ) : (
                exportList.map((item) => (
                  <tr key={item.id} className="hover:bg-white/5 transition-colors">
                    <td className="py-3 px-6 text-sm text-zinc-400">#{item.id}</td>
                    <td className="py-3 px-6 text-sm text-white">{item.source_id}</td>
                    <td className="py-3 px-6 text-sm text-zinc-300">{item.format}</td>
                    <td className="py-3 px-6 text-sm text-zinc-400">
                      {new Date(item.date_from).toLocaleDateString()} - {new Date(item.date_to).toLocaleDateString()}
                    </td>
                    <td className="py-3 px-6">
                      <div className="flex items-center gap-2">
                        {getStatusIcon(item.status)}
                        <span className={`text-sm capitalize ${getStatusColor(item.status)}`}>
                          {item.status}
                        </span>
                      </div>
                    </td>
                    <td className="py-3 px-6 text-sm text-zinc-400">
                      {new Date(item.created_at).toLocaleString()}
                    </td>
                    <td className="py-3 px-6">
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
    </motion.div>
  );
}
