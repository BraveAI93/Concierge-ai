import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import AddProfileClient from './AddProfileClient';

export default async function AddProfilePage() {
  const cookieStore = cookies();
  const token = cookieStore.get('cai_token')?.value;
  if (!token) redirect('/theconcierge/owner-auth');
  return <AddProfileClient token={token} />;
}
