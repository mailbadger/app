
import sweetalert from 'sweetalert';
import React, {Component} from 'react';
import ErrorsList from '../errors-list.jsx';
import List from '../../entities/list.js';

const l = new List();

export default class AddSubscribers extends Component {

    constructor(props) {
        super(props);
        this.state = {
            fields: [],
            hasErrors: false,
            errors: {}
        };

        this.handleSuccess = this.handleSuccess.bind(this);
        this.handleErrors = this.handleErrors.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }
    
    handleSuccess() {
        this.setState({hasErrors: false, errors: []});
        sweetalert({
            title: "Success",
            text: "The list was successfully imported!",
            type: "success"
        }, () => {
            this.props.back();
        });
    }

    handleErrors(xhr) {
        this.setState({hasErrors: true, errors: xhr.responseJSON});
    }

    handleSubmit(e) {
        e.preventDefault();

        l.createSubscribers(this.props.listId, this.refs.subscribers.getDOMNode().files[0])
            .done(this.handleSuccess)
            .fail(this.handleErrors);
    }

    componentDidMount() {
        l.allFields(this.props.listId).done((res) => {
            this.setState({fields: res});
        });
    }

    render() {
        let errors = (this.state.hasErrors) ? <ErrorsList errors={this.state.errors}/> : null;
        let fields = (field) => {
            return <th key={field.id}>{field.name}</th>;
        }
        let fieldVals = (field, key) => {
            return <td key={field.id}>Example field {key}</td>
        }

        return (
            <div>
                <div className="row">
                    <h2>Import via csv/xls file</h2>

                    <p>File format:</p>
                    <ul>
                        <li>Format your CSV the same way as the example below</li>
                        <li>The number of columns in your CSV should be the same as the example below</li>
                        <li>If you want to import more than just name and email, create custom fields first</li>
                    </ul>
                    <div className="col-lg-4">
                        <table className="table table-responsive table-hover">
                            <thead>
                            <tr>
                                <th>Name</th>
                                <th>Email</th>
                                {this.state.fields.map(fields)}
                            </tr>
                            </thead>
                            <tbody>
                            <tr>
                                <td>John Doe</td>
                                <td>john@doe.com</td>
                                {this.state.fields.map(fieldVals)}
                            </tr>
                            <tr>
                                <td>Jane Doe</td>
                                <td>jane@doe.com</td>
                                {this.state.fields.map(fieldVals)}
                            </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
                <div className="row">
                    <form onSubmit={this.handleSubmit} id="upload-form">
                        <div className="errors">{errors}</div>
                        <div className="col-lg-3">
                        <span className="btn btn-success btn-file">
                            Browse<input type="file" ref="subscribers" name="subscribers"
                                         id="subscribers" required/>
                        </span>
                        </div>
                        <button type="submit" className="btn btn-default">Import</button>
                    </form>
                    <div className="col-lg-12" style={{marginTop: '20px'}}>
                        <a href="#" onClick={this.props.back}>Back to list</a>
                    </div>
                </div>
            </div>
        );
    }
}
