/** @jsx React.DOM */

require('sweetalert');
var React = require('react');

var DeleteButton = React.createClass({
    handleSubmit: function (e) {
        e.preventDefault();
        swal({
                title: "Are you sure?",
                text: "You will not be able to recover this!",
                type: "warning",
                showCancelButton: true,
                confirmButtonColor: "#DD6B55",
                confirmButtonText: "Yes, delete it!",
                closeOnConfirm: false
            },
            function () {
                this.props.delete()
                    .done(function () {
                        swal({
                            title: "Success",
                            text: "The entity was successfully removed!",
                            type: "success"
                        }, function () {
                            location.reload();
                        });
                    })
                    .fail(function () {
                        swal('Could not delete', 'Could not delete the entity. Try again.', 'error');
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

module.exports = DeleteButton;