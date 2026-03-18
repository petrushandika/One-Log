import { useState, useRef, useEffect, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { MessageCircle, X, Send, Bot, Zap, Activity, ShieldAlert, Cpu, BarChart2, Copy, Check } from 'lucide-react';
import { chatApi } from '../lib/api';

interface Message {
  role: 'user' | 'assistant';
  text: string;
  time: string;
}

const SUGGESTIONS = [
  { icon: BarChart2,   text: 'Ada berapa ERROR hari ini?' },
  { icon: ShieldAlert, text: 'Apa itu kategori AUDIT_TRAIL?' },
  { icon: Activity,    text: 'Cara kirim PERFORMANCE log?' },
  { icon: Cpu,         text: 'Bagaimana cara debug SYSTEM_ERROR?' },
];

function now(): string {
  return new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
}

function CopyableCode({ children }: { children: React.ReactNode }) {
  const [copied, setCopied] = useState(false);

  const getText = (node: React.ReactNode): string => {
    if (typeof node === 'string') return node;
    if (typeof node === 'number') return String(node);
    if (Array.isArray(node)) return node.map(getText).join('');
    if (node && typeof node === 'object' && 'props' in (node as React.ReactElement)) {
      return getText((node as React.ReactElement).props.children);
    }
    return '';
  };

  const handleCopy = () => {
    const text = getText(children);
    navigator.clipboard.writeText(text).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  };

  return (
    <div className="relative group my-2 rounded-lg bg-[#080809] border border-white/7 overflow-hidden">
      <button
        onClick={handleCopy}
        className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 flex items-center gap-1 px-2 py-1 rounded-md bg-zinc-700/60 hover:bg-zinc-600/80 text-zinc-400 hover:text-zinc-100 transition-all text-[10px] font-medium z-10"
        title="Copy code"
      >
        {copied ? <Check size={11} className="text-emerald-400" /> : <Copy size={11} />}
        {copied ? 'Copied' : 'Copy'}
      </button>
      <div className="overflow-x-auto">
        {children}
      </div>
    </div>
  );
}

function AiMarkdown({ text }: { text: string }) {
  return (
    <ReactMarkdown
      remarkPlugins={[remarkGfm]}
      components={{
        p:          ({ children }) => <p className="my-1 leading-relaxed text-[13px] text-zinc-300">{children}</p>,
        strong:     ({ children }) => <strong className="font-semibold text-zinc-100">{children}</strong>,
        em:         ({ children }) => <em className="italic text-zinc-400">{children}</em>,
        ul:         ({ children }) => <ul className="pl-4 my-1.5 space-y-0.5 list-disc text-zinc-300 text-[13px]">{children}</ul>,
        ol:         ({ children }) => <ol className="pl-4 my-1.5 space-y-0.5 list-decimal text-zinc-300 text-[13px]">{children}</ol>,
        li:         ({ children }) => <li className="leading-relaxed">{children}</li>,
        h1:         ({ children }) => <h1 className="font-bold text-sm text-zinc-100 mt-3 mb-1 border-b border-white/7 pb-1">{children}</h1>,
        h2:         ({ children }) => <h2 className="font-semibold text-sm text-zinc-100 mt-3 mb-1">{children}</h2>,
        h3:         ({ children }) => <h3 className="font-medium text-[13px] text-zinc-200 mt-2 mb-0.5">{children}</h3>,
        a:          ({ href, children }) => <a href={href} target="_blank" rel="noopener noreferrer" className="text-purple-400 underline underline-offset-2 hover:text-purple-300">{children}</a>,
        blockquote: ({ children }) => <blockquote className="border-l-2 border-purple-500/40 pl-3 my-2 text-zinc-400 italic text-[13px]">{children}</blockquote>,
        hr:         () => <hr className="border-white/7 my-2" />,
        pre:        ({ children }) => <CopyableCode>{children}</CopyableCode>,
        code: ({ children, className }) => {
          if (className?.includes('language-')) {
            return (
              <code className="block p-3 text-[11px] font-mono text-emerald-300 whitespace-pre leading-relaxed">
                {children}
              </code>
            );
          }
          // Inline code — break-all so it wraps within the chat bubble
          return (
            <code className="bg-zinc-800/60 text-purple-300 px-1.5 py-0.5 rounded text-[11px] font-mono border border-white/5 wrap-break-word">
              {children}
            </code>
          );
        },
      }}
    >
      {text}
    </ReactMarkdown>
  );
}

function TypingDots() {
  return (
    <div className="flex items-center gap-1.5 py-0.5 px-0.5">
      {[0, 160, 320].map((d) => (
        <span
          key={d}
          className="w-2 h-2 rounded-full bg-purple-400/70 animate-bounce"
          style={{ animationDelay: `${d}ms`, animationDuration: '800ms' }}
        />
      ))}
    </div>
  );
}

export default function ChatWidget() {
  const [isOpen, setIsOpen]     = useState(false);
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput]       = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const bottomRef   = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  // Scroll to latest message
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages, isLoading]);

  // Auto-focus textarea when opened
  useEffect(() => {
    if (isOpen) setTimeout(() => textareaRef.current?.focus(), 180);
  }, [isOpen]);

  // Resize textarea to fit content
  const resizeTextarea = useCallback(() => {
    const el = textareaRef.current;
    if (!el) return;
    el.style.height = 'auto';
    el.style.height = `${Math.min(el.scrollHeight, 96)}px`;
  }, []);

  const send = useCallback(async (override?: string) => {
    const text = (override ?? input).trim();
    if (!text || isLoading) return;

    const time = now();
    setInput('');
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto';
    }
    setMessages((prev) => [...prev, { role: 'user', text, time }]);
    setIsLoading(true);

    try {
      const { data } = await chatApi.send(text);
      const reply = (data as any)?.data?.reply ?? 'Maaf, tidak ada respons.';
      setMessages((prev) => [...prev, { role: 'assistant', text: reply, time: now() }]);
    } catch {
      setMessages((prev) => [
        ...prev,
        { role: 'assistant', text: '**Gagal terhubung.** Pastikan backend aktif dan `GROQ_API_KEY` terkonfigurasi.', time: now() },
      ]);
    } finally {
      setIsLoading(false);
    }
  }, [input, isLoading]);

  const onKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); send(); }
  };

  return (
    /* Anchor: fixed bottom-right corner — panel grows upward from here */
    <div className="fixed bottom-6 right-6 z-50 flex flex-col items-end gap-3">

      {/* ── Chat panel ──────────────────────────────────────────── */}
      <AnimatePresence>
        {isOpen && (
          <motion.div
            key="panel"
            initial={{ opacity: 0, scale: 0.94, y: 8 }}
            animate={{ opacity: 1, scale: 1,    y: 0 }}
            exit={{   opacity: 0, scale: 0.94, y: 8 }}
            transition={{ ease: [0.16, 1, 0.3, 1], duration: 0.22 }}
            style={{ transformOrigin: 'bottom right' }}
            className="w-[360px] rounded-2xl border border-white/9 bg-[#0e0e11] shadow-2xl shadow-black/70 overflow-hidden"
          >
            <div className="flex flex-col" style={{ height: '460px' }}>

              {/* ── HEADER ─────────────────────────────────────── */}
              <div className="shrink-0 flex items-center justify-between px-4 py-3.5 bg-linear-to-r from-purple-900/50 via-purple-800/20 to-[#0e0e11] border-b border-white/7">
                <div className="flex items-center gap-3">
                  <div className="w-9 h-9 rounded-xl bg-purple-600/25 border border-purple-500/30 flex items-center justify-center shrink-0">
                    <Bot size={17} className="text-purple-400" />
                  </div>
                  <div>
                    <p className="text-sm font-semibold text-zinc-100 leading-tight">One Log AI</p>
                    <div className="flex items-center gap-1.5 mt-0.5">
                      <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
                      <span className="text-[10px] text-zinc-500 flex items-center gap-1">
                        <Zap size={8} className="text-purple-500" />
                        Powered by Groq AI
                      </span>
                    </div>
                  </div>
                </div>
                <button
                  onClick={() => setIsOpen(false)}
                  className="w-7 h-7 flex items-center justify-center rounded-lg hover:bg-white/10 text-zinc-500 hover:text-zinc-300 transition-colors"
                >
                  <X size={15} />
                </button>
              </div>

              {/* ── MESSAGES ───────────────────────────────────── */}
              <div className="flex-1 overflow-y-auto px-4 py-4 space-y-4 min-h-0">

                {/* Empty / welcome state */}
                {messages.length === 0 && !isLoading && (
                  <motion.div
                    initial={{ opacity: 0, y: 4 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.05 }}
                    className="space-y-4"
                  >
                    {/* Greeting bubble */}
                    <div className="flex gap-2.5">
                      <div className="w-7 h-7 rounded-lg bg-purple-600/20 border border-purple-500/20 flex items-center justify-center shrink-0 mt-0.5">
                        <Bot size={13} className="text-purple-400" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="bg-white/4 border border-white/7 rounded-2xl rounded-tl-md px-4 py-3 overflow-hidden">
                          <p className="text-[13px] text-zinc-200 leading-relaxed wrap-break-word">
                            Hai! 👋 Saya <span className="font-semibold text-purple-300">One Log AI</span>, asisten Anda untuk memantau log, error, dan performa sistem.
                          </p>
                          <p className="text-[13px] text-zinc-400 mt-1.5 leading-relaxed">
                            Tanyakan apa saja tentang sistem Anda!
                          </p>
                        </div>
                        <p className="text-[10px] text-zinc-600 mt-1 ml-1">{now()}</p>
                      </div>
                    </div>

                    {/* Quick Questions */}
                    <div>
                      <p className="text-[11px] font-semibold text-zinc-500 uppercase tracking-wider mb-2 px-1">
                        Pertanyaan Cepat
                      </p>
                      <div className="space-y-2">
                        {SUGGESTIONS.map(({ icon: Icon, text }) => (
                          <button
                            key={text}
                            onClick={() => send(text)}
                            className="w-full flex items-center gap-3 px-3.5 py-2.5 rounded-xl bg-white/3 hover:bg-purple-500/10 border border-white/7 hover:border-purple-500/20 text-left transition-all group"
                          >
                            <div className="w-7 h-7 rounded-lg bg-purple-500/10 border border-purple-500/15 flex items-center justify-center shrink-0">
                              <Icon size={13} className="text-purple-400" />
                            </div>
                            <span className="text-[13px] text-zinc-400 group-hover:text-zinc-200 transition-colors">{text}</span>
                          </button>
                        ))}
                      </div>
                    </div>
                  </motion.div>
                )}

                {/* Conversation bubbles */}
                {messages.map((msg, i) => (
                  <motion.div
                    key={i}
                    initial={{ opacity: 0, y: 5 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.16 }}
                    className={`flex gap-2.5 ${msg.role === 'user' ? 'flex-row-reverse' : ''}`}
                  >
                    {/* Avatar */}
                    {msg.role === 'assistant' && (
                      <div className="w-7 h-7 rounded-lg bg-purple-600/20 border border-purple-500/20 flex items-center justify-center shrink-0 mt-0.5">
                        <Bot size={13} className="text-purple-400" />
                      </div>
                    )}

                    <div className={`flex flex-col min-w-0 ${msg.role === 'user' ? 'items-end' : 'items-start'} max-w-[84%]`}>
                      <div
                        className={`min-w-0 w-full rounded-2xl px-3.5 py-2.5 overflow-hidden ${
                          msg.role === 'user'
                            ? 'bg-purple-600/25 border border-purple-500/25 rounded-tr-md'
                            : 'bg-white/4 border border-white/7 rounded-tl-md'
                        }`}
                      >
                        {msg.role === 'user' ? (
                          <p className="text-[13px] text-zinc-100 leading-relaxed whitespace-pre-wrap wrap-break-word">{msg.text}</p>
                        ) : (
                          <div className="min-w-0 overflow-hidden">
                            <AiMarkdown text={msg.text} />
                          </div>
                        )}
                      </div>
                      <p className="text-[10px] text-zinc-600 mt-1 mx-1">{msg.time}</p>
                    </div>
                  </motion.div>
                ))}

                {/* Loading/typing indicator */}
                {isLoading && (
                  <motion.div
                    initial={{ opacity: 0, y: 5 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="flex gap-2.5"
                  >
                    <div className="w-7 h-7 rounded-lg bg-purple-600/20 border border-purple-500/20 flex items-center justify-center shrink-0">
                      <Bot size={13} className="text-purple-400" />
                    </div>
                    <div className="bg-white/4 border border-white/7 rounded-2xl rounded-tl-md px-4 py-3">
                      <TypingDots />
                    </div>
                  </motion.div>
                )}

                <div ref={bottomRef} />
              </div>

              {/* ── INPUT BAR ──────────────────────────────────── */}
              <div className="shrink-0 border-t border-white/7 px-4 py-3 bg-[#0e0e11]">
                <div className="flex items-center gap-2 bg-white/3 border border-white/8 focus-within:border-purple-500/40 rounded-xl px-3 py-2.5 transition-colors">
                  <textarea
                    ref={textareaRef}
                    rows={1}
                    value={input}
                    onChange={(e) => { setInput(e.target.value); resizeTextarea(); }}
                    onKeyDown={onKeyDown}
                    disabled={isLoading}
                    placeholder="Tanya tentang log, error, performa..."
                    className="flex-1 resize-none bg-transparent text-[13px] text-zinc-100 placeholder:text-zinc-600 focus:outline-none disabled:opacity-40 min-h-[22px] max-h-[96px] leading-relaxed self-center"
                    style={{ minHeight: '22px', maxHeight: '96px' }}
                  />
                  <button
                    onClick={() => send()}
                    disabled={isLoading || !input.trim()}
                    className="shrink-0 w-8 h-8 flex items-center justify-center rounded-lg bg-purple-600 hover:bg-purple-500 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
                    aria-label="Kirim"
                  >
                    <Send size={14} className="text-white" />
                  </button>
                </div>
                <p className="text-[10px] text-zinc-700 text-center mt-2">
                  Enter kirim · Shift+Enter baris baru
                </p>
              </div>

            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* ── TRIGGER BUTTON ──────────────────────────────────────── */}
      <motion.button
        onClick={() => setIsOpen((v) => !v)}
        whileTap={{ scale: 0.88 }}
        whileHover={{ scale: 1.06 }}
        transition={{ type: 'spring', stiffness: 400, damping: 17 }}
        className="w-14 h-14 rounded-full flex items-center justify-center bg-purple-600 hover:bg-purple-500 shadow-lg shadow-purple-500/30 transition-colors"
        aria-label="One Log AI"
      >
        <AnimatePresence mode="wait" initial={false}>
          {isOpen ? (
            <motion.span key="x" initial={{ rotate: -80, opacity: 0, scale: 0.5 }} animate={{ rotate: 0, opacity: 1, scale: 1 }} exit={{ rotate: 80, opacity: 0, scale: 0.5 }} transition={{ duration: 0.13 }}>
              <X size={20} className="text-white" />
            </motion.span>
          ) : (
            <motion.span key="m" initial={{ rotate: 80, opacity: 0, scale: 0.5 }} animate={{ rotate: 0, opacity: 1, scale: 1 }} exit={{ rotate: -80, opacity: 0, scale: 0.5 }} transition={{ duration: 0.13 }}>
              <MessageCircle size={20} className="text-white" />
            </motion.span>
          )}
        </AnimatePresence>
      </motion.button>
    </div>
  );
}
