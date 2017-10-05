function dumpsRemove() {
    $.post('/profiling/trace/?action=delete&id=all', function() {
        location.reload();
    });
}