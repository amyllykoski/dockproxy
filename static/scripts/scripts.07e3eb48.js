"use strict";angular.module("dockuiApp",["ngAnimate","ngCookies","ngResource","ngRoute","ngSanitize","ngTouch"]).config(["$routeProvider","$httpProvider",function(a,b){b.defaults.useXDomain=!0,b.defaults.withCredentials=!1,delete b.defaults.headers.common["X-Requested-With"],b.defaults.headers.common.Accept="application/json",b.defaults.headers.common["Content-Type"]="application/json"}]),angular.module("dockuiApp").controller("MainCtrl",["$scope","$log","$interval","ImageListService","BuildMessageService",function(a,b,c,d,e){a.builds=[],a.teradataImages=[],a.customerImages=[],a.isTeamCityBusy=!0,a.isTeradataBusy=!0,a.isCustomerBusy=!0,a.teradataHostIP=d.getTeradataIP(),a.customerHostIP=d.getCustomerIP();var f,g,h=function(){e.getBuildMessages().success(function(c){e.setLatestBuildMessage(c),b.debug("Got build message",c),a.builds=c}).error(function(c){b.error("GetImages failed",c),a.status="Unable to get image list: "+c.message})},i=function(){d.getTeradataImageList().success(function(c){b.debug("Got images",c),a.teradataImages=c}).error(function(c){b.error("GetImages failed",c),a.status="Unable to get image list: "+c.message})},j=function(){d.getCustomerImageList().success(function(c){b.debug("Got images",c),a.customerImages=c}).error(function(c){b.error("GetImages failed",c),a.status="Unable to get image list: "+c.message})},k=function(){angular.isDefined(f)||(f=c(function(){h(),i(),j()},8e3))},l=function(){angular.isDefined(g)||(g=c(function(){a.isTeradataBusy=!1,a.isCustomerBusy=!1,a.isTeamCityBusy=e.isBusy()},2e3))};a.stopTick=function(){angular.isDefined(f)&&(c.cancel(f),f=void 0)},a.stopTack=function(){angular.isDefined(g)&&(c.cancel(g),g=void 0)},a.$on("$destroy",function(){a.stopTick(),a.stopTack()}),k(),l()}]),angular.module("dockuiApp").factory("ImageListService",["$http","$location","$log",function(a,b,c){var d="http://"+b.host()+":8007/",e=d+"_10.25.191.196:2375/images/json",f=d+"_153.64.104.38:2375/images/json",g=function(){return c.debug("Making AJAX request to",f),a.get(f)},h=function(){return c.debug("Making AJAX request to",e),a.get(e)},i=function(){return"10.25.191.196"},j=function(){return"153.64.10.38"};return{getTeradataImageList:g,getCustomerImageList:h,getTeradataIP:i,getCustomerIP:j}}]),angular.module("dockuiApp").service("BuildMessageService",["$http","$location","$log",function(a,b,c){var d="http://"+b.host()+":8007",e=[],f=function(){return c.debug("Getting build messages from "+d+"/build"),a.get(d+"/build")},g=function(a){e=a},h=function(){if(!e)return!1;for(var a in e)if("done"!==e[a].status&&"-"!==e[a].status)return!0;return!1};return{getBuildMessages:f,setLatestBuildMessage:g,isBusy:h}}]);