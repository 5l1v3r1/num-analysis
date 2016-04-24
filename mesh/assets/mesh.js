(function() {

  var SVG_NAMESPACE = 'http://www.w3.org/2000/svg';
  var CIRCLE_RADIUS = 0.025;
  var EDGE_THICKNESS = 0.01;

  function MeshView() {
    this._element = document.getElementById('mesh');
    this._edges = [];
    this._nodes = [];
    this._size = window.getMeshDimensions();

    var eAndN = edgesAndNodes();

    var edgesElement = document.getElementById('edges');
    for (var i = 0, len = eAndN.edges.length; i < len; i += 4) {
      var edge = document.createElementNS(SVG_NAMESPACE, 'line');
      edge.setAttribute('stroke-width', EDGE_THICKNESS);
      this._edges.push(edge);
      edgesElement.appendChild(edge);
    }

    var nodesElement = document.getElementById('nodes');
    for (var i = 0, len = eAndN.nodes.length; i < len; i += 2) {
      var node = document.createElementNS(SVG_NAMESPACE, 'circle');
      node.setAttribute('r', CIRCLE_RADIUS);
      this._nodes.push(node);
      nodesElement.appendChild(node);
    }

    this.update();
    this._registerMouseEvents();
  }

  MeshView.prototype.update = function() {
    var eAndN = edgesAndNodes();

    var edgesElement = document.getElementById('edges');
    for (var i = 0, len = eAndN.edges.length; i < len; i += 4) {
      var x1 = eAndN.edges[i];
      var y1 = eAndN.edges[i+1];
      var x2 = eAndN.edges[i+2];
      var y2 = eAndN.edges[i+3];
      var edge = this._edges[i/4];
      edge.setAttribute("x1", x1);
      edge.setAttribute("y1", y1);
      edge.setAttribute("x2", x2);
      edge.setAttribute("y2", y2);
    }

    var nodesElement = document.getElementById('nodes');
    for (var i = 0, len = eAndN.nodes.length; i < len; i += 2) {
      var x = eAndN.nodes[i];
      var y = eAndN.nodes[i+1];
      var node = this._nodes[i/2];
      node.setAttribute('cx', x);
      node.setAttribute('cy', y);
    }
  };

  MeshView.prototype._registerMouseEvents = function() {
    var draggingNode = -1;
    this._element.addEventListener('mousedown', function(e) {
      var rect = this._boundingRect();
      var x = (e.clientX - rect.left) / rect.width;
      var y = (e.clientY - rect.top) / rect.height;
      var nodes = getMeshCoords();
      var closestDist = Infinity;
      var closest = -1;
      for (var i = 0, len = nodes.length/2; i < len; ++i) {
        var nx = nodes[i * 2];
        var ny = nodes[i*2 + 1];
        var dist = Math.sqrt(Math.pow(nx-x, 2) + Math.pow(ny-y, 2));
        if (dist < closestDist) {
          closestDist = dist;
          closest = i;
        }
      }
      draggingNode = closest;
    }.bind(this));
    this._element.addEventListener('mouseup', function(e) {
      draggingNode = -1;
    });
    this._element.addEventListener('mousemove', function(e) {
      if (draggingNode < 0) {
        return;
      }
      var rect = this._boundingRect();
      var x = (e.clientX - rect.left) / rect.width;
      var y = (e.clientY - rect.top) / rect.height;
      window.moveMeshNode(draggingNode, x, y);
      this.update();
    }.bind(this));
  };

  MeshView.prototype._boundingRect = function() {
    var rect = this._element.getBoundingClientRect();
    rect = {width: rect.width, height: rect.height, left: rect.left, top: rect.top};
    if (rect.width > rect.height) {
      rect.left += (rect.width - rect.height) / 2;
      rect.width = rect.height;
    } else {
      rect.top += (rect.height - rect.width) / 2;
      rect.height = rect.width;
    }
    return rect;
  };

  function edgesAndNodes() {
    var nodes = window.getMeshCoords();
    var edgeIndices = window.getMeshEdges();
    var edges = [];
    for (var i = 0, len = edgeIndices.length; i < len; ++i) {
      var index = edgeIndices[i];
      edges.push(nodes[index*2]);
      edges.push(nodes[index*2+1]);
    }
    return {edges: edges, nodes: nodes};
  }

  window.MeshView = MeshView;

})();
