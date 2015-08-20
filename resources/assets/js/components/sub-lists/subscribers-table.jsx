/** @jsx React.DOM */

require('bootpag/lib/jquery.bootpag.min.js');

var React = require('react');
var DeleteButton = require('../delete-button.jsx');
var List = require('../../entities/list.js');

var l = new List();

var SubscriberRow = React.createClass({
    render: function () {
        return (
            <tr>
                <td>{this.props.data.name}</td>
                <td>{this.props.data.email}</td>
                <td>
                    <DeleteButton delete={l.deleteSubscriber.bind(this, this.props.listId, this.props.data.id)}/>
                </td>
            </tr>
        );
    }
});

var SubscribersTable = React.createClass({
    getInitialState: function () {
        return {subscribers: {data: []}};
    },
    componentDidMount: function () {
        l.allSubscribers(this.props.listId, true, 10, 1).done(function (response) {
            this.setState({subscribers: response});
            $('.pagination').bootpag({
                total: response.last_page,
                page: response.current_page,
                maxVisible: 5
            }).on("page", function (event, num) {
                l.allSubscribers(this.props.listId, true, 10, num).done(function (response) {
                    this.setState({subscribers: response});
                    $('.pagination').bootpag({page: response.current_page});
                }.bind(this));
            }.bind(this));
        }.bind(this));
    },
    render: function () {
        var rows = function (data) {
            return <SubscriberRow key={data.id} listId={this.props.listId} data={data}/>
        }.bind(this);
        return (
            <div>
                <table className="table table-responsive table-striped table-hover">
                    <thead>
                    <tr>
                        <th>Subscriber name</th>
                        <th>Email</th>
                        <th>Delete</th>
                    </tr>
                    </thead>
                    <tbody>
                    {this.state.subscribers.data.map(rows)}
                    </tbody>
                </table>
                <div className="col-lg-12 pagination text-center"></div>
            </div>
        );
    }
});

module.exports = SubscribersTable;
