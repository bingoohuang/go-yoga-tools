(function () {

    $.saveCaptcha = function () {
        var mobile = $("#mobileKey").val()
        var captcha = $("#captchaValue").val()

        $.ajax({
            type: 'POST',
            url: contextPath + "/saveCaptcha",
            data: {mobile: mobile, activeHomeArea: $.activeHomeArea, activeClassifier: $.activeClassifier},
            success: function (content, textStatus, request) {
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    $.setCaptcha = function () {
        var html = '<div id="captchaSettingWrapper">'
        html += '<div id="title">验证码设置 '
        html += '</div>  <div id="mobileWrapper">'
        html += ' <div>手机号码：</div>'
        html += '<input placeholder="手机号码" id="mobileKey"/> </div>'
        html += '<div id="captchaValueWrapper">'
        html += '<div>验证码：</div>'
        html += '<input id="captchaValue" readonly value="1234">'
        html += '</div> </div>'
        var $table = $(html)

        if ($("#captchaSettingWrapper").length > 0) {
            $('#captchaSettingWrapper').replaceWith($table)
        } else {
            $table.prependTo($('#captchaArea'))
        }

        $('#mobileKey').keydown(function (event) {
            var keyCode = event.keyCode || event.which
            if (keyCode == 13) {
                $.saveCaptcha()
            }
        })
    }


})()

