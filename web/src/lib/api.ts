import axios from 'axios';

const api = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor - add auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor - handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth
export const authAPI = {
  login: (username: string, password: string) =>
    api.post('/auth/login', { username, password }),
  me: () => api.get('/auth/me'),
  logout: () => api.post('/auth/logout'),
};

// Dashboard
export const dashboardAPI = {
  getStats: () => api.get('/dashboard/stats'),
};

// Users
export const usersAPI = {
  list: () => api.get('/users'),
  get: (id: number) => api.get(`/users/${id}`),
  create: (data: { username: string; email: string; password: string; role: string }) =>
    api.post('/users', data),
  update: (id: number, data: Partial<{ email: string; password: string; role: string; active: boolean }>) =>
    api.put(`/users/${id}`, data),
  delete: (id: number) => api.delete(`/users/${id}`),
};

// Packages
export const packagesAPI = {
  list: () => api.get('/packages'),
  create: (data: any) => api.post('/packages', data),
  update: (id: number, data: any) => api.put(`/packages/${id}`, data),
  delete: (id: number) => api.delete(`/packages/${id}`),
};

// Domains
export const domainsAPI = {
  list: () => api.get('/domains'),
  get: (id: number) => api.get(`/domains/${id}`),
  create: (data: { name: string; document_root?: string }) => api.post('/domains', data),
  delete: (id: number) => api.delete(`/domains/${id}`),
};

// Databases
export const databasesAPI = {
  list: () => api.get('/databases'),
  create: (data: { name: string; type?: string }) => api.post('/databases', data),
  delete: (id: number) => api.delete(`/databases/${id}`),
};

// System
export const systemAPI = {
  getStats: () => api.get('/system/stats'),
  getServices: () => api.get('/system/services'),
  restartService: (name: string) => api.post(`/system/services/${name}/restart`),
};

// Accounts (Hosting hesapları)
export const accountsAPI = {
  list: () => api.get('/accounts'),
  get: (id: number) => api.get(`/accounts/${id}`),
  create: (data: {
    username: string;
    email: string;
    password: string;
    domain: string;
    package_id: number;
  }) => api.post('/accounts', data),
  delete: (id: number) => api.delete(`/accounts/${id}`),
  suspend: (id: number) => api.post(`/accounts/${id}/suspend`),
  unsuspend: (id: number) => api.post(`/accounts/${id}/unsuspend`),
};

export default api;
