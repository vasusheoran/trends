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

Function GetSelection(Row As String) As String
    Dim RetVal
    Const TEMPLATE = "D:\wsl\data\trends\golang\grpc-client-app\trends-client-app.exe --server localhost:5001 --date |%1| --symbol |%2| --close |%3|  --high |%4| --low |%5| "
    RetVal = Shell(TemplateReplace$(TEMPLATE, "1", "2", "3", "4", "5"), vbHide)
End Function

Function SAS(rngRef As Range)
    If rngRef.Rows.Count < 2 Then
        Dim MyCellRow As String
        MyCellRow = rngRef.Row
        GetSelection (MyCellRow)
        SAS = rngRef.Value
    End If
End Function
