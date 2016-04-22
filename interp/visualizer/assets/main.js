(function() {

  var drawing = null;
  var INTERP_STEP = 0.005;
  var INTERP_COUNT = 201;

  function initialize() {
    drawing = new window.Drawing();
    drawing.onPointAdded = interpolate;
  }

  function interpolate() {
    var xs = drawing.getPointXs();
    var ys = drawing.getPointYs();
    var interpXs = [];
    for (var i = 0; i < INTERP_COUNT; ++i) {
      interpXs.push(i*INTERP_STEP);
    }
    var interpYs = window.interpolate('poly', xs, ys, 0, INTERP_STEP, INTERP_COUNT);
    console.log(interpYs);
    drawing.setInterpolatedLine(interpXs, interpYs);
  }

  $(initialize);

})();
