// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Concierge AI — Secure Backend
// API key never exposed to frontend
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

const RATE_LIMIT = new Map();
const MAX_TURNS = 20;
const MAX_TOKENS = 1000;
const WINDOW_MS = 60 * 60 * 1000;
const MAX_PER_HOUR = 60;

// In-memory analytics (resets on redeploy — use DB for production)
const ANALYTICS = {
  conversations: 0,
  messages: 0,
  byProfile: {},
  topTopics: {},
  blocked: 0,
  bookingIntents: 0,
  hourly: [],
};

function trackAnalytics(profileId, userMessage, wasBlocked) {
  ANALYTICS.messages++;
  ANALYTICS.byProfile[profileId] = (ANALYTICS.byProfile[profileId] || 0) + 1;
  if (wasBlocked) ANALYTICS.blocked++;

  // Detect booking intent
  const lower = (userMessage || "").toLowerCase();
  if (["book", "prenot", "appointment", "slot", "available", "when", "schedule"].some(w => lower.includes(w))) {
    ANALYTICS.bookingIntents++;
  }

  // Extract topic keywords
  const topics = ["price", "prezzo", "pain", "dolore", "injury", "first time", "prima volta", "menu", "allerg", "legal", "consultation", "class", "lezione", "photo", "content"];
  topics.forEach(t => {
    if (lower.includes(t)) ANALYTICS.topTopics[t] = (ANALYTICS.topTopics[t] || 0) + 1;
  });

  // Hourly tracking
  const hour = new Date().toISOString().slice(0, 13);
  const existing = ANALYTICS.hourly.find(h => h.hour === hour);
  if (existing) existing.count++;
  else ANALYTICS.hourly.push({ hour, count: 1 });
  if (ANALYTICS.hourly.length > 48) ANALYTICS.hourly.shift();
}

function getRateLimit(ip) {
  const now = Date.now();
  const entry = RATE_LIMIT.get(ip);
  if (!entry || now > entry.resetAt) {
    RATE_LIMIT.set(ip, { count: 1, resetAt: now + WINDOW_MS });
    return true;
  }
  if (entry.count >= MAX_PER_HOUR) return false;
  entry.count++;
  return true;
}

export default async function handler(req, res) {
  res.setHeader("Access-Control-Allow-Origin", "*");
  res.setHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS");
  res.setHeader("Access-Control-Allow-Headers", "Content-Type, X-Owner-Key");
  if (req.method === "OPTIONS") return res.status(200).end();

  // Analytics endpoint for owner dashboard
  if (req.method === "GET" && req.url?.includes("/api/chat")) {
    const ownerKey = req.headers["x-owner-key"];
    if (ownerKey !== process.env.OWNER_KEY) return res.status(401).json({ error: "Unauthorized" });
    return res.status(200).json(ANALYTICS);
  }

  if (req.method !== "POST") return res.status(405).json({ error: "Method not allowed" });

  const ip = req.headers["x-forwarded-for"]?.split(",")[0] || "unknown";
  if (!getRateLimit(ip)) return res.status(429).json({ error: "Too many requests. Please try again later." });

  try {
    const { messages, systemPrompt, turnCount, profileId, wasBlocked } = req.body;
    if (!messages || !Array.isArray(messages)) return res.status(400).json({ error: "Invalid request" });

    const userMessage = messages[messages.length - 1]?.content || "";
    trackAnalytics(profileId || "unknown", userMessage, wasBlocked);

    if (turnCount > MAX_TURNS) {
      return res.status(200).json({
        content: [{ text: "We've reached the end of this demo session. To continue the conversation or book directly, please use the link in the profile. Thank you!" }]
      });
    }

    const sanitized = messages.map(m => ({
      role: m.role === "assistant" ? "assistant" : "user",
      content: String(m.content).slice(0, 2000)
    }));

    const response = await fetch("https://api.anthropic.com/v1/messages", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "x-api-key": process.env.ANTHROPIC_API_KEY,
        "anthropic-version": "2023-06-01",
      },
      body: JSON.stringify({
        model: "claude-sonnet-4-20250514",
        max_tokens: MAX_TOKENS,
        system: systemPrompt,
        messages: sanitized,
      }),
    });

    if (!response.ok) {
      console.error("Anthropic error:", await response.text());
      return res.status(502).json({ error: "AI service temporarily unavailable. Please try again." });
    }

    const data = await response.json();
    return res.status(200).json(data);

  } catch (error) {
    console.error("Handler error:", error);
    return res.status(500).json({ error: "Internal server error" });
  }
}
