/** @jsx React.DOM */

var React = require('react');

var ReportsTable = require('./components/reports/reports-table.jsx');
var Campaign = require('./entities/campaign.js');

var c = new Campaign();

var Reports = React.createClass({
    getInitialState: function () {
        return {step: '', campaign: {}}
    },
    viewReport: function () {
        this.setState({step: 'view'});
    },
    back: function () {
        this.setState({step: ''});
    },
    render: function () { 
        return (
            <div className="row">
                <ReportsTable viewReport={this.viewReport} />
            </div>
        ); 
    }
});

React.render(<Reports />, document.getElementById('reports'));
