/* */

const path = require('path');
const p = require(path.join(__dirname, '_private.js'));

const mysql = require('mysql');
const con = mysql.createConnection({
  host: p.mysql.host,
  user: p.mysql.user,
  password: p.mysql.password,
  database: p.mysql.db,
  port: p.mysql.port,
  //connectTimeout: 20000, // avoid ETIMEDOUT, default 10000
  //acquireTimeout: 20000 // avoid ETIMEDOUT

});

function testDBConnection() {
  console.log('Connecting ......\n');//, con.config);
  con.connect(function (err) {
    if (err) {
      throw new Error("Error connecting to DB =>" + err);
      //console.error('Error connecting to DB => ', err);
    } else {
      console.log('Connection OK');
    }
  });
}

// AUTH METHODS

function insertNewAccount(u) {
  return new Promise(function doInsertNewAccount(resolve, reject) {
    //console.log("Create account");
    let sql = 'INSERT INTO ?? (name, hash, email, logo)';
    sql += ' VALUES (?, ?, ?, ?)';
    const inserts = [p.mysql.tableLogin, u.name, u.hash, u.email, u.logo];
    sql = mysql.format(sql, inserts);
    con.query(sql, function (err, rows) {
      if (err) {
        //console.error('Insert New Account Error =>', err);
        reject(err);
      } else {
        //console.log('New user saved on DB');
        resolve(rows);
      }
    });
  });
}

function getAccount(name) {
  return new Promise(function doGetAccount(resolve, reject) {
    //console.log('get account');
    let sql = 'SELECT * FROM ?? WHERE name=?';
    const inserts = [p.mysql.tableLogin, name];
    sql = mysql.format(sql, inserts);
    con.query(sql, function (err, rows) {
      if (err) {
        //console.error('Error Searching User =>', err);
        reject(err);
        // throw err
      } else {
        if (rows.length === 1) {
          //console.log("USER =>", rows[0].name);
          resolve([true, rows[0]]);
        } else {
          //console.log('USER =>', name, "doesnt exist");
          resolve([false, undefined]);
        }
      }
    });
  });
}

// SESSION METHODS

function saveSession() {
  console.log('save session');
  // INSERT INTO `%s` (username, sessionID) values ('%s', '%s')ON DUPLICATE KEY UPDATE sessionID = '%s'", c.Mysql.Table2, username, sessionID, sessionID
}

function deleteSession() {
  console.log('selete session');
  // DELETE FROM `%s` WHERE username = '%s'", c.Mysql.Table2, username
}

function loadAllSessions() {
  console.log('load all sessions');
  //SELECT * FROM `%s`", c.Mysql.Table2
}

module.exports = {
  testDBConnection: testDBConnection,
  insertNewAccount: insertNewAccount,
  getAccount: getAccount,
  saveSession: saveSession,
  deleteSession: deleteSession,
  loadAllSessions: loadAllSessions
};

