package service

import (
	"fmt"
	"net/smtp"
)

type EmailService struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewEmailService(host, port, username, password string) *EmailService {
	return &EmailService{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

// SendResetPasswordEmail gửi email chứa link đặt lại mật khẩu
func (e *EmailService) SendResetPasswordEmail(toEmail, resetToken string) error {
	// Link frontend để đổi pass
	resetLink := fmt.Sprintf("http://localhost:5173/reset-password?token=%s", resetToken)

	// Soạn nội dung email
	subject := "Subject: Yêu cầu đặt lại mật khẩu - AI Interview\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<h2 style="color: #2563eb;">Xin chào,</h2>
			<p>Chúng tôi nhận được yêu cầu đặt lại mật khẩu cho tài khoản của bạn.</p>
			<p>Vui lòng click vào nút dưới đây để tạo mật khẩu mới (Link có hiệu lực trong 1 giờ):</p>
			<p style="margin: 24px 0;">
				<a href="%s" style="padding: 12px 24px; background-color: #2563eb; color: white; text-decoration: none; border-radius: 6px; font-weight: bold;">Đặt lại Mật khẩu</a>
			</p>
			<p>Nếu nút bấm không hoạt động, bạn có thể copy đường dẫn sau và dán vào trình duyệt:</p>
			<p style="background-color: #f3f4f6; padding: 12px; border-radius: 4px; word-break: break-all;">
				<a href="%s">%s</a>
			</p>
			<hr style="border: none; border-top: 1px solid #eee; margin: 24px 0;" />
			<p style="font-size: 12px; color: #6b7280;">Nếu bạn không yêu cầu đặt lại mật khẩu, vui lòng bỏ qua email này.</p>
		</div>
	`, resetLink, resetLink, resetLink)

	msg := []byte(subject + mime + body)

	// Thiết lập Authentication
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)

	// Gửi mail
	addr := fmt.Sprintf("%s:%s", e.Host, e.Port)
	err := smtp.SendMail(addr, auth, e.Username, []string{toEmail}, msg)
	if err != nil {
		return err
	}
	return nil
}

// SendOTPEmail gửi email chứa mã OTP
func (e *EmailService) SendOTPEmail(toEmail, otpCode, purpose string) error {
	var subject, action string
	if purpose == "register" {
		subject = "Subject: Xác thực email đăng ký - AI Interview\r\n"
		action = "xác thực tài khoản"
	} else if purpose == "forgot_password" {
		subject = "Subject: Mã OTP đặt lại mật khẩu - AI Interview\r\n"
		action = "đặt lại mật khẩu"
	} else {
		subject = "Subject: Mã xác thực OTP - AI Interview\r\n"
		action = "xác thực"
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 500px; margin: 0 auto; padding: 20px; border: 1px solid #eaeaea; border-radius: 8px;">
			<h2 style="color: #2563eb; text-align: center;">Mã xác thực OTP</h2>
			<p>Xin chào,</p>
			<p>Bạn đã yêu cầu %s trên hệ thống AI Interview. Vui lòng sử dụng mã OTP dưới đây để tiếp tục:</p>
			<div style="text-align: center; margin: 30px 0;">
				<span style="font-size: 32px; font-weight: bold; letter-spacing: 5px; color: #111; background: #f3f4f6; padding: 15px 30px; border-radius: 8px;">%s</span>
			</div>
			<p style="color: #666; font-size: 14px;">Mã OTP này có hiệu lực trong vòng 15 phút. Tuyệt đối không chia sẻ mã này cho bất kỳ ai.</p>
			<hr style="border: none; border-top: 1px solid #eee; margin: 24px 0;" />
			<p style="font-size: 12px; color: #9ca3af; text-align: center;">Nếu bạn không yêu cầu %s, vui lòng bỏ qua email này.</p>
		</div>
	`, action, otpCode, action)

	msg := []byte(subject + mime + body)

	// Thiết lập Authentication
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)

	// Gửi mail
	addr := fmt.Sprintf("%s:%s", e.Host, e.Port)
	err := smtp.SendMail(addr, auth, e.Username, []string{toEmail}, msg)
	if err != nil {
		return err
	}
	return nil
}
