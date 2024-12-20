// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package upload

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func AddTickerInput() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"add-ticker-section\"><form class=\"flex flex-row w-full gap-4 items-center max-w-1xl\" hx-encoding=\"multipart/form-data\" hx-post=\"/upload\" _=\"on htmx:xhr:progress(loaded, total) set #progress.value to (loaded/total)*100\"><span class=\"material-symbols-outlined cursor-pointer\" hx-get=\"/select/close\" hx-trigger=\"click\" hx-swap=\"outerHTML\" hx-target=\"#add-ticker-section\">close</span> <input id=\"ticker-name\" placeholder=\"Symbol Name\" name=\"ticker-name\" type=\"text\" class=\"form-control rounded-2xl text-black font-sans text-sm flex-1\"><input id=\"add-ticker-file\" type=\"file\" name=\"file\"><button class=\"bg-white hover:opacity-50 p-3 rounded-2xl border border-black border-solid  max-w-1xl\" hx-post=\"/ticker/init\" hx-trigger=\"click\" hx-target=\"#add-ticker-section\" hx-swap=\"outerHTML\">Upload</button> <progress id=\"progress\" value=\"0\" max=\"100\"></progress></form></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
