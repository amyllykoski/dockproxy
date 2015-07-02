"use strict";angular.module("dockuiApp",["ngAnimate","ngCookies","ngResource","ngRoute","ngSanitize","ngTouch"]).config(["$routeProvider","$httpProvider",function(a,b){b.defaults.useXDomain=!0,b.defaults.withCredentials=!1,delete b.defaults.headers.common["X-Requested-With"],b.defaults.headers.common.Accept="application/json",b.defaults.headers.common["Content-Type"]="application/json"}]),angular.module("dockuiApp").controller("MainCtrl",["$scope","$log","$interval","ImageListService","BuildMessageService",function(a,b,c,d,e){a.builds=[],a.teradataImages=[],a.customerImages=[],a.isTeamCityBusy=!0,a.isTeradataBusy=!0,a.isCustomerBusy=!0,a.teradataHostIP=d.getTeradataIP(),a.customerHostIP=d.getCustomerIP();var f,g,h=function(){e.getBuildMessages().success(function(c){e.setLatestBuildMessage(c),b.debug("Got build message",c),a.builds=c}).error(function(c){b.error("GetImages failed",c),a.status="Unable to get image list: "+c.message})},i=function(){d.getTeradataImageList().success(function(c){b.debug("Got images",c),a.teradataImages=c}).error(function(c){b.error("GetImages failed",c),a.status="Unable to get image list: "+c.message})},j=function(){d.getCustomerImageList().success(function(c){b.debug("Got images",c),a.customerImages=c}).error(function(c){b.error("GetImages failed",c),a.status="Unable to get image list: "+c.message})},k=function(){angular.isDefined(f)||(f=c(function(){h(),i(),j()},8e3))},l=function(){angular.isDefined(g)||(g=c(function(){a.isTeradataBusy=!1,a.isCustomerBusy=!1,a.isTeamCityBusy=e.isBusy()},2e3))};a.stopTick=function(){angular.isDefined(f)&&(c.cancel(f),f=void 0)},a.stopTack=function(){angular.isDefined(g)&&(c.cancel(g),g=void 0)},a.$on("$destroy",function(){a.stopTick(),a.stopTack()}),k(),l()}]),angular.module("dockuiApp").factory("ImageListService",["$http","$log",function(a,b){var c="http://localhost:8007/",d=c+"_10.25.191.196:2375/images/json",e=c+"_153.64.104.38:2375/images/json",f=function(){return b.debug("Making AJAX request to",e),a.get(e)},g=function(){return b.debug("Making AJAX request to",d),a.get(d)},h=function(){return"10.25.191.196"},i=function(){return"153.64.10.38"};return{getTeradataImageList:f,getCustomerImageList:g,getTeradataIP:h,getCustomerIP:i}}]),angular.module("dockuiApp").service("BuildMessageService",["$http",function(a){var b="http://localhost:8007",c=[],d=function(){return a.get(b+"/build")},e=function(a){c=a},f=function(){if(!c)return!1;for(var a in c)if("done"!==c[a].status&&"-"!==c[a].status)return!0;return!1};return{getBuildMessages:d,setLatestBuildMessage:e,isBusy:f}}]);