var QetchQuery = angular.module('QetchQuery');

QetchQuery.controller('QetchQuery_ResultsCntrl', 
  ['$scope', '$timeout', 'QetchQuery_QueryAPI', 'Dataset_Resource', 'DatasetAPI', 'Parameters',
  function ($scope, $timeout, QetchQuery_QueryAPI, Dataset_Resource, DatasetAPI, Parameters) {

    $scope.multipleSeries = false;
    $scope.selectedSeriesNum = 0;

    $scope.matches = null;
    $scope.orderingColumns = ['+adjMatch'];

    $scope.MAX_MATCH_GOOD = Parameters.MAX_MATCH_GOOD;
    $scope.MAX_MATCH_MEDIUM = Parameters.MAX_MATCH_MEDIUM;

    $scope.$on(Parameters.DATASET_EVENTS.MATCHES_CHANGED, function (event, matches, matchIndex) {
      if (matchIndex !== undefined) return;
      $scope.matches = matches;
      $scope.adjustMatches();
      DatasetAPI.showMatches(null, DatasetAPI.smoothedDataId, $scope.selectedSeriesNum, null, null, false, null);
      $timeout(function() {
        $scope.$apply();
      });
    });

    $scope.$on(Parameters.DATASET_EVENTS.SHOW_MATCHES, function (event, matches, matchIndex, smoothIteration, minimumMatch, maximumMatch) {
      var $results = $('.result-display');
      $results.removeClass('displaying');
      for (var i in matches) {
        $results.filter('[data-match-id="' + matches[i].id + '"]').addClass('displaying');
      }
    }); 

    $scope.$on(Parameters.QUERY_EVENTS.CLEAR, function (event, matches, matchIndex) {
      $scope.matches = null;
    });

    $scope.$on(Parameters.DATASET_EVENTS.DATA_CHANGED, function (event, seriesNum, values, axes) {
      $scope.selectedSeriesNum = seriesNum;
    });

    $scope.changeColumnOrdering = function (columnName) {
      var columnIndex = -1;
      for (var i = 0; i < $scope.orderingColumns.length; i++) {
        if ($scope.orderingColumns[i].substr(1) === columnName) {
          columnIndex = i;
          break;
        }
      }

      if (columnIndex == -1) {
        $scope.orderingColumns.unshift('-' + columnName);
      } else if (columnIndex == 0) {
        $scope.orderingColumns[0] = ($scope.orderingColumns[0].substr(0, 1) == '-' ? '+' : '-') + $scope.orderingColumns[0].substr(1);
      } else {
        $scope.orderingColumns.splice(columnIndex);
        $scope.orderingColumns.unshift('-' + columnName);
      }
    };

    $scope.getLastOrderingColumn = function () {
      return $scope.orderingColumns[0].substr(1);
    };

    $scope.getLastOrderingColumnSign = function () {
      return $scope.orderingColumns[0].substr(0,1);
    };

    $scope.getSeriesName = function (snum) {
      return DatasetAPI.dataset.series[snum].desc;
    };

    $scope.$on(Parameters.DATASET_EVENTS.DATASET_LOADED, function (event, dataset) {
      $scope.multipleSeries = dataset.series.length > 1;
    });

    $scope.matchValueClass = function (matchValue) {
      if (matchValue < $scope.MAX_MATCH_GOOD) {
        return 'label-success';
      } else if (matchValue < $scope.MAX_MATCH_MEDIUM) {
        return 'label-warning';
      } else {
        return 'label-danger';
      }
    };

    // To show a particular match from the list
    $scope.showMatch = function (i, snum, smoothIteration) {
      DatasetAPI.notifyDataChanged(snum);
      DatasetAPI.showDataRepresentation(snum, smoothIteration);
      var $results = $('.result-display');
      $results.removeClass('displaying');
      DatasetAPI.showMatches(i, null, null, null, null, false, null);
    };

    $scope.showMatches = function (minimumMatch, maximumMatch) {
      DatasetAPI.showDataRepresentation($scope.selectedSeriesNum, 0);
      var $results = $('.result-display');
      $results.removeClass('displaying');
      DatasetAPI.showMatches(null, null, $scope.selectedSeriesNum, minimumMatch, maximumMatch, false, null);
    };

    document.showAllMatches = function () {
      $scope.showMatches(0, Number.MAX_VALUE);
      $scope.$apply();
    };

    $scope.adjustMatches = function () {
        var match, i;
        for (i = 0; i < $scope.matches.length; i++) {
          match = $scope.matches[i];
          match.feedback = undefined;
          match.adjMatch = match.match;
        }
      };

    }

  ]);