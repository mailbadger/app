/** @jsx React.DOM */

var React = require('react');

var TemplatesTable = require('./templates-table.jsx');
var CreateNewButton = require('../create-new-button.jsx');

var Templates = React.createClass({
    render: function () {
        return (
            <div>
                <CreateNewButton url={url_base + '/dashboard/new-template'} text="Create new template" />
                <div className="row">
                    <TemplatesTable />
                </div>
            </div>
        );
    }
});

React.render(<Templates />, document.getElementById('templates'));