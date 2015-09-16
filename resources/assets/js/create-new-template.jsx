/** @jsx React.DOM */

var React = require('react');

var TemplateForm = require('./components/templates/template-form.jsx');

React.render(<TemplateForm edit={false} />, document.getElementById('new-template'));