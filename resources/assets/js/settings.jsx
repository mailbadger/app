/** @jsx React.DOM */

var React = require('react');
var SettingsForm = require('./components/settings-form.jsx');
var User = require('./entities/user.js');
var u = new User();

u.getSettings().done(function(res) {  
    React.render(<SettingsForm data={res} />, document.getElementById('settings'));
});

