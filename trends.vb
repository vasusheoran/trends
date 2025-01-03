Sub Test()

    Dim result As String
    
    result = PostTrend("NF", _
                    "2/12/24", _
                    2, _
                    3, _
                    4, _
                    5)
    Debug.Print result
    
End Sub


' If a function's arguments are defined as follows:
Function PostTrend(ticker As String, _
                d As String, _
                c As String, _
                o As String, _
                h As String, _
                l As String, _
                Optional URL = "http://localhost:5001/api/update/index")
                
Dim data As String
data = "{""date"":""" & d & _
        """,""ticker"":""" & ticker & _
        """,""open"":" & o & _
        ",""close"":" & c & _
        ",""high"":" & h & _
        ",""low"":" & l & " }"
        
Debug.Print data
        
Dim hReq As Object
Set hReq = CreateObject("MSXML2.XMLHTTP")
    With hReq
        .Open "PUT", URL, False
        .SetRequestHeader "Content-Type", "application/json"
        '.SetTimeouts 2000, 2000, 2000, 2000
        .Send data
    End With
    
PostTrend = hReq.ResponseText
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

