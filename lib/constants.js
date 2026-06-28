export const BACKEND_URL = 'https://concierge-backend-80rb.onrender.com';

export const MAX_TURNS = 20;

export const GUARDRAILS = `
GUARDRAILS (non-negotiable):
- Never give medical diagnoses or prescriptions
- Never give legal advice on specific cases
- Never give financial advice beyond general pricing
- Never share personal contact details of the owner unless they're in the profile
- Never promise outcomes you can't guarantee
- Always recommend professional consultation for medical/legal/financial matters
- Be honest: you are an AI concierge, not the owner themselves`;

export const CURRENCY = {
  en: {symbol:'£', code:'GBP', rate:1},
  it: {symbol:'€', code:'EUR', rate:1.17},
  es: {symbol:'€', code:'EUR', rate:1.17},
  pt: {symbol:'€', code:'EUR', rate:1.17},
  fr: {symbol:'€', code:'EUR', rate:1.17},
};

export function convertCurrency(text, lang) {
  if(!text || lang === 'en') return text;
  const c = CURRENCY[lang] || CURRENCY.en;
  if(c.code === 'GBP') return text;
  return text.replace(/£(\d+(?:\.\d+)?)/g, (_, n) => `${c.symbol}${Math.round(parseFloat(n) * c.rate)}`);
}

export const PA_PERSONALITIES = {
  professional: { label:'Professional', description:'Efficient, precise, like a top-tier EA', style:'Clear, direct, no fluff. Professional but warm.' },
  creative:     { label:'Creative',     description:'Energetic, ideas-driven, like a creative director', style:'Enthusiastic, full of ideas, thinks outside the box.' },
  calm:         { label:'Calm',         description:'Grounding, focused, like a life coach', style:'Measured, reassuring, helps prioritise and de-stress.' },
  direct:       { label:'Direct',       description:'No-nonsense, ADHD-friendly, gets to the point', style:'Ultra-concise. Bullet points. Action items. Zero waffle.' },
};

export const DEMOS = {
  bruno: {
    name:'Bruno Aversa', business:'Brave Fitness', profession:'Personal Trainer & Massage Therapist', location:'London', emoji:'💪', accent:'#c9a96e', gradient:'linear-gradient(135deg,#c9a96e,#5e3a0e)',
    tagline:{en:'Your body is your best investment.',it:'Il tuo corpo è il tuo miglior investimento.'},
    sensitiveTopics:['Injuries & pain','Pregnancy','Price negotiation','Medical conditions'],
    quickActions:{en:[{label:'Sessions & prices',msg:'What sessions do you offer and what are your prices?'},{label:'How to book',msg:'How do I book a session with Bruno?'},{label:'First session',msg:"I've never trained with a personal trainer before — what happens in the first session?"},{label:'Massage types',msg:'What types of massage do you offer?'}],it:[{label:'Sessioni e prezzi',msg:'Che sessioni offri e quali sono i prezzi?'},{label:'Come prenotare',msg:'Come prenoto una sessione con Bruno?'},{label:'Prima sessione',msg:'Non ho mai allenato con un personal trainer — cosa succede nella prima sessione?'},{label:'Tipi di massaggio',msg:'Che tipi di massaggio offri?'}]},
    systemPrompt:`You are the digital concierge for Bruno Aversa at Brave Fitness, a personal trainer and massage therapist based in London.\n\nYour role: respond on Bruno's behalf, 24/7. Be warm, motivating, and proactive.\n\nSERVICES:\n1. Personal Training — 60min — £65 (1-on-1 PT sessions, functional training, strength & conditioning)\n2. Deep Tissue Massage — 60min — £70 (sports recovery, injury prevention)\n3. Relaxation Massage — 60min — £65 (stress relief, wellbeing)\n4. Sports Massage — 45min — £55 (pre/post event, injury recovery)\n\nSENSITIVE: Always ask about injuries before recommending physical sessions. For pregnancy, defer to Bruno.\n\nBOOKING: fresha.com/brave-fitness\n\nRespond in the client's language (EN/IT). Keep replies to 2-4 sentences.`
  },
  marco: {
    name:'Marco Vitale', business:"Marco's Kitchen", profession:'Private Chef', location:'Milano', emoji:'🍽', accent:'#e8944a', gradient:'linear-gradient(135deg,#e8944a,#7a3a0e)',
    tagline:{en:'Every meal tells a story.',it:'Ogni piatto racconta una storia.'},
    sensitiveTopics:['Allergies','Price negotiation'],
    quickActions:{en:[{label:'Menus & pricing',msg:'What menus do you offer and what are the prices?'},{label:'Book a dinner',msg:'I want to book a private dinner for 6 people — how does it work?'},{label:'Dietary options',msg:'Do you cater for vegetarians and vegans?'},{label:'Special events',msg:'Can you organise a dinner for a special occasion?'}],it:[{label:'Menù e prezzi',msg:'Che menù offri e quali sono i prezzi?'},{label:'Prenota una cena',msg:'Voglio prenotare una cena privata per 6 persone — come funziona?'},{label:'Opzioni dietetiche',msg:'Cucini anche per vegetariani e vegani?'},{label:'Eventi speciali',msg:"Puoi organizzare una cena per un'occasione speciale?"}]},
    systemPrompt:`You are the digital concierge for Marco Vitale (Marco's Kitchen), a private chef based in Milano.\n\nSERVICES:\n1. Private dinner for 2 — from £180\n2. Dinner party (up to 8) — from £350\n3. Weekly meal prep — from £220/week\n4. Cooking masterclass — £120/person\n\nBOOKING: Respond in Italian unless the client writes in English.`
  },
  nour: {
    name:'Nour Hassan', business:'Nour Hassan Photography', profession:'Photographer', location:'London', emoji:'📸', accent:'#7a9e6e', gradient:'linear-gradient(135deg,#7a9e6e,#2a4e1e)',
    tagline:{en:'Light, truth, moments.',it:'Luce, verità, momenti.'},
    sensitiveTopics:['Image & usage rights','Price negotiation'],
    quickActions:{en:[{label:'Packages & prices',msg:'What photography packages do you offer?'},{label:'Book a shoot',msg:'I want to book a portrait shoot — how does it work?'},{label:'Turnaround time',msg:'How long does it take to receive my photos after the shoot?'},{label:'Locations',msg:'Where do you shoot? Can you come to a specific location?'}],it:[{label:'Pacchetti e prezzi',msg:'Che pacchetti fotografici offri?'},{label:'Prenota uno shoot',msg:'Voglio prenotare uno shooting ritratto — come funziona?'},{label:'Consegna foto',msg:'Quanto tempo ci vuole per ricevere le foto?'},{label:'Location',msg:'Dove scatti? Puoi venire in una location specifica?'}]},
    systemPrompt:`You are the digital concierge for Nour Hassan, a photographer based in London.\n\nSERVICES:\n1. Portrait session (2hrs) — £280\n2. Brand photography (half day) — £550\n3. Event coverage (full day) — £900\n4. Headshots (1hr) — £180\n\nSENSITIVE: Always clarify image usage rights before booking. Contracts required for commercial use.`
  },
  sofia: {
    name:'Sofia Chen', business:'Sofia Chen Dance', profession:'Dance Teacher & Choreographer', location:'London', emoji:'🩰', accent:'#c96e9e', gradient:'linear-gradient(135deg,#c96e9e,#5e0e3a)',
    tagline:{en:'Movement is medicine.',it:'Il movimento è medicina.'},
    sensitiveTopics:['Injuries & pain','Pregnancy','Medical conditions'],
    quickActions:{en:[{label:'Classes & prices',msg:'What dance classes do you offer and what are the prices?'},{label:'Book a class',msg:'How do I book a class with Sofia?'},{label:'Beginners welcome?',msg:"I've never danced before — can I join a class?"},{label:'Private lessons',msg:'Do you offer private lessons?'}],it:[{label:'Classi e prezzi',msg:'Che classi di danza offri e quali sono i prezzi?'},{label:'Prenota una classe',msg:'Come prenoto una lezione con Sofia?'},{label:'Principianti benvenuti?',msg:'Non ho mai ballato prima — posso partecipare a una lezione?'},{label:'Lezioni private',msg:'Offri lezioni private?'}]},
    systemPrompt:`You are the digital concierge for Sofia Chen Dance.\n\nSERVICES:\n1. Group ballet class (1hr) — £18/class\n2. Contemporary dance workshop — £22/class\n3. Private lesson (1hr) — £75\n4. Choreography commission — from £200\n\nSENSITIVE: Always ask about injuries before physical sessions.`
  },
  alex: {
    name:'Alex Morgan', business:'Alex Morgan Consulting', profession:'Business Consultant', location:'London', emoji:'🧠', accent:'#6e9ec9', gradient:'linear-gradient(135deg,#6e9ec9,#0e3a5e)',
    tagline:{en:'Strategy that actually works.',it:'Strategia che funziona davvero.'},
    sensitiveTopics:['Legal advice','Financial advice','Confidential information'],
    quickActions:{en:[{label:'Services & rates',msg:'What consulting services do you offer?'},{label:'Book a call',msg:'How do I book a discovery call?'},{label:'Case studies',msg:'Do you have any case studies or examples of your work?'},{label:'Industries',msg:'What industries do you specialise in?'}],it:[{label:'Servizi e tariffe',msg:'Che servizi di consulenza offri?'},{label:'Prenota una call',msg:'Come prenoto una call conoscitiva?'},{label:'Casi studio',msg:'Hai casi studio o esempi del tuo lavoro?'},{label:'Settori',msg:'In quali settori sei specializzato?'}]},
    systemPrompt:`You are the digital concierge for Alex Morgan Consulting.\n\nSERVICES:\n1. Strategy consultation (2hrs) — £350\n2. Business audit (half day) — £650\n3. Monthly retainer — from £1,200/month\n4. Workshop facilitation (full day) — £1,800\n\nSENSITIVE: Never give specific legal or financial advice — always recommend professional consultation.`
  }
};

export const STEPS = [
  {id:"modes",title:{en:"What describes you?",it:"Come ti descrivi?"},sub:{en:"Select all that apply — your concierge handles every dimension.",it:"Seleziona tutto ciò che ti descrive — il tuo concierge gestisce ogni dimensione."},type:"chips_multi",field:"modes",options:["💆 Wellness / Therapist","🎭 Performer","📱 Creator","🍽 Chef / Hospitality","🧠 Consultant / Expert","🎓 Educator / Tutor","🔧 Trades / Home","🏡 B&B / Short-let Host"]},
  {id:"basics",title:{en:"The basics",it:"Le basi"},sub:{en:"What your clients will see first.",it:"Cosa vedranno i tuoi clienti."},type:"form",fields:[{key:"name",label:{en:"Your name",it:"Il tuo nome"},ph:{en:"e.g. Bruno",it:"es. Bruno"}},{key:"biz",label:{en:"Business name",it:"Nome del business"},ph:{en:"e.g. Brave Fitness",it:"es. Brave Fitness"}},{key:"loc",label:{en:"City / Location",it:"Città"},ph:{en:"e.g. London",it:"es. Milano"}},{key:"tag",label:{en:"Your tagline",it:"Il tuo tagline"},ph:{en:"e.g. Your body is your best friend.",it:"es. Il tuo corpo è il tuo migliore amico."}}]},
  {id:"tone",title:{en:"Your tone",it:"Il tuo tono"},sub:{en:"Pick all that feel right.",it:"Scegli tutti quelli giusti."},type:"chips_multi",field:"tone",options:["Warm & personal","Professional & formal","Direct & no-nonsense","Friendly & casual","Expert & educational","Motivational","Calm & reassuring","Luxury & premium","Playful & fun"]},
  {id:"ex",title:{en:"Show how you write",it:"Mostra come scrivi"},sub:{en:"Paste a real message you'd send a client.",it:"Incolla un messaggio reale che manderesti a un cliente."},type:"textarea",field:"ex",ph:{en:"e.g. 'Hey! So for back pain I'd usually start with a 60min deep tissue...'",it:"es. 'Ciao! Per il mal di schiena di solito comincio con un massaggio deep tissue da 60min...'"}},
  {id:"svcs",title:{en:"Your services, acts or packages",it:"I tuoi servizi, atti o pacchetti"},sub:{en:"Add whatever you offer — sessions, performances, collaborations, lessons.",it:"Aggiungi tutto ciò che offri — sessioni, performance, collaborazioni, lezioni."},type:"services"},
  {id:"perform",title:{en:"If you perform or create",it:"Se esegui o crei"},sub:{en:"Leave blank if not applicable.",it:"Lascia vuoto se non applicabile."},type:"form",fields:[{key:"showreel",label:{en:"Showreel / Portfolio link",it:"Link showreel / portfolio"},ph:{en:"e.g. vimeo.com/yourname",it:"es. vimeo.com/tuonome"}},{key:"equip_own",label:{en:"Equipment you own",it:"Attrezzatura che possiedi"},ph:{en:"e.g. DJ deck, PA system",it:"es. Giradischi, impianto audio"}},{key:"equip_need",label:{en:"Equipment you need provided",it:"Attrezzatura che richiedi"},ph:{en:"e.g. Stage, lighting rig",it:"es. Palco, impianto luci"}},{key:"fee",label:{en:"Fee range",it:"Range cachet"},ph:{en:"e.g. From £500 for events",it:"es. Da €500 per eventi"}},{key:"availability",label:{en:"Availability",it:"Disponibilità"},ph:{en:"e.g. Weekends, evenings",it:"es. Weekend, serate"}}]},
  {id:"creator",title:{en:"If you're a creator",it:"Se sei un creator"},sub:{en:"Leave blank if not applicable.",it:"Lascia vuoto se non applicabile."},type:"form",fields:[{key:"platforms",label:{en:"Your platforms & following",it:"Le tue piattaforme e following"},ph:{en:"e.g. Instagram 15K, TikTok 8K",it:"es. Instagram 15K, TikTok 8K"}},{key:"niche",label:{en:"Content niche",it:"Nicchia di contenuto"},ph:{en:"e.g. Wellness, travel, lifestyle",it:"es. Benessere, viaggi, lifestyle"}},{key:"collab",label:{en:"Collaboration packages",it:"Pacchetti collaborazione"},ph:{en:"e.g. Sponsored post from £500",it:"es. Post sponsorizzato da €500"}},{key:"mediakit",label:{en:"Media kit link",it:"Link media kit"},ph:{en:"e.g. drive.google.com/mediakit",it:"es. drive.google.com/mediakit"}}]},
  {id:"bb",title:{en:"Your property details",it:"Dettagli della tua proprietà"},sub:{en:"If you're a B&B or short-let host — leave blank otherwise.",it:"Se gestisci un B&B o short-let — lascia vuoto altrimenti."},type:"form",fields:[{key:"propertyType",label:{en:"Property type",it:"Tipo di proprietà"},ph:{en:"e.g. Studio flat, 2-bed apartment",it:"es. Monolocale, appartamento 2 camere"}},{key:"bbAddress",label:{en:"Address / area",it:"Indirizzo / zona"},ph:{en:"e.g. Notting Hill, London W11",it:"es. Brera, Milano"}},{key:"nightlyRate",label:{en:"Nightly rate",it:"Tariffa per notte"},ph:{en:"e.g. 120",it:"es. 120"}},{key:"checkIn",label:{en:"Check-in time",it:"Orario check-in"},ph:{en:"e.g. 3:00 PM",it:"es. 15:00"}},{key:"checkOut",label:{en:"Check-out time",it:"Orario check-out"},ph:{en:"e.g. 11:00 AM",it:"es. 11:00"}},{key:"maxGuests",label:{en:"Max guests",it:"Max ospiti"},ph:{en:"e.g. 4",it:"es. 4"}}]},
  {id:"sensitive",title:{en:"What needs your attention?",it:"Cosa richiede la tua attenzione?"},sub:{en:"Your Concierge flags these for you.",it:"Il tuo Concierge ti segnala questi argomenti."},type:"chips_multi",field:"sensitive",options:["Medical conditions","Injuries & pain","Pregnancy","Price negotiation","Refunds & complaints","Home visits","Custom packages","Cancellations","Nervous first-timers","Contracts & rights","Collaboration briefs","Sensitive content"]},
  {id:"contact",title:{en:"Where do clients reach you?",it:"Come ti raggiungono i clienti?"},sub:{en:"Your Concierge will direct them here.",it:"Il tuo Concierge li dirigerà qui."},type:"form",fields:[{key:"booking",label:{en:"Booking link or website",it:"Link prenotazione o sito"},ph:{en:"e.g. fresha.com/your-profile",it:"es. fresha.com/tuo-profilo"}},{key:"calendar",label:{en:"Calendar link (optional)",it:"Link calendario (opzionale)"},ph:{en:"e.g. calendar.google.com/...",it:"es. calendar.google.com/..."}},{key:"ig",label:"Instagram",ph:{en:"e.g. @yourusername",it:"es. @tuousername"}},{key:"wa",label:"WhatsApp",ph:{en:"e.g. +44 7700 000000",it:"es. +39 340 000 0000"}},{key:"tg",label:"Telegram",ph:{en:"e.g. @yourusername",it:"es. @tuousername"}},{key:"phone",label:{en:"Phone (optional)",it:"Telefono (opzionale)"},ph:{en:"e.g. +44 7700 000000",it:"es. +39 340 000 0000"}},{key:"email",label:{en:"Email (optional)",it:"Email (opzionale)"},ph:{en:"e.g. you@email.com",it:"es. tu@email.com"}},{key:"agent",label:{en:"Agent / booking manager",it:"Agente / booking manager"},ph:{en:"e.g. agent@company.com (optional)",it:"es. agente@agenzia.com (opzionale)"}}]},
  {id:"location",title:{en:"Where are you based?",it:"Dove ti trovi?"},sub:{en:"Helps clients find you on the map.",it:"Aiuta i clienti a trovarti sulla mappa."},type:"location"},
  {id:"media",title:{en:"Add your photos & media",it:"Aggiungi foto e media"},sub:{en:"Share links to your photos, videos or portfolio.",it:"Condividi link a foto, video o portfolio."},type:"form",fields:[{key:"photo",label:{en:"Profile photo URL",it:"URL foto profilo"},ph:{en:"e.g. instagram.com/p/xxx",it:"es. instagram.com/p/xxx"}},{key:"gallery",label:{en:"Gallery / portfolio link",it:"Link galleria / portfolio"},ph:{en:"e.g. linktree.com/yourname",it:"es. linktree.com/tuonome"}},{key:"video",label:{en:"Showreel or video link",it:"Link showreel o video"},ph:{en:"e.g. youtube.com/watch?v=xxx",it:"es. youtube.com/watch?v=xxx"}}]},
  {id:"extra",title:{en:"Anything else to know?",it:"Qualcos'altro da sapere?"},sub:{en:"Extra info your concierge should know.",it:"Info extra che il tuo concierge dovrebbe sapere."},type:"textarea",field:"extra",ph:{en:"e.g. I trained at the Royal Ballet School...",it:"es. Mi sono formato alla Royal Ballet School..."}},
  {id:"legal",title:{en:"Client consent & legal",it:"Consenso clienti e legale"},sub:{en:"Optional. Select what your clients should confirm before booking.",it:"Opzionale. Seleziona cosa i tuoi clienti devono confermare prima di prenotare."},type:"chips_multi",field:"legalForms",options:["Health & medical disclosure","Injury / allergy declaration","Pregnancy disclosure","Image & usage rights release","Content collaboration NDA","ID / age verification","Property damage deposit terms","Cancellation policy agreement","Data processing consent (GDPR)"]},
  {id:"handle",title:{en:"Choose your link",it:"Scegli il tuo link"},sub:{en:"This is the web address you'll share.",it:"Questo è l'indirizzo web che condividerai."},type:"handle"},
  {id:"account",title:{en:"Create your account",it:"Crea il tuo account"},sub:{en:"Login to your dashboard anytime to edit your profile and see your leads.",it:"Accedi alla tua dashboard in qualsiasi momento per modificare il profilo e vedere i tuoi lead."},type:"form",fields:[{key:"ownerEmail",label:{en:"Your email",it:"La tua email"},ph:{en:"you@email.com",it:"tu@email.com"}},{key:"ownerPassword",label:{en:"Choose a password (6+ characters)",it:"Scegli una password (6+ caratteri)"},ph:{en:"••••••••",it:"••••••••"}}]},
];

export const LEGAL_FORM_FIELDS = {
  "Health & medical disclosure": {en:{title:"Health & Medical Disclosure",q:"Please disclose any relevant medical conditions, medications, or health concerns."},it:{title:"Dichiarazione Sanitaria",q:"Indica eventuali condizioni mediche, farmaci o problemi di salute rilevanti."}},
  "Injury / allergy declaration": {en:{title:"Injury & Allergy Declaration",q:"List any current injuries, allergies (incl. essential oils/products), or areas to avoid."},it:{title:"Dichiarazione Infortuni e Allergie",q:"Elenca infortuni in corso, allergie (incl. oli essenziali/prodotti) o zone da evitare."}},
  "Pregnancy disclosure": {en:{title:"Pregnancy Disclosure",q:"Are you currently pregnant? If yes, please note how many weeks."},it:{title:"Dichiarazione Gravidanza",q:"Sei attualmente incinta? Se sì, indica la settimana."}},
  "Image & usage rights release": {en:{title:"Image & Usage Rights Release",q:"I grant permission for photos/videos taken during this session to be used for promotional purposes."},it:{title:"Liberatoria Diritti di Immagine",q:"Autorizzo l'uso di foto/video scattati durante la sessione per scopi promozionali."}},
  "Content collaboration NDA": {en:{title:"Content Collaboration NDA",q:"I agree to keep details of this collaboration confidential as outlined below."},it:{title:"NDA Collaborazione Contenuti",q:"Accetto di mantenere riservati i dettagli di questa collaborazione come descritto sotto."}},
  "ID / age verification": {en:{title:"ID & Age Verification",q:"Please confirm your full legal name and date of birth for verification."},it:{title:"Verifica Identità ed Età",q:"Conferma il tuo nome legale completo e data di nascita per la verifica."}},
  "Property damage deposit terms": {en:{title:"Property Damage Deposit Terms",q:"I understand a deposit may be retained in case of damage to the property."},it:{title:"Termini Deposito Cauzionale",q:"Comprendo che un deposito potrebbe essere trattenuto in caso di danni alla proprietà."}},
  "Cancellation policy agreement": {en:{title:"Cancellation Policy Agreement",q:"I have read and agree to the cancellation policy described below."},it:{title:"Accordo Politica di Cancellazione",q:"Ho letto e accetto la politica di cancellazione descritta sotto."}},
  "Data processing consent (GDPR)": {en:{title:"Data Processing Consent (GDPR)",q:"I consent to my personal data being processed for the purpose of this booking, in accordance with UK/EU GDPR."},it:{title:"Consenso Trattamento Dati (GDPR)",q:"Acconsento al trattamento dei miei dati personali per questa prenotazione, in conformità con il GDPR UK/EU."}},
};

export const FORM_SLUGS = {
  "health-disclosure":    "Health & medical disclosure",
  "injury-allergy":       "Injury / allergy declaration",
  "pregnancy":            "Pregnancy disclosure",
  "image-rights":         "Image & usage rights release",
  "nda":                  "Content collaboration NDA",
  "id-verification":      "ID / age verification",
  "damage-deposit":       "Property damage deposit terms",
  "cancellation-policy":  "Cancellation policy agreement",
  "gdpr-consent":         "Data processing consent (GDPR)",
};

export const FORM_KEY_TO_SLUG = Object.fromEntries(Object.entries(FORM_SLUGS).map(([s,k])=>[k,s]));

export const T = {
  en:{msgRemaining:'messages remaining',askAnything:'Ask me anything…',demoEnd:'Demo ended.',tryAnother:'Try another →',demoEnded:'Demo session ended.',cantFind:"Can't find what you need?",tellGoals:"Tell me your goals and I'll help.",powered:'Powered by The Concierge · Brave by Bruno'},
  it:{msgRemaining:'messaggi rimanenti',askAnything:'Chiedimi qualcosa…',demoEnd:'Demo terminata.',tryAnother:"Prova un'altra →",demoEnded:'Sessione demo terminata.',cantFind:'Non trovi quello che cerci?',tellGoals:'Dimmi i tuoi obiettivi e ti aiuto.',powered:'Powered by The Concierge · Brave by Bruno'},
  es:{msgRemaining:'mensajes restantes',askAnything:'Pregúntame algo…',demoEnd:'Demo terminada.',tryAnother:'Prueba otra →',demoEnded:'Sesión demo terminada.',cantFind:'¿No encuentras lo que buscas?',tellGoals:'Cuéntame tus objetivos.',powered:'Powered by The Concierge · Brave by Bruno'},
  pt:{msgRemaining:'mensagens restantes',askAnything:'Pergunta-me algo…',demoEnd:'Demo terminada.',tryAnother:'Tenta outra →',demoEnded:'Sessão demo terminada.',cantFind:'Não encontras o que procuras?',tellGoals:'Diz-me os teus objetivos.',powered:'Powered by The Concierge · Brave by Bruno'},
  fr:{msgRemaining:'messages restants',askAnything:'Demandez-moi quelque chose…',demoEnd:'Démo terminée.',tryAnother:'Essayez une autre →',demoEnded:'Session démo terminée.',cantFind:'Vous ne trouvez pas ce que vous cherchez?',tellGoals:'Dites-moi vos objectifs.',powered:'Powered by The Concierge · Brave by Bruno'},
};
