'use client';
import { useState } from 'react';

const CSS = `
*{box-sizing:border-box}
.privacy-page{background:#0c0a08;color:#e8dcc8;font-family:'Jost',sans-serif;min-height:100vh;padding:40px 16px 80px}
.privacy-page ::-webkit-scrollbar{width:3px}.privacy-page ::-webkit-scrollbar-thumb{background:rgba(201,169,110,0.2);border-radius:2px}
@keyframes privacyFadeUp{from{opacity:0;transform:translateY(10px)}to{opacity:1;transform:translateY(0)}}

.privacy-page .wrap{max-width:680px;margin:0 auto;animation:privacyFadeUp 0.5s ease forwards}

.privacy-page .page-header{text-align:center;margin-bottom:48px;padding-bottom:28px;border-bottom:1px solid rgba(201,169,110,0.12)}
.privacy-page .brand{font-size:10px;letter-spacing:0.2em;text-transform:uppercase;color:rgba(201,169,110,0.35);margin-bottom:16px}
.privacy-page .page-title{font-family:'Cormorant Garamond',serif;font-size:32px;font-weight:300;color:#e8dcc8;margin-bottom:6px}
.privacy-page .page-sub{font-size:12px;color:rgba(232,220,200,0.35);line-height:1.6}
.privacy-page .effective{display:inline-block;margin-top:10px;padding:3px 10px;background:rgba(201,169,110,0.07);border:1px solid rgba(201,169,110,0.15);border-radius:20px;font-size:10px;color:rgba(201,169,110,0.5);letter-spacing:0.05em}

.privacy-page .tabs{display:flex;gap:6px;margin-bottom:32px;background:rgba(255,255,255,0.02);border:1px solid rgba(201,169,110,0.1);border-radius:12px;padding:5px}
.privacy-page .tab{flex:1;padding:9px 8px;border-radius:8px;cursor:pointer;font-size:11px;font-weight:500;letter-spacing:0.05em;text-align:center;border:none;background:transparent;color:rgba(232,220,200,0.35);transition:all 0.2s}
.privacy-page .tab.active{background:rgba(201,169,110,0.12);color:#c9a96e;border:1px solid rgba(201,169,110,0.2)}

.privacy-page .section{display:none;animation:privacyFadeUp 0.3s ease forwards}
.privacy-page .section.active{display:block}

.privacy-page .block{background:rgba(255,255,255,0.02);border:1px solid rgba(201,169,110,0.09);border-radius:14px;padding:22px;margin-bottom:14px}
.privacy-page .block-title{font-size:11px;font-weight:600;letter-spacing:0.1em;text-transform:uppercase;color:rgba(201,169,110,0.5);margin-bottom:14px;display:flex;align-items:center;gap:8px}
.privacy-page .block-title span{font-size:15px}
.privacy-page h3{font-family:'Cormorant Garamond',serif;font-size:20px;font-weight:400;color:#e8dcc8;margin-bottom:8px}
.privacy-page p{font-size:13px;font-weight:300;color:rgba(232,220,200,0.55);line-height:1.75;margin-bottom:10px}
.privacy-page p:last-child{margin-bottom:0}
.privacy-page ul{list-style:none;margin-bottom:10px}
.privacy-page ul li{font-size:13px;font-weight:300;color:rgba(232,220,200,0.55);line-height:1.75;padding-left:16px;position:relative;margin-bottom:4px}
.privacy-page ul li::before{content:'\\2713';position:absolute;left:0;color:rgba(201,169,110,0.5);font-size:11px}
.privacy-page ul li.x::before{content:'\\2717';color:rgba(200,80,80,0.5)}
.privacy-page strong{color:rgba(232,220,200,0.8);font-weight:500}
.privacy-page a{color:rgba(201,169,110,0.6);text-decoration:underline}
.privacy-page a:hover{color:#c9a96e}

.privacy-page .highlight{background:rgba(201,169,110,0.05);border:1px solid rgba(201,169,110,0.15);border-radius:10px;padding:14px 16px;margin-bottom:10px}
.privacy-page .highlight p{color:rgba(232,220,200,0.65);margin-bottom:0}

.privacy-page .badge{display:inline-block;padding:2px 8px;border-radius:20px;font-size:9px;font-weight:600;letter-spacing:0.06em;text-transform:uppercase;margin-left:6px;vertical-align:middle}
.privacy-page .badge.gdpr{background:rgba(100,160,220,0.12);border:1px solid rgba(100,160,220,0.2);color:rgba(100,160,220,0.7)}
.privacy-page .badge.uk{background:rgba(201,169,110,0.1);border:1px solid rgba(201,169,110,0.2);color:rgba(201,169,110,0.6)}

.privacy-page .contact-box{background:rgba(201,169,110,0.04);border:1px dashed rgba(201,169,110,0.18);border-radius:12px;padding:18px;text-align:center;margin-top:20px}
.privacy-page .contact-box p{margin-bottom:6px}
.privacy-page .contact-box a{color:#c9a96e;font-weight:500}

.privacy-page .page-footer{text-align:center;margin-top:40px;padding-top:20px;border-top:1px solid rgba(201,169,110,0.08)}
.privacy-page .page-footer p{font-size:10px;color:rgba(201,169,110,0.22);letter-spacing:0.06em}
.privacy-page .back-btn{display:inline-block;margin-top:16px;padding:8px 18px;background:rgba(201,169,110,0.07);border:1px solid rgba(201,169,110,0.18);border-radius:8px;font-size:11px;color:rgba(201,169,110,0.5);text-decoration:none;transition:all 0.2s}
.privacy-page .back-btn:hover{background:rgba(201,169,110,0.14);color:#c9a96e}
`;

export default function PrivacyPage() {
  const [tab, setTab] = useState('privacy');
  const cls = (id) => `section${tab === id ? ' active' : ''}`;
  const tabCls = (id) => `tab${tab === id ? ' active' : ''}`;

  return (
    <div className="privacy-page">
      <style>{CSS}</style>
      <div className="wrap">

        <div className="page-header">
          <div className="brand">✦ Concierge AI · Brave by Bruno</div>
          <div className="page-title">Privacy & Legal</div>
          <div className="page-sub">We believe in full transparency about how your data is handled.<br/>Read this before using the AI Concierge.</div>
          <div className="effective">Effective date: 1 June 2025 · Last updated: June 2025</div>
        </div>

        <div className="tabs">
          <button className={tabCls('privacy')} onClick={() => setTab('privacy')}>Privacy Policy</button>
          <button className={tabCls('terms')} onClick={() => setTab('terms')}>Terms of Use</button>
          <button className={tabCls('cookies')} onClick={() => setTab('cookies')}>Cookies & Data</button>
          <button className={tabCls('rights')} onClick={() => setTab('rights')}>Your Rights</button>
        </div>

        {/* PRIVACY POLICY */}
        <div className={cls('privacy')}>

          <div className="block">
            <div className="block-title"><span>🛡</span> Overview <span className="badge gdpr">EU GDPR</span><span className="badge uk">UK GDPR</span></div>
            <h3>Who we are</h3>
            <p>This AI Concierge service ("Concierge AI") is operated by <strong>Bruno Aversa</strong>, trading as <strong>Brave by Bruno</strong>, based in London, United Kingdom.</p>
            <p>Contact: <a href="mailto:hello@bravebybruno.com">hello@bravebybruno.com</a></p>
            <p>We are committed to protecting your personal data and complying with the <strong>UK General Data Protection Regulation (UK GDPR)</strong> and the <strong>EU General Data Protection Regulation (EU GDPR 2016/679)</strong>.</p>
          </div>

          <div className="block">
            <div className="block-title"><span>💬</span> What data the AI Concierge processes</div>
            <h3>Conversation data</h3>
            <p>When you use the AI Concierge, your messages are sent to <strong>Anthropic, PBC</strong> (the makers of Claude AI) for processing. This is necessary to generate a response.</p>
            <ul>
              <li>Messages you type are sent to Anthropic's API for processing</li>
              <li>Conversation messages are <strong>also stored on our servers</strong>, so Bruno can review enquiries and follow up</li>
              <li>We do not link conversations to your identity unless you provide your name or email</li>
              <li className="x">We do not sell your data to third parties</li>
              <li className="x">We do not use your data for advertising</li>
              <li className="x">We do not ask for your name, email or payment details in chat</li>
            </ul>
          </div>

          <div className="block">
            <div className="block-title"><span>🤖</span> Anthropic as Data Processor</div>
            <h3>Third-party AI processing</h3>
            <p>Anthropic, PBC (548 Market St, San Francisco, CA 94104, USA) acts as a <strong>data processor</strong> on our behalf when processing your messages through the Claude API.</p>
            <p>Anthropic processes data in accordance with their <a href="https://www.anthropic.com/privacy" target="_blank" rel="noopener noreferrer">Privacy Policy</a> and applicable data protection law. Data transferred to the US is subject to appropriate safeguards under Standard Contractual Clauses (SCCs).</p>
            <div className="highlight">
              <p><strong>Important:</strong> Do not share sensitive personal data (medical history, financial details, government ID numbers) in the chat. The Concierge is designed for service enquiries only.</p>
            </div>
          </div>

          <div className="block">
            <div className="block-title"><span>⚖️</span> Legal basis for processing</div>
            <p>We process your data under the following lawful bases (UK GDPR Article 6):</p>
            <ul>
              <li><strong>Consent</strong> — You explicitly agree before starting the chat</li>
              <li><strong>Legitimate interests</strong> — Providing the AI assistant service you requested</li>
              <li><strong>Contract</strong> — Where processing is necessary to respond to your service enquiry</li>
            </ul>
          </div>

          <div className="block">
            <div className="block-title"><span>📅</span> Data retention</div>
            <p>Conversation messages are <strong>stored on our servers</strong> (via our hosting/database provider) so that Bruno can review enquiries. We do not currently apply an automatic deletion schedule; you may request erasure at any time (see "Your Rights").</p>
            <p>Anthropic may also retain API logs for safety and abuse monitoring purposes in accordance with their own retention policies (typically up to 30 days). Please refer to Anthropic's Privacy Policy for details.</p>
          </div>

          <div className="block">
            <div className="block-title"><span>🔒</span> Security</div>
            <p>All communications between your device and our service, and between our service and Anthropic, are encrypted via <strong>TLS (HTTPS)</strong>. We do not store API keys, conversation logs or user identifiers on client-side storage.</p>
          </div>

          <div className="contact-box">
            <p style={{fontSize:12,color:'rgba(232,220,200,0.5)'}}>Privacy enquiries & data requests</p>
            <a href="mailto:hello@bravebybruno.com">hello@bravebybruno.com</a>
            <p style={{fontSize:10,color:'rgba(232,220,200,0.28)',marginTop:8}}>We will respond to all data requests within 30 days as required by law.</p>
          </div>

        </div>

        {/* TERMS OF USE */}
        <div className={cls('terms')}>

          <div className="block">
            <div className="block-title"><span>📋</span> Terms of Use</div>
            <h3>Agreement to terms</h3>
            <p>By using the Concierge AI service you agree to these Terms of Use. If you do not agree, please do not use the service.</p>
            <p>These terms are governed by the laws of <strong>England and Wales</strong>.</p>
          </div>

          <div className="block">
            <div className="block-title"><span>✦</span> What this service is</div>
            <p>The AI Concierge is an <strong>informational assistant</strong> designed to answer questions about Bruno Aversa's services, availability, pricing and bookings. It:</p>
            <ul>
              <li>Responds on behalf of Bruno Aversa based on information he has provided</li>
              <li>Is available 24/7 to answer general service enquiries</li>
              <li>May route sensitive or complex questions to Bruno for personal review</li>
              <li>Is powered by Claude AI (Anthropic) and may occasionally make errors</li>
            </ul>
          </div>

          <div className="block">
            <div className="block-title"><span>⚠️</span> Limitations & disclaimers</div>
            <ul>
              <li>The AI Concierge provides <strong>general information only</strong> — it does not constitute medical, legal or financial advice</li>
              <li>Responses are generated by AI and may not always be 100% accurate — always confirm bookings directly</li>
              <li>We are not liable for decisions made based solely on AI responses</li>
              <li>The service is provided "as is" without warranty of uninterrupted availability</li>
              <li>Prices, availability and services are subject to change — always confirm with Bruno directly</li>
            </ul>
          </div>

          <div className="block">
            <div className="block-title"><span>🚫</span> Prohibited use</div>
            <p>You agree not to use this service to:</p>
            <ul>
              <li className="x">Attempt to extract confidential business information</li>
              <li className="x">Submit false, misleading or offensive content</li>
              <li className="x">Attempt to circumvent AI safety measures</li>
              <li className="x">Use automated tools or bots to interact with the service</li>
              <li className="x">Reproduce or resell this service without written permission</li>
            </ul>
            <p>We reserve the right to terminate access for users who violate these terms.</p>
          </div>

          <div className="block">
            <div className="block-title"><span>💼</span> Intellectual property</div>
            <p>All content, branding, service descriptions and responses generated by the Concierge AI on behalf of Bruno Aversa are the intellectual property of <strong>Brave by Bruno</strong>. You may not reproduce, copy or redistribute them without written consent.</p>
            <p>The underlying AI technology is owned by Anthropic, PBC and used under licence.</p>
          </div>

          <div className="block">
            <div className="block-title"><span>📝</span> Changes to these terms</div>
            <p>We may update these Terms of Use from time to time. Material changes will be communicated by updating the "Last updated" date at the top of this page. Continued use of the service constitutes acceptance of the updated terms.</p>
          </div>

        </div>

        {/* COOKIES & DATA */}
        <div className={cls('cookies')}>

          <div className="block">
            <div className="block-title"><span>🍪</span> Cookie policy</div>
            <h3>We keep it minimal</h3>
            <p>The Concierge AI is designed with a <strong>minimal data footprint</strong>. We do not use advertising cookies, tracking pixels or third-party analytics.</p>
          </div>

          <div className="block">
            <div className="block-title"><span>💾</span> What we store locally</div>
            <p>We use <strong>sessionStorage</strong> (your browser's temporary memory) only to:</p>
            <ul>
              <li>Remember that you have accepted this privacy notice (so you don't see it every message)</li>
              <li>This data is automatically deleted when you close the browser tab</li>
              <li>It never leaves your device and is never sent to our servers</li>
            </ul>
            <div className="highlight">
              <p>We use <strong>no persistent cookies</strong>. Nothing is stored after you close the tab. Nothing is shared with advertisers. Nothing is used to track you across websites.</p>
            </div>
          </div>

          <div className="block">
            <div className="block-title"><span>🔗</span> Third-party links</div>
            <p>Our page contains links to third-party services including:</p>
            <ul>
              <li><strong>Fresha</strong> — for booking management (their own privacy policy applies)</li>
              <li><strong>Instagram</strong> — social media links (Meta's privacy policy applies)</li>
              <li><strong>Anthropic</strong> — AI processing (their privacy policy applies)</li>
            </ul>
            <p>We are not responsible for the privacy practices of these third parties.</p>
          </div>

          <div className="block">
            <div className="block-title"><span>📊</span> Analytics</div>
            <p>We currently use <strong>no third-party analytics</strong> (no Google Analytics, no Meta Pixel, no tracking scripts). Basic server-side request logs may be maintained by our hosting provider (Vercel) for security purposes — please refer to <a href="https://vercel.com/legal/privacy-policy" target="_blank" rel="noopener noreferrer">Vercel's Privacy Policy</a>.</p>
          </div>

        </div>

        {/* YOUR RIGHTS */}
        <div className={cls('rights')}>

          <div className="block">
            <div className="block-title"><span>🏛</span> Your rights under UK & EU GDPR</div>
            <p>You have the following rights regarding your personal data. To exercise any of these rights, contact us at <a href="mailto:hello@bravebybruno.com">hello@bravebybruno.com</a>.</p>
          </div>

          <div className="block">
            <div className="block-title"><span>1</span> Right to access</div>
            <p>You have the right to request a copy of any personal data we hold about you. This may include your chat messages and your consent record, if you have used the AI Concierge.</p>
          </div>

          <div className="block">
            <div className="block-title"><span>2</span> Right to erasure ("Right to be forgotten")</div>
            <p>You have the right to request deletion of your personal data. As conversations are session-only and not stored, closing the chat window effectively erases the conversation. For any other data held, contact us and we will action your request within 30 days.</p>
          </div>

          <div className="block">
            <div className="block-title"><span>3</span> Right to rectification</div>
            <p>You have the right to request correction of any inaccurate personal data we hold about you.</p>
          </div>

          <div className="block">
            <div className="block-title"><span>4</span> Right to withdraw consent</div>
            <p>You may withdraw your consent to data processing at any time by simply closing the chat and clearing your browser session storage. Withdrawal of consent does not affect the lawfulness of processing before withdrawal.</p>
          </div>

          <div className="block">
            <div className="block-title"><span>5</span> Right to complain</div>
            <p>If you believe your data rights have been violated, you have the right to lodge a complaint with:</p>
            <ul>
              <li><strong>UK:</strong> Information Commissioner's Office (ICO) — <a href="https://ico.org.uk" target="_blank" rel="noopener noreferrer">ico.org.uk</a></li>
              <li><strong>EU:</strong> Your local Data Protection Authority (DPA)</li>
            </ul>
            <p>We would always prefer you contact us first so we can resolve any issues directly.</p>
          </div>

          <div className="block">
            <div className="block-title"><span>6</span> Right to data portability</div>
            <p>You have the right to receive your personal data in a structured, machine-readable format. As we hold minimal data, this right is limited in practice — but we will fulfil any reasonable request.</p>
          </div>

          <div className="contact-box">
            <p style={{fontSize:13,fontFamily:"'Cormorant Garamond',serif",fontStyle:'italic',color:'rgba(232,220,200,0.5)',marginBottom:10}}>To exercise any of your rights, get in touch:</p>
            <a href="mailto:hello@bravebybruno.com" style={{fontSize:15}}>hello@bravebybruno.com</a>
            <p style={{fontSize:10,color:'rgba(232,220,200,0.28)',marginTop:10}}>We respond to all data requests within <strong style={{color:'rgba(232,220,200,0.4)'}}>30 days</strong> as required by UK GDPR Article 12.</p>
          </div>

        </div>

        <div className="page-footer">
          <p>© 2025 Bruno Aversa · Brave by Bruno · London, UK</p>
          <p style={{marginTop:4}}>Concierge AI is powered by Claude (Anthropic) · Compliant with UK GDPR & EU GDPR 2016/679</p>
          <a href="/" className="back-btn">← Back to Bruno's page</a>
        </div>

      </div>
    </div>
  );
}
