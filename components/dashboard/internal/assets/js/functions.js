function getParameterByName(name) {
    var match = RegExp('[?&]' + name + '=([^&]*)').exec(window.location.search);
    return match && decodeURIComponent(match[1].replace(/\+/g, ' '));
}

function dateToString(d) {
    var timestamp = Date.parse(d);

    if (isNaN(timestamp) === false) {
        return new Date(timestamp).toLocaleString();
    }

    return d
}

function durationToReadableString(ns) {
    var
        hours        = Math.floor( ns / (1000 * 1000 * 1000 * 60 * 60) % 60),
        minutes      = Math.floor(ns / (1000 * 1000 * 1000 * 60) % 60),
        seconds      = Math.floor(ns / (1000 * 1000 * 1000) % 60),
        milliseconds = Math.floor(ns /  1000 * 1000 % 60)
    ;

    hours = hours < 10 ? '0' + hours : hours;
    minutes = minutes < 10 ? '0' + minutes : minutes;
    seconds = seconds < 10 ? '0' + seconds : seconds;

    return hours + ':' + minutes + ':' + seconds + '.' + milliseconds;
}