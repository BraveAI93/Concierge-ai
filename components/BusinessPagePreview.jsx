'use client';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { generateBusinessPage } from '@/lib/generateBusinessPage';

export default function BusinessPagePreview({ data, lang, onStartChat, onBack }) {
  const [copied, setCopied] = useState(false);
  const [downloaded, setDownloaded] = useState(false);
  const [cardDownloaded, setCardDownloaded] = useState(false);
  const router = useRouter();
  const L = lang || 'en';

  const html = generateBusinessPage(data, lang);

  const publicURL = data._slug
    ? `${typeof window !== 'undefined' ? window.location.origin : 'https://bravebybruno.com'}/${data._slug}`
    : (typeof window !== 'undefined' ? window.location.href : '');

  const download = () => {
    const slug = (data.name||'my-concierge').toLowerCase().replace(/\s+/g,'-').replace(/[^a-z0-9-]/g,'');
    const blob = new Blob([html], {type:'text/html'});
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url; a.download = `${slug}.html`; a.click();
    setDownloaded(true);
  };

  const downloadBusinessCard = () => {
    const slug = data._slug || (data.name||'me').toLowerCase().replace(/\s+/g,'-').replace(/[^a-z0-9-]/g,'');
    const cardURL = `https://bravebybruno.com/${slug}`;
    const name = data.name || '';
    const tagline = (data.tag && (typeof data.tag === 'object' ? data.tag.en : data.tag)) || (data.profession || '');

    const W = 600, H = 400;
    const canvas = document.createElement('canvas');
    canvas.width = W; canvas.height = H;
    const ctx = canvas.getContext('2d');

    ctx.fillStyle = '#0a0a0a';
    ctx.fillRect(0, 0, W, H);
    ctx.strokeStyle = 'rgba(201,169,110,0.35)';
    ctx.lineWidth = 1;
    ctx.strokeRect(14, 14, W - 28, H - 28);
    ctx.strokeStyle = 'rgba(201,169,110,0.6)';
    ctx.lineWidth = 1;
    ctx.beginPath(); ctx.moveTo(40, 28); ctx.lineTo(W - 40, 28); ctx.stroke();
    ctx.beginPath(); ctx.moveTo(40, H - 28); ctx.lineTo(W - 40, H - 28); ctx.stroke();
    ctx.fillStyle = '#c9a96e';
    ctx.font = '22px "Cormorant Garamond", Georgia, serif';
    ctx.textAlign = 'center';
    ctx.fillText('✦', W / 2, 60);
    ctx.fillStyle = '#ffffff';
    ctx.font = 'bold 30px "Cormorant Garamond", Georgia, serif';
    ctx.textAlign = 'center';
    ctx.fillText(name, W / 2, 98);
    const tagTrimmed = tagline.length > 52 ? tagline.slice(0, 52) + '…' : tagline;
    ctx.fillStyle = '#c9a96e';
    ctx.font = '400 14px "Jost", sans-serif';
    ctx.textAlign = 'center';
    ctx.fillText(tagTrimmed, W / 2, 120);

    const qrSize = 160;
    const qrX = (W - qrSize) / 2;
    const qrY = 140;

    const drawQR = (qrCanvas) => {
      const pad = 8;
      ctx.fillStyle = '#e8dcc8';
      ctx.beginPath();
      ctx.roundRect(qrX - pad, qrY - pad, qrSize + pad * 2, qrSize + pad * 2, 8);
      ctx.fill();
      ctx.drawImage(qrCanvas, qrX, qrY, qrSize, qrSize);
      ctx.fillStyle = 'rgba(201,169,110,0.75)';
      ctx.font = '400 11px "Jost", sans-serif';
      ctx.textAlign = 'center';
      ctx.fillText(cardURL, W / 2, qrY + qrSize + pad * 2 + 16);
      ctx.fillStyle = 'rgba(201,169,110,0.35)';
      ctx.font = '300 10px "Jost", sans-serif';
      ctx.textAlign = 'center';
      ctx.fillText('Powered by Brave by Bruno', W / 2, H - 34);
      canvas.toBlob(blob => {
        const a = document.createElement('a');
        a.href = URL.createObjectURL(blob);
        a.download = `${slug}-business-card.jpg`;
        a.click();
      }, 'image/jpeg', 0.95);
      setCardDownloaded(true);
    };

    const tmp = document.createElement('div');
    tmp.style.cssText = 'position:absolute;left:-9999px;top:-9999px;';
    document.body.appendChild(tmp);
    try {
      if (typeof QRCode !== 'undefined') {
        new QRCode(tmp, { text: cardURL, width: qrSize, height: qrSize, colorDark: '#0a0a0a', colorLight: '#e8dcc8' });
        const qrCanvas = tmp.querySelector('canvas');
        if (qrCanvas) drawQR(qrCanvas);
      }
    } finally {
      document.body.removeChild(tmp);
    }
  };

  const bg = '#0c0a08';
  const text = '#e8dcc8';
  const cardBg = 'rgba(255,255,255,0.025)';

  return (
    <div style={{minHeight:'100vh',background:bg,display:'flex',flexDirection:'column',alignItems:'center',justifyContent:'center',padding:'24px 16px',fontFamily:"'Cormorant Garamond',serif",color:text}}>
      <div style={{width:'100%',maxWidth:440}}>
        {/* Header */}
        <div style={{textAlign:'center',marginBottom:24}}>
          <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",letterSpacing:'0.18em',textTransform:'uppercase',color:'rgba(201,169,110,0.45)',marginBottom:8}}>✦ Your Concierge is ready</div>
          <div style={{fontSize:26,fontWeight:400,color:text,marginBottom:6}}>{data.name}</div>
          <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.55)',letterSpacing:'0.1em',textTransform:'uppercase'}}>{data.profession} · {data.loc||'London'}</div>
        </div>

        {/* Live link */}
        <div style={{background:'rgba(201,169,110,0.08)',border:'1px solid rgba(201,169,110,0.25)',borderRadius:20,padding:'12px 14px',marginBottom:14,textAlign:'center'}}>
          <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",letterSpacing:'0.12em',textTransform:'uppercase',color:'rgba(201,169,110,0.45)',marginBottom:5}}>Your live concierge link</div>
          <div style={{fontSize:13,fontFamily:"'Jost',sans-serif",color:'#c9a96e',wordBreak:'break-all'}}>{publicURL}</div>
        </div>

        {/* Preview card */}
        <div style={{background:cardBg,border:'1px solid rgba(201,169,110,0.15)',borderRadius:20,padding:'20px',marginBottom:14}}>
          <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.4)',marginBottom:12}}>Your page includes</div>
          {[
            ['✦','AI Concierge bubble','Opens your personal chat assistant'],
            ['📅','Booking button', data.booking ? 'Connected to your booking link' : 'Ready to add your booking link'],
            ['💆','Services listed', `${(data.services||[]).filter(s=>s.name).length} service${(data.services||[]).filter(s=>s.name).length!==1?'s':''} added`],
            ['🌍','Multilingual ready','Supports EN, IT, ES, PT, FR'],
            ['🔒','GDPR compliant','Privacy consent built in'],
          ].map(([icon,title,sub])=>(
            <div key={title} style={{display:'flex',gap:10,marginBottom:10,alignItems:'flex-start'}}>
              <div style={{fontSize:16,flexShrink:0,marginTop:1}}>{icon}</div>
              <div>
                <div style={{fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:500,color:text}}>{title}</div>
                <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",fontWeight:300,color:'rgba(201,169,110,0.45)',marginTop:1}}>{sub}</div>
              </div>
            </div>
          ))}
        </div>

        {/* Actions */}
        <button onClick={()=>{
          navigator.clipboard?.writeText(publicURL).catch(()=>{});
          setCopied(true); setTimeout(()=>setCopied(false),2500);
        }} style={{width:'100%',padding:'13px 0',background:'linear-gradient(135deg,#c9a96e,#7a4f0e)',border:'none',borderRadius:20,cursor:'pointer',color:'#0c0a08',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:700,letterSpacing:'0.08em',textTransform:'uppercase',marginBottom:9}}>
          {copied?'✓ Copied! Paste in your Instagram bio':'🔗 Copy my link for bio'}
        </button>

        <button onClick={async()=>{
          const msg = L==='it'
            ?`Ciao! Ho appena lanciato il mio assistente AI 24/7 — chiedimi qualsiasi cosa su prenotazioni, prezzi e disponibilità: ${publicURL}`
            :`Hi! I just launched my 24/7 AI concierge — ask me anything about bookings, prices and availability: ${publicURL}`;
          const waUrl=`https://wa.me/?text=${encodeURIComponent(msg)}`;
          let shared=false;
          if(navigator.share){
            try{ await navigator.share({title:data.name,text:msg,url:publicURL}); shared=true; }catch(e){ shared=false; }
          }
          if(!shared){ window.open(waUrl,'_blank'); }
        }} style={{width:'100%',padding:'12px 0',background:'rgba(201,169,110,0.06)',border:'1px solid rgba(201,169,110,0.2)',borderRadius:20,cursor:'pointer',color:'#c9a96e',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.06em',textTransform:'uppercase',marginBottom:9}}>
          📲 Share on WhatsApp
        </button>

        <button onClick={()=>{
          navigator.clipboard?.writeText(publicURL).catch(()=>{});
        }} style={{width:'100%',padding:'12px 0',background:'rgba(201,169,110,0.06)',border:'1px solid rgba(201,169,110,0.2)',borderRadius:20,cursor:'pointer',color:'#c9a96e',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.06em',textTransform:'uppercase',marginBottom:9}}>
          📸 {L==='it'?'Condividi su Instagram':'Share on Instagram'}
        </button>

        <button onClick={downloadBusinessCard} style={{width:'100%',padding:'12px 0',background:'rgba(201,169,110,0.06)',border:'1px solid rgba(201,169,110,0.2)',borderRadius:20,cursor:'pointer',color:'#c9a96e',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.06em',textTransform:'uppercase',marginBottom:9}}>
          {cardDownloaded?'✓ Business Card Downloaded':'🪪 Download Business Card'}
        </button>

        <button onClick={()=>{ window.open(publicURL,'_blank'); }} style={{width:'100%',padding:'12px 0',background:'rgba(201,169,110,0.06)',border:'1px solid rgba(201,169,110,0.2)',borderRadius:20,cursor:'pointer',color:'#c9a96e',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.06em',textTransform:'uppercase',marginBottom:9}}>
          👁 Preview as your clients see it
        </button>

        <button onClick={()=>router.push('/theconcierge/dashboard')} style={{width:'100%',padding:'12px 0',background:'rgba(201,169,110,0.08)',border:'1px solid rgba(201,169,110,0.25)',borderRadius:20,cursor:'pointer',color:'#c9a96e',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:600,letterSpacing:'0.06em',textTransform:'uppercase',marginBottom:9}}>
          🔐 My Dashboard &amp; Login
        </button>

        <button onClick={download} style={{width:'100%',padding:'12px 0',background:'none',border:'1px solid rgba(201,169,110,0.12)',borderRadius:20,cursor:'pointer',color:'rgba(201,169,110,0.5)',fontSize:12,fontFamily:"'Jost',sans-serif",fontWeight:400,letterSpacing:'0.05em',textTransform:'uppercase',marginBottom:9}}>
          {downloaded?'✓ Downloaded':'⬇ Download HTML file (optional)'}
        </button>

        <button onClick={onStartChat} style={{width:'100%',padding:'12px 0',background:'rgba(201,169,110,0.08)',border:'1px solid rgba(201,169,110,0.2)',borderRadius:20,cursor:'pointer',color:'#c9a96e',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.06em',textTransform:'uppercase',marginBottom:9}}>
          💬 Test Your Concierge
        </button>

        <div style={{textAlign:'center'}}>
          <button onClick={onBack} style={{background:'none',border:'none',cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.3)',textDecoration:'underline'}}>← Back to demos</button>
        </div>

        <div style={{marginTop:16,padding:'12px 14px',background:'rgba(201,169,110,0.04)',border:'1px dashed rgba(201,169,110,0.15)',borderRadius:10,fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.4)',lineHeight:1.65,textAlign:'center'}}>
          Upload the downloaded HTML file to your website or share directly. Your concierge is always live at the link you share.
        </div>
      </div>
    </div>
  );
}
