import { useState, useEffect, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Settings, Plus, Key, ToggleRight, CheckCircle, XCircle, X, Copy, RefreshCw, Wifi, WifiOff } from 'lucide-react';
import { sourcesApi } from '../shared/lib/api';

interface Source {
  id: string;
  name: string;
  api_key: string;
  status: string;
  health_url: string;
  created_at: string;
}

export default function Sources() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedSettingsSource, setSelectedSettingsSource] = useState<Source | null>(null);
  const [isRegenerating, setIsRegenerating] = useState<string | null>(null);
  const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null);
  const [sources, setSources] = useState<Source[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [newSource, setNewSource] = useState({ name: '', health_url: '' });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const showToast = (message: string, type: 'success' | 'error' = 'success') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 2500);
  };

  const fetchSources = useCallback(async () => {
    setIsLoading(true);
    try {
      const { data } = await sourcesApi.getAll();
      setSources(data.data ?? []);
    } catch (err) {
      console.error(err);
      showToast('Failed to load sources', 'error');
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchSources();
  }, [fetchSources]);

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newSource.name) return;
    setIsSubmitting(true);
    try {
      await sourcesApi.create({ name: newSource.name, health_url: newSource.health_url || undefined });
      showToast('Source registered successfully!');
      setNewSource({ name: '', health_url: '' });
      setIsModalOpen(false);
      fetchSources();
    } catch (err) {
      console.error(err);
      showToast('Failed to register source', 'error');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedSettingsSource) return;
    setIsSubmitting(true);
    try {
      await sourcesApi.update(selectedSettingsSource.id, {
        name: selectedSettingsSource.name,
        health_url: selectedSettingsSource.health_url,
      });
      showToast('Source updated successfully!');
      setSelectedSettingsSource(null);
      fetchSources();
    } catch (err) {
      console.error(err);
      showToast('Failed to update source', 'error');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleToggleStatus = async (source: Source) => {
    const newStatus = source.status === 'ONLINE' ? 'OFFLINE' : 'ONLINE';
    try {
      await sourcesApi.update(source.id, { status: newStatus });
      showToast(`${source.name} is now ${newStatus}`);
      fetchSources();
    } catch (err) {
      console.error(err);
      showToast('Failed to update status', 'error');
    }
  };

  const handleRegenKey = async (id: string) => {
    setIsRegenerating(id);
    try {
      await sourcesApi.rotateKey(id);
      showToast('API Key regenerated successfully!');
      fetchSources();
    } catch (err) {
      console.error(err);
      showToast('Failed to regenerate API Key', 'error');
    } finally {
      setIsRegenerating(null);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-white">Application Sources</h1>
          <p className="text-sm text-zinc-400">Manage connected applications and API keys</p>
        </div>
        <button
          onClick={() => setIsModalOpen(true)}
          className="flex items-center gap-2 px-4 py-2.5 rounded-xl bg-purple-500 text-white text-sm font-semibold shadow-lg shadow-purple-500/20 hover:bg-purple-600 transition-all active:scale-95"
        >
          <Plus size={18} />
          Register Source
        </button>
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="h-52 rounded-2xl bg-white/[0.02] border border-white/[0.05] animate-pulse" />
          ))}
        </div>
      ) : sources.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20 text-zinc-500">
          <Settings size={40} className="mb-4 opacity-30" />
          <p className="text-sm">No sources registered yet.</p>
          <button onClick={() => setIsModalOpen(true)} className="mt-4 text-purple-400 text-sm hover:underline">
            Register your first source
          </button>
        </div>
      ) : (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3"
        >
          {sources.map((source) => (
            <div
              key={source.id}
              className="p-6 rounded-2xl bg-white/[0.02] border border-white/[0.05] backdrop-blur-sm flex flex-col justify-between"
            >
              <div>
                <div className="flex items-center justify-between">
                  <h3 className="text-lg font-semibold text-zinc-100 truncate pr-2">{source.name}</h3>
                  <span
                    className={`shrink-0 flex items-center gap-1.5 text-xs font-semibold px-2 py-1 rounded-md border ${
                      source.status === 'ONLINE'
                        ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20'
                        : source.status === 'DEGRADED'
                        ? 'bg-amber-500/10 text-amber-400 border-amber-500/20'
                        : 'bg-rose-500/10 text-rose-400 border-rose-500/20'
                    }`}
                  >
                    {source.status === 'ONLINE' ? <CheckCircle size={12} /> : <XCircle size={12} />}
                    {source.status}
                  </span>
                </div>
                {source.health_url && (
                  <p className="text-xs text-zinc-500 mt-1 truncate">{source.health_url}</p>
                )}

                <div className="mt-4 p-3 rounded-xl bg-black/30 border border-white/[0.05] flex items-center justify-between gap-3">
                  <div className="flex items-center gap-2 text-xs font-mono text-zinc-300 min-w-0">
                    <Key size={14} className="text-purple-400 shrink-0" />
                    <span className="truncate">
                      {isRegenerating === source.id ? 'Regenerating...' : source.api_key}
                    </span>
                  </div>
                  <div className="flex items-center gap-1 border-l border-white/10 pl-2 shrink-0">
                    <button
                      onClick={() => {
                        navigator.clipboard.writeText(source.api_key);
                        showToast('API Key copied!');
                      }}
                      className="p-1.5 rounded-lg hover:bg-white/10 text-zinc-400 hover:text-white transition-colors"
                      title="Copy API Key"
                    >
                      <Copy size={14} />
                    </button>
                    <button
                      onClick={() => handleRegenKey(source.id)}
                      disabled={isRegenerating === source.id}
                      className="p-1.5 rounded-lg hover:bg-white/10 text-purple-400 hover:text-purple-300 transition-colors disabled:opacity-50"
                      title="Regenerate Key"
                    >
                      <RefreshCw size={14} className={isRegenerating === source.id ? 'animate-spin' : ''} />
                    </button>
                  </div>
                </div>
              </div>

              <div className="mt-6 pt-4 border-t border-white/[0.05] flex items-center justify-between gap-2">
                <button
                  onClick={() => setSelectedSettingsSource(source)}
                  className="flex items-center gap-1.5 text-xs text-zinc-400 hover:text-zinc-200 transition-colors"
                >
                  <Settings size={14} />
                  Settings
                </button>
                <button
                  onClick={() => handleToggleStatus(source)}
                  className={`flex items-center gap-1.5 text-xs font-semibold transition-colors ${
                    source.status === 'ONLINE'
                      ? 'text-rose-400 hover:text-rose-300'
                      : 'text-emerald-400 hover:text-emerald-300'
                  }`}
                >
                  {source.status === 'ONLINE' ? <WifiOff size={14} /> : <Wifi size={14} />}
                  <ToggleRight size={14} />
                  {source.status === 'ONLINE' ? 'Disable' : 'Enable'}
                </button>
              </div>
            </div>
          ))}
        </motion.div>
      )}

      {/* Edit Modal */}
      <AnimatePresence>
        {selectedSettingsSource && (
          <>
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setSelectedSettingsSource(null)}
              className="fixed inset-0 bg-black/60 backdrop-blur-md z-40"
            />
            <motion.div
              initial={{ opacity: 0, scale: 0.95, y: 20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95, y: 20 }}
              className="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full max-w-md p-6 rounded-2xl bg-[#121214] border border-white/5 z-50 shadow-2xl"
            >
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-xl font-bold text-white">Edit Source</h3>
                <button onClick={() => setSelectedSettingsSource(null)} className="text-zinc-400 hover:text-zinc-200">
                  <X size={20} />
                </button>
              </div>
              <form onSubmit={handleUpdate} className="space-y-4">
                <div>
                  <label className="block text-xs font-medium text-zinc-400 mb-1">Source Name</label>
                  <input
                    type="text"
                    required
                    value={selectedSettingsSource.name}
                    onChange={(e) => setSelectedSettingsSource({ ...selectedSettingsSource, name: e.target.value })}
                    className="w-full px-4 py-2.5 rounded-xl bg-white/[0.03] border border-white/5 text-zinc-200 focus:outline-none focus:border-purple-500/40 text-sm"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-zinc-400 mb-1">Health Check URL</label>
                  <input
                    type="url"
                    value={selectedSettingsSource.health_url ?? ''}
                    onChange={(e) => setSelectedSettingsSource({ ...selectedSettingsSource, health_url: e.target.value })}
                    placeholder="https://api.example.com/health"
                    className="w-full px-4 py-2.5 rounded-xl bg-white/[0.03] border border-white/5 text-zinc-200 placeholder-zinc-600 focus:outline-none focus:border-purple-500/40 text-sm"
                  />
                </div>
                <div className="pt-2 flex gap-3">
                  <button type="button" onClick={() => setSelectedSettingsSource(null)} className="flex-1 py-2.5 rounded-xl bg-white/5 text-zinc-300 text-sm font-medium hover:bg-white/[0.08] transition-colors">
                    Cancel
                  </button>
                  <button type="submit" disabled={isSubmitting} className="flex-1 py-2.5 rounded-xl bg-purple-500 text-white text-sm font-semibold hover:bg-purple-600 transition-colors shadow-lg shadow-purple-500/20 disabled:opacity-50">
                    {isSubmitting ? 'Saving...' : 'Save Changes'}
                  </button>
                </div>
              </form>
            </motion.div>
          </>
        )}
      </AnimatePresence>

      {/* Register Modal */}
      <AnimatePresence>
        {isModalOpen && (
          <>
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setIsModalOpen(false)}
              className="fixed inset-0 bg-black/60 backdrop-blur-md z-40"
            />
            <motion.div
              initial={{ opacity: 0, scale: 0.95, y: 20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95, y: 20 }}
              className="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full max-w-md p-6 rounded-2xl bg-[#121214] border border-white/5 z-50 shadow-2xl"
            >
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-xl font-bold text-white">Register New Source</h3>
                <button onClick={() => setIsModalOpen(false)} className="text-zinc-400 hover:text-zinc-200">
                  <X size={20} />
                </button>
              </div>
              <form onSubmit={handleRegister} className="space-y-4">
                <div>
                  <label className="block text-xs font-medium text-zinc-400 mb-1">Source Name <span className="text-rose-400">*</span></label>
                  <input
                    type="text"
                    required
                    value={newSource.name}
                    onChange={(e) => setNewSource({ ...newSource, name: e.target.value })}
                    placeholder="e.g. Payment Gateway"
                    className="w-full px-4 py-2.5 rounded-xl bg-white/[0.03] border border-white/5 text-zinc-200 placeholder-zinc-600 focus:outline-none focus:border-purple-500/40 text-sm"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-zinc-400 mb-1">Health Check URL <span className="text-zinc-600">(optional)</span></label>
                  <input
                    type="url"
                    value={newSource.health_url}
                    onChange={(e) => setNewSource({ ...newSource, health_url: e.target.value })}
                    placeholder="https://api.example.com/health"
                    className="w-full px-4 py-2.5 rounded-xl bg-white/[0.03] border border-white/5 text-zinc-200 placeholder-zinc-600 focus:outline-none focus:border-purple-500/40 text-sm"
                  />
                </div>
                <p className="text-xs text-zinc-500">An API key will be generated automatically after registration.</p>
                <div className="pt-2 flex gap-3">
                  <button type="button" onClick={() => setIsModalOpen(false)} className="flex-1 py-2.5 rounded-xl bg-white/5 text-zinc-300 text-sm font-medium hover:bg-white/[0.08] transition-colors">
                    Cancel
                  </button>
                  <button type="submit" disabled={isSubmitting} className="flex-1 py-2.5 rounded-xl bg-purple-500 text-white text-sm font-semibold hover:bg-purple-600 transition-colors shadow-lg shadow-purple-500/20 disabled:opacity-50">
                    {isSubmitting ? 'Registering...' : 'Register'}
                  </button>
                </div>
              </form>
            </motion.div>
          </>
        )}
      </AnimatePresence>

      {/* Toast */}
      <AnimatePresence>
        {toast && (
          <motion.div
            initial={{ opacity: 0, y: 50, scale: 0.9 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: 50, scale: 0.9 }}
            className={`fixed bottom-6 right-6 px-4 py-3 rounded-xl border backdrop-blur-md shadow-2xl flex items-center gap-2 z-50 text-sm font-semibold ${
              toast.type === 'success'
                ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20'
                : 'bg-rose-500/10 text-rose-400 border-rose-500/20'
            }`}
          >
            {toast.type === 'success' ? <CheckCircle size={16} /> : <XCircle size={16} />}
            {toast.message}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
