import { NextResponse } from 'next/server';
import { cookies } from 'next/headers';

const BACKEND_URL = 'https://concierge-backend-80rb.onrender.com';

export async function POST(request) {
  try {
    const body = await request.json();
    // Support direct token injection (for profile switch)
    if (body._directToken) {
      const cookieStore = cookies();
      cookieStore.set('cai_token', body._directToken, {
        httpOnly: true,
        secure: process.env.NODE_ENV === 'production',
        sameSite: 'lax',
        maxAge: 60 * 60 * 24 * 30,
        path: '/',
      });
      if (body._directSlug) {
        cookieStore.set('cai_slug', body._directSlug, {
          httpOnly: false,
          secure: process.env.NODE_ENV === 'production',
          sameSite: 'lax',
          maxAge: 60 * 60 * 24 * 30,
          path: '/',
        });
      }
      return NextResponse.json({ token: body._directToken, slug: body._directSlug });
    }

    const r = await fetch(`${BACKEND_URL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    });
    const data = await r.json();
    if (!r.ok) return NextResponse.json(data, { status: r.status });

    const cookieStore = cookies();
    cookieStore.set('cai_token', data.token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 60 * 60 * 24 * 30,
      path: '/',
    });
    cookieStore.set('cai_slug', data.slug || '', {
      httpOnly: false,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 60 * 60 * 24 * 30,
      path: '/',
    });
    return NextResponse.json(data);
  } catch(e) {
    return NextResponse.json({ error: 'Connection error' }, { status: 500 });
  }
}
