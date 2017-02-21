$(function() {
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

    function getParameterByName(name) {
        var match = RegExp('[?&]' + name + '=([^&]*)').exec(window.location.search);
        return match && decodeURIComponent(match[1].replace(/\+/g, ' '));
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