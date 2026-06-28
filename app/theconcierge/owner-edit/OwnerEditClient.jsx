'use client';
import { useRouter } from 'next/navigation';
import OwnerEditProfile from '@/components/OwnerEditProfile';
import BravePAv2 from '@/components/BravePAv2';

export default function OwnerEditClient({ token, slug }) {
  const router = useRouter();
  return (
    <>
      <OwnerEditProfile
        token={token}
        slug={slug}
        lang="en"
        onBack={() => router.push('/theconcierge/dashboard')}
        onSaved={() => router.push('/theconcierge/dashboard')}
      />
      <BravePAv2 token={token} slug={slug} profile={null} leads={[]} convCount={0} />
    </>
  );
}
