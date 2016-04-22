(function() {

  var DOT_RADIUS = 0.01;
  var SVG_NAMESPACE = 'http://www.w3.org/2000/svg';

  function Drawing() {
    this._drawing = document.getElementById('drawing');
    this._line = document.getElementById('interpolated-line');
    this._points = document.getElementById('interpolation-points');

    this._pointElements = [];
    this._pointXs = [];
    this._pointYs = [];

    this._updateSize();
    $(window).resize(this._updateSize.bind(this));

    this._registerPointerEvents();

    this.onPointAdded = null;
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

  Drawing.prototype.getPointXs = function() {
    return this._pointXs;
  };

  Drawing.prototype.getPointYs = function() {
    return this._pointYs;
  };

  Drawing.prototype.clear = function() {
    this._pointXs = [];
    this._pointYs = [];
    for (var i = 0, len = this._pointElements.length; i < len; ++i) {
      this._points.removeChild(this._pointElements[i]);
    }
    this._pointElements = [];
  };

  Drawing.prototype._updateSize = function() {
    var bbox = this._drawing.getBoundingClientRect();
    var width = bbox.width;
    var height = bbox.height;
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

  Drawing.prototype._registerPointerEvents = function() {
    this._drawing.addEventListener('click', function(e) {
      var boundingRect = this._drawing.getBoundingClientRect();
      var xCoord = e.clientX - boundingRect.left;
      var yCoord = e.clientY - boundingRect.top;
      var maxSize = Math.max(boundingRect.width, boundingRect.height);
      xCoord /= maxSize;
      yCoord /= maxSize;
      if (boundingRect.width > boundingRect.height) {
        yCoord += (1 - boundingRect.height/boundingRect.width) / 2;
      } else {
        xCoord += (1 - boundingRect.width/boundingRect.height) / 2;
      }
      this._addPoint(xCoord, yCoord);
    }.bind(this));
  };

  Drawing.prototype._addPoint = function(x, y) {
    this._pointXs.push(x);
    this._pointYs.push(y);
    var dot = document.createElementNS(SVG_NAMESPACE, 'circle');
    dot.setAttribute('cx', x.toFixed(3));
    dot.setAttribute('cy', y.toFixed(3));
    dot.setAttribute('r', DOT_RADIUS.toFixed(3));
    this._pointElements.push(dot);
    this._points.appendChild(dot);

    if (this.onPointAdded !== null) {
      this.onPointAdded();
    }
  };

  window.Drawing = Drawing;

})();
