import Layout from '../shared/components/Layout';
import Overview from '../pages/Overview';
import Logs from '../pages/Logs';
import Sources from '../pages/Sources';
import Audit from '../pages/Audit';
import Login from '../pages/Login';
import ProtectedRoute from '../shared/components/ProtectedRoute';
import { Terminal } from 'lucide-react';

import { createBrowserRouter, Link } from 'react-router-dom';

export const router = createBrowserRouter([
  {
    path: '/',
    element: (
      <ProtectedRoute>
        <Layout />
      </ProtectedRoute>
    ),
    errorElement: (
      <div className="flex flex-col items-center justify-center min-h-screen bg-[#09090b] text-white p-6 text-center">
        <div className="p-3 rounded-2xl bg-purple-500/20 text-purple-400 mb-4">
          <Terminal size={32} />
        </div>
        <h1 className="text-4xl font-bold tracking-tight">404</h1>
        <p className="text-zinc-400 mt-2">Page not found or a system error occurred.</p>
        <Link to="/" className="mt-6 px-4 py-2 bg-purple-500 hover:bg-purple-600 rounded-xl font-semibold transition-colors flex items-center justify-center gap-2 text-sm shadow-lg shadow-purple-500/20">
          Back to Dashboard
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
  {
    path: '/login',
    element: <Login />,
  },
]);
