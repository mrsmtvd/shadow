$(document).ready(function () {
    var table = $('#sql')
        .DataTable({
            language: {
                url: '/dashboard/datatables/i18n.json?locale=' + window.shadowLocale
            },
            ajax: {
                url: '/database/migrations/',
                dataSrc: 'data'
            },
            columns: [
                { data: 'id' },
                { data: 'source' },
                {
                    data: 'modified_at',
                    render: function (date) {
                        return date ? dateToString(date) : '';
                    }
                },
                {
                    data: 'applied_at',
                    render: function (date) {
                        return date ? dateToString(date) : '';
                    }
                },
                {
                    orderable: false,
                    data: null,
                    render: function (data) {
                        var content = '<div class="btn-group btn-group-xs">'
                            + '<button class="btn btn-success btn-icon show" onclick="showCode(this)"><i class="fas fa-eye" title="Show"></i></button>'
                            + '<a href="/dashboard/bindata/?path=/' + data.source + '/migrations/' + data.id + '&mode=raw" class="btn btn-info btn-icon"><i class="fas fa-file" title="Raw"></i></a>'
                            + '<a href="/dashboard/bindata/?path=/' + data.source + '/migrations/' + data.id + '&mode=file" class="btn btn-warning btn-icon"><i class="fas fa-file-download" title="Download"></i></a>';

                        if (data.applied_at) {
                            content += '<a href="javascript:void(0)" class="btn btn-danger btn-icon" data-toggle="modal" data-target="#modal" data-modal-title="Confirm rollback ' + data.id + ' for ' + data.source + ' migration" data-modal-callback="migrate(\'down\',\'' + data.id + '\',\'' + data.source + '\');">'
                                + '<i class="fa fa-backward"></i>' +
                                '</a>';
                        } else {
                            content += '<a href="javascript:void(0)" class="btn btn-danger btn-icon" data-toggle="modal" data-target="#modal" data-modal-title="Confirm apply ' + data.id + ' for ' + data.source + ' migration" data-modal-callback="migrate(\'up\',\'' + data.id + '\',\'' + data.source + '\');">'
                                + '<i class="fa fa-play"></i>' +
                                '</a>';
                        }

                        return content + '</div>';
                    }
                },
                {
                    data: 'up',
                    visible: false
                },
                {
                    data: 'down',
                    visible: false
                },
            ],
            'drawCallback': function () {
                var api = this.api();
                var rows = api.rows({page: 'current'}).nodes();
                var last = null;

                api.column(5, {page: 'current'}).data().each(function (group, i) {
                    var row = $(rows).eq(i);

                    if (last !== group) {
                        row.after('<tr class="no-hover" style="display:none"><td colspan="' + row.children().length + '">'
                            + '<pre>'
                            + '<button type="button" class="close" onclick="hideCode(this)">×</button>'
                            + "<code class=\"sql\">-- +migrate Up\n" + group + "\n\n-- +migrate Down" + api.column(6, {page:'current'} ).data()[i] + '</code>'
                            + '</pre>'
                            + '</td></tr>');
                        last = group;
                    }
                });

                hljs.initHighlightingOnLoad();
            }
        });

    window.showCode = function (e) {
        $(e).find('i')
            .toggleClass('fas fa-eye')
            .toggleClass('fas fa-eye-slash');

        $(e).closest('#sql tbody tr').next().toggle();
    };

    window.hideCode = function(e) {
        $(e).closest('#sql tbody tr').prev().find('button.show').click();
    };

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

                table.ajax.reload();
            }
        });
    };
});