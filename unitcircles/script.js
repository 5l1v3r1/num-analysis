(function() {

  var circles = [];
  var valueList = [1/3, 1/2, 1, 1.25, 1.5, 1.75, 2, 2.5, 3, 4, 5, 6, 7, 8, 9, 10, 100];
  var canvas = null;
  var slider = null;
  var numField = null;

  function registerSliderEvent() {
    slider.addEventListener('input', updateUI);
  }

  function updateUI() {
    var ctx = canvas.getContext('2d');
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    var idx = parseInt(slider.value);
    ctx.save();
    ctx.translate(1, 1);
    circles[idx].stroke(ctx, canvas.width-2);
    ctx.restore();
    numField.textContent = 'P = ' + valueList[idx];
  }

  window.addEventListener('load', function() {
    for (var i = 0, len = valueList.length; i < len; ++i) {
      circles[i] = new window.UnitCircle(valueList[i]);
    }

    canvas = document.getElementById('canvas');
    slider = document.getElementById('slider');
    numField = document.getElementById('num-field');

    slider.max = valueList.length - 1;

    registerSliderEvent();
    updateUI();
  });

})();
