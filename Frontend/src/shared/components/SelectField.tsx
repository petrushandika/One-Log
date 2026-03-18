import { ChevronDown } from 'lucide-react';

interface SelectFieldProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
  children: React.ReactNode;
  /** Extra wrapper className */
  wrapperClassName?: string;
}

/**
 * Consistent select with properly spaced custom arrow chevron.
 * Replaces every raw <select> in the app for uniform styling.
 */
export default function SelectField({ className = '', wrapperClassName = '', children, ...props }: SelectFieldProps) {
  return (
    <div className={`relative ${wrapperClassName}`}>
      <select
        className={`w-full appearance-none bg-white/3 border border-white/8 text-zinc-200 rounded-xl px-3 py-2 pr-8 text-sm focus:outline-none focus:border-purple-500/40 cursor-pointer transition-colors ${className}`}
        {...props}
      >
        {children}
      </select>
      <ChevronDown
        size={13}
        className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-500 pointer-events-none"
      />
    </div>
  );
}
