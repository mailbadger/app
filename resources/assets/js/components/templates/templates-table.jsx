

import * as bootpag from 'bootpag/lib/jquery.bootpag.min.js';
import React, {Component} from 'react';
import DeleteButton from '../delete-button.jsx';
import Template from '../../entities/template.js';

const t = new Template();

const getAllTemplates = (component) => {
    let data = {
        paginate: true,
        per_page: 10,
        page: 1
    };

    t.all(data).done((res) => {
        component.setState({templates: res});

        $('.pagination').bootpag({
            total: res.last_page,
            page: res.current_page,
            maxVisible: 5
        }).on("page", (event, num) => {
            data.page = num;
            t.all(data).done((res) => {
                component.setState({templates: res});
                $('.pagination').bootpag({page: res.current_page});
            });
        });
    });
}

class PreviewButton extends Component {
    componentDidMount() {
        $('.preview').magnificPopup({
            type: 'iframe',
            iframe: {
                markup: '<div class="mfp-iframe-scaler">' +
                    '<div class="mfp-close"></div>' +
                        '<iframe style="background: #fff none repeat scroll 0 0 !important" class="mfp-iframe" frameborder="0"></iframe>' +
                            '</div>',
                            patterns: {
                                template: {
                                    index: '',
                                    src: '%id%'
                                }
                            }
            },
            gallery: {
                enabled: true
            }
        });
    }
    
    render() {
        return (
            <a className="preview" href={url_base + '/api/templates/content/' + this.props.tid}><span
                    className="glyphicon glyphicon-eye-open"></span></a>
        )
    }
}

class TemplateRow extends Component {

    constructor(props) {
        super(props);

        this.editTemplate = this.editTemplate.bind(this);
    }

    editTemplate() {
        this.props.editTemplate(this.props.data.id);
    }

    render() {
        return (
            <tr>
                <td>{this.props.data.name}</td>
                <td><a href="#" onClick={this.editTemplate}><span className="glyphicon glyphicon-pencil"></span></a>
                </td>
                <td>
                    <PreviewButton tid={this.props.data.id}/>
                </td>
                <td>
                    <DeleteButton success={this.props.handleDelete} resourceId={this.props.data.id} entity={t}/>
                </td>
            </tr>
        );
    }
}

export default class TemplatesTable extends Component {

    constructor(props) {
        super(props);
        this.state = {
            templates: {
                data: []
            }
        };

        this.handleDelete = this.handleDelete.bind(this);
    }

    componentDidMount() {
        getAllTemplates(this);
    }

    handleDelete() {
        getAllTemplates(this);
    }

    render() {
        let rows = (data) => {
            return <TemplateRow key={data.id} data={data} handleDelete={this.handleDelete}
                editTemplate={this.props.editTemplate}/>
        }

        return (
            <div>
                <table className="table table-responsive table-striped table-hover">
                    <thead>
                        <tr>
                            <th>Template</th>
                            <th>Edit</th>
                            <th>Preview</th>
                            <th>Delete</th>
                        </tr>
                    </thead>
                    <tbody>
                        {this.state.templates.data.map(rows)}
                    </tbody>
                </table>
                <div className="col-lg-12 pagination text-center"></div>
            </div>
        );
    }
}
