'use client';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import OwnerDashboard from '@/components/OwnerDashboard';
import BravePAv2 from '@/components/BravePAv2';

export default function DashboardClient({ token, slug }) {
  const router = useRouter();
  const [paData, setPaData] = useState({ profile: null, leads: [], convCount: 0 });

  const handleLogout = async () => {
    await fetch('/api/auth/logout', { method: 'POST' });
    localStorage.removeItem('cai_owner_token');
    localStorage.removeItem('cai_owner_slug');
    localStorage.removeItem('ownerToken');
    localStorage.removeItem('ownerProfile');
    router.push('/theconcierge');
  };

  const switchProfile = async (targetSlug) => {
    const BACKEND_URL = 'https://concierge-backend-80rb.onrender.com';
    const r = await fetch(`${BACKEND_URL}/owner/switch-profile`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify({ slug: targetSlug }),
    });
    const d = await r.json().catch(() => ({}));
    if (r.ok && d.token) {
      await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ _directToken: d.token, _directSlug: d.slug }),
      });
      window.location.reload();
    }
  };

  return (
    <>
      <OwnerDashboard
        token={token}
        slug={slug}
        onBack={() => router.push('/theconcierge')}
        onLogout={handleLogout}
        onEditProfile={() => router.push('/theconcierge/owner-edit')}
        onDataLoaded={(p, l, c) => setPaData({ profile: p, leads: l, convCount: c })}
        onAddProfile={() => router.push('/theconcierge/add-profile')}
        onSwitchProfile={switchProfile}
      />
      <BravePAv2
        token={token}
        slug={slug}
        profile={paData.profile}
        leads={paData.leads}
        convCount={paData.convCount}
      />
    </>
  );
}
