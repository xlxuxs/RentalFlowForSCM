package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

// Config holds email configuration
type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// Service handles email sending
type Service struct {
	config    Config
	templates map[string]*template.Template
}

// NewService creates a new email service
func NewService(config Config) *Service {
	service := &Service{
		config:    config,
		templates: make(map[string]*template.Template),
	}
	service.loadTemplates()
	return service
}

// loadTemplates loads email templates
func (s *Service) loadTemplates() {
	// Booking Created Template
	s.templates["booking_created"] = template.Must(template.New("booking_created").Parse(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #4F46E5; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .button { background: #4F46E5; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Booking Confirmed!</h1>
        </div>
        <div class="content">
            <p>Hi {{.UserName}},</p>
            <p>Your booking request for <strong>{{.ItemTitle}}</strong> has been created successfully.</p>
            <p><strong>Booking Details:</strong></p>
            <ul>
                <li>Start Date: {{.StartDate}}</li>
                <li>End Date: {{.EndDate}}</li>
                <li>Total Amount: {{.TotalAmount}} ETB</li>
            </ul>
            <p>The owner will review your request shortly.</p>
            <p><a href="{{.BookingURL}}" class="button">View Booking</a></p>
        </div>
        <div class="footer">
            <p>© 2025 RentalFlow. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
	`))

	// Payment Successful Template
	s.templates["payment_success"] = template.Must(template.New("payment_success").Parse(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #10B981; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ Payment Successful!</h1>
        </div>
        <div class="content">
            <p>Hi {{.UserName}},</p>
            <p>Your payment of <strong>{{.Amount}} ETB</strong> has been processed successfully.</p>
            <p><strong>Transaction Details:</strong></p>
            <ul>
                <li>Reference: {{.Reference}}</li>
                <li>Amount: {{.Amount}} ETB</li>
                <li>Date: {{.Date}}</li>
            </ul>
            <p>Thank you for your payment!</p>
        </div>
        <div class="footer">
            <p>© 2025 RentalFlow. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
	`))

	// Review Received Template
	s.templates["review_received"] = template.Must(template.New("review_received").Parse(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #F59E0B; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>⭐ New Review!</h1>
        </div>
        <div class="content">
            <p>Hi {{.OwnerName}},</p>
            <p>Your item <strong>{{.ItemTitle}}</strong> received a new {{.Rating}}-star review.</p>
            <p><strong>Review:</strong></p>
            <blockquote>{{.Comment}}</blockquote>
            <p>Keep up the great work!</p>
        </div>
        <div class="footer">
            <p>© 2025 RentalFlow. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
	`))
}

// SendBookingCreated sends booking creation email
func (s *Service) SendBookingCreated(to string, data map[string]interface{}) error {
	return s.send(to, "Booking Confirmed - RentalFlow", "booking_created", data)
}

// SendPaymentSuccess sends payment confirmation email
func (s *Service) SendPaymentSuccess(to string, data map[string]interface{}) error {
	return s.send(to, "Payment Successful - RentalFlow", "payment_success", data)
}

// SendReviewReceived sends review notification email
func (s *Service) SendReviewReceived(to string, data map[string]interface{}) error {
	return s.send(to, "New Review Received - RentalFlow", "review_received", data)
}

// send sends an email using the specified template
func (s *Service) send(to, subject, templateName string, data map[string]interface{}) error {
	// Render template
	tmpl, exists := s.templates[templateName]
	if !exists {
		return fmt.Errorf("template %s not found", templateName)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Prepare email
	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", from, to, subject, body.String()))

	// Send email
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)

	err := smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
