(function() {

  function Drawing() {
    this._drawing = document.getElementById('drawing');
    this._line = document.getElementById('interpolated-line');
    this._points = document.getElementById('interpolation-points');

    this._updateSize();
    $(window).resize(this._updateSize.bind(this));
  }

  Drawing.prototype.setInterpolatedLine = function(xs, ys) {
    var data = 'M';
    if (xs.length === 0) {
      data = '';
    }
    for (var i = 0, len = xs.length; i < len; ++i) {
      var x = xs[i];
      var y = ys[i];
      data += ' ' + x.toFixed(3) + ',' + y.toFixed(3);
    }
    this._line.setAttribute('d', data);
  };

  Drawing.prototype._updateSize = function() {
    var width = this._drawing.offsetWidth;
    var height = this._drawing.offsetHeight;
    var boxStr;
    if (width > height) {
      var relHeight = height / width;
      var yStart = (1 - relHeight) / 2;
      boxStr = '0 ' + yStart.toFixed(3) + ' 1 ' + relHeight.toFixed(3);
    } else {
      var relWidth = width / height;
      var xStart = (1 - relWidth) / 2;
      boxStr = xStart.toFixed(3) + ' 0 ' + relWidth.toFixed(3) + ' 1';
    }
    this._drawing.setAttribute('viewBox', boxStr);
  };

  window.Drawing = Drawing;

})();
