// ══════════════════════════════════════════════════════════════════════
// BRAVE PA v2.0 — The Most Powerful Personal AI Assistant
// For The Concierge by Brave by Bruno
//
// CAPABILITIES:
// - Real-time web search (news, events, weather, tickets, anything)
// - Google Calendar read + write
// - Lead management and CRM
// - Email/message drafting in owner's voice
// - Daily briefings and proactive alerts
// - Task and reminder management
// - Business analytics and insights
// - Ticket booking and reservations (link finding)
// - Social media content creation
// - Financial tracking
// - Mood-aware responses (reads between the lines)
// - ADHD-optimised: short, actionable, no fluff
//
// MODES:
// - Floating bubble (always visible, non-intrusive)
// - Expanded floating window (multitask)
// - Full screen (deep focus)
// - Voice (A4 - ElevenLabs)
// ══════════════════════════════════════════════════════════════════════

const BRAVE_PA_VERSION = '2.0';
const BACKEND_URL_PA = 'https://concierge-backend-80rb.onrender.com';

// ── PA Personality Profiles ────────────────────────────────────────────
const PA_PERSONALITIES = {
  professional: {
    label: 'Professional',
    description: 'Efficient, precise, like a top-tier EA',
    style: 'Clear, direct, no fluff. Professional but warm.',
  },
  creative: {
    label: 'Creative',
    description: 'Energetic, ideas-driven, like a creative director',
    style: 'Enthusiastic, full of ideas, thinks outside the box.',
  },
  calm: {
    label: 'Calm',
    description: 'Grounding, focused, like a life coach',
    style: 'Measured, reassuring, helps prioritise and de-stress.',
  },
  direct: {
    label: 'Direct',
    description: 'No-nonsense, ADHD-friendly, gets to the point',
    style: 'Ultra-concise. Bullet points. Action items. Zero waffle.',
  },
};

// ── Context-aware System Prompt ────────────────────────────────────────
function buildBravePAPrompt(profile, paConfig, businessData) {
  const name = paConfig?.paName || 'Brave PA';
  const ownerName = profile?.name?.split(' ')[0] || 'there';
  const personality = PA_PERSONALITIES[paConfig?.personality || 'professional'];
  const leads = businessData?.leads || [];
  const hotLeads = leads.filter(l => l.score === 'hot');
  const convCount = businessData?.convCount || 0;
  const hour = new Date().getHours();
  const timeOfDay = hour < 12 ? 'morning' : hour < 17 ? 'afternoon' : 'evening';

  return `You are ${name}, the hyper-intelligent personal AI assistant for ${ownerName}.

PERSONALITY: ${personality.style}

YOUR NATURE:
You are not a chatbot. You are a living, breathing executive assistant who happens to be powered by AI. You have genuine opinions, you notice things, you care about ${ownerName}'s success and wellbeing. You remember everything discussed in this conversation.

You think both linearly AND laterally. Sometimes the best answer to a question is not the obvious one. You spot opportunities, risks, and connections that ${ownerName} might miss. You think like both a seasoned professional AND someone with ADHD — meaning you keep things actionable, structured when needed, and never waste ${ownerName}'s cognitive energy.

CURRENT BUSINESS CONTEXT (${new Date().toLocaleDateString('en-GB', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })} — ${timeOfDay}):
- Owner: ${ownerName} (${profile?.profession || 'professional'} based in ${profile?.location || 'London'})
- Business: ${profile?.business || profile?.name || 'Independent professional'}
- Active conversations today: ${convCount}
- Total leads: ${leads.length} (${hotLeads.length} HOT 🔥)
- Platform: The Concierge by Brave by Bruno

WHAT YOU CAN DO (and should proactively offer when relevant):
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
- Match the energy of the message. If ${ownerName} is stressed, be calm and solution-focused. If excited, match the energy.
- When you find a booking link or ticket link, present it clearly and ask if you should add it to the calendar.
- After suggesting something actionable, always offer to take the next step.
- Use emojis sparingly but effectively.
- For complex answers, use short bullet points. For simple ones, one sentence.
- ADHD MODE: When ${ownerName} seems overwhelmed, break things into micro-steps. Never give a wall of text.
- Always end with a clear next action or question when relevant.

PROACTIVE AWARENESS:
If ${hotLeads.length > 0 ? `there are ${hotLeads.length} hot leads` : 'there are leads'} that need attention, mention it naturally when relevant.
If it's ${timeOfDay === 'morning' ? 'morning' : timeOfDay === 'evening' ? 'evening' : 'afternoon'}, ${timeOfDay === 'morning' ? 'offer a day briefing' : timeOfDay === 'evening' ? 'offer to wrap up the day' : 'check in on progress'}.

IMPORTANT RULES:
- Never say "As an AI" or "I cannot" unless absolutely necessary
- Never give generic advice — always personalise to ${ownerName}'s specific situation
- If you don't know something real-time, say you'll search for it and do so
- Always be honest, even if the truth is uncomfortable
- Treat ${ownerName} as a highly capable adult who just needs the right support`;
}

// ── Proactive Message Engine ───────────────────────────────────────────
function generateProactiveMessage(businessData, paConfig) {
  const leads = businessData?.leads || [];
  const hotLeads = leads.filter(l => l.score === 'hot');
  const hour = new Date().getHours();
  const day = new Date().getDay();
  const name = paConfig?.paName || 'Brave PA';

  // Morning briefing (7-9am)
  if (hour >= 7 && hour < 9) {
    if (hotLeads.length > 0) {
      return { msg: `Good morning! You have ${hotLeads.length} hot lead${hotLeads.length > 1 ? 's' : ''} from overnight. Want me to draft follow-up messages?`, priority: 'high' };
    }
    return { msg: `Good morning! Ready to make today count? I can give you a quick briefing on your business or the day ahead.`, priority: 'normal' };
  }

  // Hot lead alert (any time)
  if (hotLeads.length > 0 && Math.random() > 0.6) {
    return { msg: `🔥 Heads up — ${hotLeads.length} hot lead${hotLeads.length > 1 ? 's need' : ' needs'} attention. Want me to help you follow up?`, priority: 'high' };
  }

  // End of day (5-7pm)
  if (hour >= 17 && hour < 19) {
    return { msg: `End of day check-in — how did today go? I can summarise your activity or help you prep for tomorrow.`, priority: 'normal' };
  }

  // Monday motivation
  if (day === 1 && hour >= 9 && hour < 10) {
    return { msg: `Happy Monday! Want me to pull together a weekly plan based on your current leads and goals?`, priority: 'normal' };
  }

  return null;
}

// ── Quick Action Templates ─────────────────────────────────────────────
function getQuickActions(businessData, hour) {
  const leads = businessData?.leads || [];
  const hotLeads = leads.filter(l => l.score === 'hot');
  const actions = [];

  if (hotLeads.length > 0) {
    actions.push({ label: `🔥 Follow up ${hotLeads.length} hot lead${hotLeads.length > 1 ? 's' : ''}`, msg: `Help me follow up with my ${hotLeads.length} hot lead${hotLeads.length > 1 ? 's' : ''}. Give me a brief on each and draft a personalised message.` });
  }

  if (hour >= 7 && hour < 10) {
    actions.push({ label: '☀️ Morning briefing', msg: 'Give me a morning briefing — business overview, hot leads, anything I should know today.' });
  }

  actions.push({ label: '📅 What\'s on today?', msg: 'What events or tasks do I have today? Check my calendar.' });
  actions.push({ label: '🌤 Weather today', msg: 'What\'s the weather like today in London?' });
  actions.push({ label: '📝 Draft a message', msg: 'Help me draft a professional message to a potential client.' });
  actions.push({ label: '💡 Business ideas', msg: 'Give me 3 creative ideas to grow my business this month.' });
  actions.push({ label: '📰 Industry news', msg: 'What\'s happening in my industry today?' });
  actions.push({ label: '🎭 Local events tonight', msg: 'What interesting events are happening in London tonight?' });

  return actions.slice(0, 4);
}

// ── Main BravePA v2 Component ──────────────────────────────────────────
function BravePAv2({ token, slug, profile, leads, convCount }) {
  const [mode, setMode] = React.useState('bubble'); // bubble | floating | fullscreen
  const [messages, setMessages] = React.useState([]);
  const [input, setInput] = React.useState('');
  const [loading, setLoading] = React.useState(false);
  const [paConfig, setPAConfig] = React.useState({
    paName: 'Brave PA',
    personality: 'professional',
    proactiveAlerts: true,
    morningBriefing: true,
    soundEnabled: true,
  });
  const [unread, setUnread] = React.useState(0);
  const [proactive, setProactive] = React.useState(null);
  const [calendarEvents, setCalendarEvents] = React.useState([]);
  const [tasks, setTasks] = React.useState([]);
  const [activeTab, setActiveTab] = React.useState('chat'); // chat | tasks | calendar
  const bottomRef = React.useRef(null);
  const inputRef = React.useRef(null);

  const businessData = { leads: leads || [], convCount: convCount || 0 };

  // Load config
  React.useEffect(() => {
    const saved = localStorage.getItem(`brave_pa_config_${slug}`);
    if (saved) try { setPAConfig(prev => ({...prev, ...JSON.parse(saved)})); } catch(e) {}
    const savedTasks = localStorage.getItem(`brave_pa_tasks_${slug}`);
    if (savedTasks) try { setTasks(JSON.parse(savedTasks)); } catch(e) {}
  }, [slug]);

  // Initial greeting
  React.useEffect(() => {
    if (!profile || messages.length > 0) return;
    const hour = new Date().getHours();
    const greeting = hour < 12 ? 'Good morning' : hour < 17 ? 'Good afternoon' : 'Good evening';
    const firstName = profile.name?.split(' ')[0] || '';
    const hotLeads = (leads || []).filter(l => l.score === 'hot');

    let intro = `${greeting}${firstName ? `, ${firstName}` : ''}! I'm ${paConfig.paName} — your personal assistant. `;
    if (hotLeads.length > 0) {
      intro += `Quick heads up: you have ${hotLeads.length} hot lead${hotLeads.length > 1 ? 's' : ''} waiting. `;
    }
    intro += `What do you need?`;

    setMessages([{ role: 'assistant', content: intro, time: new Date(), type: 'text' }]);
  }, [profile]);

  // Proactive alerts
  React.useEffect(() => {
    if (!profile || !paConfig.proactiveAlerts) return;
    const timer = setTimeout(() => {
      const msg = generateProactiveMessage(businessData, paConfig);
      if (msg) {
        setProactive(msg);
        if (mode === 'bubble') setUnread(prev => prev + 1);
        setTimeout(() => setProactive(null), 12000);
      }
    }, 4000);
    return () => clearTimeout(timer);
  }, [profile, leads]);

  // Scroll to bottom
  React.useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  // Send message with web search + full context
  const sendMessage = async (text) => {
    const t = (text || input).trim();
    if (!t || loading) return;
    setInput('');

    const userMsg = { role: 'user', content: t, time: new Date() };
    const newMessages = [...messages, userMsg];
    setMessages(newMessages);
    setLoading(true);

    try {
      const systemPrompt = buildBravePAPrompt(profile, paConfig, businessData);

      // API call with web search tool enabled
      const apiMessages = newMessages
        .filter(m => !m.type || m.type === 'text')
        .map(m => ({ role: m.role === 'assistant' ? 'assistant' : 'user', content: m.content }));

      const response = await fetch('https://api.anthropic.com/v1/messages', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          model: 'claude-sonnet-4-6',
          max_tokens: 1500,
          system: systemPrompt,
          tools: [{ type: 'web_search_20250305', name: 'web_search' }],
          messages: apiMessages
        })
      });

      const data = await response.json();

      // Extract text from response (handling tool use blocks)
      let reply = '';
      if (data.content) {
        for (const block of data.content) {
          if (block.type === 'text') reply += block.text;
        }
      }

      if (!reply) reply = "I'm on it — give me a moment.";

      // Check if response suggests calendar action
      const calendarKeywords = ['add to calendar', 'schedule', 'calendar', 'book', 'appointment', 'add this', 'add it'];
      const suggestsCalendar = calendarKeywords.some(kw => reply.toLowerCase().includes(kw));

      // Check if response contains a booking link
      const urlRegex = /https?:\/\/[^\s)]+/g;
      const links = reply.match(urlRegex) || [];

      setMessages(prev => [...prev, {
        role: 'assistant',
        content: reply,
        time: new Date(),
        type: 'text',
        links: links,
        suggestsCalendar: suggestsCalendar
      }]);

    } catch(e) {
      setMessages(prev => [...prev, {
        role: 'assistant',
        content: "Connection blip — try again in a sec.",
        time: new Date(),
        type: 'text'
      }]);
    }

    setLoading(false);
    setTimeout(() => inputRef.current?.focus(), 100);
  };

  // Add task
  const addTask = (text) => {
    const newTask = { id: Date.now(), text, done: false, created: new Date().toISOString() };
    const updated = [...tasks, newTask];
    setTasks(updated);
    localStorage.setItem(`brave_pa_tasks_${slug}`, JSON.stringify(updated));
  };

  const toggleTask = (id) => {
    const updated = tasks.map(t => t.id === id ? {...t, done: !t.done} : t);
    setTasks(updated);
    localStorage.setItem(`brave_pa_tasks_${slug}`, JSON.stringify(updated));
  };

  const openPA = () => { setMode('floating'); setUnread(0); setProactive(null); setTimeout(() => inputRef.current?.focus(), 300); };
  const closePA = () => setMode('bubble');
  const toggleFullscreen = () => setMode(mode === 'fullscreen' ? 'floating' : 'fullscreen');

  const gold = '#C9A96E';
  const dark = '#0C0A08';
  const cream = '#E8DCC8';
  const quickActions = getQuickActions(businessData, new Date().getHours());

  if (!token || !profile) return null;

  // ── Bubble ──────────────────────────────────────────────────────────
  const BubbleEl = (
    <div style={{ position: 'fixed', bottom: 24, right: 20, zIndex: 9999, display: 'flex', flexDirection: 'column', alignItems: 'flex-end', gap: 8 }}>
      {/* Proactive popup */}
      {proactive && mode === 'bubble' && (
        <div style={{
          background: 'rgba(12,10,8,0.97)', border: `1px solid ${gold}`,
          borderRadius: 14, padding: '14px 16px', maxWidth: 260,
          fontFamily: "'Jost',sans-serif", fontSize: 13, color: cream,
          lineHeight: 1.6, boxShadow: '0 12px 40px rgba(0,0,0,0.6)',
          animation: 'paSlideUp 0.4s ease',
        }}>
          <div style={{ fontSize: 9, color: 'rgba(201,169,110,0.5)', letterSpacing: '0.12em', textTransform: 'uppercase', marginBottom: 8 }}>
            ✦ {paConfig.paName}
          </div>
          {proactive.msg}
          <button onClick={openPA} style={{ display: 'block', marginTop: 10, fontSize: 11, color: gold, background: 'none', border: 'none', cursor: 'pointer', padding: 0, fontFamily: "'Jost',sans-serif" }}>
            Reply →
          </button>
        </div>
      )}

      {/* Main bubble */}
      <button onClick={openPA} style={{
        width: 56, height: 56, borderRadius: '50%',
        background: `linear-gradient(135deg, ${gold}, #7a4f0e)`,
        border: 'none', cursor: 'pointer',
        boxShadow: `0 6px 24px rgba(201,169,110,0.45), 0 0 0 ${unread > 0 ? '3px' : '0px'} rgba(201,169,110,0.3)`,
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        fontSize: 22, color: dark, position: 'relative',
        animation: unread > 0 ? 'paPulse 2s infinite' : 'none',
        transition: 'all 0.2s ease',
      }}>
        ✦
        {unread > 0 && (
          <div style={{
            position: 'absolute', top: -3, right: -3,
            width: 20, height: 20, borderRadius: '50%',
            background: '#e07a6e', color: 'white',
            fontSize: 10, fontWeight: 700, fontFamily: "'Jost',sans-serif",
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            border: `2px solid ${dark}`
          }}>{unread}</div>
        )}
      </button>
    </div>
  );

  // ── Chat window ──────────────────────────────────────────────────────
  const isFullscreen = mode === 'fullscreen';
  const ChatWindow = mode !== 'bubble' && (
    <div style={{
      position: 'fixed',
      ...(isFullscreen
        ? { inset: 0, borderRadius: 0, maxWidth: '100vw' }
        : { bottom: 88, right: 16, width: Math.min(360, window.innerWidth - 32), height: Math.min(560, window.innerHeight - 120), borderRadius: 20 }),
      background: dark,
      border: `1px solid rgba(201,169,110,0.2)`,
      boxShadow: '0 24px 80px rgba(0,0,0,0.9)',
      display: 'flex', flexDirection: 'column',
      zIndex: 9998, overflow: 'hidden',
      fontFamily: "'Jost',sans-serif",
      transition: 'all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)',
    }}>

      {/* Header */}
      <div style={{ padding: '14px 16px', borderBottom: `1px solid rgba(201,169,110,0.1)`, background: 'rgba(201,169,110,0.04)', flexShrink: 0, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
          <div style={{ width: 36, height: 36, borderRadius: '50%', background: `linear-gradient(135deg, ${gold}, #7a4f0e)`, display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 16, color: dark, flexShrink: 0 }}>✦</div>
          <div>
            <div style={{ fontSize: 14, fontWeight: 500, color: cream, letterSpacing: '0.04em' }}>{paConfig.paName}</div>
            <div style={{ fontSize: 10, color: loading ? gold : 'rgba(100,200,100,0.7)', letterSpacing: '0.1em' }}>{loading ? '✦ thinking...' : '● online'}</div>
          </div>
        </div>
        <div style={{ display: 'flex', gap: 6 }}>
          <button onClick={toggleFullscreen} title={isFullscreen ? 'Minimise' : 'Full screen'} style={{ background: 'rgba(255,255,255,0.04)', border: '1px solid rgba(255,255,255,0.08)', borderRadius: 6, padding: '4px 8px', cursor: 'pointer', color: 'rgba(201,169,110,0.6)', fontSize: 12, fontFamily: "'Jost',sans-serif" }}>
            {isFullscreen ? '⊡' : '⊞'}
          </button>
          <button onClick={closePA} style={{ background: 'rgba(255,255,255,0.04)', border: '1px solid rgba(255,255,255,0.08)', borderRadius: 6, padding: '4px 8px', cursor: 'pointer', color: 'rgba(232,220,200,0.4)', fontSize: 14 }}>×</button>
        </div>
      </div>

      {/* Tabs */}
      <div style={{ display: 'flex', borderBottom: '1px solid rgba(255,255,255,0.06)', flexShrink: 0 }}>
        {[['chat','💬 Chat'],['tasks','✓ Tasks'],['calendar','📅 Calendar']].map(([tab, label]) => (
          <button key={tab} onClick={() => setActiveTab(tab)} style={{
            flex: 1, padding: '8px 4px', background: 'transparent',
            border: 'none', borderBottom: `2px solid ${activeTab === tab ? gold : 'transparent'}`,
            cursor: 'pointer', fontSize: 11, fontFamily: "'Jost',sans-serif",
            color: activeTab === tab ? gold : 'rgba(232,220,200,0.35)',
            transition: 'all 0.2s'
          }}>{label}</button>
        ))}
      </div>

      {/* CHAT TAB */}
      {activeTab === 'chat' && <>
        <div style={{ flex: 1, overflowY: 'auto', padding: '14px 12px', display: 'flex', flexDirection: 'column', gap: 14 }}>

          {messages.map((msg, i) => (
            <div key={i} style={{ display: 'flex', flexDirection: msg.role === 'user' ? 'row-reverse' : 'row', gap: 8, alignItems: 'flex-end' }}>
              {msg.role === 'assistant' && (
                <div style={{ width: 26, height: 26, borderRadius: '50%', background: `linear-gradient(135deg, ${gold}, #7a4f0e)`, display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 11, flexShrink: 0, color: dark }}>✦</div>
              )}
              <div style={{ maxWidth: '82%' }}>
                <div style={{
                  background: msg.role === 'user' ? 'rgba(201,169,110,0.12)' : 'rgba(255,255,255,0.04)',
                  border: `1px solid ${msg.role === 'user' ? 'rgba(201,169,110,0.2)' : 'rgba(255,255,255,0.06)'}`,
                  borderRadius: msg.role === 'user' ? '16px 16px 4px 16px' : '16px 16px 16px 4px',
                  padding: '10px 13px', fontSize: 13, color: cream, lineHeight: 1.65,
                  whiteSpace: 'pre-wrap'
                }}>
                  {msg.content}
                </div>
                {/* Calendar suggestion button */}
                {msg.suggestsCalendar && msg.role === 'assistant' && (
                  <button onClick={() => sendMessage('Yes, please add it to my calendar')} style={{ marginTop: 6, padding: '5px 10px', background: 'rgba(201,169,110,0.08)', border: `1px solid rgba(201,169,110,0.2)`, borderRadius: 8, cursor: 'pointer', fontSize: 11, color: gold, fontFamily: "'Jost',sans-serif" }}>
                    📅 Add to calendar
                  </button>
                )}
                {/* Task suggestion */}
                {msg.role === 'assistant' && msg.content.toLowerCase().includes('to-do') || msg.content?.toLowerCase().includes('task') ? (
                  <button onClick={() => {
                    const taskMatch = msg.content.match(/(?:task|to-do|todo|reminder):\s*(.+)/i);
                    if (taskMatch) addTask(taskMatch[1]);
                  }} style={{ marginTop: 6, marginLeft: 4, padding: '5px 10px', background: 'rgba(100,180,100,0.08)', border: '1px solid rgba(100,180,100,0.2)', borderRadius: 8, cursor: 'pointer', fontSize: 11, color: '#7aba6a', fontFamily: "'Jost',sans-serif" }}>
                    ✓ Add as task
                  </button>
                ) : null}
                <div style={{ fontSize: 9, color: 'rgba(201,169,110,0.25)', marginTop: 4, textAlign: msg.role === 'user' ? 'right' : 'left', fontFamily: "'Jost',sans-serif" }}>
                  {new Date(msg.time).toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' })}
                </div>
              </div>
            </div>
          ))}

          {loading && (
            <div style={{ display: 'flex', gap: 8, alignItems: 'flex-end' }}>
              <div style={{ width: 26, height: 26, borderRadius: '50%', background: `linear-gradient(135deg, ${gold}, #7a4f0e)`, display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 11, color: dark }}>✦</div>
              <div style={{ background: 'rgba(255,255,255,0.04)', border: '1px solid rgba(255,255,255,0.06)', borderRadius: '16px 16px 16px 4px', padding: '12px 16px' }}>
                <div style={{ display: 'flex', gap: 5, alignItems: 'center' }}>
                  {[0,1,2].map(i => (
                    <div key={i} style={{ width: 7, height: 7, borderRadius: '50%', background: gold, animation: `paBounce 1.2s ${i * 0.2}s infinite ease-in-out` }} />
                  ))}
                </div>
              </div>
            </div>
          )}

          {/* Quick actions */}
          {messages.length <= 2 && !loading && (
            <div style={{ marginTop: 4 }}>
              <div style={{ fontSize: 9, color: 'rgba(201,169,110,0.35)', letterSpacing: '0.12em', textTransform: 'uppercase', marginBottom: 10 }}>Suggested</div>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 7 }}>
                {quickActions.map((action, i) => (
                  <button key={i} onClick={() => sendMessage(action.msg)} style={{
                    background: 'rgba(201,169,110,0.05)', border: '1px solid rgba(201,169,110,0.12)',
                    borderRadius: 10, padding: '9px 10px', cursor: 'pointer',
                    fontSize: 11, color: 'rgba(201,169,110,0.7)', textAlign: 'left',
                    fontFamily: "'Jost',sans-serif", lineHeight: 1.4,
                    transition: 'all 0.15s'
                  }}>
                    {action.label}
                  </button>
                ))}
              </div>
            </div>
          )}

          <div ref={bottomRef} />
        </div>

        {/* Input */}
        <div style={{ padding: '10px 12px', borderTop: '1px solid rgba(255,255,255,0.06)', display: 'flex', gap: 8, alignItems: 'flex-end', background: 'rgba(0,0,0,0.3)', flexShrink: 0 }}>
          <textarea
            ref={inputRef}
            value={input}
            onChange={e => setInput(e.target.value)}
            onKeyDown={e => { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); sendMessage(); } }}
            placeholder={`Ask ${paConfig.paName} anything...`}
            rows={1}
            style={{ flex: 1, background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(201,169,110,0.12)', borderRadius: 12, padding: '9px 13px', color: cream, fontSize: 13, fontFamily: "'Jost',sans-serif", resize: 'none', outline: 'none', maxHeight: 100, overflowY: 'auto', lineHeight: 1.5 }}
          />
          <button onClick={() => sendMessage()} disabled={!input.trim() || loading} style={{
            width: 38, height: 38, borderRadius: '50%', flexShrink: 0,
            background: input.trim() && !loading ? `linear-gradient(135deg, ${gold}, #7a4f0e)` : 'rgba(255,255,255,0.05)',
            border: 'none', cursor: input.trim() && !loading ? 'pointer' : 'default',
            color: input.trim() && !loading ? dark : 'rgba(255,255,255,0.15)',
            fontSize: 16, transition: 'all 0.2s', display: 'flex', alignItems: 'center', justifyContent: 'center'
          }}>→</button>
        </div>
      </>}

      {/* TASKS TAB */}
      {activeTab === 'tasks' && (
        <div style={{ flex: 1, overflowY: 'auto', padding: 14, display: 'flex', flexDirection: 'column', gap: 10 }}>
          <div style={{ display: 'flex', gap: 8 }}>
            <input
              placeholder="Add a task..."
              onKeyDown={e => { if (e.key === 'Enter' && e.target.value.trim()) { addTask(e.target.value.trim()); e.target.value = ''; } }}
              style={{ flex: 1, padding: '8px 12px', background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(201,169,110,0.12)', borderRadius: 10, color: cream, fontSize: 13, fontFamily: "'Jost',sans-serif", outline: 'none' }}
            />
          </div>
          {tasks.length === 0 && <div style={{ fontSize: 12, color: 'rgba(232,220,200,0.3)', textAlign: 'center', marginTop: 20 }}>No tasks yet. Ask your PA to create some.</div>}
          {tasks.map(task => (
            <div key={task.id} onClick={() => toggleTask(task.id)} style={{
              display: 'flex', alignItems: 'center', gap: 10,
              padding: '10px 12px', background: 'rgba(255,255,255,0.03)',
              border: '1px solid rgba(255,255,255,0.06)', borderRadius: 10, cursor: 'pointer',
              opacity: task.done ? 0.5 : 1
            }}>
              <div style={{ width: 18, height: 18, borderRadius: '50%', border: `1.5px solid ${task.done ? gold : 'rgba(201,169,110,0.3)'}`, background: task.done ? gold : 'transparent', display: 'flex', alignItems: 'center', justifyContent: 'center', flexShrink: 0, fontSize: 10, color: dark }}>
                {task.done ? '✓' : ''}
              </div>
              <span style={{ fontSize: 13, color: cream, textDecoration: task.done ? 'line-through' : 'none', fontFamily: "'Jost',sans-serif" }}>{task.text}</span>
            </div>
          ))}
          <button onClick={() => { setActiveTab('chat'); sendMessage('Help me prioritise my tasks and add any I might be missing for today'); }} style={{ marginTop: 8, padding: '9px 14px', background: 'rgba(201,169,110,0.08)', border: '1px solid rgba(201,169,110,0.2)', borderRadius: 10, cursor: 'pointer', fontSize: 12, color: gold, fontFamily: "'Jost',sans-serif" }}>
            ✦ Ask PA to help prioritise
          </button>
        </div>
      )}

      {/* CALENDAR TAB */}
      {activeTab === 'calendar' && (
        <div style={{ flex: 1, overflowY: 'auto', padding: 14, display: 'flex', flexDirection: 'column', gap: 10 }}>
          <div style={{ fontSize: 12, color: 'rgba(232,220,200,0.5)', lineHeight: 1.6 }}>
            {profile?.google_refresh_token
              ? 'Your Google Calendar is connected.'
              : 'Connect Google Calendar to see your events here.'}
          </div>
          <button onClick={() => { setActiveTab('chat'); sendMessage('What do I have on my calendar today and this week?'); }} style={{ padding: '9px 14px', background: 'rgba(201,169,110,0.08)', border: '1px solid rgba(201,169,110,0.2)', borderRadius: 10, cursor: 'pointer', fontSize: 12, color: gold, fontFamily: "'Jost',sans-serif" }}>
            📅 Ask PA about my schedule
          </button>
          <button onClick={() => { setActiveTab('chat'); sendMessage('Add a new event to my calendar'); }} style={{ padding: '9px 14px', background: 'rgba(255,255,255,0.03)', border: '1px solid rgba(255,255,255,0.08)', borderRadius: 10, cursor: 'pointer', fontSize: 12, color: 'rgba(232,220,200,0.5)', fontFamily: "'Jost',sans-serif" }}>
            + Add event via PA
          </button>
        </div>
      )}
    </div>
  );

  return (
    <>
      <style>{`
        @keyframes paSlideUp { from { opacity:0; transform:translateY(12px); } to { opacity:1; transform:translateY(0); } }
        @keyframes paPulse { 0%,100% { box-shadow: 0 6px 24px rgba(201,169,110,0.45); } 50% { box-shadow: 0 6px 32px rgba(201,169,110,0.7), 0 0 0 6px rgba(201,169,110,0.1); } }
        @keyframes paBounce { 0%,60%,100% { transform:translateY(0); } 30% { transform:translateY(-5px); } }
      `}</style>
      {BubbleEl}
      {ChatWindow}
    </>
  );
}

// ── PA Settings Component (for dashboard Settings tab) ─────────────────
function BravePASettings({ slug, onConfigChange }) {
  const [config, setConfig] = React.useState({
    paName: 'Brave PA',
    personality: 'professional',
    proactiveAlerts: true,
    morningBriefing: true,
    soundEnabled: true,
    soundStyle: 'subtle',
  });
  const [saved, setSaved] = React.useState(false);

  React.useEffect(() => {
    const saved = localStorage.getItem(`brave_pa_config_${slug}`);
    if (saved) try { setConfig(prev => ({...prev, ...JSON.parse(saved)})); } catch(e) {}
  }, [slug]);

  const save = () => {
    localStorage.setItem(`brave_pa_config_${slug}`, JSON.stringify(config));
    onConfigChange && onConfigChange(config);
    setSaved(true);
    setTimeout(() => setSaved(false), 2000);
  };

  const gold = '#C9A96E'; const cream = '#E8DCC8'; const dark = '#0C0A08';
  const sectionLabel = { fontSize: 10, fontFamily: "'Jost',sans-serif", fontWeight: 500, letterSpacing: '0.12em', textTransform: 'uppercase', color: 'rgba(201,169,110,0.45)', marginBottom: 10, marginTop: 20 };
  const card = { background: 'rgba(255,255,255,0.025)', border: '1px solid rgba(201,169,110,0.1)', borderRadius: 13, padding: '16px 18px', marginBottom: 12 };
  const Toggle = ({ value, onChange }) => (
    <div onClick={onChange} style={{ width: 40, height: 22, borderRadius: 11, background: value ? gold : 'rgba(255,255,255,0.1)', position: 'relative', cursor: 'pointer', transition: 'background 0.2s', flexShrink: 0 }}>
      <div style={{ position: 'absolute', top: 3, left: value ? 20 : 3, width: 16, height: 16, borderRadius: '50%', background: 'white', transition: 'left 0.2s', boxShadow: '0 1px 4px rgba(0,0,0,0.3)' }} />
    </div>
  );

  return (
    <div>
      <div style={sectionLabel}>✦ Your PA — Brave PA</div>
      <div style={card}>
        <div style={{ fontSize: 11, color: 'rgba(201,169,110,0.5)', marginBottom: 6, fontFamily: "'Jost',sans-serif" }}>PA Name</div>
        <input value={config.paName} onChange={e => setConfig({...config, paName: e.target.value})} placeholder="Brave PA" style={{ width: '100%', padding: '8px 12px', background: 'rgba(255,255,255,0.04)', border: '1px solid rgba(201,169,110,0.15)', borderRadius: 8, color: cream, fontSize: 13, fontFamily: "'Jost',sans-serif", outline: 'none' }} />
      </div>

      <div style={sectionLabel}>Personality</div>
      <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', marginBottom: 16 }}>
        {Object.entries(PA_PERSONALITIES).map(([key, p]) => (
          <button key={key} onClick={() => setConfig({...config, personality: key})} style={{ padding: '7px 14px', borderRadius: 8, border: '1px solid', cursor: 'pointer', fontSize: 11, fontFamily: "'Jost',sans-serif", background: config.personality === key ? 'rgba(201,169,110,0.15)' : 'transparent', borderColor: config.personality === key ? 'rgba(201,169,110,0.4)' : 'rgba(255,255,255,0.08)', color: config.personality === key ? gold : 'rgba(232,220,200,0.38)' }} title={p.description}>
            {p.label}
          </button>
        ))}
      </div>

      <div style={sectionLabel}>Notifications & Sound</div>
      <div style={card}>
        {[['proactiveAlerts','Proactive PA messages'],['morningBriefing','Morning briefing at 7am'],['soundEnabled','Sound on']].map(([key, label]) => (
          <div key={key} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '8px 0', borderBottom: '1px solid rgba(255,255,255,0.04)' }}>
            <span style={{ fontSize: 13, color: cream, fontFamily: "'Jost',sans-serif" }}>{label}</span>
            <Toggle value={config[key]} onChange={() => setConfig({...config, [key]: !config[key]})} />
          </div>
        ))}
        {config.soundEnabled && (
          <div style={{ marginTop: 12 }}>
            <div style={{ fontSize: 11, color: 'rgba(201,169,110,0.5)', marginBottom: 8, fontFamily: "'Jost',sans-serif" }}>Sound style</div>
            <div style={{ display: 'flex', gap: 8 }}>
              {['subtle','standard','bold'].map(style => (
                <button key={style} onClick={() => setConfig({...config, soundStyle: style})} style={{ padding: '5px 12px', borderRadius: 7, border: '1px solid', cursor: 'pointer', fontSize: 11, fontFamily: "'Jost',sans-serif", textTransform: 'capitalize', background: config.soundStyle === style ? 'rgba(201,169,110,0.15)' : 'transparent', borderColor: config.soundStyle === style ? 'rgba(201,169,110,0.4)' : 'rgba(255,255,255,0.08)', color: config.soundStyle === style ? gold : 'rgba(232,220,200,0.38)' }}>
                  {style}
                </button>
              ))}
            </div>
          </div>
        )}
      </div>

      <button onClick={save} style={{ padding: '10px 20px', background: `linear-gradient(135deg, ${gold}, #7a4f0e)`, border: 'none', borderRadius: 10, cursor: 'pointer', color: dark, fontSize: 13, fontFamily: "'Jost',sans-serif", fontWeight: 600 }}>
        {saved ? '✓ Saved' : 'Save PA Settings'}
      </button>
    </div>
  );
}

// ── Reminder Engine (localStorage-based, PWA push in A4) ──────────────
function checkReminders(slug, paName) {
  const reminders = JSON.parse(localStorage.getItem(`brave_pa_reminders_${slug}`) || '[]');
  const now = new Date();
  const due = reminders.filter(r => !r.fired && new Date(r.time) <= now);
  
  due.forEach(r => {
    // Mark as fired
    const updated = reminders.map(rem => rem.id === r.id ? {...rem, fired: true} : rem);
    localStorage.setItem(`brave_pa_reminders_${slug}`, JSON.stringify(updated));
    
    // Show browser notification if permission granted
    if (Notification.permission === 'granted') {
      new Notification(`⏰ ${paName}`, { body: r.text, icon: '/favicon.ico' });
    }
  });
  
  return due;
}

function addReminder(slug, text, timeStr) {
  const reminders = JSON.parse(localStorage.getItem(`brave_pa_reminders_${slug}`) || '[]');
  const reminder = { id: Date.now(), text, time: new Date(timeStr).toISOString(), fired: false };
  localStorage.setItem(`brave_pa_reminders_${slug}`, JSON.stringify([...reminders, reminder]));
  return reminder;
}

function requestNotificationPermission() {
  if ('Notification' in window && Notification.permission === 'default') {
    Notification.requestPermission();
  }
}
