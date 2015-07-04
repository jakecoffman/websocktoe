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
        var url = 'ws://' + $location.host() + ':' + $location.port() + '/ws';
        var ws;
        reconnect();

        var Game = {messages: [], state: {}, connected: false};

        Game.send = function (data) {
            ws.send(JSON.stringify(data));
        };

        function onopen(e) {
            $rootScope.$apply(function () {
                Game.connected = true;
            });
        }

        function onclose(e) {
            $rootScope.$apply(function () {
                Game.connected = false;
                reconnect();
            });
        }

        function reconnect() {
            ws = new WebSocket(url);
            ws.onclose = onclose;
            ws.onmessage = onmessage;
            ws.onopen = onopen;
        }

        function onmessage(e) {
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
        }

        return Game;
    }])
    .controller('RootCtrl', ['$scope', 'Game', function ($scope, Game) {
        $scope.game = Game;
    }]);
