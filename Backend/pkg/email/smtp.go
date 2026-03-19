package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/petrushandika/one-log/internal/domain"
)

// SMTPEmailService handles actual email sending via SMTP.
type SMTPEmailService struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewSMTPEmailService() *SMTPEmailService {
	return &SMTPEmailService{
		Host:     os.Getenv("MAIL_HOST"),
		Port:     os.Getenv("MAIL_PORT"),
		Username: os.Getenv("MAIL_USER"),
		Password: os.Getenv("MAIL_PASSWORD"),
		From:     os.Getenv("MAIL_FROM"),
	}
}

// SendAlertEmail sends an HTML email via SMTP for CRITICAL/ERROR logs
func (s *SMTPEmailService) SendAlertEmail(to string, logEntry *domain.LogEntry) error {
	// If SMTP is not properly configured, log a warning and return early in a graceful way
	if s.Host == "" || s.Username == "" || s.Password == "" {
		log.Println("WARNING: SMTP credentials not fully configured. Email skipped.")
		return nil
	}

	// Prepare Subject
	subject := fmt.Sprintf("Alert: [%s] %s on %s", logEntry.Level, logEntry.Category, logEntry.SourceID)

	// Prepare Template
	tmpl := `
<!DOCTYPE html>
<html>
<body>
	<h2 style="color: #d9534f;">ULAM Alert: {{.Level}} on {{.SourceID}}</h2>
	<p><strong>Category:</strong> {{.Category}}</p>
	<p><strong>Message:</strong> {{.Message}}</p>
	<p><strong>Timestamp:</strong> {{.CreatedAt}}</p>
	{{if .IPAddress}}<p><strong>IP Address:</strong> {{.IPAddress}}</p>{{end}}
	<hr/>
	{{if .StackTrace}}
	<h4>Stack Trace</h4>
	<pre style="background: #f8f9fa; padding: 10px; border-radius: 5px;">{{.StackTrace}}</pre>
	{{end}}
	<p><a href="https://ulam.your-domain.com">View Dashboard</a></p>
</body>
</html>
`
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	// MIME headers
	body.Write([]byte(fmt.Sprintf("To: %s\r\n", to)))
	body.Write([]byte(fmt.Sprintf("From: %s\r\n", s.From)))
	body.Write([]byte(fmt.Sprintf("Subject: %s\r\n", subject)))
	body.Write([]byte("MIME-Version: 1.0\r\n"))
	body.Write([]byte("Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n"))

	err = t.Execute(&body, logEntry)
	if err != nil {
		return err
	}

	return s.sendEmail(to, body.Bytes())
}

// SendRecoveryEmail sends an HTML email when a source comes back online
func (s *SMTPEmailService) SendRecoveryEmail(to string, sourceName string, downtimeDuration string) error {
	if s.Host == "" || s.Username == "" || s.Password == "" {
		log.Println("WARNING: SMTP credentials not fully configured. Recovery email skipped.")
		return nil
	}

	subject := fmt.Sprintf("✅ Recovered: %s is back online", sourceName)

	tmpl := `
<!DOCTYPE html>
<html>
<body>
	<h2 style="color: #5cb85c;">✅ Server Recovered: {{.SourceName}}</h2>
	<p>Good news! The server has recovered and is back online.</p>
	<p><strong>Downtime Duration:</strong> {{.Downtime}}</p>
	<p><strong>Recovered At:</strong> {{.RecoveredAt}}</p>
	<hr/>
	<p><a href="https://ulam.your-domain.com/status">View Status Page</a></p>
</body>
</html>
`
	t, err := template.New("recovery").Parse(tmpl)
	if err != nil {
		return err
	}

	data := struct {
		SourceName  string
		Downtime    string
		RecoveredAt string
	}{
		SourceName:  sourceName,
		Downtime:    downtimeDuration,
		RecoveredAt: time.Now().Format(time.RFC3339),
	}

	var body bytes.Buffer
	body.Write([]byte(fmt.Sprintf("To: %s\r\n", to)))
	body.Write([]byte(fmt.Sprintf("From: %s\r\n", s.From)))
	body.Write([]byte(fmt.Sprintf("Subject: %s\r\n", subject)))
	body.Write([]byte("MIME-Version: 1.0\r\n"))
	body.Write([]byte("Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n"))

	if err = t.Execute(&body, data); err != nil {
		return err
	}

	return s.sendEmail(to, body.Bytes())
}

// SendHTML sends a raw HTML email
func (s *SMTPEmailService) SendHTML(to, subject, htmlBody string) error {
	if s.Host == "" || s.Username == "" || s.Password == "" {
		log.Println("WARNING: SMTP credentials not fully configured. Email skipped.")
		return nil
	}

	var body bytes.Buffer
	body.Write([]byte(fmt.Sprintf("To: %s\r\n", to)))
	body.Write([]byte(fmt.Sprintf("From: %s\r\n", s.From)))
	body.Write([]byte(fmt.Sprintf("Subject: %s\r\n", subject)))
	body.Write([]byte("MIME-Version: 1.0\r\n"))
	body.Write([]byte("Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n"))
	body.Write([]byte(htmlBody))

	return s.sendEmail(to, body.Bytes())
}

// sendEmail is a helper to send raw email bytes via SMTP
func (s *SMTPEmailService) sendEmail(to string, body []byte) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.Host,
	}

	conn, err := tls.Dial("tcp", s.Host+":"+s.Port, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate SMTP: %v", err)
	}

	if err = client.Mail(s.From); err != nil {
		return fmt.Errorf("failed to issue MAIL command: %v", err)
	}

	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to issue RCPT command: %v", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to issue DATA command: %v", err)
	}

	_, err = w.Write(body)
	if err != nil {
		return fmt.Errorf("failed to write email body: %v", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close email writer: %v", err)
	}

	return nil
}
