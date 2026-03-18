import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Settings, Plus, Key, ToggleRight, CheckCircle, XCircle, X, Copy, RefreshCw } from 'lucide-react';

interface Source {
  id: string;
  name: string;
  apiKey: string;
  status: string;
  url: string;
}

export default function Sources() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedSettingsSource, setSelectedSettingsSource] = useState<Source | null>(null);
  const [isRegenerating, setIsRegenerating] = useState<string | null>(null);
  const [toast, setToast] = useState<{ message: string, type: 'success' | 'error' } | null>(null);

  const showToast = (message: string, type: 'success' | 'error' = 'success') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 2500);
  };

  const [sources, setSources] = useState([
    { id: '1', name: 'Auth Service', apiKey: 'ulam_live_auth_98372', status: 'ONLINE', url: 'https://auth.sample.com' },
    { id: '2', name: 'Gateway', apiKey: 'ulam_live_gate_87413', status: 'ONLINE', url: 'https://gateway.sample.com' },
    { id: '3', name: 'DB Analytics', apiKey: 'ulam_live_db_24354', status: 'OFFLINE', url: 'https://db.sample.com' },
  ]);
  const [newSource, setNewSource] = useState({ name: '', url: '' });

  const handleRegister = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newSource.name || !newSource.url) return;

    const created = {
      id: String(sources.length + 1),
      name: newSource.name,
      url: newSource.url,
      apiKey: `ulam_live_src_${Math.random().toString(36).substr(2, 5)}`,
      status: 'ONLINE'
    };

    setSources([...sources, created]);
    setNewSource({ name: '', url: '' });
    setIsModalOpen(false);
    showToast('Source application registered successfully!');
  };

  const handleUpdate = (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedSettingsSource) return;
    setSources(sources.map(s => s.id === selectedSettingsSource.id ? selectedSettingsSource : s));
    setSelectedSettingsSource(null);
    showToast('Source settings updated successfully!');
  };

  const handleRegenKey = (id: string) => {
    setIsRegenerating(id);
    setTimeout(() => {
      setSources(sources.map(s => s.id === id ? { ...s, apiKey: `ulam_live_${Math.random().toString(36).substr(2, 6)}` } : s));
      setIsRegenerating(null);
      showToast('API Key regenerated successfully!');
    }, 1000);
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
          className="flex items-center gap-2 px-4 py-2.5 rounded-xl bg-purple-500 text-white text-sm font-semibold shadow-lg shadow-purple-500/20 hover:bg-purple-600 transition-all scale-100 hover:scale-[1.02] active:scale-95"
        >
          <Plus size={18} />
          Register Source
        </button>
      </div>

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3 mt-4"
      >
        {sources.map((source) => (
          <div
            key={source.id}
            className="p-6 rounded-2xl bg-white/2 border border-white/5 backdrop-blur-sm flex flex-col justify-between"
          >
            <div>
              <div className="flex items-center justify-between">
                <h3 className="text-lg font-semibold text-zinc-100">{source.name}</h3>
                <span className={`flex items-center gap-1.5 text-xs font-semibold px-2 py-1 rounded-md border ${
                  source.status === 'ONLINE' 
                    ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' 
                    : 'bg-rose-500/10 text-rose-400 border-rose-500/20'
                }`}>
                  {source.status === 'ONLINE' ? <CheckCircle size={14}/> : <XCircle size={14}/>}
                  {source.status}
                </span>
              </div>
              <p className="text-xs text-zinc-500 mt-1">{source.url}</p>

              <div className="mt-4 p-3 rounded-xl bg-black/30 border border-white/5 flex items-center justify-between gap-3">
                <div className="flex items-center gap-2 text-xs font-mono text-zinc-300 truncate">
                  <Key size={14} className="text-purple-400" />
                  <span className="truncate">{isRegenerating === source.id ? 'Regenerating...' : source.apiKey}</span>
                </div>
                <div className="flex items-center gap-1 border-l border-white/10 pl-2">
                  <button 
                    onClick={() => {
                      navigator.clipboard.writeText(source.apiKey);
                      showToast('API Key copied to clipboard!');
                    }}
                    className="p-1.5 rounded-lg hover:bg-white/10 text-zinc-400 hover:text-white transition-colors relative"
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

            <div className="mt-6 pt-4 border-t border-white/5 flex items-center justify-between gap-2">
              <button 
                onClick={() => setSelectedSettingsSource(source)}
                className="flex items-center gap-1.5 text-xs text-zinc-400 hover:text-zinc-200 transition-colors"
              >
                <Settings size={14} />
                Settings
              </button>
              <button 
                onClick={() => {
                  const newStatus = source.status === 'ONLINE' ? 'OFFLINE' : 'ONLINE';
                  setSources(sources.map(s => s.id === source.id ? { ...s, status: newStatus } : s));
                  showToast(`Source ${source.name} is now ${newStatus}`);
                }}
                className={`flex items-center gap-1.5 text-xs font-semibold transition-colors ${source.status === 'ONLINE' ? 'text-rose-400 hover:text-rose-300' : 'text-emerald-400 hover:text-emerald-300'}`}
              >
                <ToggleRight size={14} />
                {source.status === 'ONLINE' ? 'Disable' : 'Enable'}
              </button>
            </div>
          </div>
        ))}
      </motion.div>

      {/* Settings Modal */}
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
                    className="w-full px-4 py-2.5 rounded-xl bg-white/3 border border-white/5 text-zinc-200 placeholder-zinc-600 focus:outline-none focus:border-purple-500/40 text-sm"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-zinc-400 mb-1">Endpoint URL</label>
                  <input
                    type="url"
                    required
                    value={selectedSettingsSource.url}
                    onChange={(e) => setSelectedSettingsSource({ ...selectedSettingsSource, url: e.target.value })}
                    className="w-full px-4 py-2.5 rounded-xl bg-white/3 border border-white/5 text-zinc-200 placeholder-zinc-600 focus:outline-none focus:border-purple-500/40 text-sm"
                  />
                </div>
                <div className="pt-2 flex gap-3">
                  <button
                    type="button"
                    onClick={() => setSelectedSettingsSource(null)}
                    className="flex-1 py-2.5 rounded-xl bg-white/5 text-zinc-300 text-sm font-medium hover:bg-white/8 transition-colors"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="flex-1 py-2.5 rounded-xl bg-purple-500 text-white text-sm font-semibold hover:bg-purple-600 transition-colors shadow-lg shadow-purple-500/20"
                  >
                    Save Changes
                  </button>
                </div>
              </form>
            </motion.div>
          </>
        )}
      </AnimatePresence>

      {/* Modern Modal Overlay */}
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
                  <label className="block text-xs font-medium text-zinc-400 mb-1">Source Name</label>
                  <input
                    type="text"
                    required
                    value={newSource.name}
                    onChange={(e) => setNewSource({ ...newSource, name: e.target.value })}
                    placeholder="e.g. Payment Gateway"
                    className="w-full px-4 py-2.5 rounded-xl bg-white/3 border border-white/5 text-zinc-200 placeholder-zinc-600 focus:outline-none focus:border-purple-500/40 text-sm"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-zinc-400 mb-1">Endpoint URL</label>
                  <input
                    type="url"
                    required
                    value={newSource.url}
                    onChange={(e) => setNewSource({ ...newSource, url: e.target.value })}
                    placeholder="https://api.example.com"
                    className="w-full px-4 py-2.5 rounded-xl bg-white/3 border border-white/5 text-zinc-200 placeholder-zinc-600 focus:outline-none focus:border-purple-500/40 text-sm"
                  />
                </div>
                <div className="pt-2 flex gap-3">
                  <button
                    type="button"
                    onClick={() => setIsModalOpen(false)}
                    className="flex-1 py-2.5 rounded-xl bg-white/5 text-zinc-300 text-sm font-medium hover:bg-white/8 transition-colors"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="flex-1 py-2.5 rounded-xl bg-purple-500 text-white text-sm font-semibold hover:bg-purple-600 transition-colors shadow-lg shadow-purple-500/20"
                  >
                    Register
                  </button>
                </div>
              </form>
            </motion.div>
          </>
        )}
      </AnimatePresence>

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
