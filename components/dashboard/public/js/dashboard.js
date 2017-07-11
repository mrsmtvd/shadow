$(function() {
    $('#modal').on('show.bs.modal', function (event) {
        var target = $(event.relatedTarget);
        var modal = $(this);


        var title = target.data('modal-title');
        if (title !== 'undefined') {
            modal.find('.modal-title').text(title);
        }

        var body = target.data('modal-body');
        if (body !== 'undefined') {
            modal.find('.modal-body').text(body);
        }

        var callback = target.data('modal-callback');
        if (callback !== 'undefined') {
            $('#modal').data('modal-callback', callback);
        }
    });

    $('#modal button[type=submit]').click(function (e) {
        e.preventDefault();

        var callback = $('#modal').data('modal-callback');
        if (callback !== 'undefined') {
            eval(callback);
        }
    });

    function alertUpdate() {
        $.ajax({
            type: 'GET',
            url: '/alerts/ajax/',
            success: function(r){
                if (!r.length) {
                    return;
                }

                r = r.reverse();

                var nav = $('#navbar-alerts ul');
                nav.find('li[class=alert-item],li[class=divider]').remove();

                $('#navbar-alerts>li').removeClass('disabled');

                for(var i in r) {
                    var item = $('<div></div>').text(r[i].title);

                    if (r[i].icon.length) {
                        item.prepend($('<i class="fa fa-fw"></i>').addClass('fa-' + r[i].icon));
                    }

                    item.append($('<span class="pull-right text-muted small"></span>').text(r[i].elapsed));

                    if (i < r.length) {
                        nav.prepend($('<li class="divider"></li>'));
                    }

                    nav.prepend($('<li class="alert-item"></li>').append($('<a href="javascript:void(0)"></a>').append(item)))
                }
            }
        })
    }

    if ($('#navbar-alerts').length) {
        $('#navbar-alerts>li').click(function(){
            if ($('#navbar-alerts ul li').length) {
                $(this).addClass('disabled');
            }
        });

        alertUpdate();
        setInterval(alertUpdate, 1000 * 5);
    }

    var showTab = getParameterByName('tab');
    if (showTab) {
        $('.nav-tabs a[href=#' + showTab + ']').tab('show') ;
    }
});