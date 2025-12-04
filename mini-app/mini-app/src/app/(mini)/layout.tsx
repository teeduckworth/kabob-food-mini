import { BottomNav } from '@/components/BottomNav';
import { CloseMiniAppButton } from '@/components/CloseMiniAppButton';
import { AppProviders } from '@/providers/AppProviders';

export default function MiniLayout({ children }: { children: React.ReactNode }) {
  return (
    <AppProviders>
      <div className="relative max-w-2xl mx-auto px-4 pb-32 pt-6">
        <CloseMiniAppButton />
        {children}
      </div>

      <BottomNav />
    </AppProviders>
  );
}
