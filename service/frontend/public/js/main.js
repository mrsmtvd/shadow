$(document).ready(function() {
    function getParameterByName(name) {
      var match = RegExp('[?&]' + name + '=([^&]*)').exec(window.location.search);
      return match && decodeURIComponent(match[1].replace(/\+/g, ' '));
    }

  if (location.hash.substr(0,2) == '#!') {
    $("a[href='#" + location.hash.substr(2) + "']").tab('show');
  }

  $('a[data-toggle="tab"],a[data-toggle="pill"]').on('shown.bs.tab', function (e) {
      var hash = $(e.target).attr('href');
      if (hash.substr(0,1) == "#") {
        location.replace('#!' + hash.substr(1));
      }
    });

  $('.nav-tabs a').click(function (e) {
    e.preventDefault();
    $(this).tab('show');
  });

  var showTab = getParameterByName('tab');
  if (showTab) {
    $('.nav-tabs a[href=#' + showTab + ']').tab('show') ;
  }
});
