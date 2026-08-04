package email

import (
	"fmt"
	"html"
)

// OTPTemplateData contains the dynamic values rendered into an OTP email.
type OTPTemplateData struct {
	Purpose      string
	Code         string
	ExpiryMinute int
}

// RenderOTP returns a plain-text fallback and an inline-CSS HTML email.
func RenderOTP(data OTPTemplateData) (string, string) {
	purpose := html.EscapeString(data.Purpose)
	code := html.EscapeString(data.Code)
	expiry := data.ExpiryMinute
	if expiry <= 0 {
		expiry = 10
	}
	text := fmt.Sprintf("ViệcLàm AI\n\n%s\n\nMã xác nhận của bạn là: %s\nMã có hiệu lực trong %d phút.\nNếu bạn không yêu cầu, hãy bỏ qua email này.\n\nTrân trọng,\nĐội ngũ ViệcLàm AI", data.Purpose, data.Code, expiry)
	preheader := html.EscapeString(fmt.Sprintf("Mã xác nhận %s của bạn là %s. Mã có hiệu lực trong %d phút.", data.Purpose, data.Code, expiry))
	htmlBody := fmt.Sprintf(`<!doctype html><html lang="vi"><body style="margin:0;padding:0;background:#f4f7fb;color:#1e293b;font-family:Arial,Helvetica,sans-serif;-webkit-text-size-adjust:100%%"><div style="display:none;max-height:0px;overflow:hidden">%s</div><div style="display:none;max-height:0px;overflow:hidden">&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;&#847; &zwnj;&nbsp;</div><table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="border-collapse:collapse;background:#f4f7fb"><tr><td align="center" style="padding:32px 16px"><table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="max-width:600px;border-collapse:separate;background:#fff;border-radius:16px;box-shadow:0 10px 30px rgba(15,42,79,.08);overflow:hidden"><tr><td style="height:5px;background:#1a56db;font-size:0;line-height:0">&nbsp;</td></tr><tr><td style="padding:28px 32px 12px"><table role="presentation" width="100%%" cellspacing="0" cellpadding="0"><tr><td style="font-size:24px;font-weight:700;color:#1a56db"><div style="display:flex;align-items:center;gap:8px"><span style="display:inline-flex;align-items:center;justify-content:center;width:30px;height:30px;border-radius:9px;background:#edf5ff;color:#1a56db;font-size:20px;line-height:1">🤖</span><span>ViệcLàm AI</span></div></td><td align="right"><span style="display:inline-block;padding:6px 9px;border-radius:999px;background:#e0f7ff;color:#0052cc;font-size:11px;font-weight:700">GIẢI PHÁP AI</span></td></tr></table></td></tr><tr><td style="padding:12px 32px 32px"><h1 style="margin:0 0 16px;font-size:24px;line-height:1.3;color:#1e293b">%s</h1><p style="margin:0 0 24px;font-size:16px;line-height:1.6;color:#64748b">Đây là mã xác nhận cho tài khoản ViệcLàm AI của bạn.</p><table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="border-collapse:separate;background:#f0f6ff;border:2px solid #1a56db;border-radius:12px"><tr><td align="center" style="padding:22px 16px"><div style="font-size:32px;line-height:1.2;font-weight:700;letter-spacing:8px;color:#0052cc">%s</div></td></tr></table><p style="margin:20px 0 8px;font-size:14px;line-height:1.5;color:#1e293b"><strong>Mã có hiệu lực trong %d phút.</strong></p><p style="margin:0;font-size:13px;line-height:1.6;color:#64748b">Vì lý do bảo mật, không chia sẻ mã này với bất kỳ ai. Nếu bạn không yêu cầu mã, hãy bỏ qua email này.</p></td></tr><tr><td style="padding:20px 32px;background:#f8fafc;border-top:1px solid #e2e8f0;font-size:12px;line-height:1.6;color:#888888;text-align:center">Cảm ơn bạn đã sử dụng ViệcLàm AI.<br>Đội ngũ ViệcLàm AI · Hỗ trợ tuyển dụng thông minh</td></tr></table></td></tr></table></body></html>`, preheader, purpose, code, expiry)
	return text, htmlBody
}
