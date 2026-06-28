import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import OwnerEditClient from './OwnerEditClient';

export default async function OwnerEditPage() {
  const cookieStore = cookies();
  const token = cookieStore.get('cai_token')?.value;
  const slug = cookieStore.get('cai_slug')?.value;
  if (!token) redirect('/theconcierge/owner-auth');
  return <OwnerEditClient token={token} slug={slug || ''} />;
}
