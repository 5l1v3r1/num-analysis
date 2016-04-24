(function() {

  var drawing = null;
  var $techniqueDropdown = null;

  function initialize() {
    $techniqueDropdown = $('#technique');
    $techniqueDropdown.change(interpolate);
    drawing = new window.Drawing();
    drawing.onPointAdded = interpolate;
    $('#reset-button').click(drawing.clear.bind(drawing));
  }

  function interpolate() {
    var xs = drawing.getPointXs();
    var ys = drawing.getPointYs();

    var interpCount = Math.ceil(drawing.getWidth() + 1);
    var interpStep = 1 / (interpCount - 1);

    var interpXs = [];
    for (var i = 0; i < interpCount; ++i) {
      interpXs.push(i * interpStep);
    }
    var method = $techniqueDropdown.val();
    var interpYs = window.interpolate(method, xs, ys, 0, interpStep, interpCount);
    drawing.setInterpolatedLine(interpXs, interpYs);
  }

  $(initialize);

})();
