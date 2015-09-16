/** @jsx React.DOM */

var React = require('react');

var ListForm = require('./components/sub-lists/list-form.jsx');

React.render(<ListForm edit={false} />, document.getElementById('new-sub-list'));