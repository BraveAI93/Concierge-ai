import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const token = request.cookies.get('cai_token')?.value;

  if (pathname.startsWith('/theconcierge/dashboard') && !token) {
    const url = request.nextUrl.clone();
    url.pathname = '/theconcierge/owner-auth';
    return NextResponse.redirect(url);
  }

  if (pathname === '/theconcierge' && token) {
    const url = request.nextUrl.clone();
    url.pathname = '/theconcierge/dashboard';
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/theconcierge/:path*'],
};
