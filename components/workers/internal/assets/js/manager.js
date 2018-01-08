$(document).ready(function () {
    $('#workers-show').click(function () {
        $('#workers .task-show:has(i.glyphicon-eye-open)').click();
    });

    $('#workers-hide').click(function () {
        $('#workers .task-show:has(i.glyphicon-eye-close)').click();
    });

    $('#workers-add button[type=submit]').click(function () {
        $.ajax({
            type: 'POST',
            url: '/workers/?action=workers-add',
            data: {
                count: $('#workers-add-count').val()
            },
            success: update
        });
    });

    var tableListeners = $('#listeners table')
        .on('draw.dt', function (e, settings) {
            if (settings.json) {
                $('#listeners-count').text(settings.json.recordsTotal);
            }
        })
        .DataTable({
            ajax: {
                url: '/workers/?action=stats&entity=listeners',
                dataSrc: 'data'
            },
            columns: [
                { data: 'id' },
                { data: 'name' },
                {
                    data: 'events',
                    render: function (data, type, row) {
                        var
                            content = '',
                            keys = Object.keys(data);

                        keys.sort();

                        for (var i = 0; i < keys.length; i++) {
                            if (!row.locked) {
                                content += '<a href="#" title="Removing listener" class="label label-info" data-toggle="modal" data-target="#modal" data-modal-title="Confirm remove listener #' + row.id + ' for event ' + keys[i] + '" data-modal-callback="listenersRemove(\'' + row.id + '\', \'' + data[keys[i]] + '\');">' + keys[i] + ' x</a> '
                            } else {
                                content += '<span class="label label-info">' + keys[i] + '</span> ';
                            }
                        }

                        return content;
                    }
                },
                { data: 'fires' },
                {
                    data: 'first_fired_at',
                    render: function (date) {
                        return date ? dateToString(date) : '';
                    }
                },
                {
                    data: 'last_fired_at',
                    render: function (date) {
                        return date ? dateToString(date) : '';
                    }
                },
                {
                    orderable: false,
                    data: null,
                    render: function(data, type, row) {
                        if (row.locked) {
                            return '';
                        }

                        return '<div class="btn-group btn-group-xs">'
                            + '<button type="button" class="btn btn-danger btn-icon" data-toggle="modal" data-target="#modal" data-modal-title="Confirm remove listener #' + row.id + ' for all events" data-modal-callback="listenersRemove(\'' + row.id + '\');">'
                            + '<i class="glyphicon glyphicon-trash" title="Remove listeners for all events"></i>'
                            + '</button>'
                            + '</div>';
                    }
                }
            ],
            order: [[ 1, 'asc' ], [ 2, 'asc' ]]
        });

    var tableWorkers = $('#workers table')
        .on('draw.dt', function (e, settings) {
            if (settings.json) {
                $('#workers-count').text(settings.json.recordsTotal);
            }
        })
        .DataTable({
            ajax: {
                url: '/workers/?action=stats&entity=workers',
                dataSrc: 'data'
            },
            columns: [
                { data: 'id' },
                {
                    data: 'created',
                    render: function (date) {
                        return dateToString(date);
                    }
                },
                { data: 'status' },
                {
                    data: null,
                    defaultContent: ''
                },
                {
                    orderable: false,
                    data: null,
                    render: function(data) {
                        var content = '<div class="btn-group btn-group-xs">';

                        if (data.task) {
                            content += '<button type="button" class="btn btn-success btn-circle task-show" data-task="\' + i + \'"><i class="glyphicon glyphicon-eye-open" title="Show task\'s details"></i></button>';
                        }

                        content += '<button type="button" class="btn btn-danger btn-icon" data-toggle="modal" data-target="#modal" data-modal-title="Confirm kill worker #' + data.id + '" data-modal-callback="workersRemove(\'' + data.id + '\');">'
                                 + '<i class="glyphicon glyphicon-trash" title="Remove worker"></i>'
                                 + '</button>'
                                 + '</div>';

                        return content;
                    }
                }
            ],
            order: [[ 2, 'asc' ], [ 0, 'asc' ]]
        });

    $('#workers table tbody').on('click', 'button.task-show', function (e) {
            e.preventDefault();
            var b = $(this).find('i');
            var row = tableWorkers.row($(this).closest('tr'));

            if (b.hasClass('glyphicon-eye-open')) {
                var task = row.data().task;

                b.removeClass('glyphicon-eye-open').addClass('glyphicon-eye-close');
                row.child(
                    '<table width="100%">' +
                        '<tr>' +
                            '<td>' +
                                '<ul class="list-group">' +
                                    '<li class="list-group-item"><span class="pull-right text-muted small"><em>' + task.id + '</em></span><strong>ID</strong><br /></li>' +
                                    '<li class="list-group-item"><span class="pull-right text-muted small"><em>' + task.name + '</em></span><strong>Name</strong><br /></li>' +
                                    '<li class="list-group-item"><span class="pull-right text-muted small"><em>' + task.priority + '</em></span><strong>Priority</strong><br /></li>' +
                                    '<li class="list-group-item"><span class="pull-right text-muted small"><em>' + task.repeats + '</em></span><strong>Repeats</strong><br /></li>' +
                                    '<li class="list-group-item"><span class="pull-right text-muted small"><em>' + durationToReadableString(task.repeat_interval) + '</em></span><strong>RepeatInterval</strong><br /></li>' +
                                    '<li class="list-group-item"><span class="pull-right text-muted small"><em>' + durationToReadableString(task.timeout) + '</em></span><strong>Timeout</strong><br /></li>' +
                                    '<li class="list-group-item"><span class="pull-right text-muted small"><em>' + dateToString(task.created_at) + '</em></span><strong>Created</strong><br /></li>' +
                                    '<li class="list-group-item"><span class="pull-right text-muted small"><em>' + (task.started_at ? dateToString(task.started_at) : '') + '</em></span><strong>Started</strong><br /></li>' +
                                    '<li class="list-group-item"><span class="pull-right text-muted small"><em>' + task.status + '</em></span><strong>Status</strong><br /></li>' +
                                    '<li class="list-group-item"><span class="pull-right text-muted small"><em>' + task.attempts + '</em></span><strong>Attempts</strong><br /></li>' +
                                    '<li class="list-group-item"><span class="pull-right text-muted small"><em>' + dateToString(task.allow_start_at) + '</em></span><strong>Allow start</strong><br /></li>' +
                                    '<li class="list-group-item"><span class="pull-right text-muted small"><em>' + (task.first_started_at ? dateToString(task.first_started_at) : '') + '</em></span><strong>First started</strong><br /></li>' +
                                    '<li class="list-group-item"><span class="pull-right text-muted small"><em>' + (task.last_started_at ? dateToString(task.last_started_at) : '') + '</em></span><strong>Last started</strong><br /></li>' +
                                '</ul>' +
                            '</td>' +
                        '</tr>' +
                    '</table>'
                ).show();
            } else {
                b.removeClass('glyphicon-eye-close').addClass('glyphicon-eye-open');
                row.child.hide();
            }
     });

    var tableTasks = $('#tasks table')
        .on('draw.dt', function (e, settings) {
            if (settings.json) {
                $('#tasks-count').text(settings.json.recordsTotal);
            }
        })
        .DataTable({
            ajax: {
                url: '/workers/?action=stats&entity=tasks',
                dataSrc: 'data'
            },
            columns: [
                { data: 'id' },
                { data: 'name' },
                { data: 'priority' },
                { data: 'repeats' },
                {
                    data: 'repeat_interval',
                    render: function (ns) {
                        return durationToReadableString(ns);
                    }
                },
                {
                    data: 'timeout',
                    render: function (ns) {
                        return durationToReadableString(ns);
                    }
                },
                {
                    data: 'created_at',
                    render: function (date) {
                        return dateToString(date);
                    }
                },
                {
                    data: 'started_at',
                    render: function (date) {
                        if (!date) {
                            return '';
                        }

                        return dateToString(date);
                    }
                },
                { data: 'status' },
                { data: 'attempts' },
                {
                    data: 'allow_start_at',
                    render: function (date) {
                        return dateToString(date);
                    }
                },
                {
                    data: 'first_started_at',
                    render: function (date) {
                        if (!date) {
                            return '';
                        }

                        return dateToString(date);
                    }
                },
                {
                    data: 'last_started_at',
                    render: function (date) {
                        if (!date) {
                            return '';
                        }

                        return dateToString(date);
                    }
                },
                {
                    orderable: false,
                    data: null,
                    render: function (data) {
                        return '<div class="btn-group btn-group-xs">'
                            + '<button type="button" class="btn btn-danger btn-icon task-remove" data-toggle="modal" data-target="#modal" data-modal-title="Confirm remove task #' + data.id + '" data-modal-callback="tasksRemove(\'' + data.id + '\');">'
                            + '<i class="glyphicon glyphicon-trash" title="Remove task"></i>'
                            + '</button>'
                            + '</div>';
                    }
                }
            ],
            order: [[ 2, 'asc' ], [ 3, 'asc' ]]
        });

    var update = function() {
        tableListeners.ajax.reload();
        tableWorkers.ajax.reload();
        tableTasks.ajax.reload();
    };

    var autorefresh = null;
    $('#autorefresh').click(function() {
        if (this.checked) {
            if (autorefresh === null) {
                update();
                autorefresh = window.setInterval(update, 1000 * 10);
            }
        } else if (autorefresh !== null) {
            window.clearInterval(autorefresh);
            autorefresh = null;
        }
    });

    window.listenersRemove = function(listenerId) {
        console.log(arguments);

        $.ajax({
            type: 'POST',
            url: '/workers/?action=listeners-remove',
            data: {
                id: listenerId,
                events: Array.apply(null, arguments).slice(1)
            },
            success: function() {
                tableListeners.ajax.reload();
            }
        });
    };

    window.workersRemove = function(id) {
        $.ajax({
            type: 'POST',
            url: '/workers/?action=workers-remove',
            data: {
                id: id
            },
            success: function() {
                tableWorkers.ajax.reload();
                tableTasks.ajax.reload();
            }
        });
    };

    window.tasksRemove = function(id) {
        $.ajax({
            type: 'POST',
            url: '/workers/?action=tasks-remove',
            data: {
                id: id
            },
            success: function() {
                tableWorkers.ajax.reload();
                tableTasks.ajax.reload();
            }
        });
    }
});