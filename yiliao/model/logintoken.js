/**
 * Created by wangyi on 17/8/14.
 */

var mongoose = require('mongoose');
var Schema = mongoose.Schema;

/*
 * _sdk token save*/
var Logintoken = new Schema({
    success: Boolean,
    secret: String,
    message: String,
    token: String
});

module.exports = mongoose.model('logintoken',Logintoken);