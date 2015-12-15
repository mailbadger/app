
import React, {Component} from 'react';
import ErrorsList from '../errors-list.jsx';
import DeleteButton from '../delete-button.jsx';
import Field from '../../entities/field.js';

export default class CustomFields extends Component {

    constructor(props) {
        super(props);
        this.state = {
            fields: {data: []},
            edit: false,
            name: '',
            editField: null,
            hasErrors: false,
            errors: {}
        };

        this.f = new Field(this.props.listId);

        this.handleSuccess = this.handleSuccess.bind(this);
        this.handleErrors = this.handleErrors.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleChange = this.handleChange.bind(this);
        this.handleEdit = this.handleEdit.bind(this);
        this.handleDelete = this.handleDelete.bind(this);
        this.cancelEdit = this.cancelEdit.bind(this);
    }

    handleSuccess() {
        this.setState({edit: false, editField: {name: ''}, hasErrors: false, errors: []});
        sweetalert({
            title: "Success",
            text: "The field was successfully created!",
            type: "success"
        }, () => {
            this.props.back();
        });
    }

    handleErrors(xhr) {
        this.setState({edit: false, editField: {name: ''}, hasErrors: true, errors: xhr.responseJSON});
    }

    handleSubmit(e) {
        e.preventDefault();
        this.setState({hasErrors: false, errors: []});
        let data = {
            name: this.refs.name.getDOMNode().value
        };

        if (!this.state.edit) {
            l.createField(this.props.listId, data).done(this.handleSuccess).fail(this.handleErrors);
        } else {
            l.updateField(this.props.listId, data, this.state.editField).done(this.handleSuccess).fail(this.handleErrors);
        }
    }

    getAllFields() {
        let data = {
            paginate: true,
            per_page: 10,
            page: 1
        };

        f.all(data).done((res) => {
            this.setState({fields: res});

            $('.pagination').bootpag({
                total: res.last_page,
                page: res.current_page,
                maxVisible: 5
            }).on("page", (event, num) => {
                data.page = num;
                f.all(data).done((res) => {
                    this.setState({fields: res});
                    $('.pagination').bootpag({page: res.current_page});
                });
            });
        });
    }

    handleChange(event) {
        this.setState({name: event.target.value});
    }

    handleEdit(name, id) {
        this.setState({edit: true, name: name, editField: id});
    }

    handleDelete() {
        getAllFields();
    }

    cancelEdit() {
        this.setState({edit: false, name: '', editField: null});
    }

    componentDidMount() {
        getAllFields();
    }

    render() {
        let errors = (this.state.hasErrors) ? <ErrorsList errors={this.state.errors}/> : null;
        let rows = function (field) {
            return (
                <tr key={field.id}>
                    <td>{field.name}</td>
                    <td><a href="#" onClick={this.handleEdit.bind(this, field.name, field.id)}><span
                        className="glyphicon glyphicon-pencil"></span></a>
                    </td>
                    <td><DeleteButton success={this.handleDelete} entity={this.f} resourceId={field.id} /></td>
                </tr>
            );
        };
        let cancelBtn = (this.state.edit) ? <input type="button" onClick={this.cancelEdit} value="Cancel"
                                                   className="col-lg-offset-1 col-lg-3 btn btn-default"/> : null;
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
                                <button className="col-lg-3 btn btn-default">
                                    {this.state.edit ? 'Edit field' : 'Save field'}
                                </button>
                                {this.state.edit ? cancelBtn : null}

                                <div className="col-lg-4 pull-right">
                                    <a href="#" onClick={this.props.back}>Back to lists</a>
                                </div>
                            </div>
                        </form>
                    </div>

                    <div className="row" style={{marginTop: '20px'}}>
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
}
