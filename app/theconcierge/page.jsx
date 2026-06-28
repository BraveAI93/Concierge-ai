import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';

export default async function TheConciergePage() {
  const cookieStore = cookies();
  const token = cookieStore.get('cai_token')?.value;
  if (token) redirect('/theconcierge/dashboard');
  redirect('/');
}
