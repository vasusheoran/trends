' If a function's arguments are defined as follows:
Function PostTrend(ticker As String, _
                d As String, _
                c As String, _
                o As String, _
                h As String, _
                l As String, _
                Optional URL = "http://localhost:5001/api/update/index")
                
Dim data As String
data = "{""date"":" & d & _
        ",""ticker"":""" & ticker & _
        """,""open"":" & o & _
        ",""close""," & c & _
        ",""high"":" & h & _
        ",""low"":" & l & " }"
        
Dim hReq As Object
Set hReq = CreateObject("MSXML2.XMLHTTP")
    With hReq
        .Open "PUT", URL, False
        .SetRequestHeader "Content-Type", "application/json"
        .SetTimeouts 2000, 2000, 2000, 2000
        .Send d
    End With

'.OnTimeOut = OnTimeOutMessage
Debug.Print hReq.ResponseText
Debug.Print URL
Debug.Print data

PostTrend = hReq.ResponseText

'PostTrend = Trend(URL, data)
End Function


Private Function OnTimeOutMessage()
    'Application.Caller.Value = "Server error: request time-out"
    MsgBox ("Server error: request time-out")
End Function

' If a function's arguments are defined as follows:
Function TestTrend(Optional URL = "http://localhost:5001/api/health")
Set hReq = CreateObject("MSXML2.XMLHTTP")
    With hReq
        .Open "GET", URL, False
        .Send
    End With
TestTrend = hReq.ResponseText
End Function

