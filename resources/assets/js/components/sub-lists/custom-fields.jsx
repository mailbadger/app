/** @jsx React.DOM */

var React = require('react');
var ErrorsList = require('../errors-list.jsx');
var DeleteButton = require('../delete-button.jsx');
var List = require('../../entities/list.js');
var l = new List();

var CustomFields = React.createClass({
    getInitialState: function () {
        return {
            fields: {data: []},
            edit: false,
            name: '',
            editField: null,
            hasErrors: false,
            errors: {}
        };
    },
    handleSuccess: function () {
        this.setState({edit: false, editField: {name: ''}, hasErrors: false, errors: []});
        swal({
            title: "Success",
            text: "The field was successfully created!",
            type: "success"
        }, function () {
            this.props.back();
        }.bind(this));
    },
    handleErrors: function (xhr) {
        this.setState({edit: false, editField: {name: ''}, hasErrors: true, errors: xhr.responseJSON});
    },
    handleSubmit: function (e) {
        e.preventDefault();
        this.setState({hasErrors: false, errors: []});
        var data = {
            name: this.refs.name.getDOMNode().value
        };

        if (!this.state.edit) {
            l.createField(this.props.listId, data).done(this.handleSuccess).fail(this.handleErrors);
        } else {
            l.updateField(this.props.listId, data, this.state.editField).done(this.handleSuccess).fail(this.handleErrors);
        }
    },
    handleChange: function (event) {
        this.setState({name: event.target.value});
    },
    handleEdit: function (name, id) {
        this.setState({edit: true, name: name, editField: id});
    },
    componentDidMount: function () {
        l.allFields(this.props.listId, true, 10, 1).done(function (res) {
            this.setState({fields: res});

            $('.pagination').bootpag({
                total: response.last_page,
                page: response.current_page,
                maxVisible: 5
            }).on("page", function (event, num) {
                l.allFields(this.props.listId, true, 10, num).done(function (response) {
                    this.setState({fields: response});
                    $('.pagination').bootpag({page: response.current_page});
                }.bind(this));
            }.bind(this));
        }.bind(this));
    },
    render: function () {
        var errors = (this.state.hasErrors) ? <ErrorsList errors={this.state.errors}/> : null;
        var rows = function (field) {
            return (
                <tr key={field.id}>
                    <td>{field.name}</td>
                    <td><a href="#" onClick={this.handleEdit.bind(this, field.name, field.id)}><span className="glyphicon glyphicon-pencil"></span></a>
                    </td>
                    <td><DeleteButton delete={l.deleteField.bind(this, this.props.listId, field.id)}/></td>
                </tr>
            );
        }.bind(this);
        return (
            <div>
                <div className="row">
                    <h3>Add a field</h3>

                    <div className="row">
                        <form onSubmit={this.handleSubmit}>
                            <div className="errors">{errors}</div>
                            <div className="col-lg-4">
                                <div className="form-group">
                                    <label htmlFor="name">Field name:</label>
                                    <input type="text" className="form-control" ref="name" name="name" id="name"
                                           placeholder="Name" onChange={this.handleChange}
                                           value={this.state.name} required/>
                                </div>
                                <button className="col-lg-4 btn btn-default">
                                    {this.state.edit ? 'Edit field' : 'Save field'}
                                </button>
                                {this.state.edit ? <button className="col-lg-4 btn btn-default">Cancel</button> : null}

                                <div className="col-lg-4 pull-right">
                                    <a href="#" onClick={this.props.back}>Back to lists</a>
                                </div>
                            </div>
                        </form>
                    </div>

                    <div className="row">
                        <h3 className="page-header">Existing fields</h3>
                        <table className="table table-responsive table-hover">
                            <thead>
                            <tr>
                                <th>Name</th>
                                <th>Edit</th>
                                <th>Delete</th>
                            </tr>
                            </thead>
                            <tbody>
                            {this.state.fields.data.map(rows)}
                            </tbody>
                        </table>
                        <div className="col-lg-12 pagination text-center"></div>
                    </div>
                </div>
            </div>
        );
    }
});

module.exports = CustomFields;