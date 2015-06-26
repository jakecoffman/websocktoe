'use strict';

angular.module('game.play', ['ngRoute'])

    .config(['$routeProvider', function($routeProvider) {
        $routeProvider.when('/play', {
            templateUrl: 'play/play.html',
            controller: 'PlayCtrl'
        });
    }])

    .controller('PlayCtrl', ['$scope', '$location', 'Game', function($scope, $location, Game) {
        $scope.game = Game;
        $scope.board = Game.state.board;

        $scope.move = function(x, y) {
            Game.send({x: x, y: y});
        };

        $scope.board = function(x, y) {
            return Game.state.board[x][y];
        };

        $scope.leave = function() {
            Game.send({leave: true});
            $location.path('/');
        };
    }]);
