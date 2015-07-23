/**
 * Created by filip on 22.7.15.
 */
/** @jsx React.DOM */
require('bootpag/lib/jquery.bootpag.min.js');

var React = require('react');
var Campaign = require('../lib/campaign.js');
var c = new Campaign();

var CampaignsTable = React.createClass({
    getInitialState: function () {
        return {campaigns: {data: []}};
    },
    handleData: function (data) {
        this.setState({campaigns: data});
        $('.pagination').bootpag({
            total: data.total,
            page: data.current_page,
            maxVisible: 5
        }).on("page", function (event, num) {
            c.all(num).done(this.handleData);
        });
    },
    componentDidMount: function () {
        c.all(1).done(this.handleData);
    },
    render: function () {
        var table;
        var columns = function (column) {
            return <th key={column}>{column}</th>;
        };
        var data = function (d) {
            return <tr key={d.id}>
                <td>{d.name}</td>
                <td>156</td>
                <td>{d.status}</td>
                <td>
                    <button className="btn btn-default">Delete</button>
                </td>
            </tr>;
        };
        table = <table className="table table-responsive table-striped table-hover">
            <thead>
            <tr>{this.props.columns.map(columns)}</tr>
            </thead>
            <tbody>
            {this.state.campaigns.data.map(data)}
            </tbody>
        </table>;

        return (
            <div>
                {table}
                <div className="col-lg-12 pagination text-center"></div>
            </div>

        );
    }

});

React.render(<CampaignsTable
    columns={['Campaign', 'Recipients', 'Status', 'Delete']}/>, document.getElementById('campaigns'));

