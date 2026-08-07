package main
import ("fmt";"net/smtp")
func main() { auth := smtp.PlainAuth("", "laidlpk04260@gmail.com", "ftyc zjwj lxdx dhxj", "smtp.gmail.com"); err := smtp.SendMail("smtp.gmail.com:587", auth, "laidlpk04260@gmail.com", []string{"luulai2006@gmail.com"}, []byte("Subject: Test\r\n\r\nThis is a test email")); if err != nil { fmt.Println("Error:", err) } else { fmt.Println("Success") } }
