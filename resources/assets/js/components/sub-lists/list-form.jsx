

import sweetalert from 'sweetalert';
import React, {Component} from 'react';
import List from '../../entities/list.js';
import ErrorsList from '../errors-list.jsx';

const l = new List();

export default class ListForm extends Component {

    constructor(props) {
        super(props);
        this.state = {
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
            text: "The list was successfully created!",
            type: "success"
        }, () => {
            window.location.href = url_base + '/dashboard/subscribers';
        });
    }

    handleErrors(xhr) {
        this.setState({hasErrors: true, errors: xhr.responseJSON});
    }

    handleSubmit(e) {
        e.preventDefault();
        this.setState({hasErrors: false, errors: []});
        let data;

        if(!this.props.edit) {
            data = {
                name: this.refs.name.getDOMNode().value,
                total_subscribers: 0
            };
            l.create(data).done(this.handleSuccess).fail(this.handleErrors);
        } else {
            data = {
                name: this.refs.name.getDOMNode().value
            };
            l.update(data, this.props.data.id).done(this.handleSuccess).fail(this.handleErrors);
        }
    }

    render() {
        let errors = (this.state.hasErrors) ? <ErrorsList errors={this.state.errors}/> : null;
        let backBtn = (this.props.edit) ? <a href="#" onClick={this.props.back}>Back</a> : null;
        return (
            <form onSubmit={this.handleSubmit}>
                <div className="errors">{errors}</div>
                <div className="col-lg-4">
                    <div className="form-group">
                        <label htmlFor="name">List name:</label>
                        <input type="text" className="form-control" ref="name" name="name" id="name" placeholder="Name"
                               defaultValue={this.props.edit ? this.props.data.name : ''} required/>
                    </div>
                    <button className="col-lg-4 btn btn-default">Save list</button>
                </div>
                <div className="col-lg-12">{backBtn}</div>
            </form>
        );
    }
}
