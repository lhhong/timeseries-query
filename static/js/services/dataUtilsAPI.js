var Dataset = angular.module('Dataset');

Dataset.service('Data_Utils', ['$rootScope', 'Parameters', function($rootScope, Parameters) {

  // Moving average (we iterate it multiple times) (no array copies)
  this.smooth = function (data, iterations, space) {
    var i, it;
    for (it = 0; it < iterations; it++) {
      for (i = 1; i < data.length - 1; i++) {
        var count = 1;
        var valuesSum = data[i].y;
        var origValuesSum = data[i].origY;
        for (var s = 1; s <= space ; s++) {
          if (i - s >= 0) {
            valuesSum += data[i - s].y;
            origValuesSum += data[i - s].origY;
            count += 1;
          }
          if (i + s < data.length) {
            valuesSum += data[i + s].y;
            origValuesSum += data[i + s].origY;
            count += 1;
          }
        }
        data[i].y = valuesSum / count;
        data[i].origY = origValuesSum / count;
      }
    }
  };

  this.dataHeight = function (data) {
    var miny = data[0].y, maxy = data[0].y;
    for (var i = 1; i < data.length - 1; i++) {
      if (miny > data[i].y) miny = data[i].y;
      if (maxy < data[i].y) maxy = data[i].y;
    }
    return maxy - miny;
  };

}]);
