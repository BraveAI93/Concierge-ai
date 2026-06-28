import { NextResponse } from 'next/server';
import { cookies } from 'next/headers';

export async function GET() {
  const cookieStore = cookies();
  const token = cookieStore.get('cai_token')?.value || null;
  const slug = cookieStore.get('cai_slug')?.value || null;
  return NextResponse.json({ token, slug });
}
