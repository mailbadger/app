/** @jsx React.DOM */

require('bootpag/lib/jquery.bootpag.min.js');

var React = require('react');
var DeleteButton = require('../delete-button.jsx');
var Campaign = require('../../entities/campaign.js');
var c = new Campaign();

var getSentCampaigns = function (component) {
    var data = {
        paginate: true,
        per_page: 10,
        page: 1,
        search: 'sent',
        searchFields: 'status:='
    };

    c.all(data).done(function (res) {
        component.setState({campaigns: res});

        $('.pagination').bootpag({
            total: res.last_page,
            page: res.current_page,
            maxVisible: 5
        }).on("page", function (event, num) {
            data.page = num;
            c.all(data).done(function (res) {
                component.setState({campaigns: res});
                $('.pagination').bootpag({page: res.current_page});
            });
        });
    });
};

var ReportRow = React.createClass({
    viewReport: function () {
        this.props.viewReport(this.props.data.id);
    },
    render: function () {
        return (
            <tr>
                <td><a href="#" onClick={this.viewReport}>{this.props.data.name}</a></td>
                <td>{this.props.data.recipients}</td>
                <td>{this.props.data.sent_at}</td>
                <td>{this.props.data.subject}</td>
                <td>{this.props.data.from_name}</td>
                <td>{this.props.data.from_email}</td>
            </tr>
        );
    }
});

var ReportsTable = React.createClass({
    getInitialState: function () {
        return {campaigns: {data: []}};
    },
    componentDidMount: function () {
        getSentCampaigns(this);
    },
    render: function () {
        var rows = function (data) {
            return <ReportRow key={data.id} data={data} viewReport={this.props.viewReport} />
        }.bind(this);

        return (
            <div>
                <table className="table table-responsive table-striped table-hover">
                    <thead>
                    <tr>
                        <th>Campaign</th>
                        <th>Recipients</th>
                        <th>Sent</th>
                        <th>Subject</th>
                        <th>From name</th>
                        <th>From email</th>
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

module.exports = ReportsTable;

