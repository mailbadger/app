/** @jsx React.DOM */

var React = require('react');

var ReportsTable = require('./components/reports/reports-table.jsx');
var Report = require('./components/reports/report.jsx');
var Campaign = require('./entities/campaign.js');

var c = new Campaign();

var Reports = React.createClass({
    getInitialState: function () {
        return {step: '', campaign: {}}
    },
    viewReport: function (id) {
        c.get(id).done(function (res) {
            this.setState({step: 'view', campaign: res});
        }.bind(this)); 
    },
    back: function () {
        this.setState({step: ''});
    },
    render: function () { 
        switch (this.state.step) {
            case 'view':
                return <Report data={this.state.campaign} back={this.back}/>;
            default:
                return (
                    <div className="row">
                        <ReportsTable viewReport={this.viewReport} />
                    </div>
                ); 
        }
    }
});

React.render(<Reports />, document.getElementById('reports'));
