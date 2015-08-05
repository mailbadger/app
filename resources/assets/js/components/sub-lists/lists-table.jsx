/** @jsx React.DOM */

require('bootpag/lib/jquery.bootpag.min.js');
require('sweetalert');

var React = require('react');
var List = require('../../entities/list.js');
var l = new List();

var DeleteButton = React.createClass({
    handleSubmit: function (e) {
        e.preventDefault();
        swal({
                title: "Are you sure?",
                text: "You will not be able to recover this list!",
                type: "warning",
                showCancelButton: true,
                confirmButtonColor: "#DD6B55",
                confirmButtonText: "Yes, delete it!",
                closeOnConfirm: false
            },
            function () {
                l.delete(this.props.tid)
                    .done(function () {
                        swal({
                            title: "Success",
                            text: "The list was successfully deleted!",
                            type: "success"
                        }, function () {
                            location.reload();
                        });
                    })
                    .fail(function () {
                        swal('Could not delete', 'Could not delete the list. Try again.', 'error');
                    });
            }.bind(this));
    },
    render: function () {
        return (
            <form onSubmit={this.handleSubmit}>
                <input type="hidden" name="_method" value="DELETE"/>
                <button type="submit"><span className="glyphicon glyphicon-trash"></span></button>
            </form>
        );
    }
});

var ListRow = React.createClass({
    editList: function() {
        this.props.editList(this.props.data.id);
    },
    render: function () {
        return (
            <tr>
                <td>{this.props.data.name}</td>
                <td>{this.props.data.total_subscribers}</td>
                <td>
                    <a href="#" onClick={this.editList}><span className="glyphicon glyphicon-pencil"></span></a>
                </td>
                <td>
                    <DeleteButton tid={this.props.data.id}/>
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
        l.all(true, 10, 1).done(function (response) {
            this.setState({lists: response});
            $('.pagination').bootpag({
                total: response.last_page,
                page: response.current_page,
                maxVisible: 5
            }).on("page", function (event, num) {
                l.all(true, 10, num).done(function (response) {
                    this.setState({lists: response});
                    $('.pagination').bootpag({page: response.current_page});
                }.bind(this));
            }.bind(this));
        }.bind(this));
    },
    render: function () {
        var rows = function (data) {
            return <ListRow key={data.id} data={data} editList={this.props.editList}/>
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
