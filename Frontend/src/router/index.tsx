import { createBrowserRouter, Link } from 'react-router-dom';
import Layout from '../shared/components/Layout';
import Overview from '../pages/Overview';
import Logs from '../pages/Logs';
import Sources from '../pages/Sources';
import Audit from '../pages/Audit';
import { Terminal } from 'lucide-react';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    errorElement: (
      <div className="flex flex-col items-center justify-center min-h-screen bg-[#09090b] text-white p-6 text-center">
        <div className="p-3 rounded-2xl bg-purple-500/20 text-purple-400 mb-4">
          <Terminal size={32} />
        </div>
        <h1 className="text-4xl font-bold tracking-tight">404</h1>
        <p className="text-zinc-400 mt-2">Halaman tidak ditemukan atau terjadi kesalahan sistem.</p>
        <Link to="/" className="mt-6 px-4 py-2 bg-purple-500 hover:bg-purple-600 rounded-xl font-semibold transition-colors flex items-center justify-center gap-2 text-sm shadow-lg shadow-purple-500/20">
          Kembali ke Dashboard
        </Link>
      </div>
    ),
    children: [
      { index: true, element: <Overview /> },
      { path: 'logs', element: <Logs /> },
      { path: 'sources', element: <Sources /> },
      { path: 'audit', element: <Audit /> },
    ],
  },
]);
