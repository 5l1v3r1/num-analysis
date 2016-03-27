(function() {

  var POINT_COUNT = 300;

  function UnitCircle(normNum) {
    this._normNum = normNum;
    this._points = [];
    for (var i = 0; i < POINT_COUNT; ++i) {
      var angle = (Math.PI * i * 2) / POINT_COUNT;
      var x = Math.cos(angle);
      var y = Math.sin(angle);
      var mag = this._norm(x, y);
      this._points.push([x/mag, y/mag]);
    }
  }

  UnitCircle.prototype.stroke = function(ctx, size) {
    ctx.beginPath();
    for (var i = 0, len = this._points.length; i < len; ++i) {
      var p = this._points[i];
      var x = (p[0] + 1) * (size / 2);
      var y = (p[1] + 1) * (size / 2);
      if (i === 0) {
        ctx.moveTo(x, y);
      } else {
        ctx.lineTo(x, y);
      }
    }
    ctx.closePath();
    ctx.stroke();
  };

  UnitCircle.prototype._norm = function(x, y) {
    var sum = Math.pow(Math.abs(x), this._normNum) +
      Math.pow(Math.abs(y), this._normNum);
    return Math.pow(sum, 1/this._normNum);
  };

  window.UnitCircle = UnitCircle;

})();
