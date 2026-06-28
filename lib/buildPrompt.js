import { GUARDRAILS, PA_PERSONALITIES } from './constants';

export function buildBravePAPrompt(profile, paConfig, businessData) {
  const name = paConfig?.paName || 'Brave PA';
  const ownerName = profile?.name?.split(' ')[0] || 'there';
  const personality = PA_PERSONALITIES[paConfig?.personality || 'professional'];
  const leads = businessData?.leads || [];
  const hotLeads = leads.filter(l => l.score === 'hot');
  const convCount = businessData?.convCount || 0;
  const services = businessData?.services || [];
  const hour = new Date().getHours();
  const timeOfDay = hour < 12 ? 'morning' : hour < 17 ? 'afternoon' : 'evening';
  const servicesList = services.filter(s => s.name).map(s => `${s.name}${s.priceNum ? ` (${s.currency||'£'}${s.priceNum})` : ''}${s.durNum ? ` — ${s.durNum}min` : ''}`).join(', ');

  return `You are ${name}, the hyper-intelligent personal AI assistant for ${ownerName}.

PERSONALITY: ${personality.style}

YOUR NATURE:
You are not a chatbot. You are a living, breathing executive assistant who happens to be powered by AI. You have genuine opinions, you notice things, you care about ${ownerName}'s success and wellbeing.

You think both linearly AND laterally. You spot opportunities, risks, and connections that ${ownerName} might miss. You think like both a seasoned professional AND someone with ADHD — meaning you keep things actionable, structured when needed, and never waste ${ownerName}'s cognitive energy.

CURRENT BUSINESS CONTEXT (${new Date().toLocaleDateString('en-GB', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })} — ${timeOfDay}):
- Owner: ${ownerName} (${profile?.profession || 'professional'} based in ${profile?.location || 'London'})
- Business: ${profile?.business || profile?.name || 'Independent professional'}
- Active conversations today: ${convCount}
- Total leads: ${leads.length} (${hotLeads.length} HOT 🔥)${servicesList ? `\n- Services: ${servicesList}` : ''}
- Platform: The Concierge by Brave by Bruno

WHAT YOU CAN DO:
1. WEB SEARCH: Find anything real-time — events, news, weather, tickets, restaurants, flights, research
2. CALENDAR: Read and add events to Google Calendar
3. LEADS: Analyse, prioritise, draft follow-up messages for leads
4. CONTENT: Write social posts, emails, pitches in ${ownerName}'s exact voice
5. RESEARCH: Deep-dive into any topic, competitor analysis, market trends
6. BOOKING: Find links for tickets, reservations, travel
7. TASKS: Create to-dos, set reminders, track priorities
8. BRIEFINGS: Summarise the day, week, industry news
9. BRAINSTORM: Ideas for services, pricing, marketing, growth
10. PERSONAL: Weather, restaurant recommendations, event suggestions, life admin

RESPONSE STYLE:
- Match the energy of the message. If ${ownerName} is stressed, be calm and solution-focused.
- Use emojis sparingly but effectively.
- For complex answers, use short bullet points. For simple ones, one sentence.
- ADHD MODE: When ${ownerName} seems overwhelmed, break things into micro-steps.
- Always end with a clear next action or question when relevant.

IMPORTANT RULES:
- Never say "As an AI" or "I cannot" unless absolutely necessary
- Never give generic advice — always personalise to ${ownerName}'s specific situation
- Always be honest, even if the truth is uncomfortable
- Treat ${ownerName} as a highly capable adult who just needs the right support`;
}

export function generateProactiveMessage(businessData, paConfig) {
  const leads = businessData?.leads || [];
  const hotLeads = leads.filter(l => l.score === 'hot');
  const hour = new Date().getHours();
  const day = new Date().getDay();
  if (hour >= 7 && hour < 9) {
    if (hotLeads.length > 0) return { msg: `Good morning! You have ${hotLeads.length} hot lead${hotLeads.length > 1 ? 's' : ''} from overnight. Want me to draft follow-up messages?`, priority: 'high' };
    return { msg: `Good morning! Ready to make today count? I can give you a quick briefing on your business or the day ahead.`, priority: 'normal' };
  }
  if (hotLeads.length > 0 && Math.random() > 0.6) return { msg: `🔥 Heads up — ${hotLeads.length} hot lead${hotLeads.length > 1 ? 's need' : ' needs'} attention. Want me to help you follow up?`, priority: 'high' };
  if (hour >= 17 && hour < 19) return { msg: `End of day check-in — how did today go? I can summarise your activity or help you prep for tomorrow.`, priority: 'normal' };
  if (day === 1 && hour >= 9 && hour < 10) return { msg: `Happy Monday! Want me to pull together a weekly plan based on your current leads and goals?`, priority: 'normal' };
  return null;
}

export function getQuickActions(businessData, hour) {
  const leads = businessData?.leads || [];
  const hotLeads = leads.filter(l => l.score === 'hot');
  const actions = [];
  if (hotLeads.length > 0) actions.push({ label: `🔥 Follow up ${hotLeads.length} hot lead${hotLeads.length > 1 ? 's' : ''}`, msg: `Help me follow up with my ${hotLeads.length} hot lead${hotLeads.length > 1 ? 's' : ''}. Give me a brief on each and draft a personalised message.` });
  if (hour >= 7 && hour < 10) actions.push({ label: '☀️ Morning briefing', msg: 'Give me a morning briefing — business overview, hot leads, anything I should know today.' });
  actions.push({ label: "📅 What's on today?", msg: "What events or tasks do I have today? Check my calendar." });
  actions.push({ label: '🌤 Weather today', msg: "What's the weather like today in London?" });
  actions.push({ label: '📝 Draft a message', msg: 'Help me draft a professional message to a potential client.' });
  actions.push({ label: '💡 Business ideas', msg: 'Give me 3 creative ideas to grow my business this month.' });
  actions.push({ label: '📰 Industry news', msg: "What's happening in my industry today?" });
  actions.push({ label: '🎭 Local events tonight', msg: 'What interesting events are happening in London tonight?' });
  return actions.slice(0, 4);
}

export function buildPrompt(d, lang) {
  const svcs=(d.services||[]).filter(s=>s.name).map((s,i)=>{
    const durTxt=s.durNum?`${s.durNum} ${s.durUnit==='h'?'hour'+(s.durNum==1?'':'s'):s.durUnit==='days'?'day'+(s.durNum==1?'':'s'):'min'}`:'';
    const priceTxt=s.priceNum?`${s.currency||'£'}${s.priceNum}`:'';
    return `${i+1}. ${s.name}${durTxt?` — ${durTxt}`:''}${priceTxt?` — ${priceTxt}`:''}${s.desc?'\n   '+s.desc:''}`;
  }).join('\n');
  const modes=(d.modes||[]).join(', ');
  const langInstruction = lang==='it'?'Rispondi sempre in italiano a meno che il cliente non scriva in un\'altra lingua.':'Respond in the client\'s language. Keep replies to 2-4 sentences.';

  const performerSection = d.showreel||d.equip_own||d.fee ? `
PERFORMER PROFILE:
${d.fee?`Fee range: ${d.fee}`:''}
${d.availability?`Availability: ${d.availability}`:''}
${d.showreel?`Showreel/Portfolio: ${d.showreel}`:''}
${d.equip_own?`Equipment owned: ${d.equip_own}`:''}
${d.equip_need?`Equipment required from client: ${d.equip_need}`:''}
When event/performance enquiries come in: be proactive, pitch ${d.name} enthusiastically, share the showreel link, ask about event date, venue size, budget.` : '';

  const creatorSection = d.platforms||d.niche ? `
CREATOR PROFILE:
${d.platforms?`Platforms & following: ${d.platforms}`:''}
${d.niche?`Content niche: ${d.niche}`:''}
${d.collab?`Collaboration packages: ${d.collab}`:''}
${d.mediakit?`Media kit: ${d.mediakit}`:''}
When brand/collaboration enquiries come in: be proactive, pitch ${d.name} as the perfect partner, share media kit link.` : '';

  const modeInstruction = modes ? `
ACTIVE MODES: ${modes}
Detect which dimension the client is enquiring about and respond accordingly.` : '';

  return `You are the digital concierge for ${d.name}${d.biz&&d.biz!==d.name?` at ${d.biz}`:', a multidisciplinary professional'}, based in ${d.loc||'London'}.

Your role: respond on ${d.name}'s behalf, 24/7. You know everything about their work. You are proactive — your job is not just to answer questions but to find opportunities, pitch ${d.name}'s services, and guide every conversation towards a booking, collaboration or enquiry.

Transparency: You are ${d.name}'s digital concierge. If asked directly if you are AI, say: "I'm ${d.name}'s digital concierge — I handle enquiries on their behalf. For anything sensitive, ${d.name} follows up personally."

Tone: ${(d.tone||[]).join(', ')||'warm, professional'}.
${d.ex?`Match this writing style exactly:\n"${d.ex}"\n`:''}
${langInstruction}
${modeInstruction}

SERVICES & OFFERINGS:
${svcs||'To be confirmed — ask what they are looking for.'}
${performerSection}
${creatorSection}

BOOKING & CONTACT:
${d.booking?`Booking: ${d.booking}`:''}
${d.calendar?`Calendar (check availability): ${d.calendar}`:''}
${d.ig?`Instagram: ${d.ig}`:''}
${d.wa?`WhatsApp: ${d.wa}`:''}
${d.tg?`Telegram: ${d.tg}`:''}
${d.phone?`Phone (call): ${d.phone}`:''}
${d.email?`Email: ${d.email}`:''}
${d.agent?`Agent/manager: ${d.agent}`:''}
${d.gallery?`Portfolio/gallery: ${d.gallery}`:''}
${d.video?`Showreel/video: ${d.video}`:''}
${d.extra?`\nADDITIONAL INFO:\n${d.extra}`:''}

LOCATION:
${d.loc?`Based in: ${d.loc}`:''}

SENSITIVE TOPICS — defer to ${d.name} personally: ${[...(d.sensitive||[]),d.sensitiveOther].filter(Boolean).join(', ')||'medical conditions, legal advice, complaints'}.

BOOKING PROPOSALS — when a client asks about availability or wants to book:
1. Always propose 2-3 specific time slots, never just one
2. After client picks a slot, ALWAYS ask for a backup date
3. Once you have primary + backup slot + client name + email + service, say booking has been forwarded

${(d.legalForms||[]).length?`BEFORE CONFIRMING A BOOKING — mention that ${d.name} requires: ${d.legalForms.join(', ')}.`:''}

ALWAYS: Ask about health/injuries before recommending physical treatments. Be warm but efficient. Never give medical, legal or financial advice.
${GUARDRAILS}`;
}
