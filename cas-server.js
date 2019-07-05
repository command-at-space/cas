/* */

const express = require('express');
const app = express();

const session = require('express-session');
const path = require('path');
const c = require(path.join(__dirname, '_config.js'));
const db = require(path.join(__dirname, "db.js"));
const auth = require(path.join(__dirname, "auth.js"));

if (c.app.mode === 'dev') {
  c.app.port = 3000;
}

app.disable('x-powered-by');
app.use(express.json()); // to support JSON-encoded bodies
app.use(express.urlencoded({ // to support URL-encoded bodies
  extended: true
}));
app.set('trust proxy', 1); // trust first proxy
app.use(session({
  secret: "process.env.secretSession",
  name: 'alpha',
  saveUninitialized: false,
  resave: false,
  cookie: {
    httpOnly: false, // true to production
    //secure: true // comment for localhost, only works for https
    secure: false,
  },
}));

app.use((req, res, next) => {
  next();
});

app.get("/secret", auth.isUserLogged, function (req, res) {
  res.send("SECRET KK");
});

app.post('/login', auth.isUserNotLogged, function (req, res) {
  auth.login(req, res);
});

app.post('/signup', auth.isUserNotLogged, function (req, res) {
  auth.signup(req, res);
});

app.post("/resign", function (req, res) {
  auth.resign(req, res);
});

app.post('/logout', auth.isUserLogged, function (req, res) {
  auth.logout(req, res);
});

app.use(function (req, res, next) {
  const err = new Error('Unavailable Endpoint... ' + req.path);
  err.statusCode = 403;
  next(err);
  //throw new Error('oops, error thrown!');
});

app.use(function (err, req, res, next) {
  //console.error(err.message);
  if (!err.statusCode) {
    err.statusCode = 500;
  }
  const text = { error: err.message };
  res.status(err.statusCode).send(text);
});

app.listen(c.app.port, function () {
  const time = new Date().toUTCString().split(',')[1];
  console.log('Express server on port ' + c.app.port + ' - ' + time);
  initApp();
});

function initApp() {
  db.testDBConnection();
}

module.exports = {};
