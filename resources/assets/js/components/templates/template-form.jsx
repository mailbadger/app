
import sweetalert from 'sweetalert';
import React, {Component} from 'react';
import Template from '../../entities/template.js';
import ErrorsList from '../errors-list.jsx';

const t = new Template();

export default class TemplateForm extends Component {

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
            text: "The template was successfully created!",
            type: "success"
        }, () => {
            window.location.href = url_base + '/dashboard/templates';
        });
    }

    handleErrors(xhr) {
        this.setState({hasErrors: true, errors: xhr.responseJSON});
    }

    handleSubmit(e) {
        e.preventDefault();
        this.setState({hasErrors: false, errors: []});

        for (instance in CKEDITOR.instances) {
            CKEDITOR.instances[instance].updateElement();
        }

        let data = {
            name: this.refs.name.getDOMNode().value,
            content: CKEDITOR.instances.content.getData()
        };

        if(!this.props.edit) {
            t.create(data).done(this.handleSuccess).fail(this.handleErrors);
        } else {
            t.update(data, this.props.data.id).done(this.handleSuccess).fail(this.handleErrors);
        }
    }

    componentDidMount() {
        CKEDITOR.replace('content', {
            allowedContent: {
                $1: {
                    elements: CKEDITOR.dtd,
                    attributes: true,
                    styles: true,
                    classes: true
                }
            },
            disallowedContent:  'script; *[on*]'
        });
    }

    render() {
        let errors = (this.state.hasErrors) ? <ErrorsList errors={this.state.errors}/> : null;
        let backBtn = (this.props.edit) ? <a href="#" onClick={this.props.back}>Back</a> : null;
        return (
            <form onSubmit={this.handleSubmit}>
                <div className="errors">{errors}</div>
                <div className="col-sm-4">
                    <div className="form-group">
                        <label htmlFor="name">Template name:</label>
                        <input type="text" className="form-control" ref="name" name="name" id="name" placeholder="Name"
                               defaultValue={this.props.edit ? this.props.data.name : ''} required/>
                    </div>
                    <button className="col-sm-4 btn btn-default">Save template</button>
                </div>
                <div className="col-sm-7">
                    <div className="form-group">
                        <label htmlFor="content">Content:</label>
                        <textarea name="content" id="content" defaultValue={this.props.edit ? this.props.data.content : ''}></textarea>
                    </div>
                </div>
                <div className="col-sm-4">{backBtn}</div>
                <div className="col-sm-8">
                    <div className="row col-sm-6">
                        <h3>Essential tags (HTML)</h3>
                        <p>The following tags can only be used in html</p>
                        <h4>Unsubscribe link:</h4>
                        <code>&lt;unsubscribe&gt;Unsubscribe here&lt;/unsubscribe&gt;</code>
                    </div>
                    <div className="row col-sm-6">
                        <h3>Custom tags</h3>
                        <p>You may use tags to personalize the email:</p>
                        <h4>Name:</h4>
                        <code>*|Name|*</code>
                        <h4>Email:</h4>
                        <code>*|Email|*</code>
                        <p>You can also use custom fields eg.</p>
                        <code>*|Country|*</code>
                        <p>To add them go to the subscribers list and click the 'Custom fields' button.</p>
                    </div>
                </div>
            </form>
        );
    }
}
