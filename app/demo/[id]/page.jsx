'use client';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';
import { DEMOS } from '@/lib/constants';
import Chat from '@/components/Chat';

export default function DemoPage({ params }) {
  const router = useRouter();
  const demo = DEMOS[params.id];

  useEffect(() => {
    if (!demo) router.push('/theconcierge');
    else sessionStorage.setItem('cai_consent', 'demo-' + Date.now());
  }, []);

  if (!demo) return null;

  return (
    <Chat
      profile={demo}
      systemPrompt={demo.systemPrompt}
      profileId={params.id}
      onBack={() => router.push('/theconcierge')}
      lang="en"
    />
  );
}
