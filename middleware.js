import { NextResponse } from 'next/server';

const PROTECTED = [
  '/theconcierge/dashboard',
  '/theconcierge/owner-edit',
  '/theconcierge/add-profile',
];

export function middleware(request) {
  const { pathname } = request.nextUrl;
  const isProtected = PROTECTED.some(p => pathname.startsWith(p));
  if (!isProtected) return NextResponse.next();
  const token = request.cookies.get('cai_token')?.value;
  if (!token) {
    const url = request.nextUrl.clone();
    url.pathname = '/theconcierge/owner-auth';
    return NextResponse.redirect(url);
  }
  return NextResponse.next();
}

export const config = {
  matcher: ['/theconcierge/:path*'],
};
