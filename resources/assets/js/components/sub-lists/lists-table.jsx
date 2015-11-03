/** @jsx React.DOM */

require('bootpag/lib/jquery.bootpag.min.js');

var React = require('react');
var DeleteButton = require('../delete-button.jsx');
var List = require('../../entities/list.js');
var l = new List();

var getAllLists = function (component) {
    var data = {
        paginate: true,
        per_page: 10,
        page: 1
    };

    l.all(data).done(function (res) {
        component.setState({lists: res});
        $('.pagination').bootpag({
            total: res.last_page,
            page: res.current_page,
            maxVisible: 5
        }).on("page", function (event, num) {
            data.page = num;
            l.all(data).done(function (res) {
                component.setState({lists: res});
                $('.pagination').bootpag({page: res.current_page});
            });
        });
    });
};

var ListRow = React.createClass({
    showList: function () {
        this.props.showList(this.props.data.id);
    },
    editList: function () {
        this.props.editList(this.props.data.id);
    },
    render: function () {
        return (
            <tr>
                <td><a href="#" onClick={this.showList}>{this.props.data.name}</a></td>
                <td>{this.props.data.total_subscribers}</td>
                <td>
                    <a href="#" onClick={this.editList}><span className="glyphicon glyphicon-pencil"></span></a>
                </td>
                <td>
                    <DeleteButton success={this.props.handleDelete} delete={l.delete.bind(this, this.props.data.id)}/>
                </td>
            </tr>
        );
    }
});

var ListsTable = React.createClass({
    getInitialState: function () {
        return {lists: {data: []}};
    },
    componentDidMount: function () {
        getAllLists(this);
    },
    handleDelete: function () {
        getAllLists(this);
    },
    render: function () {
        var rows = function (data) {
            return <ListRow key={data.id} data={data} handleDelete={this.handleDelete} showList={this.props.showList}
                            editList={this.props.editList}/>
        }.bind(this);
        return (
            <div>
                <table className="table table-responsive table-striped table-hover">
                    <thead>
                    <tr>
                        <th>List name</th>
                        <th>Subscribers</th>
                        <th>Edit</th>
                        <th>Delete</th>
                    </tr>
                    </thead>
                    <tbody>
                    {this.state.lists.data.map(rows)}
                    </tbody>
                </table>
                <div className="col-lg-12 pagination text-center"></div>
            </div>
        );
    }
});

module.exports = ListsTable;
