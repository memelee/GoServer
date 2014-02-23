package config

var DBHost = "127.0.0.1"
var DBName = "oj"
var DBLasting = false

var JudgePD = 0  //Pending
var JudgeRJ = 1  //Running & judging
var JudgeAC = 2  //Accepted
var JudgeCE = 3  //Compile Error
var JudgeRE = 4  //Runtime Error
var JudgeWA = 5  //Wrong Answer
var JudgeTLE = 6 //Time Limit Exceeded
var JudgeMLE = 7 //Memory Limit Exceeded
var JudgeOLE = 8 //Output Limit Exceeded

var PrivilegeNA = 0
var PrivilegePU = 1

var EncryptNA = 0
var EncryptPT = 1
var EncryptPW = 2
