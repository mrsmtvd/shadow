$(function() {
    $('.action button.action-remove').click(function() {
        var
            e = $(this),
            r = e.closest('.input-group'),
            p = r.parent();

        if (!r.data('repeated')) {
            e.addClass('hide');
            r.find('.action-add').removeClass('hide');
            r.find('.action-input').addClass('hide');
        } else {
            var rs = p.find('.input-group');

            if (rs.length > 1) {
                r.remove();

                if (rs.length === 2 && r.data('required')) {
                    r.find('.action-remove').addClass('hide');
                }
            } else {
                r.find('.action-remove').addClass('hide');

                if (!r.data('required')) {
                    r.find('.action-input').addClass('hide');
                }
            }
        }
    });

    $('.action button.action-add').click(function() {
        var
            e = $(this),
            r = e.parent().closest('.input-group'),
            v = r.find('.action-input');

        r.find('.action-remove').removeClass('hide');

        if (r.data('repeated')) {
            if (v.hasClass('hide')) {
                v.removeClass('hide');
            } else {
                var c = r.clone(true);
                r.parent().append(c);

                // TODO: select2 & icheck
            }
        } else {
            v.removeClass('hide');
            e.addClass('hide');
        }
    });

    $('.call-result button.close').click(function(){
        $(this).parent().hide();
    });

    $('form.call').submit(function () {
        event.preventDefault();
        var e = $(this),
            h = e.parent(),
            result = $('#' + e.data('result')),
            response = result.find('.response');

        var toggle = function(error) {
            response
                .removeClass(function() {
                    return (error) ? 'alert-success' : 'alert-danger';
                })
                .addClass(function() {
                    return (error) ? 'alert-danger' : 'alert-success';
                });
        };

        h.waitMe({
            effect : 'facebook'
        });

        $.ajax({
            type: e.attr('method'),
            url: e.attr('action'),
            data: e.serialize(),
            complete: function() {
                result.show();
                h.waitMe('hide');
            },
            error: function() {
                response.html('Ajax request failed');
                toggle(true)
            },
            success: function(r) {
                var p = '';

                if (r.headers) {
                    var headers = [];

                    for (var key in r.headers) {
                        for (var val in r.headers[key]) {
                            headers[headers.length] = key.toUpperCase() + ": " + val;
                        }
                    }

                    p = headers.join("\n") + "\n\n"
                }

                if (r.result) {
                    p += JSON.stringify(JSON.parse(r.result), null, 4);
                } else if (r.error) {
                    p += r.error
                }

                response.html(p);
                toggle(r.error)
            }
        });
    });
});