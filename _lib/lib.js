/* */

const bcrypt = require("bcrypt");

function sendResult(req, res, data, status) {
  res.setHeader('Content-Type', 'application/json');
  res.status(status).send(JSON.stringify(data, null, 3));
}

function saltAndHash(pass, saltRounds) {
  return new Promise(function doSaltAndHash(resolve, reject) {
    bcrypt.hash(pass, saltRounds, function (err, res) {
      if (err) {
        resolve(err);
      }
      resolve(res);
    });
  });
}

function checkPassword(userpass, dbpass) {
  return new Promise(function doSaltAndHash(resolve, reject) {
    bcrypt.compare(userpass, dbpass, function compare(err, res) {
      if (err) {
        resolve(err);
      }
      resolve(res);
    });
  });
}

module.exports = {
  sendResult: sendResult,
  saltAndHash: saltAndHash,
  checkPassword: checkPassword
};
