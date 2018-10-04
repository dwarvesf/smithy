package email

import (
	"fmt"
	"os"

	sg "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	subject = "Account recovery code"
)

type SendGrid struct {
	apiKey string
}

func New(apiKey string) *SendGrid {
	return &SendGrid{apiKey: apiKey}
}

//Send confirmation code
func (sc *SendGrid) Send(toName, toEmail, code string) error {
	from := mail.NewEmail(os.Getenv("FROM_NAME"), os.Getenv("FROM_EMAIL"))
	to := mail.NewEmail(toName, toEmail)
	contentForm := fmt.Sprintf("<div style=\"margin:0;padding:0\" dir=\"ltr\" bgcolor=\"#ffffff\"><table border=\"0\" cellspacing=\"0\" cellpadding=\"0\" align=\"center\" id=\"m_-4346728458872767771m_7919708925948187036email_table\" style=\"border-collapse:collapse\"><tbody><tr><td id=\"m_-4346728458872767771m_7919708925948187036email_content\" style=\"font-family:Helvetica Neue,Helvetica,Lucida Grande,tahoma,verdana,arial,sans-serif;background:#ffffff\"><span class=\"im\"><table border=\"0\" width=\"100 \" cellspacing=\"0\" cellpadding=\"0\" style=\"border-collapse:collapse\"><tbody><tr><td height=\"20\" style=\"line-height:20px\" colspan=\"3\"> </td></tr><tr><td height=\"1\" colspan=\"3\" style=\"line-height:1px\"><span style=\"color:#ffffff;display:none!important;font-size:1px\"></span></td></tr><tr><td width=\"15\" style=\"display:block;width:15px\">   </td><td><table border=\"0\" width=\"100\" cellspacing=\"0\" cellpadding=\"0\" style=\"border-collapse:collapse\"><tbody><tr><td height=\"16\" style=\"line-height:16px\" colspan=\"3\"> </td></tr><tr><td width=\"32\" align=\"left\" valign=\"middle\" style=\"height:32;line-height:0px\"><img src=\"https://ci3.googleusercontent.com/proxy/hCdvJs5EqIFUSoW4WskmPjwB-RwXBgVhoRaAvbXXPn7ba4r4eV4L36jLvBrSP0UusvOqnTyPOmtj04IUxxTcU6nmiMiK51wWXScyJkrYl_hQKoWskUwkxb3WgoOn7Bv6e32CaCSYYg2MFxGusCEwMOGAJJebG3QKYOprMCYw0DOndI5o9_d_5fwJpuB_ZYonrtzSwD8IDTswFy4WstzMUwMdovx8Np5oRWTTGivLHfYFt6JnxtquIMJSNKgpNi5YFTHH4Qx5YWtG5_SRTiy9_PE=s0-d-e1-ft#https://plus.google.com/u/0/_/focus/photos/public/AIbEiAIAAABECKz208iptJCMrAEiC3ZjYXJkX3Bob3RvKigxMzcyNTM1MDFjMTMzNmY2YmM4NWU5NzMyMzljNjlkZTMwNTcyMGRkMAGNzwpuZNsFK5Nw41llUsA1rPZFNg?sz=32\" width=\"32\" height=\"32\" style=\"border:0\" class=\"m_-4346728458872767771CToWUd CToWUd\"></td><td width=\"15\" style=\"display:block;width:15px\">   </td><td width=\"100\"><span style=\"font-family:Helvetica Neue,Helvetica,Lucida Grande,tahoma,verdana,arial,sans-serif;font-size:19px;line-height:32px;color:#3b5998\">Smithy</span></td></tr><tr style=\"border-bottom:solid 1px #e5e5e5\"><td height=\"16\" style=\"line-height:16px\" colspan=\"3\"> </td></tr></tbody></table></td><td width=\"15\" style=\"display:block;width:15px\">   </td></tr><tr><td width=\"15\" style=\"display:block;width:15px\">   </td><td><table border=\"0\" width=\"100\" cellspacing=\"0\" cellpadding=\"0\" style=\"border-collapse:collapse\"><tbody><tr><td height=\"28\" style=\"line-height:28px\"> </td></tr><tr><td><span class=\"m_-4346728458872767771m_7919708925948187036mb_text\" style=\"font-family:Helvetica Neue,Helvetica,Lucida Grande,tahoma,verdana,arial,sans-serif;font-size:16px;line-height:21px;color:#141823\"><table border=\"0\" cellspacing=\"0\" cellpadding=\"0\" style=\"border-collapse:collapse\"><tbody><tr><td style=\"font-size:11px;font-family:LucidaGrande,tahoma,verdana,arial,sans-serif;padding:10px;background-color:#f2f2f2;border-left:1px solid #ccc;border-right:1px solid #ccc;border-top:1px solid #ccc;border-bottom:1px solid #ccc\"><span class=\"m_-4346728458872767771m_7919708925948187036mb_text\" style=\"font-family:Helvetica Neue,Helvetica,Lucida Grande,tahoma,verdana,arial,sans-serif;font-size:16px;line-height:21px;color:#141823\">%s</span></td></tr></tbody></table><div></td></tr><tr><td height=\"14\" style=\"line-height:14px\"> </td></tr></tbody></table></td><td width=\"15\" style=\"display:block;width:15px\">   </td></tr><tr><td width=\"15\" style=\"display:block;width:15px\">   </td><td><table border=\"0\" width=\"100\" cellspacing=\"0\" cellpadding=\"0\" align=\"left\" style=\"border-collapse:collapse\"><tbody><tr style=\"border-top:solid 1px #e5e5e5\"><td height=\"16\" style=\"line-height:16px\"> </td></tr><tr><td style=\"font-family:Helvetica Neue,Helvetica,Lucida Grande,tahoma,verdana,arial,sans-serif;font-size:11px;color:#aaaaaa;line-height:16px\">This message was sent to <a href=\"mailto:thanhanmoc97@gmail.com\" style=\"color:#3b5998;text-decoration:none\" target=\"_blank\">thanhanmoc97@gmail.com</a> at your request.</td></tr></tbody></table></td><td width=\"15\" style=\"display:block;width:15px\">   </td></tr><tr><td height=\"20\" style=\"line-height:20px\" colspan=\"3\"> </td></tr></tbody></table></span><span><img src=\"https://ci3.googleusercontent.com/proxy/hCdvJs5EqIFUSoW4WskmPjwB-RwXBgVhoRaAvbXXPn7ba4r4eV4L36jLvBrSP0UusvOqnTyPOmtj04IUxxTcU6nmiMiK51wWXScyJkrYl_hQKoWskUwkxb3WgoOn7Bv6e32CaCSYYg2MFxGusCEwMOGAJJebG3QKYOprMCYw0DOndI5o9_d_5fwJpuB_ZYonrtzSwD8IDTswFy4WstzMUwMdovx8Np5oRWTTGivLHfYFt6JnxtquIMJSNKgpNi5YFTHH4Qx5YWtG5_SRTiy9_PE=s0-d-e1-ft#https://plus.google.com/u/0/_/focus/photos/public/AIbEiAIAAABECKz208iptJCMrAEiC3ZjYXJkX3Bob3RvKigxMzcyNTM1MDFjMTMzNmY2YmM4NWU5NzMyMzljNjlkZTMwNTcyMGRkMAGNzwpuZNsFK5Nw41llUsA1rPZFNg?sz=32\" style=\"border:0;width:1px;height:1px\" class=\"m_-4346728458872767771CToWUd CToWUd\"></span></td></tr></tbody></table></div>", code)
	content := mail.NewContent("text/html", contentForm)

	m := mail.NewV3MailInit(from, subject, to, content)

	body := mail.GetRequestBody(m)
	request := sg.GetRequest(sc.apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = body

	_, err := sg.API(request)
	if err != nil {
		return err
	}

	return nil
}
