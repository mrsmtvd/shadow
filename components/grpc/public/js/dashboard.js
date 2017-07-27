$(function() {
    $('form.call').submit(function () {
        event.preventDefault();
        var e = $(this);
        var h = $(this).parent();

        h.waitMe({
            effect : 'facebook'
        });

        $.ajax({
            type: e.attr('method'),
            url: e.attr('action'),
            data: e.serialize(),
            success: function(r) {
                var p = '';

                if (r.result) {
                    p = JSON.stringify(JSON.parse(r.result), null, 4);
                } else if (r.error) {
                    p = r.error
                }

                $('#' + e.data('result')).html(p).show();
                h.waitMe('hide');
            }
        });
    });
});