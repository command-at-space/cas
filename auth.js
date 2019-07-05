/* */

const path = require('path');
const lib = require(path.join(__dirname, '_lib/lib.js'));
const db = require(path.join(__dirname, "db.js"));
const c = require(path.join(__dirname, '_config.js'));

async function login(req, res) {
  console.log('/login');
  const info = {
    isLogged: false,
    error: undefined,
    profile: {
      name: req.body.user,
      hash: "",
      email: "",
      logo: ""
    }
  };
  req.session.info = info;

  const pass = req.body.pass;
  let isValidUser;
  try {
    isValidUser = await db.getAccount(info.profile.name);
  } catch (err) {
    console.error('Error Searching User =>\n', err);
    info.error = "Internal problem related to database";
    lib.sendResult(req, res, info, 500);
    return;
  }
  const isValid = isValidUser[0];
  info.profile = isValidUser[1];
  if (!isValid) {
    console.error("Non existent user " + req.body.user);
    info.error = "Non existent user " + req.body.user;
    lib.sendResult(req, res, info, 403);
    return;
  }
  const match = await lib.checkPassword(pass, info.profile.hash);
  if (!match) {
    console.error(info.profile.name, 'Invalid Password', pass);
    info.error = info.profile.name + ' Invalid Password';
    lib.sendResult(req, res, info, 403);
    return;
  }
  delete info.profile.hash;
  info.isLogged = true;
  lib.sendResult(req, res, info, 200);
}

async function signup(req, res) {
  console.log('/signup');
  const info = {
    created: false,
    error: ""
  };
  if (!req.body.user || !req.body.pass) {
    info.error = "Username and password are required";
    lib.sendResult(req, res, info, 400);
    return;
  }
  const user = {
    name: req.body.user || undefined,
    hash: await lib.saltAndHash(req.body.pass, c.app.saltRounds),
    email: req.body.mail || undefined,
    logo: req.body.logo || undefined
  };
  let existsUser;
  try {
    existsUser = await db.getAccount(user.name);
  } catch (err) {
    info.error = "Internal problem related to database";
    lib.sendResult(req, res, info, 500);
    return;
  }

  if (existsUser[0] || existsUser.name === user.name) {
    info.error = `Username ${user.name} already exists`;
    lib.sendResult(req, res, info, 400);
    return;
  }

  try {
    db.insertNewAccount(user);
    info.created = true;
    lib.sendResult(req, res, info, 200);
  } catch (err) {
    info.error = "Internal problem related to database";
    lib.sendResult(req, res, info, 500);
    return;
  }
}

function resign(req, res) {
  console.log('/resign');
  let info = req.session.info;
  if (typeof info !== "undefined") {
    if (info.isLogged === true) {
      lib.sendResult(req, res, info, 200);
      return;
    }
  }
  info = {};
  info.isLogged = false;
  lib.sendResult(req, res, info, 200);
}
function logout(req, res) {
  console.log('/logout');
  req.session.destroy();
  const info = {
    isLogged: false,
  };
  lib.sendResult(req, res, info, 200);
}

function isUserLogged(req, res, next) {
  const info = req.session.info;
  if (typeof info !== "undefined") {
    if (info.isLogged === true) {
      next();
    }
  } else {
    const info = {
      isLogged: false,
      error: "You are not authorized to view this page",
    };
    lib.sendResult(req, res, info, 200);
  }
}

function isUserNotLogged(req, res, next) {
  const info = req.session.info;
  if (typeof info === "undefined") {
    next();
  } else if (info.isLogged === false) {
    next();
  } else {
    const info = {
      isLogged: true,
      error: "You are already logged",
    };
    lib.sendResult(req, res, info, 200);
  }
}

module.exports = {
  login: login,
  signup: signup,
  resign: resign,
  logout: logout,
  isUserLogged: isUserLogged,
  isUserNotLogged: isUserNotLogged
};
