import { NextResponse } from 'next/server';
import { cookies } from 'next/headers';

export async function POST() {
  const cookieStore = cookies();
  cookieStore.delete('cai_token');
  cookieStore.delete('cai_slug');
  return NextResponse.json({ ok: true });
}
