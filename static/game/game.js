'use strict';

angular.module('game.game', ['ngRoute'])

    .config(['$routeProvider', function($routeProvider) {
        $routeProvider.when('/game', {
            templateUrl: 'game/game.html',
            controller: 'GameCtrl'
        });
    }])

    .controller('GameCtrl', ['$scope', 'Game', function($scope, Game) {
        $scope.game = Game;
    }]);
