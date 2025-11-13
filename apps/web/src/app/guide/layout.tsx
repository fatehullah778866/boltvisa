'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

const nav = [
  { href: '/guide', label: 'Overview' },
  { href: '/guide/getting-started', label: 'Getting Started' },
  { href: '/guide/dashboard', label: 'Dashboard' },
  { href: '/guide/applications', label: 'Applications' },
  { href: '/guide/visa-categories', label: 'Visa Categories' },
  { href: '/guide/troubleshooting', label: 'Troubleshooting' },
  { href: '/guide/faq', label: 'FAQ' },
];

export default function GuideLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  return (
    <main className="min-h-screen bg-gray-50">
      <div className="mx-auto max-w-6xl px-6 py-8 grid grid-cols-1 md:grid-cols-[240px_1fr] gap-6">
        <aside className="rounded-2xl bg-white p-4 shadow-sm ring-1 ring-gray-100 h-max">
          <h2 className="text-sm font-semibold text-gray-700 mb-3">Guide</h2>
          <nav className="space-y-1">
            {nav.map(item => {
              const active = pathname === item.href;
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={`block rounded-lg px-3 py-2 text-sm ${active ? 'bg-gray-100 font-medium' : 'hover:bg-gray-50 text-gray-700'}`}
                >
                  {item.label}
                </Link>
              );
            })}
          </nav>
        </aside>
        <section className="rounded-2xl bg-white p-6 shadow-sm ring-1 ring-gray-100">
          {children}
        </section>
      </div>
    </main>
  );
}
