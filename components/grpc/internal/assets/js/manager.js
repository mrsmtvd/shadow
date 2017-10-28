$(function() {
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