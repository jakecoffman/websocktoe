'use strict';

angular.module('game.setup', ['ngRoute'])

    .config(['$routeProvider', function($routeProvider) {
        $routeProvider.when('/setup', {
            templateUrl: 'setup/setup.html',
            controller: 'SetupCtrl'
        });
    }])

    .controller('SetupCtrl', ['$scope', 'Game', function($scope, Game) {
        $scope.game = Game;
        $scope.send = function(name, gameId, choice) {
            var msg = {
                name: name,
                gameId: gameId,
                choice: choice
            };
            console.log("Attempting to send", msg);
            Game.send(msg);
        }
    }]);
