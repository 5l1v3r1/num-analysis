(function() {

  var SVG_NAMESPACE = 'http://www.w3.org/2000/svg';
  var CIRCLE_RADIUS = 0.1;
  var EDGE_THICKNESS = 0.05;

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
