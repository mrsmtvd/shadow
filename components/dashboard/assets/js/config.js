$(document).ready(function () {
    // save config
    $('#configs input[id]').change(function() {
        var
            e = $(this),
            row = e.parentsUntil('tbody', 'tr'),
            current = e.val(),
            def = this.defaultValue;

        if (e.prop('type') == 'checkbox') {
            current = e.prop('checked') != '';
            def = this.defaultChecked;
        }
        
        if (current == def) {
            row.removeClass('has-error');
        } else {
            row.addClass('has-error');
        }

        $('#configs button[type=submit]').prop('disabled', $('#configs tr.has-error').length == 0);
    });

    $('#configs button[type=reset]').click(function () {
        $('#configs tr.has-error').removeClass('has-error');
        $('#configs button[type=submit]').prop('disabled', true);
    });

    $('#configs button[type=submit]').click(function () {
        var m = $('#modalConfig');
        var c = '';

        $('#configs tr.has-error td input[name][type!=hidden]').each(function(){
            var e = $(this);

            if (e.prop('type') == 'checkbox') {
                v = e.prop('checked');
            } else {
                v = e.val();
            }

            c += '<li><strong>' + e.prop('name') + '</strong>: ' + v + '</li>';
        });

        m.find('.modal-body').html('Changes options:<br /><ul>' + c + '</ul>');
        m.modal();

        return false;
    });

    $('#modalConfig button[type=submit]').click(function() {
        var data = {};

        $('#configs tr.has-error td input[id][type!=hidden]').each(function(){
            var e = $(this);
            data[e.prop('name')] = e.prop('type') == 'checkbox' ? e.prop('checked') : e.val();
        });

        $.post('#', data, function() {
            window.location.reload();
        });

        return false;
    });

    // show/hide password
    $('#configs input[type=password].password-show')
        .on('show.bs.password hide.bs.password', function () {
            var e = $('#' + $(this).prop('id') + '_value');

            if (e.length) {
                var t = e.text();

                e.text(e.data('value'));
                e.data('value', t);
            }
        })
});