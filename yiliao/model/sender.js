/**
 * Created by jimin on 17/5/3.
 */
var http=require('http');
var config = require('../config');

//post get token

function PostToken(callback){
    var data = {  
        username: 'Jim',  
        orgName: "org1"  
    };  
    data = require('querystring').stringify(data);  
    var options={
        host:config.host,
        port:config.port,
        path:"/users",
        method:'POST',
        headers:{
            "Content-Type": 'application/x-www-form-urlencoded',  
            "Content-Length": data.length  
        }
    };
     var req = http.request(options, function (serverFeedback) {  
        if (serverFeedback.statusCode == 200) {  
            var body = "";  
            serverFeedback.on('data', function (data) { body += data; })  
                          .on('end', function () { callback(null, body); });  
        }  
        else {  
            callback(err, null) 
        }  
    });  


    req.on('error',function(err){
        console.error(err);
        callback(err, null)
    });
    req.write(data);
    req.end();
}

//发送 http Post 请求

function Poster(reqCc,reToken,reqData, callback){
    var postData=JSON.stringify(reqData);
    var options={
        hostname:config.host,
        port:config.port,
        path:config.path+reqCc,
        method:'POST',
        headers:{
            "Content-Type": 'application/json',
            "authorization": 'Bearer '+reToken,
            'Content-Length':Buffer.byteLength(postData)
        }
    };
    var req=http.request(options, function(res) {
        console.log('Status:',res.statusCode);
        // console.log('headers:',JSON.stringify(res.headers));
        res.setEncoding('utf-8');
        res.on('data',function(chun){
            callback(null, chun);
        });
        res.on('end',function(){
            // console.log('No more data in response.********');
        });
    });
    req.on('error',function(err){
        console.error(err);
        callback(err, null)
    });
    req.write(postData);
    req.end();
}

//发送 http Get 请求

function Getter(reqCc,reToken,reqData, callback){
  
    var options = {
        host: config.host,
        port:config.port,
        path:config.path+reqCc+'?peer=peer1&'+reqData,
        method: 'GET',
        headers: {
            "Content-Type": 'application/json',
            "authorization": 'Bearer '+reToken
        }
    };
    var req = http.request(options, function(res) {
        console.log('Status:',res.statusCode);
        res.setEncoding('utf-8');
        var resData = "";
        res.on("data",function(data){
            resData += data;
        });
        res.on("end", function() {
            callback(null,JSON.parse(resData));
        });
    });
    req.on('error', function (e) { 
    console.log('problem with request: ' + e.message); 
    }); 
   
    req.end();
}

module.exports = {
    Poster:Poster,
    PostToken:PostToken,
    Getter:Getter
};