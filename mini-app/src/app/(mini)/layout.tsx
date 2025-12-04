import { CloseMiniAppButton } from '@/components/CloseMiniAppButton';
import { AppProviders } from '@/providers/AppProviders';

export default function MiniLayout({ children }: { children: React.ReactNode }) {
  return (
    <AppProviders>
      <div className="relative mx-auto max-w-2xl px-4 pb-12 pt-6">
        <CloseMiniAppButton />
        {children}
      </div>
    </AppProviders>
  );
}
