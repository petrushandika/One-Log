import axios from 'axios';

const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api';

const api = axios.create({
  baseURL: BASE_URL,
  headers: { 'Content-Type': 'application/json' },
  withCredentials: true,
});

// Request: attach Bearer token from localStorage (legacy support)
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

// Response: auto-logout on 401 and handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    
    // Handle rate limiting (429)
    if (error.response?.status === 429) {
      const retryAfter = error.response.headers['retry-after'] || 60;
      console.error(`Rate limit exceeded. Retry after ${retryAfter} seconds.`);
    }
    
    return Promise.reject(error);
  },
);

// ─── Auth ───────────────────────────────────────────────────────────────────
export const authApi = {
  login: (data: { email: string; password: string }) => api.post('/auth/login', data),
  refresh: () => api.post('/auth/refresh'),
  logout: () => api.post('/auth/logout'),
};

// ─── Stats ──────────────────────────────────────────────────────────────────
export const statsApi = {
  getOverview: () => api.get('/stats/overview'),
  getActivitySummary: () => api.get('/stats/activity'),
};

// ─── Sources ────────────────────────────────────────────────────────────────
export const sourcesApi = {
  getAll: () => api.get('/sources'),
  getByID: (id: string) => api.get(`/sources/${id}`),
  create: (data: { name: string; health_url?: string }) => api.post('/sources', data),
  update: (id: string, data: { name?: string; health_url?: string; status?: string }) =>
    api.patch(`/sources/${id}`, data),
  rotateKey: (id: string) => api.post(`/sources/${id}/rotate-key`),
};

// ─── Logs ────────────────────────────────────────────────────────────────────
export const logsApi = {
  getLogs: (params: {
    source_id?: string;
    level?: string;
    category?: string;
    page?: number;
    limit?: number;
    from?: string;
    to?: string;
  } = {}) => api.get('/logs', { params }),
  getByID: (id: number) => api.get(`/logs/${id}`),
  analyze: (id: number) => api.post(`/logs/${id}/analyze`),
  export: (params: {
    source_id?: string;
    level?: string;
    category?: string;
    from?: string;
    to?: string;
  } = {}) => api.get('/logs/export', { params, responseType: 'blob' }),
};

// ─── Activity ────────────────────────────────────────────────────────────────
export const activityApi = {
  list: (params: {
    source_id?: string;
    category?: string;
    page?: number;
    limit?: number;
  } = {}) => api.get('/activity', { params }),
  summary: (params: { source_id?: string; period?: string } = {}) =>
    api.get('/activity/summary', { params }),
  byUser: (userId: string) => api.get(`/activity/users/${userId}`),
  suspicious: (params: { source_id?: string; page?: number; limit?: number } = {}) =>
    api.get('/activity/suspicious', { params }),
  // Activity Analytics (Phase 2)
  getAuthMethodBreakdown: (params: { days?: number } = {}) =>
    api.get('/activity/analytics/methods', { params }),
  getLoginTimeline: (params: { days?: number } = {}) =>
    api.get('/activity/analytics/timeline', { params }),
  getFailedLoginHeatmap: (params: { days?: number } = {}) =>
    api.get('/activity/analytics/heatmap', { params }),
  getSessions: (params: { page?: number; limit?: number } = {}) =>
    api.get('/activity/sessions', { params }),
  // Activity Feed (Phase 2 Extended)
  getFeed: (params: { page?: number; limit?: number; action?: string } = {}) =>
    api.get('/activity/feed', { params }),
  getTopUsers: (params: { days?: number; limit?: number } = {}) =>
    api.get('/activity/top-users', { params }),
  // Compliance Export (Phase 2 Extended)
  requestComplianceExport: (data: { source_id: string; format: string; date_from: string; date_to: string }) =>
    api.post('/activity/compliance-export', data),
  getComplianceExports: (params: { page?: number; limit?: number } = {}) =>
    api.get('/activity/compliance-exports', { params }),
};

// ─── APM ─────────────────────────────────────────────────────────────────────
export const apmApi = {
  endpointStats: (params: { source_id?: string; period?: string } = {}) =>
    api.get('/apm/endpoints', { params }),
  timeline: (params: { source_id?: string; endpoint?: string; period?: string; interval?: string } = {}) =>
    api.get('/apm/timeline', { params }),
  // APM Thresholds (Phase 3)
  thresholds: {
    list: () => api.get('/apm/thresholds'),
    create: (data: { source_id: string; endpoint: string; p95_limit: number; p99_limit: number; email_notify: boolean }) =>
      api.post('/apm/thresholds', data),
    update: (id: string, data: { p95_limit: number; p99_limit: number; email_notify: boolean }) =>
      api.patch(`/apm/thresholds/${id}`, data),
    delete: (id: string) => api.delete(`/apm/thresholds/${id}`),
  },
  getSlowQueries: (params: { source_id?: string; threshold?: number } = {}) =>
    api.get('/apm/slow-queries', { params }),
};

// ─── Issues ──────────────────────────────────────────────────────────────────
export const issuesApi = {
  list: (params: {
    source_id?: string;
    status?: string;
    page?: number;
    limit?: number;
  } = {}) => api.get('/issues', { params }),
  get: (fingerprint: string) => api.get(`/issues/${fingerprint}`),
  updateStatus: (fingerprint: string, status: string) =>
    api.patch(`/issues/${fingerprint}`, { status }),
  logs: (fingerprint: string, params: { page?: number; limit?: number } = {}) =>
    api.get(`/issues/${fingerprint}/logs`, { params }),
  trend: (params: { source_id?: string; days?: number } = {}) =>
    api.get('/issues/analytics/trend', { params }),
  heatmap: (params: { source_id?: string; days?: number } = {}) =>
    api.get('/issues/analytics/heatmap', { params }),
};

// ─── Status (public) ─────────────────────────────────────────────────────────
export const statusApi = {
  getPublic: () => api.get('/status'),
};

// ─── Config ──────────────────────────────────────────────────────────────────
export const configApi = {
  list: (sourceId: string, params: { environment?: string; reveal?: boolean } = {}) =>
    api.get(`/sources/${sourceId}/configs`, { params }),
  save: (sourceId: string, data: { key: string; value: string; is_secret?: boolean; environment?: string }) =>
    api.post(`/sources/${sourceId}/configs`, data),
  history: (sourceId: string, params: { environment?: string } = {}) =>
    api.get(`/sources/${sourceId}/configs/history`, { params }),
};

// ─── Incidents ───────────────────────────────────────────────────────────────
export const incidentsApi = {
  list: (params: { status?: string; source_id?: string; page?: number; limit?: number } = {}) =>
    api.get('/incidents', { params }),
  timeline: (params: { source_id?: string; days?: number } = {}) =>
    api.get('/incidents/timeline', { params }),
};

// ─── Chat (AI Copilot) ───────────────────────────────────────────────────────
export const chatApi = {
  send: (message: string) => api.post<{ data: { reply: string } }>('/chat', { message }),
};

export default api;
