/** @jsx React.DOM */

var React = require('react');

var CampaignsTable = require('./campaigns-table.jsx');
var CreateNewButton = require('../create-new-button.jsx');

var Campaigns = React.createClass({
    render: function () {
        return (
            <div>
                <CreateNewButton url={url_base + '/dashboard/new-campaign'} text="Create new campaign" />
                <div className="row">
                    <CampaignsTable />
                </div>
            </div>
        );
    }
});

React.render(<Campaigns />, document.getElementById('campaigns'));