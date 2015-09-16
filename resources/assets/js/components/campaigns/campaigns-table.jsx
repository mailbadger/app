/**
 * Created by filip on 22.7.15.
 */
/** @jsx React.DOM */

require('bootpag/lib/jquery.bootpag.min.js');

var React = require('react');
var DeleteButton = require('../delete-button.jsx');
var Campaign = require('../../entities/campaign.js');
var c = new Campaign();

var getAllCampaigns = function (component) {
    c.all(true, 10, 1).done(function (res) {
        component.setState({campaigns: res});

        $('.pagination').bootpag({
            total: res.last_page,
            page: res.current_page,
            maxVisible: 5
        }).on("page", function (event, num) {
            c.all(true, 10, num).done(function (res) {
                component.setState({campaigns: res});
                $('.pagination').bootpag({page: res.current_page});
            });
        });
    });
};

var CampaignRow = React.createClass({
    editCampaign: function () {
        this.props.editCampaign(this.props.data.id);
    },
    sendCampaign: function () {
        this.props.sendCampaign(this.props.data.id);
    },
    render: function () {
        var edit = (this.props.data.status == 'draft' || this.props.data.status == 'scheduled') ?
            <span> | <a href="#" onClick={this.editCampaign}>Edit</a></span> : null;
        return (
            <tr>
                <td><a href="#" onClick={this.sendCampaign}>{this.props.data.name}</a>{edit}</td>
                <td>{this.props.data.subject}</td>
                <td>{this.props.data.from_name}</td>
                <td>{this.props.data.from_email}</td>
                <td>{this.props.data.recipients}</td>
                <td>{(() => {
                    switch (this.props.data.status) {
                        case "draft":
                            return <span className="label label-default">Draft</span>;
                        case "sent":
                            return <span className="label label-success">Sent</span>;
                        case "sending":
                            return <span className="label label-info">Sending</span>;
                    }
                })()}</td>
                <td>
                    <DeleteButton success={this.props.handleDelete} delete={c.delete.bind(this, this.props.data.id)}/>
                </td>
            </tr>
        );
    }
});

var CampaignsTable = React.createClass({
    getInitialState: function () {
        return {campaigns: {data: []}};
    },
    componentDidMount: function () {
        getAllCampaigns(this);
    },
    handleDelete: function () {
        getAllCampaigns(this);
    },
    render: function () {
        var rows = function (data) {
            return <CampaignRow key={data.id} data={data} handleDelete={this.handleDelete}
                                editCampaign={this.props.editCampaign} sendCampaign={this.props.sendCampaign}/>
        }.bind(this);
        return (
            <div>
                <table className="table table-responsive table-striped table-hover">
                    <thead>
                    <tr>
                        <th>Campaign</th>
                        <th>Subject</th>
                        <th>From name</th>
                        <th>From email</th>
                        <th>Recipients</th>
                        <th>Status</th>
                        <th>Delete</th>
                    </tr>
                    </thead>
                    <tbody>
                    {this.state.campaigns.data.map(rows)}
                    </tbody>
                </table>
                <div className="col-lg-12 pagination text-center"></div>
            </div>
        );
    }
});

module.exports = CampaignsTable;

