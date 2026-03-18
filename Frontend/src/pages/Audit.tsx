import { useState } from 'react';
import { motion } from 'framer-motion';
import { ShieldAlert, Search, ChevronLeft, ChevronRight } from 'lucide-react';

export default function Audit() {
  const [itemsPerPage, setItemsPerPage] = useState<number | 'all'>(10);
  const [currentPage, setCurrentPage] = useState(1);

  const audits = [
    { id: 1, time: '2026-03-18 14:15:22', user: 'Admin', action: 'Regnerated API Key', target: 'Auth Service', ip: '192.168.1.100' },
    { id: 2, time: '2026-03-18 14:02:10', user: 'Admin', action: 'Update Source Name', target: 'Gateway', ip: '192.168.1.100' },
    { id: 3, time: '2026-03-18 13:45:05', user: 'System', action: 'Source disconnected', target: 'DB Analytics', ip: '-' },
    { id: 4, time: '2026-03-18 13:30:00', user: 'Admin', action: 'User login successful', target: 'Auth', ip: '192.168.1.5' },
  ];

  const totalLogs = audits.length;
  const maxPage = itemsPerPage === 'all' ? 1 : Math.ceil(totalLogs / itemsPerPage);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight text-white flex items-center gap-2">
          <ShieldAlert className="text-purple-400" size={28} />
          Audit Trail
        </h1>
        <p className="text-sm text-zinc-400">Track administrative and system events history</p>
      </div>

      <div className="flex flex-col gap-4 md:flex-row md:items-center">
        <div className="flex-1 relative">
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-500" size={18} />
          <input
            type="text"
            placeholder="Search action, users, IP..."
            className="w-full pl-11 pr-4 py-2.5 rounded-xl bg-white/3 border border-white/5 text-zinc-200 placeholder-zinc-500 focus:outline-none focus:border-purple-500/30 transition-all text-sm"
          />
        </div>
      </div>

      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="rounded-2xl bg-white/2 border border-white/5 backdrop-blur-sm overflow-hidden"
      >
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead>
              <tr className="border-b border-white/5 text-xs font-semibold uppercase tracking-wider text-zinc-400">
                <th className="px-6 py-4">Timestamp</th>
                <th className="px-6 py-4">User</th>
                <th className="px-6 py-4">Action</th>
                <th className="px-6 py-4">Target</th>
                <th className="px-6 py-4">IP Address</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/3 text-sm text-zinc-300">
              {audits.map((item) => (
                <tr key={item.id} className="hover:bg-white/1 transition-colors">
                  <td className="px-6 py-4 text-xs font-mono text-zinc-500">{item.time}</td>
                  <td className="px-6 py-4 font-semibold text-purple-400">{item.user}</td>
                  <td className="px-6 py-4 text-zinc-200">{item.action}</td>
                  <td className="px-6 py-4 text-zinc-400">{item.target}</td>
                  <td className="px-6 py-4 text-xs font-mono text-zinc-500">{item.ip}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Pagination Footer */}
        <div className="p-4 border-t border-white/5 flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2 text-sm text-zinc-400">
            <span>Show</span>
            <select
              value={itemsPerPage}
              onChange={(e) => {
                const val = e.target.value === 'all' ? 'all' : Number(e.target.value);
                setItemsPerPage(val);
                setCurrentPage(1);
              }}
              className="px-2 py-1 rounded bg-white/4 border border-white/8 text-zinc-200 focus:outline-none"
            >
              <option value="10">10</option>
              <option value="50">50</option>
              <option value="100">100</option>
              <option value="all">All</option>
            </select>
            <span>entries</span>
          </div>

          <div className="flex items-center gap-4">
            <span className="text-sm text-zinc-400">
              Page <span className="text-zinc-100">{currentPage}</span> of <span className="text-zinc-100">{maxPage}</span>
            </span>
            <div className="flex items-center gap-1">
              <button
                disabled={currentPage === 1}
                onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                className="p-2 rounded-lg border border-white/4 hover:bg-white/3 disabled:opacity-40 text-zinc-300 disabled:cursor-not-allowed"
              >
                <ChevronLeft size={16} />
              </button>
              <button
                disabled={currentPage === maxPage}
                onClick={() => setCurrentPage((p) => Math.min(maxPage, p + 1))}
                className="p-2 rounded-lg border border-white/4 hover:bg-white/3 disabled:opacity-40 text-zinc-300 disabled:cursor-not-allowed"
              >
                <ChevronRight size={16} />
              </button>
            </div>
          </div>
        </div>
      </motion.div>
    </div>
  );
}
