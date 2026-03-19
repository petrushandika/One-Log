import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { SlidersHorizontal, Eye, EyeOff, Edit3, X, History, ChevronRight, RotateCcw, Lock, Unlock, AlertCircle, RefreshCw, CheckCircle } from 'lucide-react';
import SelectField from '../shared/components/SelectField';
import { sourcesApi, configApi } from '../shared/lib/api';

interface Source {
  id: string;
  name: string;
}

interface ConfigEntry {
  id: number;
  key: string;
  value: string;
  is_secret: boolean;
  environment: string;
  updated_at: string;
}

interface ConfigHistory {
  id: number;
  key: string;
  value: string;
  is_secret: boolean;
  version: number;
  created_at: string;
}

function timeAgo(dateStr: string): string {
  const diff = Date.now() - new Date(dateStr).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 1) return 'just now';
  if (mins < 60) return `${mins}m ago`;
  const hrs = Math.floor(mins / 60);
  if (hrs < 24) return `${hrs}h ago`;
  return `${Math.floor(hrs / 24)}d ago`;
}

export default function Config() {
  const queryClient = useQueryClient();
  const [selectedSource, setSelectedSource] = useState<Source | null>(null);
  const [activeTab, setActiveTab] = useState<'config' | 'history'>('config');
  const [revealedKeys, setRevealedKeys] = useState<Set<number>>(new Set());
  const [editEntry, setEditEntry] = useState<ConfigEntry | null>(null);
  const [editValue, setEditValue] = useState('');
  const [editIsSecret, setEditIsSecret] = useState(false);
  const [newKey, setNewKey] = useState('');
  const [newValue, setNewValue] = useState('');
  const [newIsSecret, setNewIsSecret] = useState(false);
  const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null);

  const showToast = (message: string, type: 'success' | 'error' = 'success') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 3000);
  };

  // Fetch sources
  const { data: sources = [], isLoading: isLoadingSources, error: sourcesError } = useQuery({
    queryKey: ['sources'],
    queryFn: async () => {
      const { data } = await sourcesApi.getAll();
      const list: Source[] = data.data ?? [];
      if (list.length > 0 && !selectedSource) {
        setSelectedSource(list[0]);
      }
      return list;
    },
  });

  // Fetch configs for selected source
  const { 
    data: configs = [], 
    isLoading: isLoadingConfigs, 
    error: configsError,
    refetch: refetchConfigs 
  } = useQuery({
    queryKey: ['configs', selectedSource?.id],
    queryFn: async () => {
      if (!selectedSource) return [];
      const { data } = await configApi.list(selectedSource.id);
      return data.data ?? [];
    },
    enabled: !!selectedSource,
  });

  // Fetch history
  const { 
    data: history = [], 
    isLoading: isLoadingHistory, 
    error: historyError,
    refetch: refetchHistory 
  } = useQuery({
    queryKey: ['config-history', selectedSource?.id],
    queryFn: async () => {
      if (!selectedSource) return [];
      const { data } = await configApi.history(selectedSource.id);
      return data.data ?? [];
    },
    enabled: !!selectedSource && activeTab === 'history',
  });

  // Reveal secret mutation
  const revealMutation = useMutation({
    mutationFn: async (entryId: number) => {
      if (!selectedSource) throw new Error('No source selected');
      const { data } = await configApi.list(selectedSource.id, { reveal: true });
      return { data: data.data ?? [], entryId };
    },
    onSuccess: ({ data, entryId }) => {
      const revealed = data.find((c: ConfigEntry) => c.id === entryId);
      if (revealed) {
        queryClient.setQueryData(['configs', selectedSource?.id], (old: ConfigEntry[] | undefined) => {
          if (!old) return [];
          return old.map((c) => c.id === entryId ? { ...c, value: revealed.value } : c);
        });
        setRevealedKeys((prev) => new Set(prev).add(entryId));
      }
    },
    onError: () => {
      showToast('Failed to reveal secret', 'error');
    },
  });

  // Save config mutation
  const saveMutation = useMutation({
    mutationFn: async (payload: { key: string; value: string; is_secret: boolean; environment?: string }) => {
      if (!selectedSource) throw new Error('No source selected');
      await configApi.save(selectedSource.id, payload);
    },
    onSuccess: () => {
      showToast('Configuration saved successfully');
      setEditEntry(null);
      setNewKey('');
      setNewValue('');
      setNewIsSecret(false);
      queryClient.invalidateQueries({ queryKey: ['configs', selectedSource?.id] });
      queryClient.invalidateQueries({ queryKey: ['config-history', selectedSource?.id] });
    },
    onError: () => {
      showToast('Failed to save configuration', 'error');
    },
  });

  // Rollback mutation
  const rollbackMutation = useMutation({
    mutationFn: async (h: ConfigHistory) => {
      if (!selectedSource) throw new Error('No source selected');
      await configApi.save(selectedSource.id, {
        key: h.key,
        value: h.value,
        is_secret: h.is_secret,
      });
    },
    onSuccess: () => {
      showToast('Configuration rolled back successfully');
      setActiveTab('config');
      queryClient.invalidateQueries({ queryKey: ['configs', selectedSource?.id] });
      queryClient.invalidateQueries({ queryKey: ['config-history', selectedSource?.id] });
    },
    onError: () => {
      showToast('Failed to rollback configuration', 'error');
    },
  });

  const revealSecret = async (entry: ConfigEntry) => {
    if (revealedKeys.has(entry.id)) {
      setRevealedKeys((prev) => { const s = new Set(prev); s.delete(entry.id); return s; });
      return;
    }
    revealMutation.mutate(entry.id);
  };

  const openEdit = (entry: ConfigEntry) => {
    setEditEntry(entry);
    setEditValue(entry.is_secret && !revealedKeys.has(entry.id) ? '' : entry.value);
    setEditIsSecret(entry.is_secret);
  };

  const handleSave = () => {
    if (!editEntry || !selectedSource) return;
    saveMutation.mutate({
      key: editEntry.key,
      value: editValue,
      is_secret: editIsSecret,
      environment: editEntry.environment,
    });
  };

  const handleAdd = () => {
    if (!selectedSource || !newKey.trim() || !newValue.trim()) return;
    saveMutation.mutate({
      key: newKey.trim(),
      value: newValue.trim(),
      is_secret: newIsSecret,
    });
  };

  const handleRollback = (h: ConfigHistory) => {
    rollbackMutation.mutate(h);
  };

  const error = sourcesError || configsError;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        <div className="p-2 rounded-xl bg-purple-500/10 text-purple-400">
          <SlidersHorizontal size={24} />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-white">Config</h1>
          <p className="text-sm text-zinc-400">Centralized configuration management per source</p>
        </div>
      </div>

      {/* Error Message */}
      {error && (
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          className="p-4 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 flex items-center gap-3"
        >
          <AlertCircle size={20} />
          <div className="flex-1">
            <p className="font-medium">Failed to load configuration data</p>
          </div>
          <button
            onClick={() => {
              queryClient.invalidateQueries({ queryKey: ['sources'] });
              refetchConfigs();
            }}
            className="px-3 py-1.5 text-sm bg-red-500/20 hover:bg-red-500/30 rounded-lg transition-colors flex items-center gap-1"
          >
            <RefreshCw size={14} />
            Retry
          </button>
        </motion.div>
      )}

      {/* Source Selector */}
      <div className="flex items-center gap-3">
        <label className="text-sm text-zinc-400 shrink-0">Source</label>
        <SelectField
          value={selectedSource?.id ?? ''}
          onChange={(e) => {
            const s = sources.find((src: Source) => src.id === e.target.value);
            if (s) { setSelectedSource(s); setRevealedKeys(new Set()); }
          }}
          wrapperClassName="flex-1 max-w-sm"
          disabled={isLoadingSources}
        >
          {sources.map((s: Source) => (
            <option key={s.id} value={s.id}>{s.name}</option>
          ))}
        </SelectField>
        {isLoadingSources && <RefreshCw size={16} className="animate-spin text-zinc-500" />}
      </div>

      {selectedSource && (
        <>
          {/* Tabs */}
          <div className="flex border-b border-white/5">
            {([
              { id: 'config', label: 'Config', icon: SlidersHorizontal },
              { id: 'history', label: 'Change History', icon: History },
            ] as const).map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center gap-2 px-5 py-3 text-sm font-medium transition-all border-b-2 -mb-px ${
                  activeTab === tab.id
                    ? 'border-purple-500 text-purple-400'
                    : 'border-transparent text-zinc-500 hover:text-zinc-300'
                }`}
              >
                <tab.icon size={14} />
                {tab.label}
              </button>
            ))}
          </div>

          {activeTab === 'config' && (
            <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="space-y-4">
              {/* Add new config */}
              <div className="bg-white/2 border border-white/5 rounded-xl p-4 space-y-3">
                <p className="text-xs font-semibold text-zinc-400 uppercase tracking-wider">Add Config Entry</p>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                  <input
                    placeholder="Key (e.g. DATABASE_URL)"
                    value={newKey}
                    onChange={(e) => setNewKey(e.target.value)}
                    className="px-3 py-2 rounded-lg bg-white/3 border border-white/5 text-zinc-200 text-sm placeholder:text-zinc-600 focus:outline-none focus:border-purple-500/40"
                  />
                  <input
                    placeholder="Value"
                    type={newIsSecret ? 'password' : 'text'}
                    value={newValue}
                    onChange={(e) => setNewValue(e.target.value)}
                    className="px-3 py-2 rounded-lg bg-white/3 border border-white/5 text-zinc-200 text-sm placeholder:text-zinc-600 focus:outline-none focus:border-purple-500/40"
                  />
                  <div className="flex gap-2">
                    <button
                      onClick={() => setNewIsSecret((v) => !v)}
                      className={`flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium border transition-all ${
                        newIsSecret
                          ? 'bg-amber-500/10 text-amber-400 border-amber-500/20'
                          : 'bg-white/2 text-zinc-500 border-white/5 hover:bg-white/4'
                      }`}
                    >
                      {newIsSecret ? <Lock size={13} /> : <Unlock size={13} />}
                      {newIsSecret ? 'Secret' : 'Plain'}
                    </button>
                    <button
                      onClick={handleAdd}
                      disabled={saveMutation.isPending || !newKey.trim() || !newValue.trim()}
                      className="flex-1 flex items-center justify-center gap-1.5 px-3 py-2 rounded-lg text-xs font-semibold bg-purple-600 hover:bg-purple-500 text-white transition-all disabled:opacity-40"
                    >
                      {saveMutation.isPending ? 'Saving...' : 'Save'}
                    </button>
                  </div>
                </div>
              </div>

              {/* Config Table */}
              <div className="bg-white/2 border border-white/5 rounded-xl overflow-hidden">
                {isLoadingConfigs ? (
                  <div className="flex items-center justify-center h-32 text-zinc-500 text-sm gap-2">
                    <RefreshCw size={16} className="animate-spin" />
                    Loading configurations...
                  </div>
                ) : configs.length === 0 ? (
                  <div className="flex flex-col items-center justify-center h-32 text-center">
                    <SlidersHorizontal size={24} className="text-zinc-600 mb-2" />
                    <p className="text-zinc-500 text-sm">No config entries yet.</p>
                    <p className="text-zinc-600 text-xs mt-1">Add your first configuration above</p>
                  </div>
                ) : (
                  <table className="w-full text-left">
                    <thead>
                      <tr className="border-b border-white/5 text-xs font-semibold uppercase text-zinc-500">
                        <th className="px-4 py-3">Key</th>
                        <th className="px-4 py-3">Value</th>
                        <th className="px-4 py-3">Type</th>
                        <th className="px-4 py-3">Updated</th>
                        <th className="px-4 py-3"></th>
                      </tr>
                    </thead>
                    <tbody>
                      {configs.map((entry: ConfigEntry) => (
                        <tr key={entry.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                          <td className="px-4 py-3">
                            <code className="text-sm text-purple-300 font-mono">{entry.key}</code>
                          </td>
                          <td className="px-4 py-3 max-w-xs">
                            <div className="flex items-center gap-2">
                              <span className="text-sm text-zinc-300 font-mono truncate">
                                {entry.is_secret && !revealedKeys.has(entry.id) ? '••••••••' : entry.value}
                              </span>
                              {entry.is_secret && (
                                <button
                                  onClick={() => revealSecret(entry)}
                                  disabled={revealMutation.isPending}
                                  className="shrink-0 text-zinc-600 hover:text-zinc-300 transition-colors disabled:opacity-50"
                                  title={revealedKeys.has(entry.id) ? 'Hide' : 'Reveal'}
                                >
                                  {revealedKeys.has(entry.id) ? <EyeOff size={13} /> : <Eye size={13} />}
                                </button>
                              )}
                            </div>
                          </td>
                          <td className="px-4 py-3">
                            {entry.is_secret ? (
                              <span className="flex items-center gap-1 text-xs text-amber-400">
                                <Lock size={11} /> Secret
                              </span>
                            ) : (
                              <span className="flex items-center gap-1 text-xs text-zinc-500">
                                <Unlock size={11} /> Plain
                              </span>
                            )}
                          </td>
                          <td className="px-4 py-3 text-xs text-zinc-600">{timeAgo(entry.updated_at)}</td>
                          <td className="px-4 py-3">
                            <button
                              onClick={() => openEdit(entry)}
                              className="p-1.5 rounded-lg hover:bg-purple-500/10 text-zinc-500 hover:text-purple-400 transition-colors"
                            >
                              <Edit3 size={14} />
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                )}
              </div>
            </motion.div>
          )}

          {activeTab === 'history' && (
            <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
              <div className="bg-white/2 border border-white/5 rounded-xl overflow-hidden">
                {isLoadingHistory ? (
                  <div className="flex items-center justify-center h-32 text-zinc-500 text-sm gap-2">
                    <RefreshCw size={16} className="animate-spin" />
                    Loading history...
                  </div>
                ) : historyError ? (
                  <div className="flex flex-col items-center justify-center h-32 text-center">
                    <AlertCircle size={24} className="text-red-500/60 mb-2" />
                    <p className="text-zinc-500 text-sm">Failed to load history</p>
                    <button 
                      onClick={() => refetchHistory()} 
                      className="text-purple-400 text-xs mt-2 hover:text-purple-300"
                    >
                      Retry
                    </button>
                  </div>
                ) : history.length === 0 ? (
                  <div className="flex flex-col items-center justify-center h-32 text-center">
                    <History size={24} className="text-zinc-600 mb-2" />
                    <p className="text-zinc-500 text-sm">No history yet.</p>
                    <p className="text-zinc-600 text-xs mt-1">Changes will appear here</p>
                  </div>
                ) : (
                  <table className="w-full text-left">
                    <thead>
                      <tr className="border-b border-white/5 text-xs font-semibold uppercase text-zinc-500">
                        <th className="px-4 py-3">Key</th>
                        <th className="px-4 py-3">Value</th>
                        <th className="px-4 py-3">Ver.</th>
                        <th className="px-4 py-3">Changed</th>
                        <th className="px-4 py-3"></th>
                      </tr>
                    </thead>
                    <tbody>
                      {history.map((h: ConfigHistory) => (
                        <tr key={h.id} className="border-b border-white/5 hover:bg-white/5 transition-colors">
                          <td className="px-4 py-3">
                            <code className="text-sm text-purple-300 font-mono">{h.key}</code>
                          </td>
                          <td className="px-4 py-3 max-w-xs">
                            <span className="text-sm text-zinc-400 font-mono truncate block">
                              {h.is_secret ? '••••••••' : h.value}
                            </span>
                          </td>
                          <td className="px-4 py-3">
                            <span className="text-xs text-zinc-500">v{h.version}</span>
                          </td>
                          <td className="px-4 py-3 text-xs text-zinc-600">{timeAgo(h.created_at)}</td>
                          <td className="px-4 py-3">
                            <button
                              onClick={() => handleRollback(h)}
                              disabled={rollbackMutation.isPending}
                              className="flex items-center gap-1 px-2.5 py-1 rounded-lg text-xs font-medium bg-white/3 border border-white/5 text-zinc-400 hover:text-zinc-200 hover:bg-white/6 transition-all disabled:opacity-50"
                              title="Rollback to this version"
                            >
                              <RotateCcw size={11} className={rollbackMutation.isPending ? 'animate-spin' : ''} />
                              {rollbackMutation.isPending ? 'Rolling back...' : 'Rollback'}
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                )}
              </div>
            </motion.div>
          )}
        </>
      )}

      {/* Edit Slide-over */}
      <AnimatePresence>
        {editEntry && (
          <div className="fixed inset-0 z-50 flex justify-end">
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setEditEntry(null)}
              className="absolute inset-0 bg-black/60 backdrop-blur-sm"
            />
            <motion.div
              initial={{ x: '100%' }}
              animate={{ x: 0 }}
              exit={{ x: '100%' }}
              transition={{ ease: 'easeInOut', duration: 0.25 }}
              className="relative w-full max-w-md bg-[#0c0c0e] border-l border-white/10 p-6 flex flex-col gap-5 shadow-2xl h-full overflow-y-auto"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2 text-zinc-200 font-semibold">
                  <ChevronRight size={16} className="text-purple-400" />
                  Edit Config
                </div>
                <button onClick={() => setEditEntry(null)} className="p-1.5 rounded-lg hover:bg-white/10 text-zinc-400">
                  <X size={18} />
                </button>
              </div>

              <div className="space-y-4">
                <div>
                  <label className="block text-xs text-zinc-400 mb-1.5 font-medium">Key</label>
                  <input
                    readOnly
                    value={editEntry.key}
                    className="w-full px-3 py-2 rounded-lg bg-white/2 border border-white/5 text-zinc-500 text-sm font-mono cursor-not-allowed"
                  />
                </div>

                <div>
                  <label className="block text-xs text-zinc-400 mb-1.5 font-medium">Value</label>
                  <textarea
                    rows={4}
                    value={editValue}
                    onChange={(e) => setEditValue(e.target.value)}
                    placeholder={editEntry.is_secret ? 'Enter new value (leave blank to keep current)' : 'Value'}
                    className="w-full px-3 py-2 rounded-lg bg-white/3 border border-white/5 text-zinc-200 text-sm font-mono placeholder:text-zinc-600 focus:outline-none focus:border-purple-500/50 resize-none"
                  />
                </div>

                <div className="flex items-center justify-between">
                  <span className="text-sm text-zinc-400">Secret value</span>
                  <button
                    onClick={() => setEditIsSecret((v) => !v)}
                    className={`relative w-11 h-6 rounded-full transition-colors ${editIsSecret ? 'bg-amber-500/60' : 'bg-white/10'}`}
                  >
                    <span className={`absolute top-1 w-4 h-4 rounded-full bg-white transition-transform ${editIsSecret ? 'translate-x-6' : 'translate-x-1'}`} />
                  </button>
                </div>
                {editIsSecret && (
                  <p className="text-xs text-amber-400/80 bg-amber-500/10 border border-amber-500/20 rounded-lg px-3 py-2">
                    Value will be encrypted with AES-256 before storage.
                  </p>
                )}
              </div>

              <div className="mt-auto flex gap-3">
                <button
                  onClick={() => setEditEntry(null)}
                  className="flex-1 px-4 py-2.5 rounded-xl border border-white/5 text-zinc-400 hover:bg-white/5 text-sm transition-all"
                >
                  Cancel
                </button>
                <button
                  onClick={handleSave}
                  disabled={saveMutation.isPending}
                  className="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-purple-600 hover:bg-purple-500 text-white text-sm font-semibold transition-all disabled:opacity-50"
                >
                  {saveMutation.isPending ? 'Saving...' : 'Save Changes'}
                </button>
              </div>
            </motion.div>
          </div>
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
            {toast.type === 'success' ? <CheckCircle size={15} /> : <AlertCircle size={15} />}
            {toast.message}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
