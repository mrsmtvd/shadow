$(document).ready(function () {
    hljs.initHighlightingOnLoad();

    $('#sql tbody tr.description button.show').click(function() {
        $(this).find('i')
            .toggleClass('fas fa-eye')
            .toggleClass('fas fa-eye-slash');

        $(this).closest('#sql tbody tr').next().toggle();
    });

    $('#sql tbody button.close').click(function() {
        $(this).closest('#sql tbody tr').prev().find('button.show').click();
    });

    window.migrate = function (action, id, source) {
        var url = '/database/migrations/' + action + '/';

        if (source !== '') {
            url += source + '/';
        }

        if (id !== '') {
            url += id;
        }

        $.post(url, function(r) {
            if (r.result === 'failed') {
                new PNotify({
                    title: 'Result operation',
                    text: r.message,
                    type: 'error',
                    hide: false,
                    styling: 'bootstrap3'
                });
                return
            }

            if (r.message !== 'undefined') {
                new PNotify({
                    title: 'Result operation',
                    text: r.message,
                    type: 'success',
                    hide: false,
                    styling: 'bootstrap3'
                });
            }
        });
    }
});