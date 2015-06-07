'use strict';

angular.module('game', [
    'ngRoute',
    'game.setup',
    'game.game'
])
    .config(['$routeProvider', function($routeProvider) {
        $routeProvider.otherwise({redirectTo: '/setup'})
    }])
    .factory('Game', ['$rootScope', '$location', function($rootScope, $location){
        var ws = new WebSocket('ws://' + $location.host() + ':' + $location.port() + '/ws');

        var Game = {messages:[], state: {}};

        Game.send = function(data){
            ws.send(JSON.stringify(data));
        };

        ws.onopen = function(e){
            $rootScope.$apply(function(){
                Game.messages.push("Connected");
            });
        };

        ws.onclose = function(e) {
            $rootScope.$apply(function(){
                Game.messages.push("Disconnected");
            });
        };

        ws.onmessage = function(e) {
            $rootScope.$apply(function(){
                var data = JSON.parse(e.data);
                Game.state = data;
                $location.path(data.view);
            });
        };

        return Game;
    }]);
