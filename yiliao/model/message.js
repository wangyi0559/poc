/**
 * Created by jimin on 17/5/2.
 */
var config = require('../config');
var BCMessageObj = function(){
    return {
        "peers":  ["localhost:7051"],
        "fcn": "调用方法",
        "args":["参数1", "参数2","参数3","参数4"]
    }
};

// 在区块链上增加病历存储位置信息。或
// 在区块链上增加病历查看记录的信息。
/**
 *
 * @param args 增加的病历编号| 病历在数据库中的位置 | 病历hash指纹 |病历疾病类型 或
 * 查看请求编号| 请求发起者编号| 医院编号 | 医院授权情况 |病人编号 |病人授权情况
 * @returns {BCMessageObj}
 * @constructor
 */
function BCMessageAdd (args) {
    var msg = new BCMessageObj();
    msg.fcn = "add";
    msg.args = args;
    return msg;
}

// 在区块链上查询病历存储位置信息。
/**
 *
 * @param args 病历编号
 * @returns {BCMessageObj}
 * @constructor
 */
function BCMessageQuery(args,callback){
    callback("fcn=query&args="+args);
}

// 转账消息
function BCMessageTrans(args){
    var msg = new BCMessageObj();
    msg.fcn = "transfer";
    msg.args = args;
    return msg;
}

// 在区块链上查询申请记录的信息。
/**
 *
 * @param args 查看请求编号
 * @returns {BCMessageObj}
 * @constructor
 */
function BCMessageQueryApplicationLog(args){
    return "fcn=test&args="+args;
}

// 授权信息
/**
 *
 * @param args 病人/医院编号| 病人/医院授权码
 * @constructor
 */
function BCMessageVerifyAdd(args) {
    var msg = new BCMessageAdd(args);
    return msg;
}
/**
 *
 * @param args 医院编号| 医院授权码| 病人编号|病人授权码
 * @constructor
 */
function BCMessageVerify(args){
    var msg = new BCMessageObj();
    return msg;
}

function BCMessage(peers, fcn, args) {
    this.peers = peers;
    this.fcn = fcn;
    this.args = args;
    if ('undefined' == typeof BCMessage._initialized) {
        BCMessage.prototype.setPeers = function (p) {
            this.peers = p;
        }
        BCMessage.prototype.setFcn = function (f) {
            this.fcn = f;
        }
        BCMessage.prototype.setArgs = function (a) {
            this.args = a;
        }
    }
    BCMessage._initialized = true;
}


exports.BCMessage = BCMessage;
exports.BCMessageObj = BCMessageObj;
exports.BCMessageAdd = BCMessageAdd;
exports.BCMessageQuery = BCMessageQuery;
exports.BCMessageQueryApplicationLog = BCMessageQueryApplicationLog;
exports.BCMessageVerifyAdd = BCMessageVerifyAdd;
exports.BCMessageVerify = BCMessageVerify;
exports.BCMessageTrans = BCMessageTrans;