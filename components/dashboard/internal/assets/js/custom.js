/**
 * Resize function without multiple trigger
 *
 * Usage:
 * $(window).smartresize(function(){  
 *     // code here
 * });
 */
(function ($, sr) {
    // debouncing function from John Hann
    // http://unscriptable.com/index.php/2009/03/20/debouncing-javascript-methods/
    var debounce = function (func, threshold, execAsap) {
        var timeout;

        return function debounced() {
            var obj = this, args = arguments;

            function delayed() {
                if (!execAsap)
                    func.apply(obj, args);
                timeout = null;
            }

            if (timeout)
                clearTimeout(timeout);
            else if (execAsap)
                func.apply(obj, args);

            timeout = setTimeout(delayed, threshold || 100);
        };
    };

    // smartresize 
    jQuery.fn[sr] = function (fn) {
        return fn ? this.bind('resize', debounce(fn)) : this.trigger(sr);
    };

})(jQuery, 'smartresize');
/**
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */

var CURRENT_URL = window.location.href.split('#')[0].split('?')[0],
    $BODY = $('body'),
    $MENU_TOGGLE = $('#menu_toggle'),
    $SIDEBAR_MENU = $('#sidebar-menu'),
    $RIGHT_COL = $('.right_col'),
    $FOOTER = $('footer');

// Table
$('table input').on('ifChecked', function () {
    checkState = '';
    $(this).parent().parent().parent().addClass('selected');
    countChecked();
});
$('table input').on('ifUnchecked', function () {
    checkState = '';
    $(this).parent().parent().parent().removeClass('selected');
    countChecked();
});

var checkState = '';

$('.bulk_action input').on('ifChecked', function () {
    checkState = '';
    $(this).parent().parent().parent().addClass('selected');
    countChecked();
});
$('.bulk_action input').on('ifUnchecked', function () {
    checkState = '';
    $(this).parent().parent().parent().removeClass('selected');
    countChecked();
});
$('.bulk_action input#check-all').on('ifChecked', function () {
    checkState = 'all';
    countChecked();
});
$('.bulk_action input#check-all').on('ifUnchecked', function () {
    checkState = 'none';
    countChecked();
});

function countChecked() {
    if (checkState === 'all') {
        $(".bulk_action input[name='table_records']").iCheck('check');
    }
    if (checkState === 'none') {
        $(".bulk_action input[name='table_records']").iCheck('uncheck');
    }

    var checkCount = $(".bulk_action input[name='table_records']:checked").length;

    if (checkCount) {
        $('.column-title').hide();
        $('.bulk-actions').show();
        $('.action-cnt').html(checkCount + ' Records Selected');
    } else {
        $('.column-title').show();
        $('.bulk-actions').hide();
    }
}

//hover and retain popover when on popover content
var originalLeave = $.fn.popover.Constructor.prototype.leave;
$.fn.popover.Constructor.prototype.leave = function (obj) {
    var self = obj instanceof this.constructor ?
        obj : $(obj.currentTarget)[this.type](this.getDelegateOptions()).data('bs.' + this.type);
    var container, timeout;

    originalLeave.call(this, obj);

    if (obj.currentTarget) {
        container = $(obj.currentTarget).siblings('.popover');
        timeout = self.timeout;
        container.one('mouseenter', function () {
            //We entered the actual popover – call off the dogs
            clearTimeout(timeout);
            //Let's monitor popover content instead
            container.one('mouseleave', function () {
                $.fn.popover.Constructor.prototype.leave.call(self, self);
            });
        });
    }
};

$('body').popover({
    selector: '[data-popover]',
    trigger: 'click hover',
    delay: {
        show: 50,
        hide: 400
    }
});

$(document).ready(function () {
    init_nprogess();
    
    // init_alerts();
    init_autosize();
    init_datatables();
    init_echarts();
    init_icheck();
    init_modals();
    init_panel_toolbox();
    init_password_show();
    init_progressbar();
    init_sidebar();
    init_select2();
    init_switchery();
    init_tabs();
    init_tagsinput();
    init_tooltip();
    init_validator();
    init_daterangepicker();
});

/**
 * Inits
 */
function init_nprogess() {
    if (typeof NProgress === 'undefined') {
        return;
    }

    $(document).ready(function () {
        NProgress.start();
    });

    $(window).load(function () {
        NProgress.done();
    });
}

/*
function init_alerts() {
    var e = $('#alerts');
    
    if (!e.length) {
        return;
    }
    
    function alertUpdate() {
        $.ajax({
            type: 'GET',
            url: '/alerts/',
            success: function(r){
                if (!r.length) {
                    e.find('.badge').hide();
                    e.find('i.fa').removeClass('green');
                    
                    return;
                }

                e.find('ul li.alert-item').remove();
                r = r.reverse();

                for(var i in r) {
                    var 
                        title = $('<span></span>').text(r[i].title),
                        time = $('<span class="time"></span>').text(r[i].elapsed);

                    if (r[i].icon.length) {
                        title.prepend($('<i class="fa fa-fw"></i>').addClass('fa-' + r[i].icon));
                    }

                    e.find('ul').prepend(
                        $('<li class="alert-item"></li>').append(
                            $('<a href="javascript:void(0)"></a>').append(
                                $('<span></span>').append(title, time)
                            )
                        )
                    );
                }

                e.find('.badge').text(r.length).show();
                e.find('i.fa').addClass('green');

                e.find('ul li').click(function(){
                    e.find('.badge').hide();
                    e.find('i.fa').removeClass('green');
                });
            }
        })
    }

    alertUpdate();
    setInterval(alertUpdate, 1000 * 5);
}
*/

function init_autosize() {
    if(typeof (autosize) === 'undefined'){
        return;
    }

    autosize($('.resizable_textarea'));
}

function init_datatables() {
    if(typeof ($.fn.DataTable) === 'undefined') {
        return;
    }

    $('.datatable').DataTable({
        paging: false,
        fixedHeader: true,
        stateSave: true,
        stateDuration: 0,
        language: {
            url: '/dashboard/datatables/i18n.json?locale=' + window.shadowLocale
        }
    });
}

var echartsTheme = {};

function init_echarts() {
    if(typeof (echarts) === 'undefined') {
        return;
    }

    echartsTheme = {
        color: [
            '#26B99A', '#34495E', '#BDC3C7', '#3498DB',
            '#9B59B6', '#8abb6f', '#759c6a', '#bfd3b7'
        ],
        title: {
            itemGap: 8,
            textStyle: {
                fontWeight: 'normal',
                color: '#408829'
            }
        },
        dataRange: {
            color: ['#1f610a', '#97b58d']
        },
        toolbox: {
            color: ['#408829', '#408829', '#408829', '#408829']
        },
        tooltip: {
            backgroundColor: 'rgba(0,0,0,0.5)',
            axisPointer: {
                type: 'line',
                lineStyle: {
                    color: '#408829',
                    type: 'dashed'
                },
                crossStyle: {
                    color: '#408829'
                },
                shadowStyle: {
                    color: 'rgba(200,200,200,0.3)'
                }
            }
        },
        dataZoom: {
            dataBackgroundColor: '#eee',
            fillerColor: 'rgba(64,136,41,0.2)',
            handleColor: '#408829'
        },
        grid: {
            borderWidth: 0
        },
        categoryAxis: {
            axisLine: {
                lineStyle: {
                    color: '#408829'
                }
            },
            splitLine: {
                lineStyle: {
                    color: ['#eee']
                }
            }
        },
        valueAxis: {
            axisLine: {
                lineStyle: {
                    color: '#408829'
                }
            },
            splitArea: {
                show: true,
                areaStyle: {
                    color: ['rgba(250,250,250,0.1)', 'rgba(200,200,200,0.1)']
                }
            },
            splitLine: {
                lineStyle: {
                    color: ['#eee']
                }
            }
        },
        timeline: {
            lineStyle: {
                color: '#408829'
            },
            controlStyle: {
                normal: {color: '#408829'},
                emphasis: {color: '#408829'}
            }
        },
        k: {
            itemStyle: {
                normal: {
                    color: '#68a54a',
                    color0: '#a9cba2',
                    lineStyle: {
                        width: 1,
                        color: '#408829',
                        color0: '#86b379'
                    }
                }
            }
        },
        map: {
            itemStyle: {
                normal: {
                    areaStyle: {
                        color: '#ddd'
                    },
                    label: {
                        textStyle: {
                            color: '#c12e34'
                        }
                    }
                },
                emphasis: {
                    areaStyle: {
                        color: '#99d2dd'
                    },
                    label: {
                        textStyle: {
                            color: '#c12e34'
                        }
                    }
                }
            }
        },
        force: {
            itemStyle: {
                normal: {
                    linkStyle: {
                        strokeColor: '#408829'
                    }
                }
            }
        },
        chord: {
            padding: 4,
            itemStyle: {
                normal: {
                    lineStyle: {
                        width: 1,
                        color: 'rgba(128, 128, 128, 0.5)'
                    },
                    chordStyle: {
                        lineStyle: {
                            width: 1,
                            color: 'rgba(128, 128, 128, 0.5)'
                        }
                    }
                },
                emphasis: {
                    lineStyle: {
                        width: 1,
                        color: 'rgba(128, 128, 128, 0.5)'
                    },
                    chordStyle: {
                        lineStyle: {
                            width: 1,
                            color: 'rgba(128, 128, 128, 0.5)'
                        }
                    }
                }
            }
        },
        gauge: {
            startAngle: 225,
            endAngle: -45,
            axisLine: {
                show: true,
                lineStyle: {
                    color: [[0.2, '#86b379'], [0.8, '#68a54a'], [1, '#408829']],
                    width: 8
                }
            },
            axisTick: {
                splitNumber: 10,
                length: 12,
                lineStyle: {
                    color: 'auto'
                }
            },
            axisLabel: {
                textStyle: {
                    color: 'auto'
                }
            },
            splitLine: {
                length: 18,
                lineStyle: {
                    color: 'auto'
                }
            },
            pointer: {
                length: '90%',
                color: 'auto'
            },
            title: {
                textStyle: {
                    color: '#333'
                }
            },
            detail: {
                textStyle: {
                    color: 'auto'
                }
            }
        },
        textStyle: {
            fontFamily: 'Arial, Verdana, sans-serif'
        }
    };
}

function init_icheck() {
    if (typeof $.fn.iCheck === 'undefined') {
        return;
    }

    if ($('input.flat')[0]) {
        $(document).ready(function () {
            $('input.flat').iCheck({
                checkboxClass: 'icheckbox_flat-green',
                radioClass: 'iradio_flat-green'
            });
        });
    }
}

function init_modals() {
    $('#modal').on('show.bs.modal', function (event) {
        var target = $(event.relatedTarget);
        var modal = $(this);

        var title = target.data('modal-title');
        if (typeof title !== 'undefined') {
            modal.find('.modal-title').text(title);
        } else {
            modal.find('.modal-title').text('');
        }

        var body = target.data('modal-body');
        var url = target.data('modal-url');

        if (typeof body !== 'undefined' || typeof url !== 'undefined') {
            if (typeof body !== 'undefined') {
                modal.find('.modal-body').text(body);
            } else {
                modal.find('.modal-body').html(
                    '<iframe src="' + url + '" ' + ' style="border:0;height:100%;width:100%">' +
                    '</iframe>'
                );
            }
        } else {
            modal.find('.modal-body').html('');
        }

        var callback = target.data('modal-callback');
        if (typeof callback !== 'undefined') {
            $('#modal').data('modal-callback', callback);
        } else {
            $('#modal').removeData('modal-callback');
        }
    });

    $('#modal button[type=submit]').click(function (e) {
        e.preventDefault();

        var callback = $('#modal').data('modal-callback');
        if (typeof callback !== 'undefined') {
            eval(callback);
        }
    });
}

function init_panel_toolbox() {
    $('.collapse-link').on('click', function () {
        var $BOX_PANEL = $(this).closest('.x_panel'),
            $ICON = $(this).find('i'),
            $BOX_CONTENT = $BOX_PANEL.find('.x_content');

        // fix for some div with hardcoded fix class
        if ($BOX_PANEL.attr('style')) {
            $BOX_CONTENT.slideToggle(200, function () {
                $BOX_PANEL.removeAttr('style');
            });
        } else {
            $BOX_CONTENT.slideToggle(200);
            $BOX_PANEL.css('height', 'auto');
        }

        $ICON.toggleClass('fa-chevron-up fa-chevron-down');
    });

    $('.close-link').click(function () {
        var $BOX_PANEL = $(this).closest('.x_panel');

        $BOX_PANEL.remove();
    });

    $('.collapsed').css('height', 'auto');
    $('.collapsed').find('.x_content').css('display', 'none');
    $('.collapsed').find('i').toggleClass('fa-chevron-up fa-chevron-down');
}

function init_password_show() {
    if (typeof $.fn.password === 'undefined') {
        return;
    }

    $('input[type=password].password-show').password();
}

function init_progressbar() {
    if (typeof $.fn.progressbar === 'undefined') {
        return;
    }

    if ($('.progress .progress-bar')[0]) {
        $('.progress .progress-bar').progressbar();
    }
}

function init_sidebar() {
    // TODO: This is some kind of easy fix, maybe we can improve this
    var setContentHeight = function () {
        $RIGHT_COL.css('min-height', $(window).height() - $FOOTER.outerHeight());
    };

    $SIDEBAR_MENU.find('a').on('click', function (ev) {
        var $li = $(this).parent();

        if ($li.is('.active')) {
            $li.removeClass('active active-sm');
            $('ul:first', $li).slideUp(function () {
                setContentHeight();
            });
        } else {
            // prevent closing menu if we are on child menu
            if (!$li.parent().is('.child_menu')) {
                $SIDEBAR_MENU.find('li').removeClass('active active-sm');
                $SIDEBAR_MENU.find('li ul').slideUp();
            } else {
                if ($BODY.is(".nav-sm")) {
                    $li.parent().find("li").removeClass("active active-sm");
                    $li.parent().find("li ul").slideUp();
                }
            }
            $li.addClass('active');

            $('ul:first', $li).slideDown(function () {
                setContentHeight();
            });
        }
    });

    var menuToggleOnClick = function () {
        var state = $BODY.hasClass('nav-md');

        if (state) {
            $SIDEBAR_MENU.find('li.active ul').hide();
            $SIDEBAR_MENU.find('li.active').addClass('active-sm').removeClass('active');
        } else {
            $SIDEBAR_MENU.find('li.active-sm ul').show();
            $SIDEBAR_MENU.find('li.active-sm').addClass('active').removeClass('active-sm');
        }
        saveStateStorage(0, 'menu_toggle_hide', state);

        $BODY.toggleClass('nav-md nav-sm');

        setContentHeight();

        $('.dataTable').each(function () {
            $(this).dataTable().fnDraw();
        });
    };

    if  (!$BODY.hasClass('nav-sm') && loadStateStorage(0, 'menu_toggle_hide')) {
        menuToggleOnClick();
    }

    // toggle small or large menu
    $MENU_TOGGLE.on('click', menuToggleOnClick);

    // check active menu
    $SIDEBAR_MENU.find('a[href="' + CURRENT_URL + '"]').parent('li').addClass('current-page');

    $SIDEBAR_MENU.find('a').filter(function () {
        return this.href == CURRENT_URL;
    }).parent('li').addClass('current-page').parents('ul').slideDown(function () {
        setContentHeight();
    }).parent().addClass('active');

    // recompute content when resizing
    $(window).smartresize(function () {
        setContentHeight();
    });

    setContentHeight();

    // fixed sidebar
    if ($.fn.mCustomScrollbar) {
        $('.menu_fixed').mCustomScrollbar({
            autoHideScrollbar: true,
            theme: 'minimal',
            mouseWheel: {preventDefault: true}
        });
    }
}

function init_select2() {
    if (typeof $.fn.select2 === 'undefined') {
        return;
    }

    $.fn.select2.defaults.set('width', '100%');

    $('.select2').select2();
}

function init_switchery() {
    if (typeof Switchery === 'undefined') {
        return;
    }

    if ($('.js-switch')[0]) {
        var elems = Array.prototype.slice.call(document.querySelectorAll('.js-switch'));
        elems.forEach(function (html) {
            var switchery = new Switchery(html, {
                color: '#26B99A'
            });
        });
    }
}

function init_tabs() {
    var showTab = getParameterByName('tab');
    if (showTab) {
        $('.nav-tabs a[href=#' + showTab + ']').tab('show') ;
    }
}

function init_tagsinput() {
    if (typeof $.fn.tagsInput === 'undefined') {
        return;
    }

    var cb = function() {
        $(this).change();
    };

    var opts = {
        width: 'auto',
        onAddTag: cb,
        onRemoveTag: cb
    };

    $('input[type=text].tags,textarea.tags').each(function() {
        var el = $(this);

        el.tagsInput(jQuery.extend(opts, {
            defaultText: el.data('default-text') || 'add a tag'
        }));
    });
}

function init_tooltip() {
    $('[data-toggle="tooltip"]').tooltip({
        container: 'body'
    });
}

function init_validator() {
    if(typeof (validator) === 'undefined'){
        return;
    }

    validator.message.date = 'not a real date';

    $('form')
        .on('blur', 'input[required], input.optional, select.required', validator.checkField)
        .on('change', 'select.required', validator.checkField)
        .on('keypress', 'input[required][pattern]', validator.keypress)
        .submit(function(e) {
            if (validator.checkAll(this)) {
                $(this).data('valid', true);
                return true;
            }

            $(this).data('valid', false);
            return false;
        });

    $('.multi.required').on('keyup blur', 'input', function() {
        validator.checkField.apply($(this).siblings().last()[0]);
    });
}

function init_daterangepicker() {
    if(typeof ($.fn.daterangepicker) === 'undefined') {
        return;
    }

    $.fn.daterangepicker.defaultOptions = {
        timePicker: true,
        timePicker24Hour: true,
        applyButtonClasses: 'btn-success',
        locale: {
            format: 'YYYY.MM.DD HH:mm:ss',
            applyLabel: window.shadowI18n.labels.apply,
            cancelLabel: window.shadowI18n.labels.cancel,
            fromLabel: window.shadowI18n.labels.from,
            toLabel: window.shadowI18n.labels.to,
            customRangeLabel: window.shadowI18n.labels.customRange,
            weekLabel: window.shadowI18n.labels.week,
            daysOfWeek: window.shadowI18n.daysOfWeek,
            monthNames: window.shadowI18n.monthNames,
            firstDay: 1
        }
    };
}

function loadStateStorage(iStateDuration, path) {
    try {
        return JSON.parse(
            (iStateDuration === -1 ? sessionStorage : localStorage).getItem('shadow_'+path)
        );
    } catch (e) {
        return {};
    }
}

function saveStateStorage(iStateDuration, path, data) {
    try {
        (iStateDuration === -1 ? sessionStorage : localStorage).setItem('shadow_'+path, JSON.stringify(data));
    } catch (e) {}
}