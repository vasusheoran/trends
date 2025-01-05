var sse;

document.body.addEventListener('htmx:sseOpen', function (e)  {
    sse = e.detail.source
    console.log("Saved the eventSource", sse)
})

document.body.addEventListener('htmx:sseError', function (e)  {
    sse.close()
})

window.onbeforeunload = function () {
    sse.close();
};