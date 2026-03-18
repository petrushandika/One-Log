import { useState } from 'react';
import { Link, useLocation, Outlet, useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { LayoutGrid, FileText, Settings, ShieldAlert, LogOut, Terminal, Menu, X } from 'lucide-react';

export default function Layout() {
  const location = useLocation();
  const navigate = useNavigate();
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  const handleLogout = () => {
     localStorage.removeItem('token');
     navigate('/login');
  };

  const menuItems = [
    { title: 'Overview', icon: LayoutGrid, path: '/' },
    { title: 'Logs', icon: FileText, path: '/logs' },
    { title: 'Sources', icon: Settings, path: '/sources' },
    { title: 'Audit Trail', icon: ShieldAlert, path: '/audit' },
  ];

  return (
    <div className="flex bg-[#09090b] min-h-screen">
      {/* Sidebar - Desktop */}
      <aside className="hidden lg:flex fixed inset-y-0 left-0 w-64 border-r border-white/5 bg-[#0c0c0e]/80 backdrop-blur-md p-6 flex-col justify-between z-30">
        <div>
          <div className="flex items-center gap-2 mb-10">
            <div className="p-2 rounded-xl bg-purple-500/20 text-purple-400">
              <Terminal size={24} />
            </div>
            <h1 className="text-xl font-bold tracking-tight text-white">One Log</h1>
          </div>

          <nav className="space-y-1">
            {menuItems.map((item) => {
              const isActive = location.pathname === item.path;
              return (
                <Link
                  key={item.path}
                  to={item.path}
                  className={`flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium transition-all duration-200 ${
                    isActive
                      ? 'bg-purple-500/10 text-purple-400 border border-purple-500/20'
                      : 'text-zinc-400 hover:bg-white/3 hover:text-zinc-200 border border-transparent'
                  }`}
                >
                  <item.icon size={18} />
                  <span>{item.title}</span>
                </Link>
              );
            })}
          </nav>
        </div>

        <div>
          <div className="flex items-center gap-3 p-3 mb-4 rounded-xl bg-white/2 border border-white/5">
            <img src="https://avatar.vercel.sh/admin" alt="Admin" className="w-9 h-9 rounded-lg" />
            <div>
              <p className="text-sm font-semibold text-zinc-100">Administrator</p>
              <p className="text-xs text-zinc-500">Full Access</p>
            </div>
          </div>
          <button 
            onClick={handleLogout}
            className="flex items-center justify-center w-full gap-2 px-4 py-3 text-sm font-semibold transition-all duration-200 border rounded-xl bg-red-500/10 hover:bg-red-500 text-red-500 hover:text-white border-red-500/20 shadow-lg shadow-red-500/5"
          >
            <LogOut size={16} />
            Sign Out
          </button>
        </div>
      </aside>

      {/* Sidebar - Mobile Drawer */}
      <AnimatePresence>
        {isSidebarOpen && (
          <>
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setIsSidebarOpen(false)}
              className="fixed inset-0 bg-black/60 backdrop-blur-sm z-40 lg:hidden"
            />
            <motion.aside
              initial={{ x: -250 }}
              animate={{ x: 0 }}
              exit={{ x: -250 }}
              transition={{ ease: 'easeInOut', duration: 0.3 }}
              className="fixed inset-y-0 left-0 w-64 border-r border-white/5 bg-[#0c0c0e] p-6 flex flex-col justify-between z-50 lg:hidden"
            >
              <div>
                <div className="flex items-center justify-between mb-10">
                  <div className="flex items-center gap-2">
                    <div className="p-2 rounded-xl bg-purple-500/20 text-purple-400">
                      <Terminal size={22} />
                    </div>
                    <span className="text-lg font-bold text-white">One Log</span>
                  </div>
                  <button onClick={() => setIsSidebarOpen(false)} className="text-zinc-400 hover:text-zinc-200">
                    <X size={20} />
                  </button>
                </div>

                <nav className="space-y-1">
                  {menuItems.map((item) => {
                    const isActive = location.pathname === item.path;
                    return (
                      <Link
                        key={item.path}
                        to={item.path}
                        onClick={() => setIsSidebarOpen(false)}
                        className={`flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium ${
                          isActive
                            ? 'bg-purple-500/10 text-purple-400 border border-purple-500/20'
                            : 'text-zinc-400 hover:bg-white/3'
                        }`}
                      >
                        <item.icon size={18} />
                        <span>{item.title}</span>
                      </Link>
                    );
                  })}
                </nav>
              </div>

              <div>
                <button 
                  onClick={handleLogout}
                  className="flex items-center justify-center w-full gap-2 px-4 py-3 text-sm font-semibold transition-all duration-200 border rounded-xl bg-red-500/10 hover:bg-red-500 text-red-500 hover:text-white border-red-500/20 shadow-lg shadow-red-500/5"
                >
                  <LogOut size={16} />
                  Sign Out
                </button>
              </div>
            </motion.aside>
          </>
        )}
      </AnimatePresence>

      {/* Main Content Pane */}
      <div className="flex-1 lg:ml-64 w-full">
        {/* Top Navbar */}
        <header className="sticky top-0 z-20 p-4 md:p-6 border-b border-white/5 bg-[#09090b]/40 backdrop-blur-md flex items-center justify-between">
          <div className="flex items-center gap-3">
            <button onClick={() => setIsSidebarOpen(true)} className="p-2 -ml-2 rounded-lg hover:bg-white/5 text-zinc-400 lg:hidden">
              <Menu size={20} />
            </button>
            <div className="flex items-center gap-1.5 text-xs text-zinc-500 font-medium">
              <Link to="/" className="hover:text-zinc-200 transition-colors">Home</Link>
              {location.pathname !== '/' && (
                <>
                  <span className="text-zinc-600">&gt;</span>
                  <span className="text-zinc-200 font-semibold">
                    {location.pathname === '/logs' ? 'Logs' : 
                     location.pathname === '/sources' ? 'Sources' : 
                     location.pathname === '/audit' ? 'Audit Trail' : 'Page'}
                  </span>
                </>
              )}
            </div>
          </div>
          <div className="flex items-center gap-4">
            <span className="px-3 py-1 text-xs font-semibold rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
              Live
            </span>
          </div>
        </header>

        {/* Routers Outlet View */}
        <main className="p-4 md:p-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
