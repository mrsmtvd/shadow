$(function() {
    $('.call-result button.close').click(function(){
        $(this).parent().hide();
    });
    
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

                var r = $('#' + e.data('result'));
                
                r.find('.response').html(p);
                r.show();
                
                h.waitMe('hide');
            }
        });
    });
});