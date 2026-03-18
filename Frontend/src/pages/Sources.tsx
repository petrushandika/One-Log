import { useState, useEffect, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Server, Plus, Key, CheckCircle, XCircle, X, Copy, RefreshCw, Wifi, WifiOff, Eye, EyeOff, Settings } from 'lucide-react';
import { sourcesApi } from '../shared/lib/api';

interface Source {
  id: string;
  name: string;
  status: string;
  health_url: string;
  created_at: string;
  updated_at: string;
}

interface RevealedKey {
  key: string;
  visible: boolean;
}

function maskKey(key: string): string {
  // Show prefix up to first underscore after "ulam_live_", mask the rest
  const prefix = key.startsWith('ulam_live_') ? 'ulam_live_' : key.slice(0, 8);
  return `${prefix}${'•'.repeat(12)}`;
}

export default function Sources() {
  const [sources, setSources] = useState<Source[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedSettingsSource, setSelectedSettingsSource] = useState<Source | null>(null);
  const [newSource, setNewSource] = useState({ name: '', health_url: '' });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isRegenerating, setIsRegenerating] = useState<string | null>(null);
  const [isTogglingStatus, setIsTogglingStatus] = useState<string | null>(null);
  const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null);

  // Map of sourceId → { key: plaintext, visible: bool }
  // Only populated on create or rotate — never from list API (key is server-masked)
  const [revealedKeys, setRevealedKeys] = useState<Map<string, RevealedKey>>(new Map());

  const showToast = (message: string, type: 'success' | 'error' = 'success') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 3000);
  };

  const fetchSources = useCallback(async () => {
    setIsLoading(true);
    try {
      const { data } = await sourcesApi.getAll();
      setSources(data.data ?? []);
    } catch {
      showToast('Failed to load sources', 'error');
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => { fetchSources(); }, [fetchSources]);

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newSource.name.trim()) return;
    setIsSubmitting(true);
    try {
      const { data } = await sourcesApi.create({
        name: newSource.name.trim(),
        health_url: newSource.health_url.trim() || undefined,
      });
      const rawKey: string = data.data?.api_key ?? '';
      const newSourceId: string = data.data?.id ?? '';

      setNewSource({ name: '', health_url: '' });
      setIsModalOpen(false);
      await fetchSources();

      // Reveal the key immediately after creation
      if (rawKey && newSourceId) {
        setRevealedKeys((prev) => new Map(prev).set(newSourceId, { key: rawKey, visible: true }));
      }
      showToast('Source registered! Copy the API key below.');
    } catch {
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
        // only send health_url when it has a value — empty string would fail URL validation
        ...(selectedSettingsSource.health_url?.trim()
          ? { health_url: selectedSettingsSource.health_url.trim() }
          : {}),
      });
      showToast('Source updated');
      setSelectedSettingsSource(null);
      fetchSources();
    } catch {
      showToast('Failed to update source', 'error');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleToggleStatus = async (source: Source) => {
    const newStatus = source.status === 'ONLINE' ? 'OFFLINE' : 'ONLINE';
    setIsTogglingStatus(source.id);
    try {
      await sourcesApi.update(source.id, { status: newStatus });
      showToast(`${source.name} is now ${newStatus}`);
      fetchSources();
    } catch {
      showToast('Failed to update status', 'error');
    } finally {
      setIsTogglingStatus(null);
    }
  };

  const handleRotateKey = async (id: string) => {
    setIsRegenerating(id);
    try {
      const { data } = await sourcesApi.rotateKey(id);
      const newKey: string = data.data?.new_api_key ?? '';
      if (newKey) {
        setRevealedKeys((prev) => new Map(prev).set(id, { key: newKey, visible: true }));
        showToast('New API key generated — save it now!');
      } else {
        showToast('Key rotated', 'success');
      }
      fetchSources();
    } catch {
      showToast('Failed to rotate API key', 'error');
    } finally {
      setIsRegenerating(null);
    }
  };

  const toggleKeyVisibility = (id: string) => {
    setRevealedKeys((prev) => {
      const next = new Map(prev);
      const entry = next.get(id);
      if (entry) next.set(id, { ...entry, visible: !entry.visible });
      return next;
    });
  };

  const copyKey = (id: string, fallback?: string) => {
    const key = revealedKeys.get(id)?.key ?? fallback ?? '';
    if (!key) { showToast('No key to copy — rotate to reveal', 'error'); return; }
    navigator.clipboard.writeText(key);
    showToast('API key copied!');
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-white flex items-center gap-2.5">
            <Server size={22} className="text-purple-400" />
            Sources
          </h1>
          <p className="text-sm text-zinc-400 mt-0.5">
            {sources.length} registered application{sources.length !== 1 ? 's' : ''}
          </p>
        </div>
        <button
          onClick={() => setIsModalOpen(true)}
          className="flex items-center gap-2 px-4 py-2.5 rounded-xl bg-purple-600 text-white text-sm font-semibold shadow-lg shadow-purple-500/20 hover:bg-purple-500 transition-all active:scale-95"
        >
          <Plus size={16} />
          Register Source
        </button>
      </div>

      {/* Source Grid */}
      {isLoading ? (
        <div className="grid grid-cols-1 gap-5 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="h-56 rounded-2xl bg-white/2 border border-white/5 animate-pulse" />
          ))}
        </div>
      ) : sources.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20 text-zinc-500">
          <div className="p-4 rounded-2xl bg-white/2 border border-white/5 mb-4">
            <Server size={32} className="opacity-40" />
          </div>
          <p className="text-sm font-medium">No sources registered yet</p>
          <button onClick={() => setIsModalOpen(true)} className="mt-3 text-purple-400 text-sm hover:text-purple-300 transition-colors">
            Register your first source →
          </button>
        </div>
      ) : (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="grid grid-cols-1 gap-5 md:grid-cols-2 lg:grid-cols-3"
        >
          {sources.map((source, i) => {
            const revealed = revealedKeys.get(source.id);
            const isOnline = source.status === 'ONLINE';
            const isDegraded = source.status === 'DEGRADED';

            return (
              <motion.div
                key={source.id}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.04 }}
                className="p-5 rounded-2xl bg-white/2 border border-white/5 flex flex-col gap-4"
              >
                {/* Top row: name + status badge */}
                <div className="flex items-start justify-between gap-2">
                  <div className="min-w-0">
                    <h3 className="font-semibold text-zinc-100 truncate">{source.name}</h3>
                    {source.health_url && (
                      <p className="text-xs text-zinc-500 truncate mt-0.5">{source.health_url}</p>
                    )}
                  </div>
                  <span
                    className={`shrink-0 flex items-center gap-1 text-[11px] font-semibold px-2 py-1 rounded-lg border ${
                      isOnline
                        ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20'
                        : isDegraded
                        ? 'bg-amber-500/10 text-amber-400 border-amber-500/20'
                        : 'bg-rose-500/10 text-rose-400 border-rose-500/20'
                    }`}
                  >
                    {isOnline ? <CheckCircle size={11} /> : <XCircle size={11} />}
                    {source.status}
                  </span>
                </div>

                {/* API Key section */}
                <div className="rounded-xl bg-black/30 border border-white/6 overflow-hidden">
                  <div className="flex items-center gap-1.5 px-3 py-2 border-b border-white/5">
                    <Key size={12} className="text-purple-400 shrink-0" />
                    <span className="text-[10px] font-semibold uppercase tracking-wider text-zinc-500">API Key</span>
                  </div>

                  {revealed ? (
                    /* Key was just generated/rotated — show it (eye toggleable) */
                    <div className="px-3 py-2.5 flex items-center gap-2">
                      <code className="flex-1 text-[11px] font-mono text-zinc-300 truncate select-all">
                        {revealed.visible ? revealed.key : maskKey(revealed.key)}
                      </code>
                      <div className="flex items-center gap-0.5 shrink-0">
                        <button
                          onClick={() => toggleKeyVisibility(source.id)}
                          className="p-1.5 rounded-lg hover:bg-white/10 text-zinc-500 hover:text-zinc-200 transition-colors"
                          title={revealed.visible ? 'Hide key' : 'Show key'}
                        >
                          {revealed.visible ? <EyeOff size={13} /> : <Eye size={13} />}
                        </button>
                        <button
                          onClick={() => copyKey(source.id)}
                          className="p-1.5 rounded-lg hover:bg-white/10 text-zinc-500 hover:text-zinc-200 transition-colors"
                          title="Copy key"
                        >
                          <Copy size={13} />
                        </button>
                        <button
                          onClick={() => handleRotateKey(source.id)}
                          disabled={isRegenerating === source.id}
                          className="p-1.5 rounded-lg hover:bg-white/10 text-purple-400 hover:text-purple-300 transition-colors disabled:opacity-40"
                          title="Rotate key"
                        >
                          <RefreshCw size={13} className={isRegenerating === source.id ? 'animate-spin' : ''} />
                        </button>
                      </div>
                    </div>
                  ) : (
                    /* No key revealed — show masked + rotate CTA */
                    <div className="px-3 py-2.5 flex items-center gap-2">
                      <code className="flex-1 text-[11px] font-mono text-zinc-600 select-none">
                        ulam_live_••••••••••••
                      </code>
                      <button
                        onClick={() => handleRotateKey(source.id)}
                        disabled={isRegenerating === source.id}
                        className="shrink-0 flex items-center gap-1 px-2 py-1 rounded-lg text-[11px] font-medium bg-purple-600/15 hover:bg-purple-600/25 border border-purple-500/20 text-purple-400 hover:text-purple-300 transition-all disabled:opacity-40"
                        title="Generate new key to reveal"
                      >
                        <RefreshCw size={11} className={isRegenerating === source.id ? 'animate-spin' : ''} />
                        {isRegenerating === source.id ? 'Generating...' : 'Rotate & Reveal'}
                      </button>
                    </div>
                  )}
                </div>

                {/* Footer */}
                <div className="flex items-center justify-between pt-1">
                  <button
                    onClick={() => setSelectedSettingsSource(source)}
                    className="flex items-center gap-1.5 text-xs text-zinc-500 hover:text-zinc-200 transition-colors"
                  >
                    <Settings size={13} />
                    Settings
                  </button>
                  <button
                    onClick={() => handleToggleStatus(source)}
                    disabled={isTogglingStatus === source.id}
                    className={`flex items-center gap-1.5 text-xs font-semibold transition-colors disabled:opacity-50 ${
                      isOnline
                        ? 'text-rose-400 hover:text-rose-300'
                        : 'text-emerald-400 hover:text-emerald-300'
                    }`}
                  >
                    {isOnline ? <WifiOff size={13} /> : <Wifi size={13} />}
                    {isTogglingStatus === source.id ? 'Updating...' : isOnline ? 'Disable' : 'Enable'}
                  </button>
                </div>
              </motion.div>
            );
          })}
        </motion.div>
      )}

      {/* Edit Modal */}
      <AnimatePresence>
        {selectedSettingsSource && (
          <>
            <motion.div
              initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}
              onClick={() => setSelectedSettingsSource(null)}
              className="fixed inset-0 bg-black/60 backdrop-blur-sm z-40"
            />
            <motion.div
              initial={{ opacity: 0, scale: 0.95, y: 20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95, y: 20 }}
              className="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full max-w-md p-6 rounded-2xl bg-[#111113] border border-white/8 z-50 shadow-2xl"
            >
              <div className="flex items-center justify-between mb-5">
                <h3 className="text-lg font-bold text-white">Edit Source</h3>
                <button onClick={() => setSelectedSettingsSource(null)} className="p-1.5 rounded-lg hover:bg-white/10 text-zinc-400 hover:text-zinc-200 transition-colors">
                  <X size={18} />
                </button>
              </div>
              <form onSubmit={handleUpdate} className="space-y-4">
                <div>
                  <label className="block text-xs font-medium text-zinc-400 mb-1.5">Source Name</label>
                  <input
                    type="text" required
                    value={selectedSettingsSource.name}
                    onChange={(e) => setSelectedSettingsSource({ ...selectedSettingsSource, name: e.target.value })}
                    className="w-full px-3 py-2.5 rounded-xl bg-white/3 border border-white/8 text-zinc-200 focus:outline-none focus:border-purple-500/40 text-sm"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-zinc-400 mb-1.5">Health Check URL</label>
                  <input
                    type="url"
                    value={selectedSettingsSource.health_url ?? ''}
                    onChange={(e) => setSelectedSettingsSource({ ...selectedSettingsSource, health_url: e.target.value })}
                    placeholder="https://api.example.com/health"
                    className="w-full px-3 py-2.5 rounded-xl bg-white/3 border border-white/8 text-zinc-200 placeholder-zinc-600 focus:outline-none focus:border-purple-500/40 text-sm"
                  />
                </div>
                <div className="flex gap-3 pt-1">
                  <button type="button" onClick={() => setSelectedSettingsSource(null)} className="flex-1 py-2.5 rounded-xl bg-white/5 text-zinc-300 text-sm hover:bg-white/8 transition-colors">
                    Cancel
                  </button>
                  <button type="submit" disabled={isSubmitting} className="flex-1 py-2.5 rounded-xl bg-purple-600 text-white text-sm font-semibold hover:bg-purple-500 transition-colors disabled:opacity-50">
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
              initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}
              onClick={() => setIsModalOpen(false)}
              className="fixed inset-0 bg-black/60 backdrop-blur-sm z-40"
            />
            <motion.div
              initial={{ opacity: 0, scale: 0.95, y: 20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95, y: 20 }}
              className="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full max-w-md p-6 rounded-2xl bg-[#111113] border border-white/8 z-50 shadow-2xl"
            >
              <div className="flex items-center justify-between mb-5">
                <div>
                  <h3 className="text-lg font-bold text-white">Register Source</h3>
                  <p className="text-xs text-zinc-500 mt-0.5">An API key will be generated automatically.</p>
                </div>
                <button onClick={() => setIsModalOpen(false)} className="p-1.5 rounded-lg hover:bg-white/10 text-zinc-400 hover:text-zinc-200 transition-colors">
                  <X size={18} />
                </button>
              </div>
              <form onSubmit={handleRegister} className="space-y-4">
                <div>
                  <label className="block text-xs font-medium text-zinc-400 mb-1.5">
                    Source Name <span className="text-rose-400">*</span>
                  </label>
                  <input
                    type="text" required
                    value={newSource.name}
                    onChange={(e) => setNewSource({ ...newSource, name: e.target.value })}
                    placeholder="e.g. Payment Gateway"
                    className="w-full px-3 py-2.5 rounded-xl bg-white/3 border border-white/8 text-zinc-200 placeholder-zinc-600 focus:outline-none focus:border-purple-500/40 text-sm"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-zinc-400 mb-1.5">
                    Health Check URL <span className="text-zinc-600">(optional)</span>
                  </label>
                  <input
                    type="url"
                    value={newSource.health_url}
                    onChange={(e) => setNewSource({ ...newSource, health_url: e.target.value })}
                    placeholder="https://api.example.com/health"
                    className="w-full px-3 py-2.5 rounded-xl bg-white/3 border border-white/8 text-zinc-200 placeholder-zinc-600 focus:outline-none focus:border-purple-500/40 text-sm"
                  />
                </div>
                <div className="flex gap-3 pt-1">
                  <button type="button" onClick={() => setIsModalOpen(false)} className="flex-1 py-2.5 rounded-xl bg-white/5 text-zinc-300 text-sm hover:bg-white/8 transition-colors">
                    Cancel
                  </button>
                  <button type="submit" disabled={isSubmitting} className="flex-1 py-2.5 rounded-xl bg-purple-600 text-white text-sm font-semibold hover:bg-purple-500 transition-colors disabled:opacity-50">
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
            initial={{ opacity: 0, y: 40, scale: 0.95 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: 40, scale: 0.95 }}
            className={`fixed bottom-24 right-6 px-4 py-3 rounded-xl border shadow-2xl flex items-center gap-2 z-60 text-sm font-medium ${
              toast.type === 'success'
                ? 'bg-[#111113] text-emerald-400 border-emerald-500/20'
                : 'bg-[#111113] text-rose-400 border-rose-500/20'
            }`}
          >
            {toast.type === 'success' ? <CheckCircle size={15} /> : <XCircle size={15} />}
            {toast.message}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
