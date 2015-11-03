/** @jsx React.DOM */

require('bootpag/lib/jquery.bootpag.min.js');

var React = require('react');
var DeleteButton = require('../delete-button.jsx');
var Template = require('../../entities/template.js');
var t = new Template();

var getAllTemplates = function (component) {
    var data = {
        paginate: true,
        per_page: 10,
        page: 1
    };

    t.all(data).done(function (res) {
        component.setState({templates: res});

        $('.pagination').bootpag({
            total: res.last_page,
            page: res.current_page,
            maxVisible: 5
        }).on("page", function (event, num) {
            data.page = num;
            t.all(data).done(function (res) {
                component.setState({templates: res});
                $('.pagination').bootpag({page: res.current_page});
            });
        });
    });
};

var PreviewButton = React.createClass({
    componentDidMount: function () {
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
    },
    render: function () {
        return (
            <a className="preview" href={url_base + '/api/templates/content/' + this.props.tid}><span
                    className="glyphicon glyphicon-eye-open"></span></a>
        )
    }
});

var TemplateRow = React.createClass({
    editTemplate: function () {
        this.props.editTemplate(this.props.data.id);
    },
    render: function () {
        return (
            <tr>
                <td>{this.props.data.name}</td>
                <td><a href="#" onClick={this.editTemplate}><span className="glyphicon glyphicon-pencil"></span></a>
                </td>
                <td>
                    <PreviewButton tid={this.props.data.id}/>
                </td>
                <td>
                    <DeleteButton success={this.props.handleDelete} delete={t.delete.bind(this, this.props.data.id)}/>
                </td>
            </tr>
        );
    }
});

var TemplatesTable = React.createClass({
    getInitialState: function () {
        return {templates: {data: []}};
    },
    componentDidMount: function () {
        getAllTemplates(this);
    },
    handleDelete: function () {
        getAllTemplates(this);
    },
    render: function () {
        var rows = function (data) {
            return <TemplateRow key={data.id} data={data} handleDelete={this.handleDelete}
                editTemplate={this.props.editTemplate}/>
        }.bind(this);
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
});

module.exports = TemplatesTable;
