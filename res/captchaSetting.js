(function () {

    $.saveCaptcha = function () {
        var mobile = $("#mobileKey").val()
        var captcha = $("#captchaValue").val()

        $.ajax({
            type: 'POST',
            url: contextPath + "/saveCaptcha",
            data: {mobile: mobile, captcha: captcha},
            success: function (content, textStatus, request) {
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }
})()

