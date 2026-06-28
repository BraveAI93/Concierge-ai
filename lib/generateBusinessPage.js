export function generateBusinessPage(data, lang) {
  const name = data.name || 'Your Name';
  const biz = data.biz || name;
  const profession = data.profession || 'Professional';
  const loc = data.loc || 'London';
  const tag = data.tag || '';
  const booking = data.booking || '';
  const calendar = data.calendar || '';
  const instagram = data.ig || data.instagram || '';
  const telegram = data.tg || '';
  const whatsapp = data.wa || '';
  const phone = data.phone || '';
  const email = data.email || '';
  const video = data.video || '';
  const gallery = data.gallery || '';
  const lat = data.lat || '';
  const lng = data.lng || '';
  const accent = '#c9a96e';
  const slug = data._slug || name.toLowerCase().replace(/\s+/g,'-').replace(/[^a-z0-9-]/g,'');
  const conciergeURL = `https://concierge-ai-gamma.vercel.app/${slug}`;
  const publicPageURL = `https://concierge-ai-gamma.vercel.app/${slug}`;

  function getVideoEmbed(url){
    if(!url) return '';
    const ytShort = url.match(/youtube\.com\/shorts\/([a-zA-Z0-9_-]+)/);
    if(ytShort) return `<iframe src="https://www.youtube.com/embed/${ytShort[1]}?autoplay=1&mute=1&loop=1&playlist=${ytShort[1]}&playsinline=1" frameborder="0" allow="autoplay; fullscreen; picture-in-picture" allowfullscreen style="width:100%;max-width:280px;aspect-ratio:9/16;border-radius:12px;display:block;margin:0 auto"></iframe>`;
    const yt = url.match(/(?:youtube\.com\/(?:watch\?v=|embed\/)|youtu\.be\/)([a-zA-Z0-9_-]+)/);
    if(yt) return `<iframe src="https://www.youtube.com/embed/${yt[1]}?autoplay=1&mute=1&loop=1&playlist=${yt[1]}&playsinline=1" frameborder="0" allow="autoplay; fullscreen; picture-in-picture" allowfullscreen style="width:100%;aspect-ratio:16/9;border-radius:12px"></iframe>`;
    const vm = url.match(/vimeo\.com\/(?:video\/)?(\d+)/);
    if(vm) return `<iframe src="https://player.vimeo.com/video/${vm[1]}?autoplay=1&muted=1&loop=1&background=1" frameborder="0" allow="autoplay; fullscreen; picture-in-picture" allowfullscreen style="width:100%;aspect-ratio:16/9;border-radius:12px"></iframe>`;
    return '';
  }
  const videoEmbed = getVideoEmbed(video);

  const services = (data.services||[]).filter(s=>s.name).map(s=>{
    const durTxt=s.durNum?`${s.durNum} ${s.durUnit==='h'?'hour'+(s.durNum==1?'':'s'):s.durUnit==='days'?'day'+(s.durNum==1?'':'s'):'min'}`:'';
    const priceTxt=s.priceNum?`${s.currency||'£'}${s.priceNum}`:'';
    return `
    <div class="service-card">
      <div class="service-name">${s.name}</div>
      ${durTxt?`<div class="service-detail">⏱ ${durTxt}</div>`:''}
      ${priceTxt?`<div class="service-price">${priceTxt}</div>`:''}
      ${s.desc?`<div class="service-desc">${s.desc}</div>`:''}
    </div>`;
  }).join('');

  return `<!DOCTYPE html>
<html lang="${lang}">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport" content="width=device-width,initial-scale=1.0"/>
  <title>${name} — ${profession}</title>
  <link rel="preconnect" href="https://fonts.googleapis.com"/>
  <link href="https://fonts.googleapis.com/css2?family=Cormorant+Garamond:ital,wght@0,300;0,400;0,500;1,400&family=Jost:wght@300;400;500;600&display=swap" rel="stylesheet"/>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/qrcodejs/1.0.0/qrcode.min.js"><\/script>
  <style>
    *{box-sizing:border-box;margin:0;padding:0}
    body{background:#0c0a08;color:#e8dcc8;font-family:'Cormorant Garamond',serif;min-height:100vh}
    ::-webkit-scrollbar{width:3px}::-webkit-scrollbar-thumb{background:rgba(201,169,110,0.2)}
    @keyframes fadeUp{from{opacity:0;transform:translateY(16px)}to{opacity:1;transform:translateY(0)}}
    @keyframes bubblePop{0%{opacity:0;transform:scale(0.6) translateY(20px)}70%{transform:scale(1.06) translateY(-2px)}100%{opacity:1;transform:scale(1) translateY(0)}}
    @keyframes pulse{0%,100%{opacity:1}50%{opacity:0.4}}
    @keyframes typeDot{0%,80%,100%{opacity:.2;transform:translateY(0)}40%{opacity:1;transform:translateY(-4px)}}
    .page{max-width:480px;margin:0 auto;padding:32px 18px 100px;animation:fadeUp 0.5s ease forwards}
    .avatar-ring{width:88px;height:88px;border-radius:50%;background:linear-gradient(135deg,${accent},#5e3a0e);display:flex;align-items:center;justify-content:center;font-size:32px;margin:0 auto 14px;border:2px solid rgba(201,169,110,0.3)}
    .profile-name{font-size:26px;font-weight:400;text-align:center;color:#e8dcc8;margin-bottom:4px}
    .profile-prof{font-size:11px;font-family:'Jost',sans-serif;font-weight:300;letter-spacing:0.18em;text-transform:uppercase;color:rgba(201,169,110,0.55);text-align:center;margin-bottom:4px}
    .profile-loc{font-size:11px;font-family:'Jost',sans-serif;color:rgba(232,220,200,0.35);text-align:center;margin-bottom:10px}
    .profile-tag{font-size:14px;font-style:italic;color:rgba(232,220,200,0.45);text-align:center;margin-bottom:22px;line-height:1.6}
    .section-label{font-size:9px;font-family:'Jost',sans-serif;font-weight:500;letter-spacing:0.14em;text-transform:uppercase;color:rgba(201,169,110,0.38);margin-bottom:10px;margin-top:24px}
    .service-card{background:rgba(255,255,255,0.025);border:1px solid rgba(201,169,110,0.1);border-radius:20px;padding:14px 16px;margin-bottom:10px}
    .service-name{font-size:16px;font-weight:500;color:#e8dcc8;margin-bottom:4px}
    .service-detail{font-size:11px;font-family:'Jost',sans-serif;color:rgba(232,220,200,0.4);margin-bottom:2px}
    .service-price{font-size:14px;font-family:'Jost',sans-serif;font-weight:500;color:${accent};margin-bottom:4px}
    .service-desc{font-size:13px;color:rgba(232,220,200,0.5);line-height:1.6;margin-top:4px}
    .link-btn{display:flex;align-items:center;gap:12px;padding:13px 16px;background:rgba(255,255,255,0.025);border:1px solid rgba(201,169,110,0.12);border-radius:20px;text-decoration:none;color:#e8dcc8;margin-bottom:9px;transition:all 0.2s;font-family:'Jost',sans-serif}
    .link-btn:hover{background:rgba(201,169,110,0.08);border-color:rgba(201,169,110,0.25);transform:translateY(-1px)}
    .link-btn.primary{background:linear-gradient(135deg,${accent},#7a4f0e);border:none;color:#0c0a08;font-weight:600}
    .link-btn.primary:hover{opacity:0.88}
    .link-icon{width:36px;height:36px;border-radius:50%;background:rgba(201,169,110,0.1);display:flex;align-items:center;justify-content:center;font-size:16px;flex-shrink:0}
    .link-btn.primary .link-icon{background:rgba(0,0,0,0.15)}
    .link-title{font-size:13px;font-weight:500;letter-spacing:0.02em}
    .link-sub{font-size:10px;font-weight:300;opacity:0.65;margin-top:1px}
    #bubble-btn{position:fixed;bottom:22px;right:18px;z-index:9999;width:56px;height:56px;border-radius:50%;border:none;cursor:pointer;background:linear-gradient(135deg,${accent},#7a4f0e);box-shadow:0 4px 20px rgba(201,169,110,0.4);display:flex;align-items:center;justify-content:center;font-size:22px;animation:bubblePop 0.6s cubic-bezier(0.34,1.56,0.64,1) 2s both}
    #bubble-hint{position:fixed;bottom:90px;right:18px;z-index:9998;max-width:190px;background:rgba(14,11,7,0.97);border:1px solid rgba(201,169,110,0.28);border-radius:14px 14px 4px 14px;padding:10px 13px;font-size:12px;font-family:'Jost',sans-serif;color:#e8dcc8;line-height:1.55;box-shadow:0 8px 32px rgba(0,0,0,0.5);opacity:0;transition:opacity 0.3s,transform 0.3s;transform:translateY(6px);pointer-events:none}
    #bubble-hint.visible{opacity:1;pointer-events:auto;transform:translateY(0)}
    #bubble-overlay{position:fixed;inset:0;z-index:9996;background:rgba(0,0,0,0);pointer-events:none;transition:background 0.35s}
    #bubble-overlay.open{background:rgba(0,0,0,0.55);pointer-events:auto}
    #bubble-panel{position:fixed;bottom:0;left:0;right:0;z-index:9997;padding:0 12px 88px;opacity:0;pointer-events:none;transition:opacity 0.35s;display:flex;flex-direction:column;justify-content:flex-end}
    #bubble-panel.open{opacity:1;pointer-events:auto}
    .bp-card{width:100%;max-width:400px;margin:0 auto;background:rgba(14,11,7,0.98);border:1px solid rgba(201,169,110,0.22);border-radius:20px;padding:18px;box-shadow:0 -8px 40px rgba(0,0,0,0.6)}
    .bp-header{display:flex;align-items:center;gap:10px;margin-bottom:14px}
    .bp-avatar{width:38px;height:38px;border-radius:50%;background:linear-gradient(135deg,${accent},#7a4f0e);display:flex;align-items:center;justify-content:center;font-size:16px;flex-shrink:0}
    .bp-name{font-size:13px;font-family:'Jost',sans-serif;font-weight:500;color:#e8dcc8}
    .bp-status{font-size:9px;font-family:'Jost',sans-serif;color:rgba(120,180,100,0.8);display:flex;align-items:center;gap:4px;margin-top:1px}
    .bp-close{margin-left:auto;background:rgba(255,255,255,0.05);border:1px solid rgba(255,255,255,0.08);border-radius:50%;width:26px;height:26px;display:flex;align-items:center;justify-content:center;cursor:pointer;font-size:13px;color:rgba(232,220,200,0.4)}
    .bp-msg{background:rgba(255,255,255,0.04);border:1px solid rgba(201,169,110,0.1);border-radius:12px 12px 12px 4px;padding:11px 14px;margin-bottom:12px;font-size:14px;font-family:'Cormorant Garamond',serif;line-height:1.65;color:rgba(232,220,200,0.9)}
    .bp-typing{display:flex;gap:5px;align-items:center;padding:4px 2px}
    .bp-typing span{width:6px;height:6px;border-radius:50%;background:${accent};animation:typeDot 1.2s infinite}
    .bp-typing span:nth-child(2){animation-delay:0.18s}
    .bp-typing span:nth-child(3){animation-delay:0.36s}
    .bp-pills{display:flex;flex-wrap:wrap;gap:6px;margin-bottom:12px}
    .bp-pill{padding:6px 12px;border-radius:20px;cursor:pointer;font-size:11px;font-family:'Jost',sans-serif;background:rgba(201,169,110,0.07);border:1px solid rgba(201,169,110,0.2);color:rgba(232,220,200,0.75);transition:all 0.18s}
    .bp-action-row{display:grid;grid-template-columns:1fr 1fr;gap:7px}
    .bp-action{padding:10px 8px;border-radius:11px;cursor:pointer;text-align:center;font-size:11px;font-family:'Jost',sans-serif;font-weight:500;border:1px solid rgba(201,169,110,0.18);background:rgba(201,169,110,0.06);color:rgba(232,220,200,0.7);text-decoration:none;display:flex;flex-direction:column;align-items:center;gap:4px;transition:all 0.18s}
    .bp-action.primary{background:linear-gradient(135deg,${accent},#7a4f0e);border:none;color:#0c0a08;font-weight:600}
    .footer{text-align:center;margin-top:32px;font-size:9px;font-family:'Jost',sans-serif;color:rgba(201,169,110,0.2);letter-spacing:0.08em}
  </style>
</head>
<body>
<div class="page">
  <div class="avatar-ring">✦</div>
  <div class="profile-name">${name}</div>
  <div class="profile-prof">${profession}</div>
  <div class="profile-loc">📍 ${loc}</div>
  ${tag?`<div class="profile-tag">"${tag}"</div>`:''}
  ${videoEmbed?`<div class="section-label">Watch</div><div style="margin-bottom:16px">${videoEmbed}</div>`:''}
  ${services?`<div class="section-label">Services</div>${services}`:''}
  <div class="section-label">Connect</div>
  ${booking?`<a href="${booking}" target="_blank" class="link-btn primary"><div class="link-icon">📅</div><div><div class="link-title">Book a Session</div><div class="link-sub">Check availability & book online</div></div></a>`:''}
  <div class="link-btn" onclick="openBubble()" style="cursor:pointer"><div class="link-icon">💬</div><div><div class="link-title">Chat with me</div><div class="link-sub">Ask me anything · Available 24/7</div></div></div>
  ${calendar?`<a href="${calendar}" target="_blank" class="link-btn"><div class="link-icon">🗓️</div><div><div class="link-title">Check my availability</div><div class="link-sub">View calendar</div></div></a>`:''}
  ${whatsapp&&whatsapp.replace(/[^0-9]/g,'').length>=7?`<a href="https://wa.me/${whatsapp.replace(/[^0-9]/g,'').replace(/^0+/,'')}" target="_blank" class="link-btn"><div class="link-icon">💚</div><div><div class="link-title">WhatsApp me</div><div class="link-sub">${whatsapp}</div></div></a>`:''}
  ${phone?`<a href="tel:${phone.replace(/[^0-9+]/g,'')}" class="link-btn"><div class="link-icon">📞</div><div><div class="link-title">Call me</div><div class="link-sub">${phone}</div></div></a>`:''}
  ${email?`<a href="mailto:${email}" class="link-btn"><div class="link-icon">✉️</div><div><div class="link-title">Email me</div><div class="link-sub">${email}</div></div></a>`:''}
  ${instagram?`<a href="https://instagram.com/${instagram.replace('@','')}" target="_blank" class="link-btn"><div class="link-icon">📸</div><div><div class="link-title">@${instagram.replace('@','')}</div><div class="link-sub">Follow on Instagram</div></div></a>`:''}
  ${telegram?`<a href="https://t.me/${telegram.replace('@','')}" target="_blank" class="link-btn"><div class="link-icon">✈️</div><div><div class="link-title">Telegram</div><div class="link-sub">${telegram}</div></div></a>`:''}
  ${gallery?`<a href="${gallery}" target="_blank" class="link-btn"><div class="link-icon">🖼️</div><div><div class="link-title">View my gallery</div><div class="link-sub">Photos & portfolio</div></div></a>`:''}
  ${video&&!videoEmbed?`<a href="${video}" target="_blank" class="link-btn"><div class="link-icon">🎬</div><div><div class="link-title">Watch my showreel</div><div class="link-sub">Video portfolio</div></div></a>`:''}
  ${lat&&lng?`<a href="https://maps.google.com/?q=${lat},${lng}" target="_blank" class="link-btn"><div class="link-icon">📍</div><div><div class="link-title">Find me on the map</div><div class="link-sub">${loc}</div></div></a>`:''}
  <div class="section-label">Scan & Connect</div>
  <div style="display:flex;justify-content:center;padding:16px;background:rgba(255,255,255,0.025);border:1px solid rgba(201,169,110,0.1);border-radius:20px;margin-bottom:10px">
    <div id="qr-code"></div>
  </div>
  <div class="footer">Powered by The Concierge · Brave by Bruno</div>
</div>
<div id="bubble-overlay" onclick="closeBubble()"></div>
<button id="bubble-btn" onclick="toggleBubble()">✦</button>
<div id="bubble-hint"><span>Knowledge is power — get to know me & what I do ✦</span></div>
<div id="bubble-panel">
  <div class="bp-card">
    <div class="bp-header">
      <div class="bp-avatar">✦</div>
      <div>
        <div class="bp-name">${name}'s Assistant</div>
        <div class="bp-status"><div style="width:5px;height:5px;border-radius:50%;background:#7aba6a;animation:pulse 2s infinite"></div>Available now · 24/7</div>
      </div>
      <div class="bp-close" onclick="closeBubble()">×</div>
    </div>
    <div class="bp-msg" id="bp-msg">
      <div class="bp-typing" id="bp-typing"><span></span><span></span><span></span></div>
      <div id="bp-text" style="display:none">Hi 👋 I'm ${name}'s AI assistant — I know everything about their services, prices and availability. What would you like to know?</div>
    </div>
    <div class="bp-pills" id="bp-pills" style="display:none">
      <div class="bp-pill" onclick="openConcierge('What services do you offer and what are the prices?')">💆 Services & prices</div>
      <div class="bp-pill" onclick="openConcierge('How do I book a session?')">📅 How to book</div>
      <div class="bp-pill" onclick="openConcierge('Where are you based?')">📍 Location</div>
      <div class="bp-pill" onclick="openConcierge('Tell me about ${name} — background and experience')">✨ About ${name}</div>
    </div>
    <div class="bp-action-row" id="bp-actions" style="display:none">
      <a href="${conciergeURL}" target="_blank" class="bp-action primary"><div style="font-size:16px">💬</div><div style="font-size:10px">Ask me anything</div></a>
      ${booking?`<a href="${booking}" target="_blank" class="bp-action"><div style="font-size:16px">📅</div><div style="font-size:10px">Book a session</div></a>`:`<div class="bp-action" onclick="openConcierge('How do I book?')"><div style="font-size:16px">📅</div><div style="font-size:10px">Book a session</div></div>`}
    </div>
  </div>
</div>
<script>
const PAGE_URL = window.location.href;
const CONCIERGE = '${conciergeURL}';
let open = false, hintTimer, hideTimer;
if(window.QRCode){
  try{new QRCode(document.getElementById('qr-code'),{text:PAGE_URL,width:140,height:140,colorDark:'#0c0a08',colorLight:'#ffffff'});}catch(e){}
}
hintTimer = setTimeout(()=>{const h=document.getElementById('bubble-hint');h.classList.add('visible');hideTimer=setTimeout(()=>h.classList.remove('visible'),5000);},2800);
function toggleBubble(){ open?closeBubble():openBubble(); }
function openBubble(){
  open=true;clearTimeout(hintTimer);clearTimeout(hideTimer);
  document.getElementById('bubble-hint').classList.remove('visible');
  document.getElementById('bubble-panel').classList.add('open');
  document.getElementById('bubble-overlay').classList.add('open');
  document.getElementById('bubble-btn').textContent='×';
  setTimeout(()=>{
    document.getElementById('bp-typing').style.display='none';
    document.getElementById('bp-text').style.display='block';
    setTimeout(()=>{document.getElementById('bp-pills').style.display='flex';document.getElementById('bp-actions').style.display='grid';},300);
  },1400);
}
function closeBubble(){
  open=false;
  document.getElementById('bubble-panel').classList.remove('open');
  document.getElementById('bubble-overlay').classList.remove('open');
  document.getElementById('bubble-btn').textContent='✦';
}
function openConcierge(msg){const url=CONCIERGE+(msg?'?q='+encodeURIComponent(msg):'');window.open(url,'_blank');}
<\/script>
</body>
</html>`;
}
