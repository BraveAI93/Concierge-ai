'use client';
import { useState, useEffect } from 'react';
import { BACKEND_URL, STEPS } from '@/lib/constants';

function getActiveStepIds(modes) {
  if (!modes || !modes.length) return null;
  const hasBB = modes.some(m => m.includes('B&B'));
  const hasPerformer = modes.some(m => m.includes('Performer'));
  const hasCreator = modes.some(m => m.includes('Creator'));
  const hasWellness = modes.some(m => m.includes('Wellness'));
  const hasChef = modes.some(m => m.includes('Chef') || m.includes('Hospitality'));
  const hasGeneral = modes.some(m => m.includes('Consultant') || m.includes('Educator') || m.includes('Trades'));
  const common = ['modes', 'basics', 'tone', 'ex'];
  const tail = ['handle', 'account'];
  let mid = [];
  if (hasBB && !hasWellness && !hasPerformer && !hasCreator) {
    mid = ['bb', 'location', 'legal', 'contact', 'media', 'extra'];
  } else if (hasWellness && !hasPerformer && !hasCreator) {
    mid = ['svcs', 'location', 'sensitive', 'legal', 'contact', 'media', 'extra'];
  } else if (hasPerformer && !hasCreator && !hasWellness) {
    mid = ['svcs', 'perform', 'media', 'sensitive', 'contact', 'extra'];
  } else if (hasCreator && !hasPerformer && !hasWellness) {
    mid = ['creator', 'media', 'sensitive', 'contact', 'extra'];
  } else {
    const s = new Set();
    if (hasWellness || hasChef || hasGeneral || hasPerformer) s.add('svcs');
    if (hasPerformer) s.add('perform');
    if (hasCreator) s.add('creator');
    if (hasBB) s.add('bb');
    if (hasWellness || hasChef || hasGeneral || hasBB) s.add('location');
    s.add('sensitive');
    if (hasWellness || hasBB) s.add('legal');
    s.add('contact'); s.add('media'); s.add('extra');
    ['svcs', 'perform', 'creator', 'bb', 'location', 'sensitive', 'legal', 'contact', 'media', 'extra']
      .forEach(id => { if (s.has(id)) mid.push(id); });
  }
  return [...common, ...mid, ...tail];
}

export default function Onboarding({ onComplete, onBack, lang }) {
  const [step, setStep] = useState(0);
  const [data, setData] = useState({});
  const [services, setServices] = useState([{name:'',durNum:'',durUnit:'min',priceNum:'',currency:'£',desc:''}]);
  const [handle, setHandle] = useState('');
  const [handleStatus, setHandleStatus] = useState('');
  const [handleTouched, setHandleTouched] = useState(false);
  const [geoStatus, setGeoStatus] = useState('');
  const [obUploading, setObUploading] = useState(false);
  const [obUploadPct, setObUploadPct] = useState(0);
  const [mediaUploads, setMediaUploads] = useState([]);
  const [svcImporting, setSvcImporting] = useState(false);
  const [svcPreview, setSvcPreview] = useState(null);
  const L = lang || 'en';

  const _activeIds = getActiveStepIds(data.modes || []);
  const _stepDefs = _activeIds ? _activeIds.map(id => STEPS.find(s => s.id === id)).filter(Boolean) : STEPS;
  const total = _stepDefs.length;
  const safeStep = Math.min(step, total - 1);
  const cur = _stepDefs[safeStep];
  const set = (k, v) => setData(d => ({...d, [k]: v}));
  const tog = (f, v) => { const a = data[f] || []; set(f, a.includes(v) ? a.filter(x => x !== v) : [...a, v]); };

  const slugify = (s) => s.toLowerCase().trim().replace(/\s+/g,'-').replace(/[^a-z0-9-]/g,'').replace(/-+/g,'-').replace(/^-|-$/g,'');

  const obUploadFile = async (file) => {
    setObUploading(true); setObUploadPct(0);
    const fd = new FormData(); fd.append('file', file);
    try {
      await new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.upload.onprogress = e => { if (e.lengthComputable) setObUploadPct(Math.round(e.loaded / e.total * 100)); };
        xhr.onload = () => {
          if (xhr.status === 200) { const d = JSON.parse(xhr.responseText); setMediaUploads(u => [...u, d.url]); resolve(); }
          else reject(new Error('Upload failed'));
        };
        xhr.onerror = () => reject(new Error('Upload failed'));
        xhr.open('POST', `${BACKEND_URL}/media/upload`);
        xhr.send(fd);
      });
    } catch(e) { alert('Upload failed: ' + e.message); }
    setObUploading(false); setObUploadPct(0);
  };

  const obImportServices = async (file) => {
    setSvcImporting(true); setSvcPreview(null);
    try {
      const b64 = await new Promise((res, rej) => { const r = new FileReader(); r.onload = e => res(e.target.result.split(',')[1]); r.onerror = rej; r.readAsDataURL(file); });
      const mime = file.type || 'image/jpeg';
      const ownerTok = localStorage.getItem('cai_owner_token') || localStorage.getItem('ownerToken') || '';
      const resp = await fetch(`${BACKEND_URL}/ai/import-services`, { method: 'POST', headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + ownerTok }, body: JSON.stringify({ image_base64: b64, mime_type: mime }) });
      const d = await resp.json();
      if (resp.ok && d.services) {
        const mapped = d.services.map(s => ({ name: s.name || '', durNum: s.duration_minutes ? String(s.duration_minutes) : '', durUnit: 'min', priceNum: s.price ? String(s.price) : '', currency: s.currency || '£', desc: s.description || '' }));
        setSvcPreview(mapped);
      } else { alert(d.error || 'Could not extract services'); }
    } catch(e) { alert('Import failed: ' + e.message); }
    setSvcImporting(false);
  };

  const useMyLocation = () => {
    if (!navigator.geolocation) { setGeoStatus('unsupported'); return; }
    setGeoStatus('locating');
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        set('lat', pos.coords.latitude); set('lng', pos.coords.longitude);
        setGeoStatus('done');
        fetch(`https://nominatim.openstreetmap.org/reverse?format=json&lat=${pos.coords.latitude}&lon=${pos.coords.longitude}`)
          .then(r => r.json())
          .then(d => { const place = d.address?.city || d.address?.town || d.address?.village || d.display_name?.split(',')[0] || ''; if (place) set('loc', place); })
          .catch(() => {});
      },
      () => setGeoStatus('denied'),
      { timeout: 8000 }
    );
  };

  useEffect(() => {
    if (cur && cur.id === 'handle' && !handleTouched) {
      const base = slugify(data.biz || data.name || '');
      if (base) { setHandle(base); checkHandle(base); }
    }
  }, [step]);

  const checkHandle = async (h) => {
    const s = slugify(h);
    if (!s || s.length < 3) { setHandleStatus('short'); return; }
    setHandleStatus('checking');
    try {
      const r = await fetch(`${BACKEND_URL}/check-slug/${s}`);
      const d = await r.json();
      setHandleStatus(d.available ? 'available' : (d.reason === 'reserved' ? 'reserved' : 'taken'));
    } catch(e) { setHandleStatus('available'); }
  };

  let _handleTimer = null;
  const onHandleChange = (v) => {
    setHandleTouched(true);
    const s = slugify(v);
    setHandle(s);
    if (typeof window !== 'undefined') {
      if (window._handleTimer) clearTimeout(window._handleTimer);
      window._handleTimer = setTimeout(() => checkHandle(s), 500);
    }
  };

  const canNext = () => {
    if (!cur) return false;
    if (cur.type === 'chips') return !!data[cur.field];
    if (cur.type === 'form') {
      if (cur.id === 'perform' || cur.id === 'creator' || cur.id === 'contact' || cur.id === 'media' || cur.id === 'account' || cur.id === 'bb') return true;
      return !!data[cur.fields[0].key]?.trim();
    }
    if (cur.type === 'chips_multi') {
      if (cur.id === 'sensitive' || cur.id === 'tone' || cur.id === 'legal') return true;
      return (data[cur.field] || []).length > 0;
    }
    if (cur.type === 'textarea') {
      if (cur.id === 'extra') return true;
      return data[cur.field]?.trim().length > 10;
    }
    if (cur.type === 'services') return services.some(s => s.name.trim());
    if (cur.type === 'handle') return handleStatus === 'available';
    if (cur.type === 'location') return true;
    return true;
  };

  const next = () => {
    if (cur.type === 'services') set('services', services.filter(s => s.name.trim()));
    if (cur.id === 'media' && mediaUploads.length > 0) set('media_gallery', [...(data.media_gallery || []), ...mediaUploads]);
    if (safeStep < total - 1) setStep(safeStep + 1);
    else onComplete({...data, services: services.filter(s => s.name.trim()), handle, media_gallery: [...(data.media_gallery || []), ...mediaUploads]});
  };

  const chip = (label, active, onClick) => (
    <button key={label} onClick={onClick} style={{padding:'7px 12px',borderRadius:20,cursor:'pointer',fontSize:12,fontFamily:"'Jost',sans-serif",fontWeight:active?500:300,border:active?'1px solid #c9a96e':'1px solid rgba(201,169,110,0.16)',background:active?'rgba(201,169,110,0.15)':'rgba(201,169,110,0.04)',color:active?'#e8c878':'rgba(232,220,200,0.52)',transition:'all 0.15s',marginBottom:4}}>{label}</button>
  );

  if (!cur) return null;

  return (
    <div style={{minHeight:'100vh',background:'#0c0a08',display:'flex',flexDirection:'column',alignItems:'center',justifyContent:'center',padding:'20px 16px'}}>
      <div style={{width:'100%',maxWidth:490}}>
        <div style={{textAlign:'center',marginBottom:24}}>
          <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",fontWeight:300,letterSpacing:'0.18em',textTransform:'uppercase',color:'rgba(201,169,110,0.42)',marginBottom:5}}>✦ The Concierge</div>
          <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.28)',letterSpacing:'0.07em'}}>Step {safeStep+1} / {total}</div>
        </div>
        <div style={{height:2,background:'rgba(201,169,110,0.1)',borderRadius:2,marginBottom:22}}>
          <div style={{height:'100%',width:`${(safeStep/Math.max(total-1,1))*100}%`,background:'linear-gradient(90deg,#c9a96e,#e8c878)',borderRadius:2,transition:'width 0.4s ease'}}/>
        </div>

        <div style={{animation:'fadeUp 0.3s ease forwards',background:'rgba(12,10,8,0.78)',border:'1px solid rgba(201,169,110,0.2)',borderRadius:20,padding:'24px 22px',marginBottom:11,backdropFilter:'blur(18px)',boxShadow:'0 0 32px rgba(201,169,110,0.07),0 12px 40px rgba(0,0,0,0.5)'}}>
          <div style={{fontSize:20,fontWeight:500,color:'#e8dcc8',marginBottom:5}}>{typeof cur.title==='object'?cur.title[L]:cur.title}</div>
          <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",fontWeight:300,color:'rgba(232,220,200,0.4)',lineHeight:1.6,marginBottom:17}}>{typeof cur.sub==='object'?cur.sub[L]:cur.sub}</div>

          {cur.type==='chips'&&<div style={{display:'flex',flexWrap:'wrap',gap:7}}>{cur.options.map(o=>chip(o,data[cur.field]===o,()=>set(cur.field,o)))}</div>}

          {cur.type==='chips_multi'&&<div>
            <div style={{display:'flex',flexWrap:'wrap',gap:7}}>{cur.options.map(o=>chip(o,(data[cur.field]||[]).includes(o),()=>tog(cur.field,o)))}</div>
            {cur.id==='sensitive'&&<div style={{marginTop:14}}>
              <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.08em',textTransform:'uppercase',color:'rgba(201,169,110,0.4)',marginBottom:6}}>{L==='it'?'Altro (specifico per il tuo lavoro)':'Other (specific to your work)'}</div>
              <input value={data.sensitiveOther||''} onChange={e=>set('sensitiveOther',e.target.value)} placeholder={L==='it'?'es. Allergie ad oli essenziali...':'e.g. Essential oil allergies...'} style={{width:'100%',padding:'9px 12px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.16)',borderRadius:8,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif",outline:'none'}}/>
            </div>}
          </div>}

          {cur.type==='form'&&cur.fields.map(f=>(
            <div key={f.key} style={{marginBottom:11}}>
              <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.42)',marginBottom:6}}>{typeof f.label==='object'?f.label[L]:f.label}</div>
              <input value={data[f.key]||''} onChange={e=>set(f.key,e.target.value)} placeholder={typeof f.ph==='object'?f.ph[L]:f.ph} style={{width:'100%',padding:'9px 12px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.16)',borderRadius:8,color:'#e8dcc8',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:300,outline:'none'}}/>
            </div>
          ))}

          {cur.type==='form'&&cur.id==='media'&&<div style={{marginTop:4}}>
            <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.42)',marginBottom:8}}>Or upload directly</div>
            <button onClick={()=>document.getElementById('ob-media-upload').click()} disabled={obUploading} style={{width:'100%',padding:'9px 0',background:'rgba(201,169,110,0.07)',border:'1px dashed rgba(201,169,110,0.25)',borderRadius:20,cursor:obUploading?'default':'pointer',color:'rgba(201,169,110,0.7)',fontSize:12,fontFamily:"'Jost',sans-serif"}}>
              {obUploading?`Uploading… ${obUploadPct}%`:'⬆ Upload from device'}
            </button>
            <input id="ob-media-upload" type="file" accept="image/*,video/*" style={{display:'none'}} onChange={e=>{if(e.target.files[0])obUploadFile(e.target.files[0]);e.target.value='';}}/>
            {obUploading&&<div style={{marginTop:6,height:3,background:'rgba(201,169,110,0.1)',borderRadius:2}}><div style={{height:'100%',width:obUploadPct+'%',background:'linear-gradient(90deg,#c9a96e,#e8c878)',borderRadius:2,transition:'width 0.3s'}}/></div>}
            {mediaUploads.length>0&&<div style={{display:'grid',gridTemplateColumns:'repeat(3,1fr)',gap:6,marginTop:10}}>
              {mediaUploads.map((url,i)=>(
                <div key={i} style={{position:'relative',paddingBottom:'100%',background:'rgba(201,169,110,0.06)',borderRadius:8,overflow:'hidden'}}>
                  <img src={url} style={{position:'absolute',top:0,left:0,width:'100%',height:'100%',objectFit:'cover',borderRadius:8}} onError={e=>e.target.style.display='none'}/>
                  <button onClick={()=>setMediaUploads(u=>u.filter((_,idx)=>idx!==i))} style={{position:'absolute',top:3,right:3,background:'rgba(0,0,0,0.6)',border:'none',borderRadius:'50%',width:18,height:18,cursor:'pointer',color:'#e8dcc8',fontSize:10,lineHeight:'18px',textAlign:'center'}}>×</button>
                </div>
              ))}
            </div>}
          </div>}

          {cur.type==='textarea'&&<textarea value={data[cur.field]||''} onChange={e=>set(cur.field,e.target.value)} placeholder={typeof cur.ph==='object'?cur.ph[L]:cur.ph} rows={4} style={{width:'100%',padding:'10px 12px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.16)',borderRadius:8,color:'#e8dcc8',fontSize:13,fontFamily:"'Cormorant Garamond',serif",lineHeight:1.65,outline:'none'}}/>}

          {cur.type==='services'&&<div>
            {services.map((s,i)=>(
              <div key={i} style={{marginBottom:11,padding:'11px 13px',background:'rgba(201,169,110,0.04)',border:'1px solid rgba(201,169,110,0.09)',borderRadius:10}}>
                <div style={{display:'flex',justifyContent:'space-between',marginBottom:8}}>
                  <span style={{fontSize:9,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.38)',letterSpacing:'0.08em',textTransform:'uppercase'}}>Service {i+1}</span>
                  {services.length>1&&<button onClick={()=>setServices(services.filter((_,idx)=>idx!==i))} style={{background:'none',border:'none',cursor:'pointer',color:'rgba(200,80,80,0.42)',fontSize:14,lineHeight:1}}>×</button>}
                </div>
                <div style={{marginBottom:6}}>
                  <input value={s.name} onChange={e=>{const n=[...services];n[i].name=e.target.value;setServices(n);}} placeholder="e.g. Deep Tissue Massage" style={{width:'100%',padding:'6px 9px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:6,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif",outline:'none'}}/>
                </div>
                <div style={{display:'flex',gap:6,marginBottom:6}}>
                  <div style={{flex:2}}>
                    <input value={s.durNum||''} onChange={e=>{const n=[...services];n[i].durNum=e.target.value;setServices(n);}} placeholder="60" type="number" style={{width:'100%',padding:'6px 9px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:6,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif",outline:'none'}}/>
                  </div>
                  <div style={{flex:1}}>
                    <select value={s.durUnit||'min'} onChange={e=>{const n=[...services];n[i].durUnit=e.target.value;setServices(n);}} style={{width:'100%',padding:'6px 4px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:6,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif"}}>
                      <option value="min">min</option>
                      <option value="h">hours</option>
                      <option value="days">days</option>
                    </select>
                  </div>
                </div>
                <div style={{display:'flex',gap:6,marginBottom:6}}>
                  <div style={{flex:1}}>
                    <select value={s.currency||'£'} onChange={e=>{const n=[...services];n[i].currency=e.target.value;setServices(n);}} style={{width:'100%',padding:'6px 4px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:6,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif"}}>
                      <option value="£">£ GBP</option>
                      <option value="€">€ EUR</option>
                      <option value="$">$ USD</option>
                    </select>
                  </div>
                  <div style={{flex:2}}>
                    <input value={s.priceNum||''} onChange={e=>{const n=[...services];n[i].priceNum=e.target.value;setServices(n);}} placeholder="70" type="number" style={{width:'100%',padding:'6px 9px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:6,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif",outline:'none'}}/>
                  </div>
                </div>
                <textarea value={s.desc} onChange={e=>{const n=[...services];n[i].desc=e.target.value;setServices(n);}} placeholder="Description" rows={2} style={{width:'100%',padding:'6px 9px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:6,color:'#e8dcc8',fontSize:12,fontFamily:"'Cormorant Garamond',serif",lineHeight:1.5,outline:'none'}}/>
              </div>
            ))}
            <button onClick={()=>setServices([...services,{name:'',durNum:'',durUnit:'min',priceNum:'',currency:'£',desc:''}])} style={{width:'100%',padding:'8px 0',background:'none',border:'1px dashed rgba(201,169,110,0.18)',borderRadius:20,cursor:'pointer',color:'rgba(201,169,110,0.38)',fontSize:12,fontFamily:"'Jost',sans-serif"}}>+ Add service</button>
            <div style={{marginTop:8}}>
              <button onClick={()=>document.getElementById('ob-svc-screenshot').click()} disabled={svcImporting} style={{width:'100%',padding:'8px 0',background:'none',border:'1px dashed rgba(110,143,201,0.3)',borderRadius:20,cursor:'pointer',color:svcImporting?'rgba(110,143,201,0.35)':'rgba(110,143,201,0.65)',fontSize:12,fontFamily:"'Jost',sans-serif"}}>
                {svcImporting?'Analysing screenshot…':'📸 Import services from screenshot'}
              </button>
              <input id="ob-svc-screenshot" type="file" accept="image/*" style={{display:'none'}} onChange={e=>{if(e.target.files[0])obImportServices(e.target.files[0]);e.target.value='';}}/>
            </div>
            {svcPreview&&<div style={{marginTop:10,padding:'11px 13px',background:'rgba(110,143,201,0.06)',border:'1px solid rgba(110,143,201,0.18)',borderRadius:10}}>
              <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(110,143,201,0.7)',marginBottom:8}}>AI found {svcPreview.length} service{svcPreview.length!==1?'s':''} — confirm to add</div>
              {svcPreview.map((s,i)=>(
                <div key={i} style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.7)',marginBottom:4,paddingLeft:8,borderLeft:'2px solid rgba(110,143,201,0.3)'}}>
                  <b style={{color:'#e8dcc8'}}>{s.name}</b>{s.priceNum?` · ${s.currency}${s.priceNum}`:''}
                </div>
              ))}
              <div style={{display:'flex',gap:8,marginTop:10}}>
                <button onClick={()=>{setServices(sv=>[...sv.filter(s=>s.name.trim()),...svcPreview]);setSvcPreview(null);}} style={{flex:1,padding:'7px 0',background:'rgba(110,143,201,0.15)',border:'1px solid rgba(110,143,201,0.3)',borderRadius:20,cursor:'pointer',color:'rgba(110,143,201,0.9)',fontSize:12,fontFamily:"'Jost',sans-serif"}}>✓ Add all</button>
                <button onClick={()=>setSvcPreview(null)} style={{padding:'7px 14px',background:'none',border:'1px solid rgba(201,169,110,0.15)',borderRadius:20,cursor:'pointer',color:'rgba(201,169,110,0.4)',fontSize:12,fontFamily:"'Jost',sans-serif"}}>Dismiss</button>
              </div>
            </div>}
          </div>}

          {cur.type==='location'&&<div>
            <button onClick={useMyLocation} style={{width:'100%',padding:'11px 0',background:'rgba(201,169,110,0.1)',border:'1px solid rgba(201,169,110,0.25)',borderRadius:20,cursor:'pointer',color:'#c9a96e',fontSize:12,fontFamily:"'Jost',sans-serif",fontWeight:600,letterSpacing:'0.05em',textTransform:'uppercase',marginBottom:10}}>
              📍 {geoStatus==='locating'?(L==='it'?'Localizzazione...':'Locating...'):(L==='it'?'Usa la mia posizione attuale':'Use my current location')}
            </button>
            {geoStatus==='done'&&<div style={{fontSize:11,fontFamily:"'Jost',sans-serif",color:'#7aba6a',marginBottom:10}}>✓ {L==='it'?'Posizione rilevata':'Location detected'}: {data.loc}</div>}
            {geoStatus==='denied'&&<div style={{fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(200,100,100,0.7)',marginBottom:10}}>{L==='it'?'Permesso negato — inserisci manualmente':'Permission denied — enter manually below'}</div>}
            {geoStatus==='unsupported'&&<div style={{fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.4)',marginBottom:10}}>{L==='it'?'Non supportato — inserisci manualmente':'Not supported — enter manually below'}</div>}
            <div style={{marginBottom:6}}>
              <div style={{fontSize:8,fontFamily:"'Jost',sans-serif",letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.32)',marginBottom:4}}>{L==='it'?'Città / zona':'City / area'}</div>
              <input value={data.loc||''} onChange={e=>set('loc',e.target.value)} placeholder={L==='it'?'es. Londra':'e.g. London'} style={{width:'100%',padding:'9px 12px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.16)',borderRadius:8,color:'#e8dcc8',fontSize:13,fontFamily:"'Jost',sans-serif",outline:'none'}}/>
            </div>
            <div style={{marginBottom:6}}>
              <div style={{fontSize:8,fontFamily:"'Jost',sans-serif",letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.32)',marginBottom:4}}>{L==='it'?'Raggio di copertura (opzionale)':'Coverage range (optional)'}</div>
              <input value={data.coverageRange||''} onChange={e=>set('coverageRange',e.target.value)} placeholder={L==='it'?'es. entro 10km':'e.g. within 10km'} style={{width:'100%',padding:'9px 12px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.16)',borderRadius:8,color:'#e8dcc8',fontSize:13,fontFamily:"'Jost',sans-serif",outline:'none'}}/>
            </div>
            <div>
              <div style={{fontSize:8,fontFamily:"'Jost',sans-serif",letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.32)',marginBottom:4}}>{L==='it'?'Altre città disponibili (opzionale)':'Other cities you cover (optional)'}</div>
              <input value={data.otherCities||''} onChange={e=>set('otherCities',e.target.value)} placeholder={L==='it'?'es. Milano (estate)':'e.g. Milan (summer)'} style={{width:'100%',padding:'9px 12px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.16)',borderRadius:8,color:'#e8dcc8',fontSize:13,fontFamily:"'Jost',sans-serif",outline:'none'}}/>
            </div>
          </div>}

          {cur.type==='bb_details'&&<div>
            <div style={{marginBottom:11}}>
              <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.42)',marginBottom:6}}>Property type</div>
              <select value={data.propertyType||''} onChange={e=>set('propertyType',e.target.value)} style={{width:'100%',padding:'9px 12px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.16)',borderRadius:8,color:data.propertyType?'#e8dcc8':'rgba(232,220,200,0.3)',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:300}}>
                <option value="" style={{background:'#1a1410'}}>Select type...</option>
                {['B&B','Boutique Hotel','Short-let','Apartment','Villa'].map(t=><option key={t} value={t} style={{background:'#1a1410'}}>{t}</option>)}
              </select>
            </div>
            <div style={{display:'flex',gap:8,marginBottom:11}}>
              {[['checkIn','Check-in','e.g. 3:00 PM'],['checkOut','Check-out','e.g. 11:00 AM']].map(([k,label,ph])=>(
                <div key={k} style={{flex:1}}>
                  <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.42)',marginBottom:6}}>{label}</div>
                  <input value={data[k]||''} onChange={e=>set(k,e.target.value)} placeholder={ph} style={{width:'100%',padding:'9px 12px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.16)',borderRadius:8,color:'#e8dcc8',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:300,outline:'none'}}/>
                </div>
              ))}
            </div>
            <div style={{display:'flex',gap:8,marginBottom:11}}>
              <div style={{flex:1}}>
                <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.42)',marginBottom:6}}>Currency</div>
                <select value={data.bbCurrency||'£'} onChange={e=>set('bbCurrency',e.target.value)} style={{width:'100%',padding:'9px 8px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.16)',borderRadius:8,color:'#e8dcc8',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:300}}>
                  <option value="£" style={{background:'#1a1410'}}>£ GBP</option>
                  <option value="€" style={{background:'#1a1410'}}>€ EUR</option>
                  <option value="$" style={{background:'#1a1410'}}>$ USD</option>
                </select>
              </div>
              <div style={{flex:2}}>
                <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.42)',marginBottom:6}}>Nightly rate</div>
                <input value={data.nightlyRate||''} onChange={e=>set('nightlyRate',e.target.value)} placeholder="e.g. 120" type="number" style={{width:'100%',padding:'9px 12px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.16)',borderRadius:8,color:'#e8dcc8',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:300,outline:'none'}}/>
              </div>
            </div>
            <div style={{marginBottom:11}}>
              <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.42)',marginBottom:8}}>Amenities</div>
              <div style={{display:'flex',flexWrap:'wrap',gap:6}}>
                {['WiFi','Parking','Breakfast included','Kitchen access','Garden/Terrace','Pool','Hot tub','Air conditioning','Heating','TV/Netflix','Washing machine','Towels & linen','Late check-in','Luggage storage'].map(a=>{
                  const active=(data.amenities||[]).includes(a);
                  return <button key={a} onClick={()=>tog('amenities',a)} style={{padding:'6px 11px',borderRadius:20,cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",fontWeight:active?500:300,border:active?'1px solid #c9a96e':'1px solid rgba(201,169,110,0.16)',background:active?'rgba(201,169,110,0.15)':'rgba(201,169,110,0.04)',color:active?'#e8c878':'rgba(232,220,200,0.52)',transition:'all 0.15s',marginBottom:2}}>{a}</button>;
                })}
              </div>
            </div>
          </div>}

          {cur.type==='handle'&&<div>
            <div style={{display:'flex',alignItems:'center',background:'rgba(255,255,255,0.04)',border:`1px solid ${handleStatus==='available'?'rgba(100,180,80,0.4)':handleStatus==='taken'||handleStatus==='reserved'?'rgba(200,80,80,0.4)':'rgba(201,169,110,0.16)'}`,borderRadius:20,overflow:'hidden'}}>
              <span style={{padding:'9px 4px 9px 12px',fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.3)',whiteSpace:'nowrap'}}>bravebybruno.com/</span>
              <input value={handle} onChange={e=>onHandleChange(e.target.value)} placeholder="your-name" style={{flex:1,padding:'9px 12px 9px 0',background:'none',border:'none',color:'#e8dcc8',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:500,outline:'none'}}/>
            </div>
            <div style={{marginTop:8,fontSize:11,fontFamily:"'Jost',sans-serif",minHeight:16}}>
              {handleStatus==='checking'&&<span style={{color:'rgba(201,169,110,0.5)'}}>Checking...</span>}
              {handleStatus==='available'&&<span style={{color:'#7aba6a'}}>✓ Available!</span>}
              {handleStatus==='taken'&&<span style={{color:'rgba(200,100,100,0.8)'}}>✗ Taken — try adding your city or profession</span>}
              {handleStatus==='reserved'&&<span style={{color:'rgba(200,100,100,0.8)'}}>✗ Reserved — pick another</span>}
              {handleStatus==='short'&&<span style={{color:'rgba(201,169,110,0.4)'}}>At least 3 characters</span>}
            </div>
            {(handleStatus==='taken'||handleStatus==='reserved')&&<div style={{marginTop:10,display:'flex',flexWrap:'wrap',gap:6}}>
              {[slugify((data.name||'')+'-'+(data.loc||'')),slugify((data.biz||data.name||'')+'-'+(data.profession||'').split(' ')[0]),slugify((data.name||'me'))+'-'+Math.random().toString(36).slice(2,5)].filter(s=>s&&s.length>3).map(sug=>(
                <button key={sug} onClick={()=>{setHandleTouched(true);setHandle(sug);checkHandle(sug);}} style={{padding:'5px 10px',borderRadius:20,cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",background:'rgba(201,169,110,0.08)',border:'1px solid rgba(201,169,110,0.2)',color:'rgba(232,220,200,0.7)'}}>{sug}</button>
              ))}
            </div>}
          </div>}
        </div>

        <div style={{display:'flex',gap:8}}>
          {safeStep>0&&<button onClick={()=>setStep(Math.max(0,safeStep-1))} style={{padding:'10px 16px',background:'none',border:'1px solid rgba(201,169,110,0.12)',borderRadius:20,cursor:'pointer',color:'rgba(232,220,200,0.35)',fontSize:12,fontFamily:"'Jost',sans-serif"}}>{L==='it'?'Indietro':'Back'}</button>}
          <button onClick={next} disabled={!canNext()} style={{flex:1,padding:'11px 0',background:canNext()?'linear-gradient(135deg,#c9a96e,#8c5e14)':'rgba(201,169,110,0.09)',border:'none',borderRadius:20,cursor:canNext()?'pointer':'default',color:canNext()?'#0c0a08':'rgba(201,169,110,0.25)',fontSize:12,fontFamily:"'Jost',sans-serif",fontWeight:600,letterSpacing:'0.07em',textTransform:'uppercase',transition:'all 0.2s'}}>
            {safeStep===total-1?(L==='it'?'Lancia il tuo Concierge ✦':'Launch My Concierge ✦'):(L==='it'?'Continua →':'Continue →')}
          </button>
        </div>
        <div style={{textAlign:'center',marginTop:12}}>
          <button onClick={onBack} style={{background:'none',border:'none',cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.26)',textDecoration:'underline'}}>{L==='it'?'← Torna alle demo':'← Back to demos'}</button>
        </div>
      </div>
    </div>
  );
}
