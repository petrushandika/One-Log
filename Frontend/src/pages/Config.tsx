import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { SlidersHorizontal, Eye, EyeOff, Edit3, X, History, ChevronRight, RotateCcw, Lock, Unlock } from 'lucide-react';
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
  const [sources, setSources] = useState<Source[]>([]);
  const [selectedSource, setSelectedSource] = useState<Source | null>(null);
  const [configs, setConfigs] = useState<ConfigEntry[]>([]);
  const [history, setHistory] = useState<ConfigHistory[]>([]);
  const [activeTab, setActiveTab] = useState<'config' | 'history'>('config');
  const [isLoading, setIsLoading] = useState(false);
  const [revealedKeys, setRevealedKeys] = useState<Set<number>>(new Set());
  const [editEntry, setEditEntry] = useState<ConfigEntry | null>(null);
  const [editValue, setEditValue] = useState('');
  const [editIsSecret, setEditIsSecret] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [newKey, setNewKey] = useState('');
  const [newValue, setNewValue] = useState('');
  const [newIsSecret, setNewIsSecret] = useState(false);
  const [isAdding, setIsAdding] = useState(false);

  useEffect(() => {
    sourcesApi.getAll().then(({ data }) => {
      const list: Source[] = data.data ?? [];
      setSources(list);
      if (list.length > 0) setSelectedSource(list[0]);
    }).catch(console.error);
  }, []);

  useEffect(() => {
    if (!selectedSource) return;
    fetchConfigs();
    if (activeTab === 'history') fetchHistory();
  }, [selectedSource, activeTab]);

  const fetchConfigs = async () => {
    if (!selectedSource) return;
    setIsLoading(true);
    try {
      const { data } = await configApi.list(selectedSource.id);
      setConfigs(data.data ?? []);
    } catch (err) {
      console.error('Failed to fetch configs', err);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchHistory = async () => {
    if (!selectedSource) return;
    try {
      const { data } = await configApi.history(selectedSource.id);
      setHistory(data.data ?? []);
    } catch (err) {
      console.error('Failed to fetch history', err);
    }
  };

  const revealSecret = async (entry: ConfigEntry) => {
    if (revealedKeys.has(entry.id)) {
      setRevealedKeys((prev) => { const s = new Set(prev); s.delete(entry.id); return s; });
      return;
    }
    try {
      const { data } = await configApi.list(selectedSource!.id, { reveal: true });
      const revealed = (data.data ?? []).find((c: ConfigEntry) => c.id === entry.id);
      if (revealed) {
        setConfigs((prev) => prev.map((c) => c.id === entry.id ? { ...c, value: revealed.value } : c));
        setRevealedKeys((prev) => new Set(prev).add(entry.id));
      }
    } catch (err) {
      console.error('Failed to reveal secret', err);
    }
  };

  const openEdit = (entry: ConfigEntry) => {
    setEditEntry(entry);
    setEditValue(entry.is_secret && !revealedKeys.has(entry.id) ? '' : entry.value);
    setEditIsSecret(entry.is_secret);
  };

  const handleSave = async () => {
    if (!editEntry || !selectedSource) return;
    setIsSaving(true);
    try {
      await configApi.save(selectedSource.id, {
        key: editEntry.key,
        value: editValue,
        is_secret: editIsSecret,
        environment: editEntry.environment,
      });
      setEditEntry(null);
      fetchConfigs();
      fetchHistory();
    } catch (err) {
      console.error('Failed to save config', err);
    } finally {
      setIsSaving(false);
    }
  };

  const handleAdd = async () => {
    if (!selectedSource || !newKey.trim() || !newValue.trim()) return;
    setIsAdding(true);
    try {
      await configApi.save(selectedSource.id, { key: newKey.trim(), value: newValue.trim(), is_secret: newIsSecret });
      setNewKey('');
      setNewValue('');
      setNewIsSecret(false);
      fetchConfigs();
    } catch (err) {
      console.error('Failed to add config', err);
    } finally {
      setIsAdding(false);
    }
  };

  const handleRollback = async (h: ConfigHistory) => {
    if (!selectedSource) return;
    try {
      await configApi.save(selectedSource.id, {
        key: h.key,
        value: h.value,
        is_secret: h.is_secret,
      });
      fetchConfigs();
      setActiveTab('config');
    } catch (err) {
      console.error('Failed to rollback', err);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-white tracking-tight flex items-center gap-2">
          <SlidersHorizontal size={22} className="text-purple-400" />
          Config
        </h1>
        <p className="text-sm text-zinc-400 mt-0.5">Centralized configuration management per source</p>
      </div>

      {/* Source Selector */}
      <div className="flex items-center gap-3">
        <label className="text-sm text-zinc-400 shrink-0">Source</label>
        <SelectField
          value={selectedSource?.id ?? ''}
          onChange={(e) => {
            const s = sources.find((src) => src.id === e.target.value);
            if (s) { setSelectedSource(s); setRevealedKeys(new Set()); }
          }}
          wrapperClassName="flex-1 max-w-sm"
        >
          {sources.map((s) => (
            <option key={s.id} value={s.id}>{s.name}</option>
          ))}
        </SelectField>
      </div>

      {selectedSource && (
        <>
          {/* Tabs */}
          <div className="flex border-b border-white/[0.06]">
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
              <div className="p-4 rounded-xl bg-white/[0.02] border border-white/[0.05] space-y-3">
                <p className="text-xs font-semibold text-zinc-400 uppercase tracking-wider">Add Config Entry</p>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                  <input
                    placeholder="Key (e.g. DATABASE_URL)"
                    value={newKey}
                    onChange={(e) => setNewKey(e.target.value)}
                    className="px-3 py-2 rounded-lg bg-white/[0.03] border border-white/10 text-zinc-200 text-sm placeholder:text-zinc-600 focus:outline-none focus:border-purple-500/40"
                  />
                  <input
                    placeholder="Value"
                    type={newIsSecret ? 'password' : 'text'}
                    value={newValue}
                    onChange={(e) => setNewValue(e.target.value)}
                    className="px-3 py-2 rounded-lg bg-white/[0.03] border border-white/10 text-zinc-200 text-sm placeholder:text-zinc-600 focus:outline-none focus:border-purple-500/40"
                  />
                  <div className="flex gap-2">
                    <button
                      onClick={() => setNewIsSecret((v) => !v)}
                      className={`flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium border transition-all ${
                        newIsSecret
                          ? 'bg-amber-500/10 text-amber-400 border-amber-500/20'
                          : 'bg-white/[0.02] text-zinc-500 border-white/10 hover:bg-white/[0.04]'
                      }`}
                    >
                      {newIsSecret ? <Lock size={13} /> : <Unlock size={13} />}
                      {newIsSecret ? 'Secret' : 'Plain'}
                    </button>
                    <button
                      onClick={handleAdd}
                      disabled={isAdding || !newKey.trim() || !newValue.trim()}
                      className="flex-1 flex items-center justify-center gap-1.5 px-3 py-2 rounded-lg text-xs font-semibold bg-purple-600 hover:bg-purple-500 text-white transition-all disabled:opacity-40"
                    >
                      Save
                    </button>
                  </div>
                </div>
              </div>

              {/* Config Table */}
              <div className="rounded-2xl bg-white/[0.02] border border-white/[0.05] overflow-hidden">
                {isLoading ? (
                  <div className="flex items-center justify-center h-32 text-zinc-500 text-sm">Loading...</div>
                ) : configs.length === 0 ? (
                  <div className="flex flex-col items-center justify-center h-32 text-center">
                    <SlidersHorizontal size={24} className="text-zinc-600 mb-2" />
                    <p className="text-zinc-500 text-sm">No config entries yet.</p>
                  </div>
                ) : (
                  <table className="w-full text-left">
                    <thead>
                      <tr className="border-b border-white/[0.05] text-xs font-semibold uppercase tracking-wider text-zinc-400">
                        <th className="px-5 py-3">Key</th>
                        <th className="px-5 py-3">Value</th>
                        <th className="px-5 py-3">Type</th>
                        <th className="px-5 py-3">Updated</th>
                        <th className="px-5 py-3"></th>
                      </tr>
                    </thead>
                    <tbody>
                      {configs.map((entry) => (
                        <tr key={entry.id} className="border-b border-white/[0.03] hover:bg-white/[0.02] transition-colors">
                          <td className="px-5 py-3">
                            <code className="text-sm text-purple-300 font-mono">{entry.key}</code>
                          </td>
                          <td className="px-5 py-3 max-w-xs">
                            <div className="flex items-center gap-2">
                              <span className="text-sm text-zinc-300 font-mono truncate">
                                {entry.is_secret && !revealedKeys.has(entry.id) ? '••••••••' : entry.value}
                              </span>
                              {entry.is_secret && (
                                <button
                                  onClick={() => revealSecret(entry)}
                                  className="shrink-0 text-zinc-600 hover:text-zinc-300 transition-colors"
                                  title={revealedKeys.has(entry.id) ? 'Hide' : 'Reveal'}
                                >
                                  {revealedKeys.has(entry.id) ? <EyeOff size={13} /> : <Eye size={13} />}
                                </button>
                              )}
                            </div>
                          </td>
                          <td className="px-5 py-3">
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
                          <td className="px-5 py-3 text-xs text-zinc-600">{timeAgo(entry.updated_at)}</td>
                          <td className="px-5 py-3">
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
              <div className="rounded-2xl bg-white/[0.02] border border-white/[0.05] overflow-hidden">
                {history.length === 0 ? (
                  <div className="flex flex-col items-center justify-center h-32 text-center">
                    <History size={24} className="text-zinc-600 mb-2" />
                    <p className="text-zinc-500 text-sm">No history yet.</p>
                  </div>
                ) : (
                  <table className="w-full text-left">
                    <thead>
                      <tr className="border-b border-white/[0.05] text-xs font-semibold uppercase tracking-wider text-zinc-400">
                        <th className="px-5 py-3">Key</th>
                        <th className="px-5 py-3">Value</th>
                        <th className="px-5 py-3">Ver.</th>
                        <th className="px-5 py-3">Changed</th>
                        <th className="px-5 py-3"></th>
                      </tr>
                    </thead>
                    <tbody>
                      {history.map((h) => (
                        <tr key={h.id} className="border-b border-white/[0.03] hover:bg-white/[0.02] transition-colors">
                          <td className="px-5 py-3">
                            <code className="text-sm text-purple-300 font-mono">{h.key}</code>
                          </td>
                          <td className="px-5 py-3 max-w-xs">
                            <span className="text-sm text-zinc-400 font-mono truncate block">
                              {h.is_secret ? '••••••••' : h.value}
                            </span>
                          </td>
                          <td className="px-5 py-3">
                            <span className="text-xs text-zinc-500">v{h.version}</span>
                          </td>
                          <td className="px-5 py-3 text-xs text-zinc-600">{timeAgo(h.created_at)}</td>
                          <td className="px-5 py-3">
                            <button
                              onClick={() => handleRollback(h)}
                              className="flex items-center gap-1 px-2.5 py-1 rounded-lg text-xs font-medium bg-white/[0.03] border border-white/10 text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.06] transition-all"
                              title="Rollback to this version"
                            >
                              <RotateCcw size={11} />
                              Rollback
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
                    className="w-full px-3 py-2 rounded-lg bg-white/[0.02] border border-white/10 text-zinc-500 text-sm font-mono cursor-not-allowed"
                  />
                </div>

                <div>
                  <label className="block text-xs text-zinc-400 mb-1.5 font-medium">Value</label>
                  <textarea
                    rows={4}
                    value={editValue}
                    onChange={(e) => setEditValue(e.target.value)}
                    placeholder={editEntry.is_secret ? 'Enter new value (leave blank to keep current)' : 'Value'}
                    className="w-full px-3 py-2 rounded-lg bg-white/[0.03] border border-white/10 text-zinc-200 text-sm font-mono placeholder:text-zinc-600 focus:outline-none focus:border-purple-500/50 resize-none"
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
                  className="flex-1 px-4 py-2.5 rounded-xl border border-white/10 text-zinc-400 hover:bg-white/5 text-sm transition-all"
                >
                  Cancel
                </button>
                <button
                  onClick={handleSave}
                  disabled={isSaving}
                  className="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-purple-600 hover:bg-purple-500 text-white text-sm font-semibold transition-all disabled:opacity-50"
                >
                  {isSaving ? 'Saving...' : 'Save Changes'}
                </button>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>
    </div>
  );
}
