package config

const DBHost = "127.0.0.1"
const DBName = "oj"
const DBLasting = false

const CookieExpires = 1800

const ModuleNA = 0 //None
const ModuleP = 1  //Problem
const ModuleC = 2  //Contest
const ModuleE = 3  //Exercise

const JudgeNA = 0  //None
const JudgePD = 1  //Pending
const JudgeRJ = 2  //Running & judging
const JudgeAC = 3  //Accepted
const JudgeCE = 4  //Compile Error
const JudgeRE = 5  //Runtime Error
const JudgeWA = 6  //Wrong Answer
const JudgeTLE = 7 //Time Limit Exceeded
const JudgeMLE = 8 //Memory Limit Exceeded
const JudgeOLE = 9 //Output Limit Exceeded

const LanguageNA = 0   //None
const LanguageC = 1    //C
const LanguageCPP = 2  //C++
const LanguageJAVA = 3 //Java

const PrivilegeNA = 0 //None
const PrivilegePU = 1 //Primary User
const PrivilegeSB = 2 //Source Broswer
const PrivilegeAD = 3 //Admin

const EncryptNA = 0 //None
const EncryptPB = 1 //Public
const EncryptPT = 2 //Private
const EncryptPW = 3 //Password
