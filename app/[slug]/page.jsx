import { notFound } from 'next/navigation';
import PublicProfileClient from './PublicProfileClient';

const BACKEND_URL = 'https://concierge-backend-80rb.onrender.com';
const RESERVED = ['api', 'theconcierge', '_next', 'favicon.ico', 'robots.txt', 'sitemap.xml', 'demo'];

export default async function SlugPage({ params }) {
  const { slug } = params;
  if (RESERVED.includes(slug)) notFound();

  let profileData = null;
  try {
    const r = await fetch(`${BACKEND_URL}/profile/${slug}`, { next: { revalidate: 60 } });
    if (r.ok) profileData = await r.json();
  } catch(e) {}

  if (!profileData || !profileData.slug) notFound();

  let pdata = {};
  try { pdata = JSON.parse(profileData.profile_data || '{}'); } catch(e) {}

  const data = {
    ...pdata,
    _slug: profileData.slug,
    name: profileData.name,
    profession: profileData.profession,
    loc: profileData.location,
    business: profileData.business,
    accent: profileData.accent || '#c9a96e',
  };

  return <PublicProfileClient data={data} systemPrompt={profileData.system_prompt} slug={slug} />;
}
