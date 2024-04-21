// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package template

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Body(code string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\" xmlns=\"http://www.w3.org/1999/xhtml\" xmlns:o=\"urn:schemas-microsoft-com:office:office\"><head><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><meta http-equiv=\"Content-type\" content=\"text/html; charset=UTF-8\"><!--[if mso]> \n        <noscript> \n            <xml> \n                <o:OfficeDocumentSettings> \n                <o:PixelsPerInch>96</o:PixelsPerInch> \n                </o:OfficeDocumentSettings> \n            </xml> \n        </noscript> \n    <![endif]--><style>\n    table, td, div, h1, p {font-family: Arial, sans-serif;}\n\ttable, td {border:none !important;}\n    </style></head><body style=\"margin:0;padding:0;\"><table role=\"presentation\" style=\"width:100%;border-collapse:collapse;border:0;border-spacing:0;background:#ffffff;\"><tr><td align=\"center\" style=\"padding:0;\"><span style=\"font-size:36px;font-weight:bold;color:#facc14;\">SAY</span> <span style=\"font-size:36px;font-weight:bold;color:#15903d;\">WUB</span></td></tr><tr><td align=\"center\" style=\"padding:16px 0 0 0;font-size:24px;\">Login code</td></tr><tr><td align=\"center\" style=\"padding:0;font-size:24px;\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(code)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `template/email.templ`, Line: 37, Col: 9}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</td></tr></table></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}