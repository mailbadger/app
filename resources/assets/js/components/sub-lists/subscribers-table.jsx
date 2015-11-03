/** @jsx React.DOM */

require('bootpag/lib/jquery.bootpag.min.js');

var React = require('react');
var DeleteButton = require('../delete-button.jsx');
var List = require('../../entities/list.js');

var l = new List();

var getAllSubscribers = function (component, listId) {
    var data = {
        paginate: true,
        per_page: 10,
        page: 1
    };

    l.allSubscribers(listId, data).done(function (res) {
        component.setState({subscribers: res});

        $('.pagination').bootpag({
            total: res.last_page,
            page: res.current_page,
            maxVisible: 5
        }).on("page", function (event, num) {
            data.page = num;
            l.allSubscribers(listId, data).done(function (res) {
                component.setState({subscribers: res});
                $('.pagination').bootpag({page: res.current_page});
            });
        });
    });
};

var SubscriberRow = React.createClass({
    render: function () {
        return (
            <tr>
                <td>{this.props.data.name}</td>
                <td>{this.props.data.email}</td>
                <td>
                    <DeleteButton success={this.props.handleDelete} delete={l.deleteSubscriber.bind(this, this.props.listId, this.props.data.id)}/>
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
        getAllSubscribers(this, this.props.listId);
    },
    handleDelete: function () {
        getAllSubscribers(this, this.props.listId);
    },
    render: function () {
        var rows = function (data) {
            return <SubscriberRow key={data.id} handleDelete={this.handleDelete} listId={this.props.listId} data={data}/>
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
