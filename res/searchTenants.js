(function () {
    $('#tenantSearchKey').keydown(function (event) {
        var keyCode = event.keyCode || event.which
        if (keyCode == 13) {
            var tenantSearchKey = $.trim($('#tenantSearchKey').val())
            if (tenantSearchKey === '') {
                alert("please input tid/tcode/tname")
                return
            }
            $.searchTenants(tenantSearchKey)
        }
    }).focus(function () {
        $(this).select()
    })

    $.searchTenants = function (tenantSearchKey) {
        $.ajax({
            type: 'POST',
            url: contextPath + "/searchTenants",
            data: {tenantSearchKey: tenantSearchKey},
            success: function (content, textStatus, request) {
                var tenantsHtml = ''

                var hasContent = content && content.length
                if (hasContent) {
                    for (var j = 0; j < content.length; j++) {
                        tenantsHtml += '<span tid="' + content[j].MerchantId
                            + '" tcode="' + content[j].MerchantCode
                            + '" homeArea="' + content[j].HomeArea
                            + '" classifier="' + content[j].Classifier
                            + '">'
                            + content[j].MerchantName + '</span>'
                    }
                }

                $('#toolsDiv').toggle(!!hasContent)
                clearActiveMerchant()
                $('#resultArea').html('')
                $('#tenantsList').html(tenantsHtml)
                $('#tenantsList span:first-child').click()
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    function clearActiveMerchant() {
        $.activeMerchantId = null
        $.activeMerchantCode = null
        $.activeHomeArea = null
        $.activeClassifier = null
        $.activeMerchantName = null
    }

    $('#tenantsList').on('click', 'span', function () {
        $('#tenantsList span').removeClass('active')
        var $this = $(this).addClass('active')

        $.activeMerchantId = $this.attr('tid')
        $.activeMerchantCode = $this.attr('tcode')
        $.activeHomeArea = $this.attr('homeArea')
        $.activeClassifier = $this.attr('classifier')
        $.activeMerchantName = $this.text()

        $('#resultArea').html('')
        $('#tidtcodeSpan').html('&nbsp;<span title="tcode" >' + $.activeMerchantCode + '</span>' +
            '&nbsp;<span title="home area">' + $.activeHomeArea + '</span>' +
            '&nbsp;<span><a href="' + createUrl() + '" target="_blank">Home</a></span>')
    })

    var createUrl = function () {
        var url = ''

        if (location.hostname.indexOf("test.") >= 0) {
            url += 'http://test.go.easy-hi.com'
        } else {
            if ($.activeHomeArea === 'south-center') {
                url += 'https://app.easy-hi.com'
            } else if ($.activeHomeArea === 'north-center') {
                url += 'https://appn.easy-hi.com'
            }
        }

        if ($.activeClassifier === 'yoga') {
            url += '/yoga-system/center/' + $.activeMerchantCode + '/index/vision#'
        } else if ($.activeClassifier === 'et') {
            url += '/et-server/center/' + $.activeMerchantCode + '/vision#'
        }

        return url
    }

})()