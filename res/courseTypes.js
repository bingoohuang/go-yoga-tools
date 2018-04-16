(function () {
    function attachRadioChangeEvent(radioParent, tagChangedCell) {
        radioParent.find(':radio').change(function () {
            var $td = $(this).parent('td')
            var newVal = $(this).val()
            if (tagChangedCell) {
                $td.toggleClass('changedCell', newVal !== $td.attr('oldValue'))
            }
            $td.attr('newValue', newVal)
        })
    }

    $.maintainCourseTypes = function (oldResultId) {
        $.ajax({
            type: 'POST',
            url: contextPath + "/queryCourseTypes",
            data: {tid: $.activeMerchantId},
            success: function (content, textStatus, request) {
                resultId = oldResultId || ++$.resultId
                var html = '<div id="resultDiv' + resultId + '">'
                var hasContent = content && content.length
                var maxShowOrder = 0
                if (hasContent) {
                    html += "<table><thead><td>显示顺序</td><td>课种编码</td><td>课种名称</td><td>预订类型</td><td>最少预订人数</td></thead><tbody>"
                    for (var j = 0; j < content.length; j++) {
                        var r = content[j]

                        if (r.ShowOrder > maxShowOrder) maxShowOrder = r.ShowOrder

                        html += '<tr class="existingRows" CourseTypeId="' + r.CourseTypeId + '">'
                        html += '<td pname="SHOW_ORDER" class="contentEditable">' + r.ShowOrder + '</td>'
                        html += '<td pname="COURSE_TYPE_ID">' + r.CourseTypeId + '</td>'
                        html += '<td pname="COURSE_TYPE_NAME" class="contentEditable">' + r.CourseTypeName + '</td>'
                        html += '<td pname="SUBSCRIBE_TYPE"  oldValue="' + r.SubscribeType + '">'
                        html += '<input type="radio" name="SubscribeType_' + r.CourseTypeId + '" value="1" ' + (r.SubscribeType === '1' ? 'checked' : '') + '>小班课'
                        html += '<input type="radio" name="SubscribeType_' + r.CourseTypeId + '" value="2" ' + (r.SubscribeType === '2' ? 'checked' : '') + '>私教课'
                        html += '</td>'
                        html += '<td pname="MIN_SUBSCRIPTIONS"  class="contentEditable">' + r.MinSubscriptions + '</td>'
                        html += '</tr>'
                    }
                    html += "</tbody></table>"
                }
                html += "<button class='new'>新增</button><button class='save'>保存</button></div>"
                var $table = $(html)

                if (oldResultId !== undefined) {
                    $('#resultDiv' + oldResultId).replaceWith($table)
                } else {
                    $table.prependTo($('#resultArea'))
                }

                $table.find('td.contentEditable').attr('contenteditable', true)
                    .each(function (index, td) {
                        var $td = $(td)
                        $td.attr('oldValue', $td.text())
                    })
                    .blur(function () {
                        var $td = $(this)
                        $td.toggleClass('changedCell', $td.text() !== $td.attr('oldValue'))
                    })

                attachRadioChangeEvent($table, true)

                $table.find('.new').click(function () {
                    var showOrder = ++maxShowOrder
                    var newRowHtml = '<tr class="newRows"><td class="contentEditable">' + showOrder + '</td><td class="contentEditable"></td><td class="contentEditable"></td>'
                    newRowHtml += '<td newValue="1">'
                    newRowHtml += '<input type="radio" name="SubscribeType_' + showOrder + '" value="1" checked>小班课'
                    newRowHtml += '<input type="radio" name="SubscribeType_' + showOrder + '" value="2">私教课'
                    newRowHtml += '</td>'
                    newRowHtml += '<td class="contentEditable">0</td></tr>'
                    var $newRow = $(newRowHtml)
                    $newRow.insertAfter($table.find('tbody tr:last'))
                        .find('td.contentEditable').attr('contenteditable', true)
                    attachRadioChangeEvent($newRow, false)
                })
                $table.find('.save').click(function () {
                    var sqls = []
                    $table.find('tbody tr.existingRows').each(function (index, tr) {
                        var $tr = $(tr)
                        var changedCells = $tr.find('td.changedCell')
                        if (changedCells.length) {
                            var updateItems = []
                            changedCells.each(function (j, td) {
                                var $td = $(td)
                                var pname = $td.attr('pname')
                                if (pname === 'SubscribeType') {
                                    $td.find(':radio')
                                    updateItems.push(pname + "='" + $td.attr('newValue') + "'")
                                } else {
                                    updateItems.push(pname + "='" + $td.text() + "'")
                                }
                            })

                            sqls.push('update tt_d_course_type set ' + updateItems.join(',') + " where COURSE_TYPE_ID = '" + $tr.attr('CourseTypeId') + "'")
                        }
                    })

                    $table.find('tbody tr.newRows').each(function (index, tr) {
                        var $tds = $(tr).find('td')
                        var courseTypeId = $tds.eq(1).text();
                        sqls.push('insert into tt_d_course_type(SHOW_ORDER, COURSE_TYPE_ID, COURSE_TYPE_NAME, SUBSCRIBE_TYPE, MIN_SUBSCRIPTIONS, STATE, CREATE_TIME, UPDATE_TIME)' +
                            " values('" + $tds.eq(0).text() + "', '" + courseTypeId + "', '" + $tds.eq(2).text() + "', '" + $tds.eq(3).attr('newValue') + "','" + $tds.eq(4).text() + "', '1', now(), now())")
                        sqls.push("insert into tt_f_course_category select '" + courseTypeId + "',COURSE_ID,STATE,NOW(),NOW() from tt_f_course_category where COURSE_TYPE_ID = '2001'")
                        sqls.push("insert into tt_f_cardtype_rule select CARD_TYPE_ID, '" + courseTypeId + "',DEDU_TIMES, PRIORITY_ORDER, STATE,NOW(),NOW() from tt_f_cardtype_rule where COURSE_TYPE_ID = '2001'")
                    })

                    $.ajax({
                        type: 'POST',
                        url: contextPath + "/updateCourseTypes",
                        data: {
                            tid: $.activeMerchantId,
                            tcode: $.activeMerchantCode,
                            courseTypesSqls: sqls.join(';'),
                            activeHomeArea: $.activeHomeArea
                        },
                        success: function (content, textStatus, request) {
                            $.maintainCourseTypes(resultId)
                        },
                        error: function (jqXHR, textStatus, errorThrown) {
                            alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                        }
                    })
                })

            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }
})()

