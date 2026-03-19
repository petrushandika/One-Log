
import { motion, AnimatePresence } from 'framer-motion';
import { AlertTriangle, X } from 'lucide-react';

interface ConfirmDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  type?: 'danger' | 'warning' | 'info';
}

export default function ConfirmDialog({
  isOpen,
  onClose,
  onConfirm,
  title,
  message,
  confirmText = 'Confirm',
  cancelText = 'Cancel',
  type = 'danger',
}: ConfirmDialogProps) {
  const handleConfirm = () => {
    onConfirm();
    onClose();
  };

  const getColors = () => {
    switch (type) {
      case 'danger':
        return {
          icon: 'text-red-400 bg-red-500/10 border-red-500/20',
          confirm: 'bg-red-600 hover:bg-red-500',
        };
      case 'warning':
        return {
          icon: 'text-amber-400 bg-amber-500/10 border-amber-500/20',
          confirm: 'bg-amber-600 hover:bg-amber-500',
        };
      default:
        return {
          icon: 'text-blue-400 bg-blue-500/10 border-blue-500/20',
          confirm: 'bg-blue-600 hover:bg-blue-500',
        };
    }
  };

  const colors = getColors();

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50"
          />
          
          {/* Modal */}
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 20 }}
            className="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full max-w-md p-6 rounded-2xl bg-[#111113] border border-white/5 z-50 shadow-2xl"
          >
            {/* Header */}
            <div className="flex items-start gap-4 mb-5">
              <div className={`p-3 rounded-xl border ${colors.icon}`}>
                <AlertTriangle size={24} />
              </div>
              <div className="flex-1">
                <h3 className="text-lg font-semibold text-white mb-1">{title}</h3>
                <p className="text-sm text-zinc-400">{message}</p>
              </div>
              <button
                onClick={onClose}
                className="p-1.5 rounded-lg hover:bg-white/5 text-zinc-500 hover:text-zinc-300 transition-colors"
              >
                <X size={18} />
              </button>
            </div>

            {/* Actions */}
            <div className="flex gap-3 justify-end">
              <button
                onClick={onClose}
                className="px-4 py-2 rounded-xl bg-white/5 text-zinc-300 text-sm font-medium hover:bg-white/10 transition-colors"
              >
                {cancelText}
              </button>
              <button
                onClick={handleConfirm}
                className={`px-4 py-2 rounded-xl text-white text-sm font-medium transition-colors ${colors.confirm}`}
              >
                {confirmText}
              </button>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
}
