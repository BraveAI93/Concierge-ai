'use client';
import { useState, useEffect } from 'react';
import { BACKEND_URL } from '@/lib/constants';
import { buildPrompt } from '@/lib/buildPrompt';

export default function OwnerEditProfile({ token, slug, lang, onBack, onSaved }) {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [profile, setProfile] = useState(null);
  const [pdata, setPdata] = useState({});
  const [saved, setSaved] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [uploadPct, setUploadPct] = useState(0);
  const [editServices, setEditServices] = useState([]);
  const [svcImporting, setSvcImporting] = useState(false);
  const [svcPreview, setSvcPreview] = useState(null);
  const L = lang || 'en';

  useEffect(() => {
    fetch(`${BACKEND_URL}/owner/profile`, { headers: { Authorization: `Bearer ${token}` } })
      .then(r => r.ok ? r.json() : null)
      .then(p => {
        if (p) {
          setProfile(p);
          let parsed = {}; try { parsed = JSON.parse(p.profile_data || '{}'); } catch(e) {}
          setPdata(parsed);
          setEditServices(parsed.services || []);
        }
        setLoading(false);
      }).catch(() => setLoading(false));
  }, [token]);

  const set = (k, v) => setPdata(d => ({...d, [k]: v}));

  const save = async () => {
    setSaving(true);
    const merged = {...pdata, services: editServices};
    const sp = buildPrompt(merged, lang);
    try {
      await fetch(`${BACKEND_URL}/owner/profile`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify({ name: merged.name, business: merged.biz || merged.name, profession: merged.profession || profile.profession, location: merged.loc || '', system_prompt: sp, profile_data: JSON.stringify(merged), active: true })
      });
      setSaved(true);
      setTimeout(() => { onSaved && onSaved(); }, 1200);
    } catch(e) { alert('Could not save — please try again.'); }
    setSaving(false);
  };

  const uploadFile = async (file) => {
    setUploading(true); setUploadPct(0);
    const fd = new FormData(); fd.append('file', file);
    try {
      await new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.upload.onprogress = e => { if (e.lengthComputable) setUploadPct(Math.round(e.loaded / e.total * 100)); };
        xhr.onload = () => {
          if (xhr.status === 200) {
            const d = JSON.parse(xhr.responseText);
            const cur = (pdata.media_gallery || []);
            setPdata(p => ({...p, media_gallery: [...cur, d.url]}));
            resolve();
          } else reject(new Error('Upload failed'));
        };
        xhr.onerror = () => reject(new Error('Upload failed'));
        xhr.open('POST', `${BACKEND_URL}/media/upload`);
        xhr.setRequestHeader('Authorization', 'Bearer ' + token);
        xhr.send(fd);
      });
    } catch(e) { alert('Upload failed: ' + e.message); }
    setUploading(false); setUploadPct(0);
  };

  const importServices = async (file) => {
    setSvcImporting(true); setSvcPreview(null);
    try {
      const b64 = await new Promise((res, rej) => { const r = new FileReader(); r.onload = e => res(e.target.result.split(',')[1]); r.onerror = rej; r.readAsDataURL(file); });
      const mimeType = file.type || 'image/jpeg';
      const resp = await fetch(`${BACKEND_URL}/ai/import-services`, { method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: 'Bearer ' + token }, body: JSON.stringify({ image_base64: b64, mime_type: mimeType }) });
      const d = await resp.json();
      if (resp.ok && d.services) {
        const mapped = d.services.map(s => ({ name: s.name || '', durNum: s.duration_minutes ? String(s.duration_minutes) : '', durUnit: 'min', priceNum: s.price ? String(s.price) : '', currency: s.currency || '£', desc: s.description || '' }));
        setSvcPreview(mapped);
      } else { alert(d.error || 'Could not extract services'); }
    } catch(e) { alert('Import failed: ' + e.message); }
    setSvcImporting(false);
  };

  if (loading) return (
    <div style={{minHeight:'100vh',background:'#0c0a08',display:'flex',alignItems:'center',justifyContent:'center'}}>
      <div style={{fontSize:28,color:'#c9a96e',animation:'bravePulse 1.8s ease-in-out infinite',textShadow:'0 0 20px rgba(201,169,110,0.6)'}}>✦</div>
    </div>
  );

  const fieldStyle = {width:'100%',padding:'10px 13px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.18)',borderRadius:20,color:'#e8dcc8',fontSize:13,fontFamily:"'Jost',sans-serif",marginBottom:14,outline:'none'};
  const labelStyle = {fontSize:10,fontFamily:"'Jost',sans-serif",letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.45)',marginBottom:6};

  return (
    <div style={{minHeight:'100vh',background:'#0c0a08',padding:'24px 16px',fontFamily:"'Cormorant Garamond',Georgia,serif",color:'#e8dcc8'}}>
      <div style={{maxWidth:480,margin:'0 auto'}}>
        <div style={{display:'flex',alignItems:'center',justifyContent:'space-between',marginBottom:24}}>
          <div style={{fontFamily:"'Jost',sans-serif",fontWeight:500,fontSize:14,letterSpacing:'0.08em',textTransform:'uppercase',color:'#e8dcc8'}}>Edit My Profile</div>
          <button onClick={onBack} style={{padding:'7px 13px',background:'rgba(255,255,255,0.03)',border:'1px solid rgba(255,255,255,0.09)',borderRadius:20,cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.38)'}}>← Cancel</button>
        </div>

        {[
          ['name','Name','Your name'],
          ['biz','Business name','Business name'],
          ['tag','Tagline','A one-line tagline'],
          ['loc','Location','City, area'],
          ['booking','Booking link','fresha.com/...'],
          ['calendar','Calendar link','calendly.com/...'],
          ['wa','WhatsApp','+44 7700 000000'],
          ['phone','Phone','+44 7700 000000'],
          ['email','Email (public)','you@email.com'],
          ['ig','Instagram','@your_handle'],
          ['tg','Telegram','@your_handle'],
          ['video','Video / showreel link','youtube.com/...'],
          ['gallery','Gallery link','instagram.com/...'],
        ].map(([key, label, ph]) => (
          <div key={key}>
            <div style={labelStyle}>{label}</div>
            <input value={pdata[key]||''} onChange={e=>set(key,e.target.value)} placeholder={ph} style={fieldStyle}/>
          </div>
        ))}

        <div style={{marginBottom:14}}>
          <div style={labelStyle}>Upload media</div>
          <button onClick={()=>document.getElementById('ep-media-upload').click()} disabled={uploading} style={{width:'100%',padding:'10px 0',background:'rgba(201,169,110,0.07)',border:'1px dashed rgba(201,169,110,0.28)',borderRadius:20,cursor:uploading?'default':'pointer',color:'rgba(201,169,110,0.72)',fontSize:12,fontFamily:"'Jost',sans-serif",marginBottom:6}}>
            {uploading?`Uploading… ${uploadPct}%`:'⬆ Upload from device (images & videos)'}
          </button>
          <input id="ep-media-upload" type="file" accept="image/*,video/*" style={{display:'none'}} onChange={e=>{if(e.target.files[0])uploadFile(e.target.files[0]);e.target.value='';}}/>
          {uploading && <div style={{height:3,background:'rgba(201,169,110,0.1)',borderRadius:2,marginBottom:8}}><div style={{height:'100%',width:uploadPct+'%',background:'linear-gradient(90deg,#c9a96e,#e8c878)',borderRadius:2,transition:'width 0.3s'}}/></div>}
          {(pdata.media_gallery||[]).length>0 && (
            <div style={{display:'grid',gridTemplateColumns:'repeat(3,1fr)',gap:6}}>
              {(pdata.media_gallery||[]).map((url,i)=>(
                <div key={i} style={{position:'relative',paddingBottom:'100%',background:'rgba(201,169,110,0.06)',borderRadius:8,overflow:'hidden'}}>
                  <img src={url} style={{position:'absolute',top:0,left:0,width:'100%',height:'100%',objectFit:'cover',borderRadius:8}} onError={e=>e.target.style.display='none'}/>
                  <button onClick={()=>setPdata(p=>({...p,media_gallery:(p.media_gallery||[]).filter((_,idx)=>idx!==i)}))} style={{position:'absolute',top:3,right:3,background:'rgba(0,0,0,0.65)',border:'none',borderRadius:'50%',width:18,height:18,cursor:'pointer',color:'#e8dcc8',fontSize:10,lineHeight:'18px',textAlign:'center'}}>×</button>
                </div>
              ))}
            </div>
          )}
        </div>

        <div style={{marginBottom:14}}>
          <div style={labelStyle}>Services</div>
          {editServices.map((s,i)=>(
            <div key={i} style={{marginBottom:8,padding:'9px 11px',background:'rgba(201,169,110,0.04)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:10}}>
              <div style={{display:'flex',justifyContent:'space-between',marginBottom:6}}>
                <span style={{fontSize:9,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.38)',letterSpacing:'0.08em',textTransform:'uppercase'}}>Service {i+1}</span>
                <button onClick={()=>setEditServices(editServices.filter((_,idx)=>idx!==i))} style={{background:'none',border:'none',cursor:'pointer',color:'rgba(200,80,80,0.4)',fontSize:14,lineHeight:1}}>×</button>
              </div>
              <input value={s.name||''} onChange={e=>{const n=[...editServices];n[i]={...n[i],name:e.target.value};setEditServices(n);}} placeholder="Service name" style={{width:'100%',padding:'6px 9px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:6,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif",marginBottom:4,outline:'none'}}/>
              <div style={{display:'flex',gap:5,marginBottom:4}}>
                <input value={s.currency||'£'} onChange={e=>{const n=[...editServices];n[i]={...n[i],currency:e.target.value};setEditServices(n);}} placeholder="£" style={{width:40,padding:'5px 6px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:6,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif",outline:'none'}}/>
                <input value={s.priceNum||''} onChange={e=>{const n=[...editServices];n[i]={...n[i],priceNum:e.target.value};setEditServices(n);}} placeholder="Price" type="number" style={{flex:1,padding:'5px 8px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:6,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif",outline:'none'}}/>
                <input value={s.durNum||''} onChange={e=>{const n=[...editServices];n[i]={...n[i],durNum:e.target.value};setEditServices(n);}} placeholder="Mins" type="number" style={{width:56,padding:'5px 6px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:6,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif",outline:'none'}}/>
              </div>
              <input value={s.desc||''} onChange={e=>{const n=[...editServices];n[i]={...n[i],desc:e.target.value};setEditServices(n);}} placeholder="Description" style={{width:'100%',padding:'5px 9px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:6,color:'#e8dcc8',fontSize:12,fontFamily:"'Cormorant Garamond',serif",outline:'none'}}/>
            </div>
          ))}
          <button onClick={()=>setEditServices([...editServices,{name:'',durNum:'',durUnit:'min',priceNum:'',currency:'£',desc:''}])} style={{width:'100%',padding:'7px 0',background:'none',border:'1px dashed rgba(201,169,110,0.18)',borderRadius:20,cursor:'pointer',color:'rgba(201,169,110,0.38)',fontSize:12,fontFamily:"'Jost',sans-serif",marginBottom:6}}>+ Add service</button>
          <button onClick={()=>document.getElementById('ep-svc-screenshot').click()} disabled={svcImporting} style={{width:'100%',padding:'7px 0',background:'none',border:'1px dashed rgba(110,143,201,0.3)',borderRadius:20,cursor:'pointer',color:svcImporting?'rgba(110,143,201,0.35)':'rgba(110,143,201,0.65)',fontSize:12,fontFamily:"'Jost',sans-serif"}}>
            {svcImporting?'Analysing screenshot…':'📸 Import services from screenshot'}
          </button>
          <input id="ep-svc-screenshot" type="file" accept="image/*" style={{display:'none'}} onChange={e=>{if(e.target.files[0])importServices(e.target.files[0]);e.target.value='';}}/>
          {svcPreview && (
            <div style={{marginTop:8,padding:'10px 13px',background:'rgba(110,143,201,0.06)',border:'1px solid rgba(110,143,201,0.2)',borderRadius:10}}>
              <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(110,143,201,0.7)',marginBottom:8}}>AI found {svcPreview.length} service{svcPreview.length!==1?'s':''}</div>
              {svcPreview.map((s,i)=>(
                <div key={i} style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.7)',marginBottom:3,paddingLeft:8,borderLeft:'2px solid rgba(110,143,201,0.3)'}}>
                  <b style={{color:'#e8dcc8'}}>{s.name}</b>{s.priceNum?` · ${s.currency}${s.priceNum}`:''}
                </div>
              ))}
              <div style={{display:'flex',gap:8,marginTop:8}}>
                <button onClick={()=>{setEditServices(sv=>[...sv.filter(s=>s.name.trim()),...svcPreview]);setSvcPreview(null);}} style={{flex:1,padding:'6px 0',background:'rgba(110,143,201,0.15)',border:'1px solid rgba(110,143,201,0.3)',borderRadius:20,cursor:'pointer',color:'rgba(110,143,201,0.9)',fontSize:12,fontFamily:"'Jost',sans-serif"}}>✓ Add all</button>
                <button onClick={()=>setSvcPreview(null)} style={{padding:'6px 12px',background:'none',border:'1px solid rgba(201,169,110,0.15)',borderRadius:20,cursor:'pointer',color:'rgba(201,169,110,0.4)',fontSize:12,fontFamily:"'Jost',sans-serif"}}>Dismiss</button>
              </div>
            </div>
          )}
        </div>

        <div style={labelStyle}>Additional info for your concierge</div>
        <textarea value={pdata.extra||''} onChange={e=>set('extra',e.target.value)} rows={3} style={{...fieldStyle,fontFamily:"'Cormorant Garamond',serif",lineHeight:1.6,borderRadius:12}}/>

        <button onClick={save} disabled={saving} style={{width:'100%',padding:'13px 0',background:saved?'rgba(100,180,80,0.8)':'linear-gradient(135deg,#c9a96e,#7a4f0e)',border:'none',borderRadius:20,cursor:'pointer',color:'#0c0a08',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:700,letterSpacing:'0.08em',textTransform:'uppercase',marginTop:8}}>
          {saved?'✓ Saved!':saving?'Saving...':'Save Changes'}
        </button>
      </div>
    </div>
  );
}
