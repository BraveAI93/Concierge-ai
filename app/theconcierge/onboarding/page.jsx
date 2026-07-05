'use client';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import Onboarding from '@/components/Onboarding';
import BusinessPagePreview from '@/components/BusinessPagePreview';
import { buildPrompt } from '@/lib/buildPrompt';

const BACKEND_URL = 'https://concierge-backend-80rb.onrender.com';

export default function OnboardingPage() {
  const router = useRouter();
  const [onboardData, setOnboardData] = useState(null);
  const [showPreview, setShowPreview] = useState(false);

  const handleComplete = async (data) => {
    const sp = buildPrompt(data, 'en');
    const slug = data.handle || ((data.name || 'me').toLowerCase().replace(/\s+/g, '-').replace(/[^a-z0-9-]/g, '') + '-' + Math.random().toString(36).slice(2, 6));

    let saveOk = false;
    try {
      const saveResp = await fetch(`${BACKEND_URL}/profile`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ slug, name: data.name, business: data.biz || data.name, profession: data.profession || 'Professional', location: data.loc || '', system_prompt: sp, profile_data: JSON.stringify(data), accent: '#c9a96e', active: true }),
      });
      if (saveResp.ok) { data._slug = slug; saveOk = true; }
    } catch(e) {}

    if (saveOk && data.ownerEmail && data.ownerPassword && data.ownerPassword.length >= 6) {
      try {
        const r = await fetch('/api/auth/signup', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ slug, email: data.ownerEmail, password: data.ownerPassword }),
        });
        const d = await r.json().catch(() => ({}));
        if (r.ok) {
          localStorage.setItem('cai_owner_token', d.token);
          localStorage.setItem('cai_owner_slug', d.slug);
          localStorage.setItem('ownerToken', d.token);
          localStorage.setItem('ownerProfile', JSON.stringify({ slug: d.slug }));
          router.push('/theconcierge/dashboard');
          return;
        }
      } catch(e) {}
    }

    setOnboardData({ ...data, _slug: slug });
    setShowPreview(true);
  };

  if (showPreview && onboardData) {
    return (
      <BusinessPagePreview
        data={onboardData}
        lang="en"
        onStartChat={() => router.push('/theconcierge')}
        onBack={() => router.push('/theconcierge')}
      />
    );
  }

  return (
    <Onboarding
      lang="en"
      onBack={() => router.push('/theconcierge')}
      onComplete={handleComplete}
    />
  );
}
