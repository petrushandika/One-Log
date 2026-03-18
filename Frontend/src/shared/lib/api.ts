import axios from 'axios';

// 1. Create Axios Instance set base config
const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

// 2. Setup JWT Auth Interceptor
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// 3. API Methods Index
export const authApi = {
  login: (data: { email?: string; password?: string }) => api.post('/auth/login', data),
};

export const statsApi = {
  getOverview: () => api.get('/stats/overview'),
};

export const sourcesApi = {
  getAll: () => api.get('/sources'),
  getByID: (id: string) => api.get(`/sources/${id}`),
  create: (data: { name: string; url: string }) => api.post('/sources', data),
  rotateKey: (id: string) => api.post(`/sources/${id}/rotate-key`),
};

export const logsApi = {
  getLogs: (params: { level?: string; source?: string; search?: string } = {}) => 
    api.get('/logs', { params }),
  getByID: (id: number) => api.get(`/logs/${id}`),
  analyze: (id: number) => api.post(`/logs/${id}/analyze`),
};

export default api;
