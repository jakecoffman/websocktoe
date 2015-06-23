'use strict';

angular.module('game', [
    'ngRoute',
    'game.setup',
    'game.play'
])
    .config(['$routeProvider', function ($routeProvider) {
        $routeProvider.otherwise({redirectTo: '/setup'})
    }])
    .run(['$location', function ($location) {
        $location.path("/");
    }])
    .factory('Game', ['$rootScope', '$location', function ($rootScope, $location) {
        var ws = new WebSocket('ws://' + $location.host() + ':' + $location.port() + '/ws');

        var Game = {messages: [], state: {}};

        Game.send = function (data) {
            ws.send(JSON.stringify(data));
        };

        ws.onopen = function (e) {
            $rootScope.$apply(function () {
                Game.messages.unshift("Connected");
            });
        };

        ws.onclose = function (e) {
            $rootScope.$apply(function () {
                Game.messages.unshift("Disconnected");
            });
        };

        ws.onmessage = function (e) {
            $rootScope.$apply(function () {
                var data = JSON.parse(e.data);
                console.log(data);
                switch (data.type) {
                    case "state":
                        Game.state = data;
                        $location.path(data.view.toLowerCase());
                        break;
                    case "message":
                        Game.messages.unshift(data.message);
                        break;
                    default:
                        console.log("Unknown message type", data.type);
                }
            });
        };

        return Game;
    }])
    .controller('RootCtrl', ['$scope', 'Game', function ($scope, Game) {
        $scope.game = Game;
    }]);
