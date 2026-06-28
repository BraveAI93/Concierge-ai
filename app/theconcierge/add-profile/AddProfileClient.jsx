'use client';
import { useRouter } from 'next/navigation';
import Onboarding from '@/components/Onboarding';
import { buildPrompt } from '@/lib/buildPrompt';

const BACKEND_URL = 'https://concierge-backend-80rb.onrender.com';

export default function AddProfileClient({ token }) {
  const router = useRouter();

  const handleComplete = async (data) => {
    const sp = buildPrompt(data, 'en');
    const newSlug = data.handle || ((data.name || 'me').toLowerCase().replace(/\s+/g, '-').replace(/[^a-z0-9-]/g, '') + '-' + Math.random().toString(36).slice(2, 6));
    try {
      const r = await fetch(`${BACKEND_URL}/owner/profiles`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify({ slug: newSlug, name: data.name, business: data.biz || data.name, profession: data.profession || 'Professional', location: data.loc || '', system_prompt: sp, profile_data: JSON.stringify(data), accent: '#c9a96e', active: true }),
      });
      const d = await r.json().catch(() => ({}));
      if (r.ok && d.token) {
        await fetch('/api/auth/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ _directToken: d.token, _directSlug: d.slug || newSlug }),
        });
        localStorage.setItem('cai_owner_token', d.token);
        localStorage.setItem('ownerToken', d.token);
      }
    } catch(e) { alert(`Could not save profile: ${e.message}`); }
    router.push('/theconcierge/dashboard');
  };

  return (
    <Onboarding
      lang="en"
      onBack={() => router.push('/theconcierge/dashboard')}
      onComplete={handleComplete}
    />
  );
}
