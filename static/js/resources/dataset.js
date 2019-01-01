var Dataset = angular.module('Dataset');

Dataset.factory('Dataset_Resource', function($resource) {
  return $resource('/datasets/:key/:snum', null, {
    definition: {
      method: 'GET',
      url: '/datasets/definition'
    },
  });
});