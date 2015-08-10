var oncourseApp = angular.module('oncourse');

oncourseApp.controller('oncourseCtrl', ['$scope', '$http', '$filter', '$interval',
  function($scope, $http, $filter, $interval) {
    $http.get('/data').success(function(data) {
      $scope.data = data;
    });
  }
]);
