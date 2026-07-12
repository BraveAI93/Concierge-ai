import { useEffect, useState } from 'react';
import { BACKEND_URL } from './constants';

export function useFeatureFlags() {
  const [flags, setFlags] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;
    fetch(`${BACKEND_URL}/flags`)
      .then(r => r.ok ? r.json() : { flags: [] })
      .then(d => { if (!cancelled) setFlags(d.flags || []); })
      .catch(() => { if (!cancelled) setFlags([]); })
      .finally(() => { if (!cancelled) setLoading(false); });
    return () => { cancelled = true; };
  }, []);

  const getState = (name) => flags?.find(f => f.name === name)?.state || 'UNKNOWN';

  return { flags, loading, getState };
}
