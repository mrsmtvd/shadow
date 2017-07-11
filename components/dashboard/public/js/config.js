$(function() {
    // save config
    $('#configs input').change(function() {
        $(this).parentsUntil('tbody', 'tr').addClass('has-error');
        $('#configs button[type=submit]').prop('disabled', false);
    });

    $('#configs button[type=reset]').click(function () {
        $('#configs tr.has-error').removeClass('has-error');
        $('#configs button[type=submit]').prop('disabled', true);
    });

    $('#configs button[type=submit]').click(function () {
        var m = $('#modalConfig');
        var c = '';

        $('#configs tr.has-error td input[type!=hidden]').each(function(){
            var e = $(this);

            if (e.prop('type') == 'checkbox') {
                v = e.prop('checked');
            } else {
                v = e.val();
            }

            c += '<li><strong>' + e.prop('id') + '</strong>: ' + v + '</li>';
        });

        m.find('.modal-body').html('Changes options:<br /><ul>' + c + '</ul>');
        m.modal();

        return false;
    });

    $('#modalConfig button[type=submit]').click(function() {
        var data = {};

        $('#configs tr.has-error td input[type!=hidden]').each(function(){
            var e = $(this);
            data[e.prop('id')] = e.prop('type') == 'checkbox' ? e.prop('checked') : e.val();
        });

        $.post('#', data, function() {
            window.location.reload();
        });

        return false;
    });
});