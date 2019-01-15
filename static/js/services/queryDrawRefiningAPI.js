var QetchQueryDrawRefining = angular.module('QetchQuery');

QetchQueryDrawRefining.service('QetchQuery_DrawRefining', ['$rootScope', 'DatasetAPI', 'QetchQuery_QueryAPI', 'Parameters',
    function ($rootScope, DatasetAPI, QetchQuery_QueryAPI, Parameters) {
  var self = this;

  this.predefinedQueries = null;
  this.queryHistory = [];

  this.queryFunctions = [];

  // Transforms a list of points in a string format in a real list of points
  this.pointsListToPointArray = function (ptsStr) {
    var ptsStrLst = ptsStr.substring(1, ptsStr.length - 1).split(')(');
    var ptsLst = [];
    for (var j = 0; j < ptsStrLst.length; j++) {
      var pts = ptsStrLst[j].split(',');
      var x = parseFloat(pts[0]);
      var y = parseFloat(pts[1]);
      var pt = new Qetch.Point(x, y, x, y);
      ptsLst.push(pt);
    }
    return ptsLst;
  };

  this.queryUpdated = function (points) {
    var qtangents = QetchQuery_QueryAPI.extractTangents(points);
    var qsections = QetchQuery_QueryAPI.findCurveSections(qtangents, points, Parameters.DIVIDE_SECTION_MIN_HEIGHT_QUERY);
    this.addQueryInHistory({points: points, bounds: this.pointsBounds(points), tangents: qtangents, sections: qsections});

  };

  this.addQueryInHistory = function (query) {
    this.queryHistory.push(query);
    $rootScope.$broadcast(Parameters.QUERY_REFINEMENT_EVENTS.QUERY_HISTORY_UPDATE, this.queryHistory);
  };

  // get the bounds of a list of points
  this.pointsBounds = function (points) {
    var bounds = {
      minY: Number.MAX_SAFE_INTEGER,
      minX: Number.MAX_SAFE_INTEGER,
      maxY: Number.MIN_SAFE_INTEGER,
      maxX: Number.MIN_SAFE_INTEGER
    };
    for (var j = 0; j < points.length; j++) {
      var pt = points[j];
      if (bounds.minX > pt.x) bounds.minX = pt.x;
      if (bounds.maxX < pt.x) bounds.maxX = pt.x;
      if (bounds.minY > pt.y) bounds.minY = pt.y;
      if (bounds.maxY < pt.y) bounds.maxY = pt.y;
    }
    return bounds;
  };

}]);  