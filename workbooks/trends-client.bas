Attribute VB_Name = "Module1"
Function TemplateReplace$(TEMPLATE$, ParamArray replacements())
    Dim c&, s$, t$, e
    s = TEMPLATE
    For Each e In replacements
        c = c + 1
        t = "|%" & c & "|"
        If InStrB(e, "~~") Then e = Replace(e, "~~", Chr(34))
        If InStrB(s, t) Then s = Replace(s, t, e)
    Next
    TemplateReplace = s
End Function

Function GetSelection(row)
    Dim RetVal
    Const TEMPLATE = "D:\wsl\data\trends\golang\client-app\http\trends-client-app.exe --server localhost:5000 --date ""|%1|"" --symbol ""|%2|"" --close |%3|  --high |%4| --low |%5| "
    Cmd = TemplateReplace$(TEMPLATE, Format(ActiveSheet.Range("U" & row), "m:d:yyyy h:m:s"), ActiveSheet.Range("A" & row).Value, ActiveSheet.Range("B" & row).Value, ActiveSheet.Range("I" & row).Value, ActiveSheet.Range("J" & row).Value)
    RetVal = Shell(Cmd, vbHide)
    'MsgBox Cmd
End Function

Function SAS(rngRef As Range)
    If rngRef.Rows.Count < 2 Then
        MyCellRow = rngRef.row
        GetSelection (MyCellRow)
        SAS = rngRef.Value
    End If
End Function

Sub test()
    row = "2"
    GetSelection (row)
End Sub

