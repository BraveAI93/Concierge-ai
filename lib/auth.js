import { cookies } from 'next/headers';

export const TOKEN_COOKIE = 'cai_token';
export const SLUG_COOKIE = 'cai_slug';

export function setAuthCookies(token, slug) {
  const cookieStore = cookies();
  cookieStore.set(TOKEN_COOKIE, token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'lax',
    maxAge: 60 * 60 * 24 * 30,
    path: '/',
  });
  if (slug) {
    cookieStore.set(SLUG_COOKIE, slug, {
      httpOnly: false,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 60 * 60 * 24 * 30,
      path: '/',
    });
  }
}

export function getAuthFromCookies() {
  const cookieStore = cookies();
  const token = cookieStore.get(TOKEN_COOKIE)?.value || null;
  const slug = cookieStore.get(SLUG_COOKIE)?.value || null;
  return { token, slug };
}

export function clearAuthCookies() {
  const cookieStore = cookies();
  cookieStore.delete(TOKEN_COOKIE);
  cookieStore.delete(SLUG_COOKIE);
}
