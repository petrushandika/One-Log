import Layout from '../shared/components/Layout';
import Overview from '../pages/Overview';
import Logs from '../pages/Logs';
import Sources from '../pages/Sources';
import Audit from '../pages/Audit';
import Issues from '../pages/Issues';
import APM from '../pages/APM';
import Status from '../pages/Status';
import Config from '../pages/Config';
import Incidents from '../pages/Incidents';
import ActivityAnalytics from '../pages/ActivityAnalytics';
import ActivityFeed from '../pages/ActivityFeed';
import ComplianceExport from '../pages/ComplianceExport';
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
      { path: 'issues', element: <Issues /> },
      { path: 'apm', element: <APM /> },
      { path: 'status', element: <Status /> },
      { path: 'sources', element: <Sources /> },
      { path: 'config', element: <Config /> },
      { path: 'incidents', element: <Incidents /> },
      { path: 'activity', element: <ActivityAnalytics /> },
      { path: 'activity/feed', element: <ActivityFeed /> },
      { path: 'audit', element: <Audit /> },
      { path: 'compliance', element: <ComplianceExport /> },
    ],
  },
  {
    path: '/login',
    element: <Login />,
  },
]);
