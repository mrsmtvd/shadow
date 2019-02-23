$(document).ready(function () {
    var groupColumn = 0;

    var tableLoggers = $('#loggers table')
        .DataTable({
            pageLength: 50,
            language: {
                url: '/dashboard/datatables/i18n.json'
            },
            order: [[groupColumn, 'asc']],
            columnDefs: [{
                visible: false,
                targets: groupColumn
            }],
            'drawCallback': function (settings) {
                var api = this.api();
                var rows = api.rows({page: 'current'}).nodes();
                var last = null;

                api.column(groupColumn, {page: 'current'}).data().each(function (group, i) {
                    var $aRow = $(rows).eq(i);

                    if (last !== group) {
                        $aRow.before('<tr class="group"></td><td colspan="' + $aRow.children().length + '">' + group + '</td></tr>');
                        last = group;
                    }
                });
            }
        });
    tableLoggers.on('click', 'tr.group', function () {
        var currentOrder = tableLoggers.order()[0];

        if (currentOrder[0] === groupColumn && currentOrder[1] === 'asc') {
            tableLoggers.order([groupColumn, 'desc']).draw();
        }
        else {
            tableLoggers.order([groupColumn, 'asc']).draw();
        }
    });
});