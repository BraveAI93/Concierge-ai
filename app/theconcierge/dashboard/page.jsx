import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import DashboardClient from './DashboardClient';

export default async function DashboardPage() {
  const cookieStore = cookies();
  const token = cookieStore.get('cai_token')?.value;
  const slug = cookieStore.get('cai_slug')?.value;
  if (!token) redirect('/theconcierge/owner-auth');
  return <DashboardClient token={token} slug={slug || ''} />;
}
