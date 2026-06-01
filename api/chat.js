const RATE_LIMIT = new Map();
const MAX_TURNS = 20;
const MAX_TOKENS = 1000;
const WINDOW_MS = 60 * 60 * 1000;
const MAX_PER_HOUR = 60;

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

module.exports = async function handler(req, res) {
  res.setHeader("Access-Control-Allow-Origin", "*");
  res.setHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS");
  res.setHeader("Access-Control-Allow-Headers", "Content-Type, X-Owner-Key");
  if (req.method === "OPTIONS") return res.status(200).end();

  const ip = req.headers["x-forwarded-for"]?.split(",")[0] || "unknown";
  if (!getRateLimit(ip)) return res.status(429).json({ error: "Too many requests." });

  if (req.method !== "POST") return res.status(405).json({ error: "Method not allowed" });

  try {
    const { messages, systemPrompt, turnCount } = req.body;
    if (!messages || !Array.isArray(messages)) return res.status(400).json({ error: "Invalid request" });

    if (turnCount > MAX_TURNS) {
      return res.status(200).json({
        content: [{ text: "We've reached the end of this demo session. Please book directly or get in touch — thank you!" }]
      });
    }

    const sanitized = messages.map(m => ({
      role: m.role === "assistant" ? "assistant" : "user",
      content: String(m.content).slice(0, 2000)
    }));

    const apiKey = process.env.ANTHROPIC_API_KEY;
    console.log("API key present:", !!apiKey, "Length:", apiKey ? apiKey.length : 0);

    const response = await fetch("https://api.anthropic.com/v1/messages", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "x-api-key": apiKey,
        "anthropic-version": "2023-06-01",
      },
      body: JSON.stringify({
        model: "claude-haiku-4-5-20251001",
        max_tokens: MAX_TOKENS,
        system: systemPrompt,
        messages: sanitized,
      }),
    });

    const responseText = await response.text();
    console.log("Anthropic status:", response.status);
    console.log("Anthropic response:", responseText.slice(0, 200));

    if (!response.ok) {
      return res.status(502).json({ error: "AI service temporarily unavailable." });
    }

    const data = JSON.parse(responseText);
    return res.status(200).json(data);

  } catch (error) {
    console.error("Handler error:", error.message);
    return res.status(500).json({ error: "Internal server error" });
  }
};
