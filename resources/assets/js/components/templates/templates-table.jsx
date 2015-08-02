/** @jsx React.DOM */

require('bootpag/lib/jquery.bootpag.min.js');
require('sweetalert');

var React = require('react');
var Template = require('../../entities/template.js');
var t = new Template();

var DeleteButton = React.createClass({
    handleSubmit: function (e) {
        e.preventDefault();
        swal({
                title: "Are you sure?",
                text: "You will not be able to recover this template!",
                type: "warning",
                showCancelButton: true,
                confirmButtonColor: "#DD6B55",
                confirmButtonText: "Yes, delete it!",
                closeOnConfirm: false
            },
            function () {
                t.delete(this.props.tid)
                    .done(function () {
                        swal({
                            title: "Success",
                            text: "The template was successfully deleted!",
                            type: "success"
                        }, function () {
                            location.reload();
                        });
                    })
                    .fail(function () {
                        swal('Could not delete', 'Check if the template belongs to a campaign.', 'error');
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
    editTemplate: function() {
        this.props.editTemplate(this.props.data.id);
    },
    render: function () {
        return (
            <tr>
                <td>{this.props.data.name}</td>
                <td><a href="#" onClick={this.editTemplate}>Edit</a></td>
                <td>
                    <PreviewButton tid={this.props.data.id}/>
                </td>
                <td>
                    <DeleteButton tid={this.props.data.id}/>
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
        t.all(true, 10, 1).done(function (response) {
            this.setState({templates: response});

            $('.pagination').bootpag({
                total: response.last_page,
                page: response.current_page,
                maxVisible: 5
            }).on("page", function (event, num) {
                t.all(true, 10, num).done(function (response) {
                    this.setState({templates: response});
                    $('.pagination').bootpag({page: response.current_page});
                }.bind(this));
            }.bind(this));
        }.bind(this));
    },
    render: function () {
        var rows = function (data) {
            return <TemplateRow key={data.id} data={data} editTemplate={this.props.editTemplate}/>
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
