# 🔗 Third-Party Application Integration Guide

## Unified Log & Activity Monitor (ULAM) - External Integration Architecture

> Dokumen ini menjelaskan bagaimana ULAM dapat mengintegrasikan dan memonitor aktivitas dari berbagai aplikasi third-party seperti Google Workspace, GitHub, Instagram, Twitter, dan lainnya.

---

## 📋 Daftar Isi

1. [Konsep Integrasi](#konsep-integrasi)
2. [Aplikasi yang Didukung](#aplikasi-yang-didukung)
3. [Data yang Bisa Di-Track](#data-yang-bisa-di-track)
4. [Cara Integrasi per Aplikasi](#cara-integrasi-per-aplikasi)
5. [Kasus Penggunaan](#kasus-penggunaan)
6. [Keamanan & Privasi](#keamanan--privasi)
7. [Roadmap Integrasi](#roadmap-integrasi)

---

## 🎯 Konsep Integrasi

### Bagaimana ULAM Bekerja sebagai Aggregator

```
┌─────────────────────────────────────────────────────────────┐
│                     ULAM DASHBOARD                          │
│  ┌──────────────┬──────────────┬──────────────┐            │
│  │  Google      │  GitHub      │  Instagram   │            │
│  │  Workspace   │  Enterprise  │  Business    │            │
│  └──────┬───────┴──────┬───────┴──────┬───────┘            │
│         │              │              │                      │
│  ┌──────▼───────┬──────▼───────┬──────▼───────┐            │
│  │  Twitter/X   │  Slack       │  Dropbox     │            │
│  │  API         │  Enterprise  │  Business    │            │
│  └──────┬───────┴──────┬───────┴──────┬───────┘            │
│         │              │              │                      │
└─────────┼──────────────┼──────────────┼──────────────────────┘
          │              │              │
          ▼              ▼              ▼
    ┌─────────┐    ┌─────────┐    ┌─────────┐
    │ Connect │    │ Connect │    │ Connect │
    │  ors    │    │  ors    │    │  ors    │
    └────┬────┘    └────┬────┘    └────┬────┘
         │              │              │
         └──────────────┴──────────────┘
                         │
                    ┌────▼────┐
                    │  ULAM   │
                    │  CORE   │
                    └─────────┘
```

**Prinsip Kerja:**
1. **Connector** - Background worker yang pull data dari API third-party
2. **Normalizer** - Transform data ke format ULAM standard
3. **Aggregator** - Gabungkan semua data di satu dashboard
4. **Analyzer** - Deteksi pola mencurigakan cross-platform
5. **Alerter** - Kirim notifikasi real-time

---

## 📱 Aplikasi yang Didukung

### Tier 1: Enterprise/Business (High Priority)

| Aplikasi | Tipe Data | Status | Keterangan |
|----------|-----------|--------|------------|
| **Google Workspace** | Login, Drive, Gmail, Calendar, Docs | 🔴 Planned | Admin SDK Reports API |
| **GitHub Enterprise** | Repo access, SSH keys, OAuth apps | 🔴 Planned | Audit Log API |
| **Slack Enterprise** | Channel access, file sharing, DMs | 🔴 Planned | Audit Logs API |
| **Microsoft 365** | Azure AD, SharePoint, Teams | 🔴 Planned | Microsoft Graph API |
| **AWS** | IAM, S3, EC2 activities | 🔴 Planned | CloudTrail + CloudWatch |

### Tier 2: Social Media (Medium Priority)

| Aplikasi | Tipe Data | Status | Keterangan |
|----------|-----------|--------|------------|
| **Instagram Business** | Login, post, DM, story | 🔴 Planned | Instagram Graph API |
| **Twitter/X** | Login, tweet, DM, followers | 🔴 Planned | Twitter API v2 |
| **LinkedIn** | Login, post, connection | 🔴 Planned | LinkedIn API |
| **Facebook** | Login, page admin, groups | 🔴 Planned | Facebook Graph API |

### Tier 3: Productivity Tools (Medium Priority)

| Aplikasi | Tipe Data | Status | Keterangan |
|----------|-----------|--------|------------|
| **Notion** | Page access, edits, shares | 🔴 Planned | Notion API |
| **Figma** | File access, comments, edits | 🔴 Planned | Figma API |
| **Dropbox Business** | File access, share, delete | 🔴 Planned | Dropbox Business API |
| **Zoom** | Meeting access, recordings | 🔴 Planned | Zoom API |

### Tier 4: Development Tools (Low Priority)

| Aplikasi | Tipe Data | Status | Keterangan |
|----------|-----------|--------|------------|
| **Jira** | Issue access, project changes | 🔴 Planned | Jira REST API |
| **GitLab** | Repo access, merge requests | 🔴 Planned | GitLab API |
| **Vercel** | Deployment, team access | 🔴 Planned | Vercel API |
| **Docker Hub** | Image pull/push, team access | 🔴 Planned | Docker Hub API |

---

## 📊 Data yang Bisa Di-Track

### 1. Authentication Events

```json
{
  "category": "AUTH_EVENT",
  "level": "INFO",
  "message": "User login detected",
  "context": {
    "platform": "google_workspace",
    "user_email": "john.doe@company.com",
    "ip_address": "203.0.113.45",
    "location": "Jakarta, Indonesia",
    "device": "iPhone 14 Pro",
    "browser": "Chrome 120.0",
    "timestamp": "2026-03-19T14:30:00Z",
    "auth_method": "oauth2",
    "is_suspicious": false
  }
}
```

**Yang bisa dideteksi:**
- ✅ Login dari multiple device berbeda
- ✅ Login dari lokasi berbeda dalam waktu singkat (impossible travel)
- ✅ Login di luar jam kerja
- ✅ Failed login attempts (brute force)
- ✅ Login dari IP yang tidak dikenal
- ✅ Password changes
- ✅ 2FA enabled/disabled
- ✅ Recovery email/phone changes

### 2. File/Asset Access

```json
{
  "category": "FILE_ACCESS",
  "level": "WARN",
  "message": "Sensitive file accessed",
  "context": {
    "platform": "google_drive",
    "user_email": "john.doe@company.com",
    "file_name": "Financial_2026_Q1.pdf",
    "file_id": "1ABC123xyz",
    "action": "DOWNLOAD",
    "ip_address": "203.0.113.45",
    "timestamp": "2026-03-19T15:45:00Z",
    "is_shared_externally": true,
    "sensitivity_level": "HIGH"
  }
}
```

**Yang bisa dideteksi:**
- ✅ File download/upload
- ✅ File share (internal/external)
- ✅ File delete/permanent delete
- ✅ File edit/version changes
- ✅ Bulk download (data exfiltration)
- ✅ Access ke file sensitif
- ✅ Share link creation

### 3. Repository/Code Access (GitHub/GitLab)

```json
{
  "category": "REPO_ACCESS",
  "level": "CRITICAL",
  "message": "SSH key added to repository",
  "context": {
    "platform": "github",
    "user": "john.doe",
    "repo": "company/payment-gateway",
    "action": "SSH_KEY_ADDED",
    "ip_address": "203.0.113.45",
    "timestamp": "2026-03-19T16:00:00Z",
    "ssh_key_fingerprint": "SHA256:abc123...",
    "is_suspicious": true,
    "reason": "New SSH key from unknown device"
  }
}
```

**Yang bisa dideteksi:**
- ✅ Repo access (clone, pull, push)
- ✅ SSH key additions
- ✅ OAuth app installations
- ✅ Personal access token creation
- ✅ Branch protection changes
- ✅ Admin privilege escalation
- ✅ Secret/credential exposure

### 4. Social Media Activities

```json
{
  "category": "SOCIAL_ACTIVITY",
  "level": "INFO",
  "message": "Post published on company account",
  "context": {
    "platform": "twitter",
    "account": "@company_official",
    "user_email": "social@company.com",
    "action": "POST_CREATED",
    "content_hash": "sha256:xyz789...",
    "ip_address": "203.0.113.45",
    "timestamp": "2026-03-19T10:00:00Z",
    "engagement": {
      "likes": 150,
      "retweets": 45,
      "replies": 23
    }
  }
}
```

**Yang bisa dideteksi:**
- ✅ Post/tweet published
- ✅ Post deleted
- ✅ Direct messages
- ✅ Follower/following changes
- ✅ Login dari device baru
- ✅ Password/email changes
- ✅ Privacy setting changes

---

## 🔧 Cara Integrasi per Aplikasi

### 1. Google Workspace

#### Setup Steps:
1. **Create Service Account** di Google Cloud Console
2. **Enable Admin SDK API**
3. **Domain-wide Delegation** untuk service account
4. **Add scopes:**
   - `https://www.googleapis.com/auth/admin.reports.audit.readonly`
   - `https://www.googleapis.com/auth/admin.reports.usage.readonly`
5. **Download JSON key file**

#### Data yang Diambil:
```go
// ULAM Connector
func SyncGoogleWorkspace() {
    // Reports API
    activities := reportsService.Activities.List("all", "login")
    
    // Events yang di-track:
    // - login_success
    // - login_failure  
    // - 2sv_disable
    // - password_change
    // - email_forwarding_enabled
    // - calendar_event_deleted
    // - drive_file_download
    // - gmail_forwarding_enabled
    // - admin_user_deleted
}
```

#### Contoh Alert:
```
🚨 SUSPICIOUS ACTIVITY DETECTED

User: john.doe@company.com
Platform: Google Workspace
Activity: Login from 3 different locations in 1 hour

Location 1: Jakarta, Indonesia (14:00 WIB)
Location 2: Singapore (14:30 WIB)  
Location 3: New York, USA (14:45 WIB)

⚠️ Impossible travel pattern detected!
Account might be compromised or shared.
```

### 2. GitHub Enterprise

#### Setup Steps:
1. **Enable Audit Log** di Organization Settings
2. **Generate Personal Access Token** dengan scope `admin:org`
3. **Enable Audit Log Streaming** (opsional untuk real-time)

#### Data yang Diambil:
```go
// ULAM Connector
func SyncGitHubAudit() {
    // Audit Log API
    events := githubClient.Orgs.GetAuditLog("your-org")
    
    // Events yang di-track:
    // - repo.access
    // - repo.create
    // - repo.destroy
    // - team.add_member
    // - team.remove_member
    // - ssh_key.create
    // - oauth_authorization.create
    // - protected_branch.policy_override
}
```

#### Contoh Alert:
```
🚨 CRITICAL: Repository Access Anomaly

User: developer-john
Repository: company/payment-gateway
Activity: SSH key added at 03:00 AM

Details:
- SSH Key: SHA256:abc123...
- Device: Unknown
- IP: 185.220.101.45 (Tor Exit Node)
- Time: 03:15 AM (Outside working hours)

🔒 Recommendation: Immediately review and revoke if unauthorized!
```

### 3. Instagram Business

#### Setup Steps:
1. **Create Facebook App** di Meta for Developers
2. **Add Instagram Graph API** product
3. **Request permissions:**
   - `instagram_basic`
   - `instagram_manage_insights`
   - `pages_read_engagement`
4. **Generate Access Token** untuk Business Account

#### Data yang Diambil:
```javascript
// ULAM Connector
async function syncInstagram() {
  const insights = await fetch(`https://graph.instagram.com/me/insights?...`)
  
  // Events yang di-track:
  // - media_publish
  // - media_delete
  // - story_publish
  // - comment_create
  // - mention
  // - profile_update
}
```

#### Contoh Alert:
```
⚠️ Instagram Activity Alert

Account: @company_official
User: social@company.com
Activity: Post deleted

Details:
- Post ID: 123456789
- Posted: 2026-03-19 09:00 AM
- Deleted: 2026-03-19 15:30 PM (6.5 hours later)
- Reason: Unknown

Note: Post was gaining traction (150 likes, 45 shares)
🔍 Review if deletion was intentional or compromised account.
```

### 4. Twitter/X API

#### Setup Steps:
1. **Apply for Developer Account** di Twitter Developer Portal
2. **Create Project & App**
3. **Generate Bearer Token**
4. **Elevated Access** untuk analytics data

#### Data yang Diambil:
```python
# ULAM Connector
import tweepy

def sync_twitter():
    client = tweepy.Client(bearer_token="YOUR_TOKEN")
    
    # Events yang di-track:
    # - tweet.create
    # - tweet.delete
    # - dm.send
    # - follow.create
    # - block.create
    # - password_change
    # - email_change
```

#### Contoh Alert:
```
🚨 Twitter Security Alert

Account: @company_official
Activity: Email address changed

Details:
- Old Email: social@company.com
- New Email: social-backup@gmail.com
- Changed by: admin@company.com
- Time: 2026-03-19 11:00 AM

⚠️ This change bypassed normal approval workflow!
🔒 Recommendation: Verify with admin team.
```

### 5. Slack Enterprise

#### Setup Steps:
1. **Enable Audit Logs** di Slack Admin Dashboard
2. **Create API Token** dengan scope `auditlogs:read`
3. **Configure Log Retention** (90 days default)

#### Data yang Diambil:
```go
// ULAM Connector
func SyncSlackAudit() {
    // Audit Logs API
    logs := slackClient.GetAuditLogs(
        slack.AuditLogFilter{
            Oldest: time.Now().Add(-1 * time.Hour),
        },
    )
    
    // Events yang di-track:
    // - user_login
    // - user_logout
    // - file_download
    // - file_shared
    // - channel_created
    // - channel_deleted
    // - user_channel_join
    // - app_installed
}
```

#### Contoh Alert:
```
⚠️ Slack Data Exfiltration Risk

User: developer-john@company.com
Activity: Bulk file download detected

Details:
- Files Downloaded: 150 files
- Total Size: 2.3 GB
- Time Window: 15 minutes
- Sensitive Files: 12 (containing "confidential" or "password")

🚨 This exceeds normal usage pattern!
🔍 Possible data exfiltration attempt.
```

---

## 💼 Kasus Penggunaan

### Kasus 1: Deteksi Account Sharing

**Skenario:** 1 akun Google Workspace dibagi 10 orang

**ULAM Detection:**
```
📊 ANOMALY DETECTED: Account Sharing

Account: admin@company.com
Detection Method: Multiple concurrent logins

Timeline:
├─ 08:00 - Login from Jakarta (IP: 203.0.113.10) ✅
├─ 08:05 - Login from Surabaya (IP: 203.0.113.20) ⚠️
├─ 08:07 - Login from Bandung (IP: 203.0.113.30) ⚠️
├─ 08:10 - Login from Singapore (IP: 203.0.113.40) 🚨
└─ 08:15 - Login from Malaysia (IP: 203.0.113.50) 🚨

Analysis:
- 5 logins dalam 15 menit
- 3 lokasi berbeda (impossible travel)
- Device berbeda: Chrome, Firefox, Safari, Mobile App, Tablet

🚨 HIGH RISK: Account sharing atau compromise detected!

Recommendation:
1. Disable account immediately
2. Force password reset
3. Enable mandatory 2FA
4. Audit all recent activities
5. Create individual accounts for each user
```

### Kasus 2: Insider Threat Detection

**Skenario:** Karyawan resign download semua file sensitif

**ULAM Detection:**
```
🚨 CRITICAL: Potential Data Theft

User: john.resign@company.com
Status: Resignation submitted (Last day: 2026-03-31)

Activity Pattern (Last 7 Days):
├─ Google Drive: Downloaded 450 files (normally 10/day)
├─ GitHub: Cloned 12 private repositories
├─ Slack: Exported all channel history
├─ Email: Forwarded 200 emails to personal Gmail
└─ Time: Mostly after 6 PM and weekends

Files Accessed:
🔴 Financial_Forecast_2026.xlsx
🔴 Client_List_Confidential.csv
🔴 Product_Roadmap_Q3.pdf
🔴 Employee_Salaries_2026.pdf
🔴 Acquisition_Plan_Draft.docx

🚨 THREAT LEVEL: CRITICAL
📊 Confidence: 95%

Immediate Actions:
1. Disable all access immediately
2. Revoke OAuth tokens
3. Force logout from all devices
4. Audit file downloads
5. Check for data already exfiltrated
6. Legal review for NDA violation
```

### Kasus 3: Compromised Account

**Skenario:** Akun Instagram business di-hack

**ULAM Detection:**
```
🚨 SECURITY BREACH: Account Compromise

Account: @company_official (500K followers)
Platform: Instagram

Anomalous Activities:
├─ 02:30 AM - Login from Russia (IP: 91.203.164.x)
├─ 02:35 AM - Email changed to: recover.xyz@gmail.com
├─ 02:37 AM - 2FA disabled
├─ 02:40 AM - Password changed
├─ 02:45 AM - Post published: "🚀 Crypto investment opportunity..."
└─ 02:50 AM - DM sent to 1000 followers with phishing link

Impacts:
- Reputational damage: HIGH
- Follower loss: -5K (and growing)
- Legal risk: Phishing content posted
- Revenue impact: Estimated $50K loss

🔒 Emergency Response:
1. Contact Instagram support immediately
2. Freeze account
3. Post disclaimer from backup account
4. Notify legal and PR team
5. Forensic investigation
6. Customer notification

Recovery Time: 24-72 hours
```

---

## 🔒 Keamanan & Privasi

### Data yang Tidak Di-track (Privacy-First)

❌ **Tidak pernah di-track:**
- Passwords
- Credit card numbers
- Private messages content (hanya metadata)
- Personal photos/videos
- Medical records
- Banking details

✅ **Hanya di-track:**
- Login events (who, when, from where)
- File access metadata (nama, ukuran, waktu, bukan isi)
- Permission changes
- Security settings changes
- Admin activities

### Compliance

ULAM integrasi mematuhi:
- ✅ **GDPR** - Right to be forgotten, data minimization
- ✅ **SOC 2** - Security controls
- ✅ **ISO 27001** - Information security
- ✅ **CCPA** - California privacy rights
- ✅ **HIPAA** - Healthcare data (jika applicable)

---

## 📅 Roadmap Integrasi

### Q2 2026
- [ ] Google Workspace Connector
- [ ] GitHub Enterprise Connector
- [ ] Slack Enterprise Connector

### Q3 2026
- [ ] Microsoft 365 Connector
- [ ] AWS CloudTrail Connector
- [ ] Instagram Business Connector

### Q4 2026
- [ ] Twitter/X API Connector
- [ ] LinkedIn Connector
- [ ] Notion Connector

### Q1 2027
- [ ] Dropbox Business Connector
- [ ] Zoom Connector
- [ ] Jira Connector

---

## 🛠️ Technical Requirements

### API Keys Needed per Platform

```yaml
Google Workspace:
  - Service Account JSON
  - Admin SDK enabled
  - Domain-wide delegation

GitHub Enterprise:
  - Personal Access Token (PAT)
  - Organization admin access
  - Audit log access enabled

Slack Enterprise:
  - OAuth Token
  - Audit logs:read scope
  - Enterprise Grid subscription

Instagram Business:
  - Facebook App ID
  - Instagram Business Account
  - Graph API access token

Twitter/X:
  - API Key & Secret
  - Bearer Token
  - Elevated access level
```

### Storage Requirements

| Platform | Events/Day | Storage/Day | Retention |
|----------|------------|-------------|-----------|
| Google Workspace | ~10,000 | ~50 MB | 90 days |
| GitHub Enterprise | ~5,000 | ~25 MB | 90 days |
| Slack Enterprise | ~50,000 | ~200 MB | 30 days |
| Instagram Business | ~1,000 | ~5 MB | 30 days |
| **Total** | **~66,000** | **~280 MB** | **-** |

**Rekomendasi:** 100GB storage untuk 1 tahun data

---

## 📞 Support & Troubleshooting

### Common Issues

**1. API Rate Limits**
```
Error: 429 Too Many Requests
Solution: Implement exponential backoff
```

**2. Token Expiration**
```
Error: 401 Unauthorized
Solution: Auto-refresh OAuth tokens
```

**3. Missing Permissions**
```
Error: 403 Forbidden
Solution: Check OAuth scopes
```

### Contact

- 📧 Email: support@ulam.io
- 💬 Slack: #ulam-support
- 📚 Docs: https://docs.ulam.io

---

## 📝 Summary

ULAM dapat menjadi **Single Pane of Glass** untuk monitoring aktivitas di SEMUA aplikasi:

✅ **100+ aplikasi** bisa di-integrasikan  
✅ **Real-time detection** anomali dan threat  
✅ **Cross-platform correlation** (connect the dots)  
✅ **Immutable audit trail** untuk compliance  
✅ **AI-powered analysis** untuk deteksi pola canggih  

**Value Proposition:**
- 🔐 **Security:** Deteksi account compromise, insider threats, data exfiltration
- 📊 **Compliance:** Audit trail untuk GDPR, SOC2, ISO 27001
- 🚨 **Alerting:** Real-time notification via Email, Telegram, Slack
- 🤖 **Automation:** Auto-respon untuk security incidents

---

**Dokumen ini akan di-update seiring dengan rilis connector baru.**

*Last Updated: 2026-03-19*  
*Version: 1.0*  
*Author: ULAM Team*
